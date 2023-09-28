package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service/sd"
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
	"time"
)

type SdJobHandler struct {
	BaseHandler
	redis   *redis.Client
	db      *gorm.DB
	service *sd.Service
}

func NewSdJobHandler(app *core.AppServer, redisCli *redis.Client, db *gorm.DB, service *sd.Service) *SdJobHandler {
	h := SdJobHandler{
		redis:   redisCli,
		db:      db,
		service: service,
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
	h.service.Clients.Put(sessionId, client)
	logger.Infof("New websocket connected, IP: %s", c.ClientIP())
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
	if !h.App.Config.SdConfig.Enabled {
		resp.ERROR(c, "Stable Diffusion service is disabled")
		return
	}

	if !h.checkLimits(c) {
		return
	}

	var data struct {
		SessionId string `json:"session_id"`
		types.SdTaskParams
	}
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Width <= 0 {
		data.Width = 512
	}
	if data.Height <= 0 {
		data.Height = 512
	}
	if data.CfgScale <= 0 {
		data.CfgScale = 7
	}
	if data.Seed == 0 {
		data.Seed = -1
	}
	if data.Steps <= 0 {
		data.Steps = 20
	}
	if data.Sampler == "" {
		data.Sampler = "Euler a"
	}
	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	params := types.SdTaskParams{
		TaskId:         fmt.Sprintf("task(%s)", utils.RandString(15)),
		Prompt:         data.Prompt,
		NegativePrompt: data.NegativePrompt,
		Steps:          data.Steps,
		Sampler:        data.Sampler,
		FaceFix:        data.FaceFix,
		CfgScale:       data.CfgScale,
		Seed:           data.Seed,
		Height:         data.Height,
		Width:          data.Width,
		HdFix:          data.HdFix,
		HdRedrawRate:   data.HdRedrawRate,
		HdScale:        data.HdScale,
		HdScaleAlg:     data.HdScaleAlg,
		HdSteps:        data.HdSteps,
	}
	job := model.SdJob{
		UserId:    userId,
		Type:      types.TaskImage.String(),
		TaskId:    params.TaskId,
		Params:    utils.JsonEncode(params),
		Prompt:    data.Prompt,
		Progress:  0,
		Started:   false,
		CreatedAt: time.Now(),
	}
	res := h.db.Create(&job)
	if res.Error != nil {
		resp.ERROR(c, "error with save job: "+res.Error.Error())
		return
	}

	h.service.PushTask(types.SdTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Src:       types.TaskSrcImg,
		Type:      types.TaskImage,
		Prompt:    data.Prompt,
		Params:    params,
		UserId:    userId,
	})
	var jobVo vo.SdJob
	err := utils.CopyObject(job, &jobVo)
	if err == nil {
		// 推送任务到前端
		client := h.service.Clients.Get(data.SessionId)
		if client != nil {
			utils.ReplyChunkMessage(client, jobVo)
		}
	}
	resp.SUCCESS(c)
}

// JobList 获取 MJ 任务列表
func (h *SdJobHandler) JobList(c *gin.Context) {
	status := h.GetInt(c, "status", 0)
	var items []model.SdJob
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

	var jobs = make([]vo.SdJob, 0)
	for _, item := range items {
		var job vo.SdJob
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
