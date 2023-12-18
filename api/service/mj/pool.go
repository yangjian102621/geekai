package mj

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ServicePool Mj service pool
type ServicePool struct {
	services  []*Service
	taskQueue *store.RedisQueue
}

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, appConfig *types.AppConfig) *ServicePool {
	services := make([]*Service, 0)
	queue := store.NewRedisQueue("MidJourney_Task_Queue", redisCli)
	// create mj client and service
	for k, config := range appConfig.MjConfigs {
		if config.Enabled == false {
			continue
		}
		// create mj client
		client := NewClient(config, appConfig.ProxyURL)

		name := fmt.Sprintf("MjService-%d", k)
		// create mj service
		service := NewService(name, queue, 4, 600, db, client, manager, appConfig.ProxyURL)
		botName := fmt.Sprintf("MjBot-%d", k)
		bot, err := NewBot(botName, appConfig.ProxyURL, &config, service)
		if err != nil {
			continue
		}

		err = bot.Run()
		if err != nil {
			continue
		}

		// run mj service
		go func() {
			service.Run()
		}()

		services = append(services, service)
	}

	return &ServicePool{
		taskQueue: queue,
		services:  services,
	}
}

// PushTask push a new mj task in to task queue
func (p *ServicePool) PushTask(task types.MjTask) {
	logger.Debugf("add a new MidJourney task to the task list: %+v", task)
	p.taskQueue.RPush(task)
}

// HasAvailableService check if it has available mj service in pool
func (p *ServicePool) HasAvailableService() bool {
	return len(p.services) > 0
}
