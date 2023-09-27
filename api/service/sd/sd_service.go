package sd

import (
	"chatplus/core/types"
	"chatplus/service/mj"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

// SD 绘画服务

const RunningJobKey = "StableDiffusion_Running_Job"

type Service struct {
	taskQueue *store.RedisQueue
	redis     *redis.Client
	db        *gorm.DB
	Client    *Client
}

func NewService(redisCli *redis.Client, db *gorm.DB, client *Client) *Service {
	return &Service{
		redis:     redisCli,
		db:        db,
		Client:    client,
		taskQueue: store.NewRedisQueue("stable_diffusion_task_queue", redisCli),
	}
}

func (s *Service) Run() {
	logger.Info("Starting StableDiffusion job consumer.")
	ctx := context.Background()
	for {
		_, err := s.redis.Get(ctx, RunningJobKey).Result()
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
		err = s.Client.Txt2Img(task.Params)
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
		s.redis.Set(ctx, mj.RunningJobKey, utils.JsonEncode(task), time.Minute*5)
	}
}

func (s *Service) PushTask(task types.SdTask) {
	logger.Infof("add a new MidJourney Task: %+v", task)
	s.taskQueue.RPush(task)
}
