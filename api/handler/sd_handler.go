package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/service/sd"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SdJobHandler struct {
	BaseHandler
	redis    *redis.Client
	db       *gorm.DB
	pool     *sd.ServicePool
	uploader *oss.UploaderManager
}

func NewSdJobHandler(app *core.AppServer, db *gorm.DB, pool *sd.ServicePool, manager *oss.UploaderManager) *SdJobHandler {
	h := SdJobHandler{
		db:       db,
		pool:     pool,
		uploader: manager,
	}
	h.App = app
	return &h
}

func (h *SdJobHandler) checkLimits(c *gin.Context) bool {
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if !h.pool.HasAvailableService() {
		resp.ERROR(c, "Stable-Diffusion 池子中没有没有可用的服务！")
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
		CreatedAt: time.Now(),
	}
	res := h.db.Create(&job)
	if res.Error != nil {
		resp.ERROR(c, "error with save job: "+res.Error.Error())
		return
	}

	h.pool.PushTask(types.SdTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Type:      types.TaskImage,
		Prompt:    data.Prompt,
		Params:    params,
		UserId:    userId,
	})

	// update user's img calls
	h.db.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls - ?", 1))

	resp.SUCCESS(c)
}

// JobList 获取 stable diffusion 任务列表
func (h *SdJobHandler) JobList(c *gin.Context) {
	status := h.GetInt(c, "status", 0)
	userId := h.GetInt(c, "user_id", 0)
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)

	session := h.db.Session(&gorm.Session{})
	if status == 1 {
		session = session.Where("progress = ?", 100).Order("id DESC")
	} else {
		session = session.Where("progress < ?", 100).Order("id ASC")
	}
	if userId > 0 {
		session = session.Where("user_id = ?", userId)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}

	var items []model.SdJob
	res := session.Find(&items)
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

		if job.Progress == -1 {
			h.db.Delete(&model.SdJob{Id: job.Id})
		}

		if item.Progress < 100 {
			// 5 分钟还没完成的任务直接删除
			if time.Now().Sub(item.CreatedAt) > time.Minute*5 {
				h.db.Delete(&item)
				// 退回绘图次数
				h.db.Model(&model.User{}).Where("id = ?", item.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls + ?", 1))
				continue
			}
			// 正在运行中任务使用代理访问图片
			image, err := utils.DownloadImage(item.ImgURL, "")
			if err == nil {
				job.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
			}
		}
		jobs = append(jobs, job)
	}
	resp.SUCCESS(c, jobs)
}

// Remove remove task image
func (h *SdJobHandler) Remove(c *gin.Context) {
	var data struct {
		Id     uint   `json:"id"`
		ImgURL string `json:"img_url"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// remove job recode
	res := h.db.Delete(&model.SdJob{Id: data.Id})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	// remove image
	err := h.uploader.GetUploadHandler().Delete(data.ImgURL)
	if err != nil {
		logger.Error("remove image failed: ", err)
	}

	resp.SUCCESS(c)
}
