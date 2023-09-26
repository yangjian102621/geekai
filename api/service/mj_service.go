package service

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
	"time"
)

var logger = logger2.GetLogger()

// MJ 绘画服务

const MjRunningJobKey = "MidJourney_Running_Job"

type MjService struct {
	config    types.ChatPlusExtConfig
	client    *req.Client
	taskQueue *store.RedisQueue
	redis     *redis.Client
	db        *gorm.DB
}

func NewMjService(appConfig *types.AppConfig, client *redis.Client, db *gorm.DB) *MjService {
	return &MjService{
		config:    appConfig.ExtConfig,
		redis:     client,
		db:        db,
		taskQueue: store.NewRedisQueue("midjourney_task_queue", client),
		client:    req.C().SetTimeout(30 * time.Second)}
}

func (s *MjService) Run() {
	logger.Info("Starting MidJourney job consumer.")
	ctx := context.Background()
	for {
		_, err := s.redis.Get(ctx, MjRunningJobKey).Result()
		if err == nil { // 队列串行执行
			time.Sleep(time.Second * 3)
			continue
		}
		var task types.MjTask
		err = s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}
		logger.Infof("Consuming Task: %+v", task)
		switch task.Type {
		case types.TaskImage:
			err = s.image(task.Prompt)
			break
		case types.TaskUpscale:
			err = s.upscale(MjUpscaleReq{
				Index:       task.Index,
				MessageId:   task.MessageId,
				MessageHash: task.MessageHash,
			})
			break
		case types.TaskVariation:
			err = s.variation(MjVariationReq{
				Index:       task.Index,
				MessageId:   task.MessageId,
				MessageHash: task.MessageHash,
			})
		}
		if err != nil {
			logger.Error("绘画任务执行失败：", err)
			if task.RetryCount <= 5 {
				s.taskQueue.RPush(task)
			}
			task.RetryCount += 1
			time.Sleep(time.Second * 3)
			continue
		}

		// 更新任务的执行状态
		s.db.Model(&model.MidJourneyJob{}).Where("id = ?", task.Id).UpdateColumn("started", true)
		// 锁定任务执行通道，直到任务超时（5分钟）
		s.redis.Set(ctx, MjRunningJobKey, utils.JsonEncode(task), time.Minute*5)
	}
}

func (s *MjService) PushTask(task types.MjTask) {
	logger.Infof("add a new MidJourney Task: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *MjService) image(prompt string) error {
	logger.Infof("MJ 绘画参数：%+v", prompt)
	body := map[string]string{"prompt": prompt}
	url := fmt.Sprintf("%s/api/mj/image", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}

type MjUpscaleReq struct {
	Index       int32  `json:"index"`
	MessageId   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

func (s *MjService) upscale(upReq MjUpscaleReq) error {
	url := fmt.Sprintf("%s/api/mj/upscale", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(upReq).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}

type MjVariationReq struct {
	Index       int32  `json:"index"`
	MessageId   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

func (s *MjService) variation(upReq MjVariationReq) error {
	url := fmt.Sprintf("%s/api/mj/variation", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(upReq).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}
