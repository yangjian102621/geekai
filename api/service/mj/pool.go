package mj

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ServicePool Mj service pool
type ServicePool struct {
	services  []Service
	taskQueue *store.RedisQueue
}

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, appConfig *types.AppConfig) *ServicePool {
	// create mj client and service
	for _, config := range appConfig.MjConfigs {
		if config.Enabled == false {
			continue
		}
		// create mj client
		client := NewClient(&config, appConfig.ProxyURL)

		// create mj service
		service := NewService()
	}

	return &ServicePool{
		taskQueue: store.NewRedisQueue("MidJourney_Task_Queue", redisCli),
	}
}

func (p *ServicePool) PushTask(task types.MjTask) {
	logger.Debugf("add a new MidJourney task to the task list: %+v", task)
	p.taskQueue.RPush(task)
}
