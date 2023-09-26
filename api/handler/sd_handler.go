package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SdJobHandler struct {
	BaseHandler
	redis           *redis.Client
	db              *gorm.DB
	mjService       *service.MjService
	uploaderManager *oss.UploaderManager
	lock            sync.Mutex
	clients         *types.LMap[string, *types.WsClient]
}

func NewSdJobHandler(
	app *core.AppServer,
	client *redis.Client,
	db *gorm.DB,
	manager *oss.UploaderManager,
	mjService *service.MjService) *MidJourneyHandler {
	h := MidJourneyHandler{
		redis:           client,
		db:              db,
		uploaderManager: manager,
		lock:            sync.Mutex{},
		mjService:       mjService,
		clients:         types.NewLMap[string, *types.WsClient](),
	}
	h.App = app
	return &h
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *SdJobHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	sessionId := c.Query("session_id")
	client := types.NewWsClient(ws)
	// 删除旧的连接
	h.clients.Delete(sessionId)
	h.clients.Put(sessionId, client)
	logger.Infof("New websocket connected, IP: %s", c.ClientIP())
}

type sdNotifyData struct {
	TaskId    string
	ImageName string
	ImageData string
	Progress  int
	Seed      string
	Success   bool
	Message   string
}

func (h *SdJobHandler) Notify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != h.App.Config.ExtConfig.Token {
		resp.NotAuth(c)
		return
	}
	var data sdNotifyData
	if err := c.ShouldBindJSON(&data); err != nil || data.TaskId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	logger.Debugf("收到 MidJourney 回调请求：%+v", data)

	h.lock.Lock()
	defer h.lock.Unlock()

	err, finished := h.notifyHandler(c, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 解除任务锁定
	if finished && (data.Progress == 100) {
		h.redis.Del(c, service.MjRunningJobKey)
	}
	resp.SUCCESS(c)

}

