package sd

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServicePool struct {
	services  []*Service
	taskQueue *store.RedisQueue
}

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, appConfig *types.AppConfig) *ServicePool {
	services := make([]*Service, 0)
	queue := store.NewRedisQueue("StableDiffusion_Task_Queue", redisCli)
	// create mj client and service
	for k, config := range appConfig.SdConfigs {
		if config.Enabled == false {
			continue
		}

		// create sd service
		name := fmt.Sprintf("StableDifffusion Service-%d", k)
		service := NewService(name, 1, 300, config, queue, db, manager)
		// run sd service
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
func (p *ServicePool) PushTask(task types.SdTask) {
	logger.Debugf("add a new MidJourney task to the task list: %+v", task)
	p.taskQueue.RPush(task)
}

// HasAvailableService check if it has available mj service in pool
func (p *ServicePool) HasAvailableService() bool {
	return len(p.services) > 0
}
