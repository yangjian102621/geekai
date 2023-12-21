package mj

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ServicePool Mj service pool
type ServicePool struct {
	services        []*Service
	taskQueue       *store.RedisQueue
	db              *gorm.DB
	uploaderManager *oss.UploaderManager
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
		service := NewService(name, queue, 4, 600, db, client)
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
		taskQueue:       queue,
		services:        services,
		uploaderManager: manager,
		db:              db,
	}
}

func (p *ServicePool) DownloadImages() {
	go func() {
		var items []model.MidJourneyJob
		for {
			res := p.db.Where("img_url = ? AND progress = ?", "", 100).Find(&items)
			if res.Error != nil {
				continue
			}

			// download images
			for _, item := range items {
				imgURL, err := p.uploaderManager.GetUploadHandler().PutImg(item.OrgURL, true)
				if err != nil {
					logger.Error("error with download image: ", err)
					continue
				}

				item.ImgURL = imgURL
				p.db.Updates(&item)
			}

			time.Sleep(time.Second * 5)
		}
	}()
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
