package sd

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

type ServicePool struct {
	services    []*Service
	taskQueue   *store.RedisQueue
	notifyQueue *store.RedisQueue
	db          *gorm.DB
	Clients     *types.LMap[uint, *types.WsClient] // UserId => Client
}

func NewServicePool(db *gorm.DB, redisCli *redis.Client, manager *oss.UploaderManager, appConfig *types.AppConfig, levelDB *store.LevelDB) *ServicePool {
	services := make([]*Service, 0)
	taskQueue := store.NewRedisQueue("StableDiffusion_Task_Queue", redisCli)
	notifyQueue := store.NewRedisQueue("StableDiffusion_Queue", redisCli)
	// create mj client and service
	for _, config := range appConfig.SdConfigs {
		if config.Enabled == false {
			continue
		}

		// create sd service
		name := fmt.Sprintf("StableDifffusion Service-%s", config.Model)
		service := NewService(name, config, taskQueue, notifyQueue, db, manager, levelDB)
		// run sd service
		go func() {
			service.Run()
		}()

		services = append(services, service)
	}

	return &ServicePool{
		taskQueue:   taskQueue,
		notifyQueue: notifyQueue,
		services:    services,
		db:          db,
		Clients:     types.NewLMap[uint, *types.WsClient](),
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
			var userId uint
			err := p.notifyQueue.LPop(&userId)
			if err != nil {
				continue
			}
			client := p.Clients.Get(userId)
			if client == nil {
				continue
			}
			err = client.Send([]byte("Task Updated"))
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
				// 5 分钟还没完成的任务直接删除
				if time.Now().Sub(job.CreatedAt) > time.Minute*5 || job.Progress == -1 {
					p.db.Delete(&job)
					var user model.User
					p.db.Where("id = ?", job.UserId).First(&user)
					// 退回绘图次数
					res = p.db.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("power", gorm.Expr("power + ?", job.Power))
					if res.Error == nil && res.RowsAffected > 0 {
						p.db.Create(&model.PowerLog{
							UserId:    user.Id,
							Username:  user.Username,
							Type:      types.PowerConsume,
							Amount:    job.Power,
							Balance:   user.Power + job.Power,
							Mark:      types.PowerAdd,
							Model:     "stable-diffusion",
							Remark:    fmt.Sprintf("任务失败，退回算力。任务ID：%s", job.TaskId),
							CreatedAt: time.Now(),
						})
					}
					continue
				}
			}

		}
	}()
}

// HasAvailableService check if it has available mj service in pool
func (p *ServicePool) HasAvailableService() bool {
	return len(p.services) > 0
}
