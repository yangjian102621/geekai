package plus

import (
	"chatplus/core/types"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// Service MJ 绘画服务
type Service struct {
	name             string  // service name
	Client           *Client // MJ Client
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
		Client:           client,
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
		if task.ChannelId != "" && task.ChannelId != s.Client.Config.Name {
			s.taskQueue.RPush(task)
			time.Sleep(time.Second)
			continue
		}

		logger.Infof("%s handle a new MidJourney task: %+v", s.name, task)
		var res ImageRes
		switch task.Type {
		case types.TaskImage:
			index := strings.Index(task.Prompt, " ")
			res, err = s.Client.Imagine(task.Prompt[index+1:])
			break
		case types.TaskUpscale:
			res, err = s.Client.Upscale(task.Index, task.MessageId, task.MessageHash)
			break
		case types.TaskVariation:
			res, err = s.Client.Variation(task.Index, task.MessageId, task.MessageHash)
		}

		if err != nil || (res.Code != 1 && res.Code != 22) {
			logger.Error("绘画任务执行失败：", err)
			// update the task progress
			s.db.Model(&model.MidJourneyJob{Id: uint(task.Id)}).UpdateColumn("progress", -1)
			// 任务失败，通知前端
			s.notifyQueue.RPush(task.UserId)
			// restore img_call quota
			s.db.Model(&model.User{}).Where("id = ?", task.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls + ?", 1))

			// TODO: 任务提交失败，加入队列重试
			continue
		}
		logger.Infof("任务提交成功：%+v", res)
		// lock the task until the execute timeout
		s.taskStartTimes[task.Id] = time.Now()
		atomic.AddInt32(&s.handledTaskNum, 1)
		// 更新任务 ID/频道
		s.db.Model(&model.MidJourneyJob{}).Where("id = ?", task.Id).UpdateColumns(map[string]interface{}{
			"task_id":    res.Result,
			"channel_id": s.Client.Config.Name,
		})

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

type CBReq struct {
	Id          string      `json:"id"`
	Action      string      `json:"action"`
	Status      string      `json:"status"`
	Prompt      string      `json:"prompt"`
	PromptEn    string      `json:"promptEn"`
	Description string      `json:"description"`
	SubmitTime  int64       `json:"submitTime"`
	StartTime   int64       `json:"startTime"`
	FinishTime  int64       `json:"finishTime"`
	Progress    string      `json:"progress"`
	ImageUrl    string      `json:"imageUrl"`
	FailReason  interface{} `json:"failReason"`
	Properties  struct {
		FinalPrompt string `json:"finalPrompt"`
	} `json:"properties"`
}

func (s *Service) Notify(data CBReq, job model.MidJourneyJob) error {

	job.Progress = utils.IntValue(strings.Replace(data.Progress, "%", "", 1), 0)
	job.Prompt = data.Properties.FinalPrompt
	if data.ImageUrl != "" {
		job.OrgURL = data.ImageUrl
	}
	job.UseProxy = true
	job.MessageId = data.Id
	logger.Debugf("JOB: %+v", job)
	res := s.db.Updates(&job)
	if res.Error != nil {
		return fmt.Errorf("error with update job: %v", res.Error)
	}

	if data.Status == "SUCCESS" {
		// release lock task
		atomic.AddInt32(&s.handledTaskNum, -1)
	}

	s.notifyQueue.RPush(job.UserId)
	return nil
}