func (h *SdJobHandler) notifyHandler(c *gin.Context, data sdNotifyData) (error, bool) {
	taskString, err := h.redis.Get(c, service.MjRunningJobKey).Result()
	if err != nil { // 过期任务，丢弃
		logger.Warn("任务已过期：", err)
		return nil, true
	}

	var task types.SdTask
	err = utils.JsonDecode(taskString, &task)
	if err != nil { // 非标准任务，丢弃
		logger.Warn("任务解析失败：", err)
		return nil, false
	}

	var job model.SdJob
	res := h.db.Where("id = ?", task.Id).First(&job)
	if res.Error != nil {
		logger.Warn("非法任务：", res.Error)
		return nil, false
	}
	job.Params = utils.JsonEncode(task.Params)
	job.ReferenceId = data.ImageData
	job.Progress = data.Progress
	job.Prompt = data.Prompt
	job.Hash = data.Image.Hash

	// 任务完成，将最终的图片下载下来
	if data.Progress == 100 {
		imgURL, err := h.uploaderManager.GetUploadHandler().PutImg(data.Image.URL)
		if err != nil {
			logger.Error("error with download img: ", err.Error())
			return err, false
		}
		job.ImgURL = imgURL
	} else {
		// 临时图片直接保存，访问的时候使用代理进行转发
		job.ImgURL = data.Image.URL
	}
	res = h.db.Updates(&job)
	if res.Error != nil {
		logger.Error("error with update job: ", res.Error)
		return res.Error, false
	}

	var jobVo vo.MidJourneyJob
	err := utils.CopyObject(job, &jobVo)
	if err == nil {
		if data.Progress < 100 {
			image, err := utils.DownloadImage(jobVo.ImgURL, h.App.Config.ProxyURL)
			if err == nil {
				jobVo.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
			}
		}

		// 推送任务到前端
		client := h.clients.Get(task.SessionId)
		if client != nil {
			utils.ReplyChunkMessage(client, jobVo)
		}
	}

	// 更新用户剩余绘图次数
	if data.Progress == 100 {
		h.db.Model(&model.User{}).Where("id = ?", task.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls - ?", 1))
	}

	return nil, true
}

func (h *SdJobHandler) checkLimits(c *gin.Context) bool {
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if user.ImgCalls <= 0 {
		resp.ERROR(c, "您的绘图次数不足，请联系管理员充值！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *SdJobHandler) Image(c *gin.Context) {
	var data struct {
		SessionId string  `json:"session_id"`
		Prompt    string  `json:"prompt"`
		Rate      string  `json:"rate"`
		Model     string  `json:"model"`
		Chaos     int     `json:"chaos"`
		Raw       bool    `json:"raw"`
		Seed      int64   `json:"seed"`
		Stylize   int     `json:"stylize"`
		Img       string  `json:"img"`
		Weight    float32 `json:"weight"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if !h.checkLimits(c) {
		return
	}

	var prompt = data.Prompt
	if data.Rate != "" && !strings.Contains(prompt, "--ar") {
		prompt += " --ar " + data.Rate
	}
	if data.Seed > 0 && !strings.Contains(prompt, "--seed") {
		prompt += fmt.Sprintf(" --seed %d", data.Seed)
	}
	if data.Stylize > 0 && !strings.Contains(prompt, "--s") && !strings.Contains(prompt, "--stylize") {
		prompt += fmt.Sprintf(" --s %d", data.Stylize)
	}
	if data.Chaos > 0 && !strings.Contains(prompt, "--c") && !strings.Contains(prompt, "--chaos") {
		prompt += fmt.Sprintf(" --c %d", data.Chaos)
	}
	if data.Img != "" {
		prompt = fmt.Sprintf("%s %s", data.Img, prompt)
		if data.Weight > 0 {
			prompt += fmt.Sprintf(" --iw %f", data.Weight)
		}
	}
	if data.Raw {
		prompt += " --style raw"
	}
	if data.Model != "" && !strings.Contains(prompt, "--v") && !strings.Contains(prompt, "--niji") {
		prompt += data.Model
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	job := model.MidJourneyJob{
		Type:      service.Image.String(),
		UserId:    userId,
		Progress:  0,
		Prompt:    prompt,
		CreatedAt: time.Now(),
	}
	if res := h.db.Create(&job); res.Error != nil {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	h.mjService.PushTask(service.MjTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Src:       service.TaskSrcImg,
		Type:      service.Image,
		Prompt:    prompt,
		UserId:    userId,
	})

	var jobVo vo.MidJourneyJob
	err := utils.CopyObject(job, &jobVo)
	if err == nil {
		// 推送任务到前端
		client := h.clients.Get(data.SessionId)
		if client != nil {
			utils.ReplyChunkMessage(client, jobVo)
		}
	}
	resp.SUCCESS(c)
}

// JobList 获取 MJ 任务列表
func (h *SdJobHandler) JobList(c *gin.Context) {
	status := h.GetInt(c, "status", 0)
	var items []model.MidJourneyJob
	var res *gorm.DB
	userId, _ := c.Get(types.LoginUserID)
	if status == 1 {
		res = h.db.Where("user_id = ? AND progress = 100", userId).Order("id DESC").Find(&items)
	} else {
		res = h.db.Where("user_id = ? AND progress < 100", userId).Order("id ASC").Find(&items)
	}
	if res.Error != nil {
		resp.ERROR(c, types.NoData)
		return
	}

	var jobs = make([]vo.MidJourneyJob, 0)
	for _, item := range items {
		var job vo.MidJourneyJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}
		if item.Progress < 100 {
			// 30 分钟还没完成的任务直接删除
			if time.Now().Sub(item.CreatedAt) > time.Minute*30 {
				h.db.Delete(&item)
				continue
			}
			if item.ImgURL != "" { // 正在运行中任务使用代理访问图片
				image, err := utils.DownloadImage(item.ImgURL, h.App.Config.ProxyURL)
				if err == nil {
					job.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}
		}
		jobs = append(jobs, job)
	}
	resp.SUCCESS(c, jobs)
}
