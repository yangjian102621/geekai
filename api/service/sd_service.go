package service

import (
	"chatplus/core/types"
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

// SD 绘画服务

const SdRunningJobKey = "StableDiffusion_Running_Job"

type SdService struct {
	config    types.ChatPlusExtConfig
	client    *req.Client
	taskQueue *store.RedisQueue
	redis     *redis.Client
	db        *gorm.DB
}

func NewSdService(appConfig *types.AppConfig, client *redis.Client, db *gorm.DB) *SdService {
	return &SdService{
		config:    appConfig.ExtConfig,
		redis:     client,
		db:        db,
		taskQueue: store.NewRedisQueue("stable_diffusion_task_queue", client),
		client:    req.C().SetTimeout(30 * time.Second)}
}

func (s *SdService) Run() {
	logger.Info("Starting StableDiffusion job consumer.")
	ctx := context.Background()
	for {
		_, err := s.redis.Get(ctx, SdRunningJobKey).Result()
		if err == nil { // 队列串行执行
			time.Sleep(time.Second * 3)
			continue
		}
		var task types.SdTask
		err = s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}
		logger.Infof("Consuming Task: %+v", task)
		err = s.txt2img(task.Params)
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

func (s *SdService) PushTask(task types.SdTask) {
	logger.Infof("add a new MidJourney Task: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *SdService) txt2img(params types.SdParams) error {
	logger.Infof("SD 绘画参数：%+v", params)
	url := fmt.Sprintf("%s/api/mj/image", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}
