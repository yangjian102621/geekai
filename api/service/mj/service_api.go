package mj

import (
	"chatplus/core/types"
	"chatplus/store"
	"chatplus/store/model"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// Service MJ 绘画服务
type MjApiService struct {
	name             string     // service name
	client           *ApiClient // MJ api client
	taskQueue        *store.RedisQueue
	notifyQueue      *store.RedisQueue
	db               *gorm.DB
	maxHandleTaskNum int32 // max task number current service can handle
	handledTaskNum   int32 // already handled task number
	// key: task id, value: MjTaskStatus
	// to check if the task is timeout and check the task progress
	taskStatus      map[int]*MjTaskStatus
	taskStatusMutex *sync.Mutex
	taskTimeout     int64
}

func NewApiService(name string, taskQueue *store.RedisQueue, notifyQueue *store.RedisQueue, maxTaskNum int32, timeout int64, db *gorm.DB, client *ApiClient) *MjApiService {
	return &MjApiService{
		name:             name,
		db:               db,
		taskQueue:        taskQueue,
		notifyQueue:      notifyQueue,
		client:           client,
		taskTimeout:      timeout,
		maxHandleTaskNum: maxTaskNum,
		taskStatusMutex:  &sync.Mutex{},

		taskStatus: make(map[int]*MjTaskStatus, 0),
	}
}

func (s *MjApiService) Run() {
	logger.Infof("Starting MidJourney job consumer for %s", s.name)

	go func() {
		for {
			s.checkTasks()
			time.Sleep(time.Second * 10)
		}
	}()

	for {
		if !s.canHandleTask() {
			// current service is full, can not handle more task
			// waiting for running task finish
			time.Sleep(time.Second * 3)
			continue
		}

		var task types.MjTask
		var referenceId string
		err := s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}

		logger.Infof("%s handle a new MidJourney task: %+v", s.name, task)
		switch task.Type {
		case types.TaskImage:
			err, referenceId = s.client.Imagine(task.Prompt)
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
		s.Notify(MjTaskStatus{
			Progress:  0,
			MessageId: referenceId,
			JobId:     task.Id,
		})
		s.taskStatusMutex.Lock()
		s.taskStatus[task.Id] = &MjTaskStatus{
			CreatedAt: time.Now(),
			Progress:  0,
			MessageId: referenceId,
			JobId:     task.Id,
			State:     InQueue,
		}
		s.taskStatusMutex.Unlock()
		atomic.AddInt32(&s.handledTaskNum, 1)

	}
}

// check if current service instance can handle more task
func (s *MjApiService) canHandleTask() bool {
	handledNum := atomic.LoadInt32(&s.handledTaskNum)
	return handledNum < s.maxHandleTaskNum
}

func (s *MjApiService) safeDeleteFromTaskStatus(id int) {
	s.taskStatusMutex.Lock()
	delete(s.taskStatus, id)
	s.taskStatusMutex.Unlock()
}

// remove the expired tasks
func (s *MjApiService) checkTasks() {
	logger.Info("check tasks")
	for k, t := range s.taskStatus {
		if time.Now().Unix()-t.CreatedAt.Unix() > s.taskTimeout {
			s.safeDeleteFromTaskStatus(k)
			atomic.AddInt32(&s.handledTaskNum, -1)
			// delete task from database
			s.db.Delete(&model.MidJourneyJob{Id: uint(k)}, "progress < 100")
		} else {
			// check task progress
			messageId := t.MessageId
			if messageId != "" {
				err, res := s.client.CheckStatus(messageId)
				if err != nil {
					logger.Error("error with check status: ", err)
					continue
				}
				progress := -1
				if res.State != Fail {
					progress = extractPercentage(res.Progress)
				}
				if res.State != Success && progress == 100 {
					progress = 99
					logger.Warn("progress is 100 but state is not success")
				}

				mjTaskStatus := MjTaskStatus{
					JobId:    k,
					Progress: progress,
					State:    res.State,
					OrgURL:   res.ImgURL,
				}

				s.Notify(mjTaskStatus)

			}
		}
	}
}

const (
	InQueue TaskStatus = "in_queue"
	Success TaskStatus = "success"
	Fail    TaskStatus = "fail"
)

type MjTaskStatus struct {
	CreatedAt time.Time
	Progress  int
	MessageId string
	JobId     int
	OrgURL    string
	State     TaskStatus
}

func (s *MjApiService) Notify(data MjTaskStatus) {
	var job model.MidJourneyJob

	tx := s.db.Session(&gorm.Session{}).Where("id = ?", data.JobId)
	res := tx.First(&job)
	if job.Progress == 100 && data.State == Success {
		logger.Info("已完成不用再检测状态", data.JobId)
		s.safeDeleteFromTaskStatus(data.JobId)
		return
	}

	if res.Error != nil {
		logger.Warn("非法任务：", res.Error)
		return
	}
	logger.Infof("task id %d, progress %d imageUrl %s", data.JobId, data.Progress, data.OrgURL)

	job.MessageId = data.MessageId
	job.Progress = data.Progress
	job.OrgURL = data.OrgURL

	res = s.db.Updates(&job)
	if res.Error != nil {
		logger.Error("error with update job: ", res.Error)
		return
	}

	if data.State == Success || data.State == Fail {
		// release lock task
		atomic.AddInt32(&s.handledTaskNum, -1)
	}

	s.notifyQueue.RPush(job.UserId)

}
