package mj

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"fmt"
	"github.com/go-redis/redis/v8"
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
}

var logger = logger2.GetLogger()

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, appConfig *types.AppConfig) *ServicePool {
	services := make([]*Service, 0)
	taskQueue := store.NewRedisQueue("MidJourney_Task_Queue", redisCli)
	notifyQueue := store.NewRedisQueue("MidJourney_Notify_Queue", redisCli)

	for k, config := range appConfig.MjPlusConfigs {
		if config.Enabled == false {
			continue
		}
		cli := NewPlusClient(config)
		name := fmt.Sprintf("mj-plus-service-%d", k)
		service := NewService(name, taskQueue, notifyQueue, 4, 600, db, cli)
		go func() {
			service.Run()
		}()
		services = append(services, service)
	}

	for k, config := range appConfig.MjProxyConfigs {
		if config.Enabled == false {
			continue
		}
		cli := NewProxyClient(config)
		name := fmt.Sprintf("mj-proxy-service-%d", k)
		service := NewService(name, taskQueue, notifyQueue, 4, 600, db, cli)
		go func() {
			service.Run()
		}()
		services = append(services, service)
	}

	return &ServicePool{
		taskQueue:       taskQueue,
		notifyQueue:     notifyQueue,
		services:        services,
		uploaderManager: manager,
		db:              db,
		Clients:         types.NewLMap[uint, *types.WsClient](),
	}
}

func (p *ServicePool) CheckTaskNotify() {
	go func() {
		for {
			var userId uint
			err := p.notifyQueue.LPop(&userId)
			if err != nil {
				continue
			}
			cli := p.Clients.Get(userId)
			if cli == nil {
				continue
			}
			err = cli.Send([]byte("Task Updated"))
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
				var imgURL string
				var err error
				if servicePlus := p.getService(v.ChannelId); servicePlus != nil {
					task, _ := servicePlus.Client.QueryTask(v.TaskId)
					if len(task.Buttons) > 0 {
						v.Hash = GetImageHash(task.Buttons[0].CustomId)
					}
					imgURL, err = p.uploaderManager.GetUploadHandler().PutImg(v.OrgURL, false)
				} else {
					imgURL, err = p.uploaderManager.GetUploadHandler().PutImg(v.OrgURL, true)
				}
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
				err = cli.Send([]byte("Task Updated"))
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
		var items []model.MidJourneyJob
		for {
			res := p.db.Where("progress < ?", 100).Find(&items)
			if res.Error != nil {
				continue
			}

			for _, job := range items {
				// 失败或者 30 分钟还没完成的任务删除并退回算力
				if time.Now().Sub(job.CreatedAt) > time.Minute*30 || job.Progress == -1 {
					// 删除任务
					p.db.Delete(&job)
					// 退回算力
					tx := p.db.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("power", gorm.Expr("power + ?", job.Power))
					if tx.Error == nil && tx.RowsAffected > 0 {
						var user model.User
						p.db.Where("id = ?", job.UserId).First(&user)
						p.db.Create(&model.PowerLog{
							UserId:    user.Id,
							Username:  user.Username,
							Type:      types.PowerConsume,
							Amount:    job.Power,
							Balance:   user.Power + job.Power,
							Mark:      types.PowerAdd,
							Model:     "mid-journey",
							Remark:    fmt.Sprintf("绘画任务失败，退回算力。任务ID：%s", job.TaskId),
							CreatedAt: time.Now(),
						})
					}
					continue
				}

				if servicePlus := p.getService(job.ChannelId); servicePlus != nil {
					_ = servicePlus.Notify(job)
				}
			}

			time.Sleep(time.Second)
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
