package mj

import (
	"chatplus/core/types"
	"chatplus/service/mj/plus"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ServicePool Mj service pool
type ServicePool struct {
	services        []interface{}
	taskQueue       *store.RedisQueue
	notifyQueue     *store.RedisQueue
	db              *gorm.DB
	uploaderManager *oss.UploaderManager
	Clients         *types.LMap[uint, *types.WsClient] // UserId => Client
}

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, appConfig *types.AppConfig) *ServicePool {
	services := make([]interface{}, 0)
	taskQueue := store.NewRedisQueue("MidJourney_Task_Queue", redisCli)
	notifyQueue := store.NewRedisQueue("MidJourney_Notify_Queue", redisCli)

	for k, config := range appConfig.MjPlusConfigs {
		if config.Enabled == false {
			continue
		}
		client := plus.NewClient(config)
		name := fmt.Sprintf("MidJourney Plus Service-%d", k)
		servicePlus := plus.NewService(name, taskQueue, notifyQueue, 10, 600, db, client)
		go func() {
			servicePlus.Run()
		}()
		services = append(services, servicePlus)
	}

	if len(services) == 0 {
		// create mj client and service
		for k, config := range appConfig.MjConfigs {
			if config.Enabled == false {
				continue
			}
			// create mj client
			client := NewClient(config, appConfig.ProxyURL, appConfig.ImgCdnURL)

			name := fmt.Sprintf("MjService-%d", k)
			// create mj service
			service := NewService(name, taskQueue, notifyQueue, 4, 600, db, client)
			botName := fmt.Sprintf("MjBot-%d", k)
			bot, err := NewBot(botName, appConfig.ProxyURL, config, service)
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
			client := p.Clients.Get(userId)
			err = client.Send([]byte("Task Updated"))
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
				var imgURL string
				var err error
				if v.UseProxy {
					if servicePlus := p.getServicePlus(v.ChannelId); servicePlus != nil {
						task, _ := servicePlus.Client.QueryTask(v.TaskId)
						if task.ImageUrl != "" {
							imgURL, err = p.uploaderManager.GetUploadHandler().PutImg(task.ImageUrl, false)
						}
						if len(task.Buttons) > 0 {
							v.Hash = getImageHash(task.Buttons[0].CustomId)
						}
					}
				} else {
					imgURL, err = p.uploaderManager.GetUploadHandler().PutImg(v.OrgURL, true)
				}
				if err != nil {
					logger.Error("error with download image: ", err)
					continue
				}

				v.ImgURL = imgURL
				p.db.Updates(&v)

				client := p.Clients.Get(uint(v.UserId))
				err = client.Send([]byte("Task Updated"))
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

func (p *ServicePool) Notify(data plus.CBReq) error {
	logger.Infof("收到任务回调：%+v", data)
	var job model.MidJourneyJob
	res := p.db.Where("task_id = ?", data.Id).First(&job)
	if res.Error != nil {
		return fmt.Errorf("非法任务：%s", data.Id)
	}

	if servicePlus := p.getServicePlus(job.ChannelId); servicePlus != nil {
		return servicePlus.Notify(data, job)
	}

	return nil
}

func (p *ServicePool) getServicePlus(name string) *plus.Service {
	for _, s := range p.services {
		if servicePlus, ok := s.(*plus.Service); ok {
			if servicePlus.Client.Config.Name == name {
				return servicePlus
			}
		}
	}
	return nil
}

func getImageHash(action string) string {
	split := strings.Split(action, "::")
	if len(split) > 5 {
		return split[4]
	}
	return split[len(split)-1]
}
