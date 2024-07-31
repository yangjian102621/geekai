package mj

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/oss"
	"geekai/service/sd"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ServicePool Mj service pool
type ServicePool struct {
	services        []*Service
	taskQueue       *store.RedisQueue
	notifyQueue     *store.RedisQueue
	db              *gorm.DB
	uploaderManager *oss.UploaderManager
	Clients         *types.LMap[uint, *types.WsClient] // UserId => Client
	licenseService  *service.LicenseService
}

var logger = logger2.GetLogger()

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, licenseService *service.LicenseService) *ServicePool {
	services := make([]*Service, 0)
	taskQueue := store.NewRedisQueue("MidJourney_Task_Queue", redisCli)
	notifyQueue := store.NewRedisQueue("MidJourney_Notify_Queue", redisCli)
	return &ServicePool{
		taskQueue:       taskQueue,
		notifyQueue:     notifyQueue,
		services:        services,
		uploaderManager: manager,
		db:              db,
		Clients:         types.NewLMap[uint, *types.WsClient](),
		licenseService:  licenseService,
	}
}

func (p *ServicePool) InitServices(plusConfigs []types.MjPlusConfig, proxyConfigs []types.MjProxyConfig) {
	// stop old service
	for _, s := range p.services {
		s.Stop()
	}
	p.services = make([]*Service, 0)

	for _, config := range plusConfigs {
		if config.Enabled == false {
			continue
		}

		cli := NewPlusClient(config, p.licenseService)
		name := utils.Md5(config.ApiURL)
		plusService := NewService(name, p.taskQueue, p.notifyQueue, p.db, cli)
		go func() {
			plusService.Run()
		}()
		p.services = append(p.services, plusService)
	}

	// for mid-journey proxy
	for _, config := range proxyConfigs {
		if config.Enabled == false {
			continue
		}
		cli := NewProxyClient(config)
		name := utils.Md5(config.ApiURL)
		proxyService := NewService(name, p.taskQueue, p.notifyQueue, p.db, cli)
		go func() {
			proxyService.Run()
		}()
		p.services = append(p.services, proxyService)
	}
}

func (p *ServicePool) CheckTaskNotify() {
	go func() {
		for {
			var message sd.NotifyMessage
			err := p.notifyQueue.LPop(&message)
			if err != nil {
				continue
			}
			cli := p.Clients.Get(uint(message.UserId))
			if cli == nil {
				continue
			}
			err = cli.Send([]byte(message.Message))
			if err != nil {
				continue
			}
		}
	}()
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
			for _, v := range items {
				if v.OrgURL == "" {
					continue
				}

				logger.Infof("try to download image: %s", v.OrgURL)
				mjService := p.getService(v.ChannelId)
				if mjService == nil {
					logger.Errorf("Invalid task: %+v", v)
					continue
				}

				task, _ := mjService.Client.QueryTask(v.TaskId)
				if len(task.Buttons) > 0 {
					v.Hash = GetImageHash(task.Buttons[0].CustomId)
				}
				// 如果是返回的是 discord 图片地址，则使用代理下载
				proxy := false
				if strings.HasPrefix(v.OrgURL, "https://cdn.discordapp.com") {
					proxy = true
				}
				imgURL, err := p.uploaderManager.GetUploadHandler().PutUrlFile(v.OrgURL, proxy)

				if err != nil {
					logger.Errorf("error with download image %s, %v", v.OrgURL, err)
					continue
				} else {
					logger.Infof("download image %s successfully.", v.OrgURL)
				}

				v.ImgURL = imgURL
				p.db.Updates(&v)

				cli := p.Clients.Get(uint(v.UserId))
				if cli == nil {
					continue
				}
				err = cli.Send([]byte(sd.Finished))
				if err != nil {
					continue
				}
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

// SyncTaskProgress 异步拉取任务
func (p *ServicePool) SyncTaskProgress() {
	go func() {
		var jobs []model.MidJourneyJob
		for {
			res := p.db.Where("progress < ?", 100).Find(&jobs)
			if res.Error != nil {
				continue
			}

			for _, job := range jobs {
				// 5 分钟还没完成的任务标记为失败
				if time.Now().Sub(job.CreatedAt) > time.Minute*5 {
					job.Progress = 101
					job.ErrMsg = "任务超时"
					p.db.Updates(&job)
					continue
				}

				if servicePlus := p.getService(job.ChannelId); servicePlus != nil {
					_ = servicePlus.Notify(job)
				}
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

func (p *ServicePool) getService(name string) *Service {
	for _, s := range p.services {
		if s.Name == name {
			return s
		}
	}
	return nil
}
