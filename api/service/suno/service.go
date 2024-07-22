package dalle

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service/oss"
	"geekai/store"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

type Service struct {
	httpClient    *req.Client
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	taskQueue     *store.RedisQueue
	notifyQueue   *store.RedisQueue
	Clients       *types.LMap[uint, *types.WsClient] // UserId => Client
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client) *Service {
	return &Service{
		httpClient:    req.C().SetTimeout(time.Minute * 3),
		db:            db,
		taskQueue:     store.NewRedisQueue("Suno_Task_Queue", redisCli),
		notifyQueue:   store.NewRedisQueue("Suno_Notify_Queue", redisCli),
		Clients:       types.NewLMap[uint, *types.WsClient](),
		uploadManager: manager,
	}
}

func (s *Service) PushTask(task types.SunoTask) {
	logger.Infof("add a new Suno task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Run() {
	logger.Info("Starting Suno job consumer...")
	go func() {
		for {
			var task types.SunoTask
			err := s.taskQueue.LPop(&task)
			if err != nil {
				logger.Errorf("taking task with error: %v", err)
				continue
			}

		}
	}()
}

func (s *Service) Create(task types.SunoTask) {

}
