package service

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/store"
	"chatplus/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"time"
)

var logger = logger2.GetLogger()

// MJ 绘画服务

const MjRunningJobKey = "MidJourney_Running_Job"

type TaskType string

const (
	Image     = TaskType("image")
	Upscale   = TaskType("upscale")
	Variation = TaskType("variation")
)

type TaskSrc string

const (
	TaskSrcChat = TaskSrc("chat")
	TaskSrcImg  = TaskSrc("img")
)

type MjTask struct {
	Id          string   `json:"id"`
	Src         TaskSrc  `json:"src"`
	Type        TaskType `json:"type"`
	UserId      int      `json:"user_id"`
	Prompt      string   `json:"prompt,omitempty"`
	ChatId      string   `json:"chat_id,omitempty"`
	RoleId      int      `json:"role_id,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Index       int32    `json:"index,omitempty"`
	MessageId   string   `json:"message_id,omitempty"`
	MessageHash string   `json:"message_hash,omitempty"`
	RetryCount  int      `json:"retry_count"`
}

type MjService struct {
	config    types.ChatPlusExtConfig
	client    *req.Client
	taskQueue *store.RedisQueue
	redis     *redis.Client
}

func NewMjService(config types.ChatPlusExtConfig, client *redis.Client) *MjService {
	return &MjService{
		config:    config,
		redis:     client,
		taskQueue: store.NewRedisQueue("midjourney_task_queue", client),
		client:    req.C().SetTimeout(30 * time.Second)}
}

func (s *MjService) Run() {
	ctx := context.Background()
	for {
		_, err := s.redis.Get(ctx, MjRunningJobKey).Result()
		if err == nil { // a task is running, waiting for finish
			time.Sleep(time.Second * 3)
			continue
		}
		var task MjTask
		err = s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}

		switch task.Type {
		case Image:
			err = s.image(task.Prompt)
			break
		case Upscale:
			err = s.upscale(MjUpscaleReq{
				Index:       task.Index,
				MessageId:   task.MessageId,
				MessageHash: task.MessageHash,
			})
			break
		case Variation:
			err = s.variation(MjVariationReq{
				Index:       task.Index,
				MessageId:   task.MessageId,
				MessageHash: task.MessageHash,
			})
		}
		if err != nil {
			if task.RetryCount > 5 {
				continue
			}
			task.RetryCount += 1
			time.Sleep(time.Second)
			s.taskQueue.RPush(task)
			// TODO: 执行失败通知聊天客户端
			continue
		}

		// 锁定任务执行通道，直到任务超时（10分钟）
		s.redis.Set(ctx, MjRunningJobKey, utils.JsonEncode(task), time.Second*600)
	}
}

func (s *MjService) PushTask(task MjTask) {
	s.taskQueue.RPush(task)
}

func (s *MjService) image(prompt string) error {
	logger.Infof("MJ 绘画参数：%+v", prompt)
	body := map[string]string{"prompt": prompt}
	url := fmt.Sprintf("%s/api/mj/image", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}

type MjUpscaleReq struct {
	Index       int32  `json:"index"`
	MessageId   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

func (s *MjService) upscale(upReq MjUpscaleReq) error {
	url := fmt.Sprintf("%s/api/mj/upscale", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(upReq).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}

type MjVariationReq struct {
	Index       int32  `json:"index"`
	MessageId   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

func (s *MjService) variation(upReq MjVariationReq) error {
	url := fmt.Sprintf("%s/api/mj/variation", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", s.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(upReq).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return errors.New(res.Message)
	}

	return nil
}
