package sd

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServicePool struct {
	services    []*Service
	taskQueue   *store.RedisQueue
	notifyQueue *store.RedisQueue
	db          *gorm.DB
	Clients     *types.LMap[uint, *types.WsClient] // UserId => Client
	uploader    *oss.UploaderManager
	levelDB     *store.LevelDB
}

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, levelDB *store.LevelDB) *ServicePool {
	services := make([]*Service, 0)
	taskQueue := store.NewRedisQueue("StableDiffusion_Task_Queue", redisCli)
	notifyQueue := store.NewRedisQueue("StableDiffusion_Queue", redisCli)

	return &ServicePool{
		taskQueue:   taskQueue,
		notifyQueue: notifyQueue,
		services:    services,
		db:          db,
		Clients:     types.NewLMap[uint, *types.WsClient](),
		uploader:    manager,
		levelDB:     levelDB,
	}
}

func (p *ServicePool) InitServices(configs []types.StableDiffusionConfig) {
	// stop old service
	for _, s := range p.services {
		s.Stop()
	}
	p.services = make([]*Service, 0)

	for k, config := range configs {
		if config.Enabled == false {
			continue
		}

		// create sd service
		name := fmt.Sprintf(" sd-service-%d", k)
		service := NewService(name, config, p.taskQueue, p.notifyQueue, p.db, p.uploader, p.levelDB)
		// run sd service
		go func() {
			service.Run()
		}()

		p.services = append(p.services, service)
	}
}

// PushTask push a new mj task in to task queue
func (p *ServicePool) PushTask(task types.SdTask) {
	logger.Debugf("add a new MidJourney task to the task list: %+v", task)
	p.taskQueue.RPush(task)
}

func (p *ServicePool) CheckTaskNotify() {
	go func() {
		logger.Info("Running Stable-Diffusion task notify checking ...")
		for {
			var message NotifyMessage
			err := p.notifyQueue.LPop(&message)
			if err != nil {
				continue
			}
			client := p.Clients.Get(uint(message.UserId))
			if client == nil {
				continue
			}
			err = client.Send([]byte(message.Message))
			if err != nil {
				continue
			}
		}
	}()
}

// CheckTaskStatus 检查任务状态，自动删除过期或者失败的任务
func (p *ServicePool) CheckTaskStatus() {
	go func() {
		logger.Info("Running Stable-Diffusion task status checking ...")
		for {
			var jobs []model.SdJob
			res := p.db.Where("progress < ?", 100).Find(&jobs)
			if res.Error != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, job := range jobs {
				// 5 分钟还没完成的任务标记为失败
				if time.Now().Sub(job.CreatedAt) > time.Minute*5 {
					job.Progress = 101
					job.ErrMsg = "任务超时"
					p.db.Updates(&job)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

// HasAvailableService check if it has available mj service in pool
func (p *ServicePool) HasAvailableService() bool {
	return len(p.services) > 0
}
