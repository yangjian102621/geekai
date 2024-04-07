package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/oss"
	"chatplus/service/sd"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SdJobHandler struct {
	BaseHandler
	redis     *redis.Client
	pool      *sd.ServicePool
	uploader  *oss.UploaderManager
	snowflake *service.Snowflake
	leveldb   *store.LevelDB
}

func NewSdJobHandler(app *core.AppServer, db *gorm.DB, pool *sd.ServicePool, manager *oss.UploaderManager, snowflake *service.Snowflake, levelDB *store.LevelDB) *SdJobHandler {
	return &SdJobHandler{
		pool:      pool,
		uploader:  manager,
		snowflake: snowflake,
		leveldb:   levelDB,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *SdJobHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return
	}

	userId := h.GetInt(c, "user_id", 0)
	if userId == 0 {
		logger.Info("Invalid user ID")
		c.Abort()
		return
	}

	client := types.NewWsClient(ws)
	h.pool.Clients.Put(uint(userId), client)
	logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
}

func (h *SdJobHandler) checkLimits(c *gin.Context) bool {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if !h.pool.HasAvailableService() {
		resp.ERROR(c, "Stable-Diffusion 池子中没有没有可用的服务！")
		return false
	}

	if user.Power < h.App.SysConfig.SdPower {
		resp.ERROR(c, "当前用户剩余算力不足以完成本次绘画！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *SdJobHandler) Image(c *gin.Context) {
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
	taskId, err := h.snowflake.Next(true)
	if err != nil {
		resp.ERROR(c, "error with generate task id: "+err.Error())
		return
	}
	params := types.SdTaskParams{
		TaskId:       taskId,
		Prompt:       data.Prompt,
		NegPrompt:    data.NegPrompt,
		Steps:        data.Steps,
		Sampler:      data.Sampler,
		FaceFix:      data.FaceFix,
		CfgScale:     data.CfgScale,
		Seed:         data.Seed,
		Height:       data.Height,
		Width:        data.Width,
		HdFix:        data.HdFix,
		HdRedrawRate: data.HdRedrawRate,
		HdScale:      data.HdScale,
		HdScaleAlg:   data.HdScaleAlg,
		HdSteps:      data.HdSteps,
	}

	job := model.SdJob{
		UserId:    userId,
		Type:      types.TaskImage.String(),
		TaskId:    params.TaskId,
		Params:    utils.JsonEncode(params),
		Prompt:    data.Prompt,
		Progress:  0,
		Power:     h.App.SysConfig.SdPower,
		CreatedAt: time.Now(),
	}
	res := h.DB.Create(&job)
	if res.Error != nil {
		resp.ERROR(c, "error with save job: "+res.Error.Error())
		return
	}

	h.pool.PushTask(types.SdTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Type:      types.TaskImage,
		Params:    params,
		UserId:    userId,
	})

	client := h.pool.Clients.Get(uint(job.UserId))
	if client != nil {
		_ = client.Send([]byte("Task Updated"))
	}

	// update user's power
	tx := h.DB.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("power", gorm.Expr("power - ?", job.Power))
	// 记录算力变化日志
	if tx.Error == nil && tx.RowsAffected > 0 {
		user, _ := h.GetLoginUser(c)
		h.DB.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      types.PowerConsume,
			Amount:    job.Power,
			Balance:   user.Power - job.Power,
			Mark:      types.PowerSub,
			Model:     "stable-diffusion",
			Remark:    fmt.Sprintf("绘图操作，任务ID：%s", job.TaskId),
			CreatedAt: time.Now(),
		})
	}

	resp.SUCCESS(c)
}

// ImgWall 照片墙
func (h *SdJobHandler) ImgWall(c *gin.Context) {
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)
	err, jobs := h.getData(true, 0, page, pageSize, true)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, jobs)
}

// JobList 获取 SD 任务列表
func (h *SdJobHandler) JobList(c *gin.Context) {
	status := h.GetBool(c, "status")
	userId := h.GetLoginUserId(c)
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)
	publish := h.GetBool(c, "publish")

	err, jobs := h.getData(status, userId, page, pageSize, publish)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, jobs)
}

// JobList 获取 MJ 任务列表
func (h *SdJobHandler) getData(finish bool, userId uint, page int, pageSize int, publish bool) (error, []vo.SdJob) {

	session := h.DB.Session(&gorm.Session{})
	if finish {
		session = session.Where("progress = ?", 100).Order("id DESC")
	} else {
		session = session.Where("progress < ?", 100).Order("id ASC")
	}
	if userId > 0 {
		session = session.Where("user_id = ?", userId)
	}
	if publish {
		session = session.Where("publish", publish)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}

	var items []model.SdJob
	res := session.Find(&items)
	if res.Error != nil {
		return res.Error, nil
	}

	var jobs = make([]vo.SdJob, 0)
	for _, item := range items {
		var job vo.SdJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}

		if item.Progress < 100 {
			// 从 leveldb 中获取图片预览数据
			imageData, err := h.leveldb.Get(item.TaskId)
			if err == nil {
				job.ImgURL = "data:image/png;base64," + string(imageData)
			}
		}
		jobs = append(jobs, job)
	}

	return nil, jobs
}

// Remove remove task image
func (h *SdJobHandler) Remove(c *gin.Context) {
	var data struct {
		Id     uint   `json:"id"`
		UserId uint   `json:"user_id"`
		ImgURL string `json:"img_url"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// remove job recode
	res := h.DB.Delete(&model.SdJob{Id: data.Id})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	// remove image
	err := h.uploader.GetUploadHandler().Delete(data.ImgURL)
	if err != nil {
		logger.Error("remove image failed: ", err)
	}

	client := h.pool.Clients.Get(data.UserId)
	if client != nil {
		_ = client.Send([]byte("Task Updated"))
	}

	resp.SUCCESS(c)
}

// Publish 发布/取消发布图片到画廊显示
func (h *SdJobHandler) Publish(c *gin.Context) {
	var data struct {
		Id     uint `json:"id"`
		Action bool `json:"action"` // 发布动作，true => 发布，false => 取消分享
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.SdJob{Id: data.Id}).UpdateColumn("publish", true)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c)
}
