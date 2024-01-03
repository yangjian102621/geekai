package mj

import (
	"chatplus/core/types"
	"chatplus/store"
	"chatplus/store/model"
	"gorm.io/gorm"
	"strings"
	"sync/atomic"
	"time"
)

// Service MJ 绘画服务
type Service struct {
	name             string  // service name
	client           *Client // MJ client
	taskQueue        *store.RedisQueue
	notifyQueue      *store.RedisQueue
	db               *gorm.DB
	maxHandleTaskNum int32             // max task number current service can handle
	handledTaskNum   int32             // already handled task number
	taskStartTimes   map[int]time.Time // task start time, to check if the task is timeout
	taskTimeout      int64
}

func NewService(name string, taskQueue *store.RedisQueue, notifyQueue *store.RedisQueue, maxTaskNum int32, timeout int64, db *gorm.DB, client *Client) *Service {
	return &Service{
		name:             name,
		db:               db,
		taskQueue:        taskQueue,
		notifyQueue:      notifyQueue,
		client:           client,
		taskTimeout:      timeout,
		maxHandleTaskNum: maxTaskNum,
		taskStartTimes:   make(map[int]time.Time, 0),
	}
}

func (s *Service) Run() {
	logger.Infof("Starting MidJourney job consumer for %s", s.name)
	for {
		s.checkTasks()
		if !s.canHandleTask() {
			// current service is full, can not handle more task
			// waiting for running task finish
			time.Sleep(time.Second * 3)
			continue
		}

		var task types.MjTask
		err := s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}

		// if it's reference message, check if it's this channel's  message
		if task.ChannelId != "" && task.ChannelId != s.client.Config.ChanelId {
			s.taskQueue.RPush(task)
			s.db.Model(&model.MidJourneyJob{Id: uint(task.Id)}).UpdateColumn("progress", -1)
			s.notifyQueue.RPush(task.UserId)
			time.Sleep(time.Second)
			continue
		}

		logger.Infof("%s handle a new MidJourney task: %+v", s.name, task)
		switch task.Type {
		case types.TaskImage:
			err = s.client.Imagine(task.Prompt)
			break
		case types.TaskUpscale:
			err = s.client.Upscale(task.Index, task.MessageId, task.MessageHash)

			break
		case types.TaskVariation:
			err = s.client.Variation(task.Index, task.MessageId, task.MessageHash)
		}

		if err != nil {
			logger.Error("绘画任务执行失败：", err)
			// update the task progress
			s.db.Model(&model.MidJourneyJob{Id: uint(task.Id)}).UpdateColumn("progress", -1)
			s.notifyQueue.RPush(task.UserId)
			// restore img_call quota
			s.db.Model(&model.User{}).Where("id = ?", task.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls + ?", 1))
			continue
		}

		// lock the task until the execute timeout
		s.taskStartTimes[task.Id] = time.Now()
		atomic.AddInt32(&s.handledTaskNum, 1)

	}
}

// check if current service instance can handle more task
func (s *Service) canHandleTask() bool {
	handledNum := atomic.LoadInt32(&s.handledTaskNum)
	return handledNum < s.maxHandleTaskNum
}

// remove the expired tasks
func (s *Service) checkTasks() {
	for k, t := range s.taskStartTimes {
		if time.Now().Unix()-t.Unix() > s.taskTimeout {
			delete(s.taskStartTimes, k)
			atomic.AddInt32(&s.handledTaskNum, -1)
			// delete task from database
			s.db.Delete(&model.MidJourneyJob{Id: uint(k)}, "progress < 100")
		}
	}
}

func (s *Service) Notify(data CBReq) {
	// extract the task ID
	split := strings.Split(data.Prompt, " ")
	var job model.MidJourneyJob
	res := s.db.Where("message_id = ?", data.MessageId).First(&job)
	if res.Error == nil && data.Status == Finished {
		logger.Warn("重复消息：", data.MessageId)
		return
	}

	tx := s.db.Session(&gorm.Session{}).Order("id ASC")
	if data.ReferenceId != "" {
		tx = tx.Where("reference_id = ?", data.ReferenceId)
	} else {
		tx = tx.Where("task_id = ?", split[0])
	}
	res = tx.First(&job)
	if res.Error != nil {
		logger.Warn("非法任务：", res.Error)
		return
	}

	job.ChannelId = data.ChannelId
	job.MessageId = data.MessageId
	job.ReferenceId = data.ReferenceId
	job.Progress = data.Progress
	job.Prompt = data.Prompt
	job.Hash = data.Image.Hash
	job.OrgURL = data.Image.URL
	if s.client.Config.UseCDN {
		job.UseProxy = true
		job.ImgURL = strings.ReplaceAll(data.Image.URL, "https://cdn.discordapp.com", s.client.Config.DiscordCDN)
	}

	res = s.db.Updates(&job)
	if res.Error != nil {
		logger.Error("error with update job: ", res.Error)
		return
	}

	if data.Status == Finished {
		// release lock task
		atomic.AddInt32(&s.handledTaskNum, -1)
	}

	s.notifyQueue.RPush(job.UserId)

}
