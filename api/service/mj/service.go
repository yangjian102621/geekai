package mj

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

// MJ 绘画服务

const RunningJobKey = "MidJourney_Running_Job"

type Service struct {
	client        *Client // MJ 客户端
	taskQueue     *store.RedisQueue
	redis         *redis.Client
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	Clients       *types.LMap[string, *types.WsClient] // MJ 绘画页面 websocket 连接池，用户推送绘画消息
	ChatClients   *types.LMap[string, *types.WsClient] // 聊天页面 websocket 连接池，用于推送绘画消息
	proxyURL      string
}

func NewService(redisCli *redis.Client, db *gorm.DB, client *Client, manager *oss.UploaderManager, config *types.AppConfig) *Service {
	return &Service{
		redis:         redisCli,
		db:            db,
		taskQueue:     store.NewRedisQueue("MidJourney_Task_Queue", redisCli),
		client:        client,
		uploadManager: manager,
		Clients:       types.NewLMap[string, *types.WsClient](),
		ChatClients:   types.NewLMap[string, *types.WsClient](),
		proxyURL:      config.ProxyURL,
	}
}

func (s *Service) Run() {
	logger.Info("Starting MidJourney job consumer.")
	ctx := context.Background()
	for {
		_, err := s.redis.Get(ctx, RunningJobKey).Result()
		if err == nil { // 队列串行执行
			time.Sleep(time.Second * 3)
			continue
		}
		var task types.MjTask
		err = s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}
		logger.Infof("Consuming Task: %+v", task)
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
			if task.RetryCount <= 5 {
				s.taskQueue.RPush(task)
			}
			task.RetryCount += 1
			time.Sleep(time.Second * 3)
			continue
		}

		// 更新任务的执行状态
		s.db.Model(&model.MidJourneyJob{}).Where("id = ?", task.Id).UpdateColumn("started", true)
		// 锁定任务执行通道，直到任务超时（5分钟）
		s.redis.Set(ctx, RunningJobKey, utils.JsonEncode(task), time.Minute*5)
	}
}

func (s *Service) PushTask(task types.MjTask) {
	logger.Infof("add a new MidJourney Task: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Notify(data CBReq) {
	taskString, err := s.redis.Get(context.Background(), RunningJobKey).Result()
	if err != nil { // 过期任务，丢弃
		logger.Warn("任务已过期：", err)
		return
	}

	var task types.MjTask
	err = utils.JsonDecode(taskString, &task)
	if err != nil { // 非标准任务，丢弃
		logger.Warn("任务解析失败：", err)
		return
	}

	var job model.MidJourneyJob
	res := s.db.Where("message_id = ?", data.MessageId).First(&job)
	if res.Error == nil && data.Status == Finished {
		logger.Warn("重复消息：", data.MessageId)
		return
	}

	if task.Src == types.TaskSrcImg { // 绘画任务
		var job model.MidJourneyJob
		res := s.db.Where("id = ?", task.Id).First(&job)
		if res.Error != nil {
			logger.Warn("非法任务：", res.Error)
			return
		}
		job.MessageId = data.MessageId
		job.ReferenceId = data.ReferenceId
		job.Progress = data.Progress
		job.Prompt = data.Prompt
		job.Hash = data.Image.Hash

		// 任务完成，将最终的图片下载下来
		if data.Progress == 100 {
			imgURL, err := s.uploadManager.GetUploadHandler().PutImg(data.Image.URL, true)
			if err != nil {
				logger.Error("error with download img: ", err.Error())
				return
			}
			job.ImgURL = imgURL
		} else {
			// 临时图片直接保存，访问的时候使用代理进行转发
			job.ImgURL = data.Image.URL
		}
		res = s.db.Updates(&job)
		if res.Error != nil {
			logger.Error("error with update job: ", res.Error)
			return
		}

		var jobVo vo.MidJourneyJob
		err := utils.CopyObject(job, &jobVo)
		if err == nil {
			if data.Progress < 100 {
				image, err := utils.DownloadImage(jobVo.ImgURL, s.proxyURL)
				if err == nil {
					jobVo.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}

			// 推送任务到前端
			client := s.Clients.Get(task.SessionId)
			if client != nil {
				utils.ReplyChunkMessage(client, jobVo)
			}
		}

	} else if task.Src == types.TaskSrcChat { // 聊天任务
		wsClient := s.ChatClients.Get(task.SessionId)
		if data.Status == Finished {
			if wsClient != nil && data.ReferenceId != "" {
				content := fmt.Sprintf("**%s** 任务执行成功，正在从 MidJourney 服务器下载图片，请稍后...", data.Prompt)
				utils.ReplyMessage(wsClient, content)
			}
			// download image
			imgURL, err := s.uploadManager.GetUploadHandler().PutImg(data.Image.URL, true)
			if err != nil {
				logger.Error("error with download image: ", err)
				if wsClient != nil && data.ReferenceId != "" {
					content := fmt.Sprintf("**%s** 图片下载失败：%s", data.Prompt, err.Error())
					utils.ReplyMessage(wsClient, content)
				}
				return
			}

			tx := s.db.Begin()
			data.Image.URL = imgURL
			message := model.HistoryMessage{
				UserId:     uint(task.UserId),
				ChatId:     task.ChatId,
				RoleId:     uint(task.RoleId),
				Type:       types.MjMsg,
				Icon:       task.Icon,
				Content:    utils.JsonEncode(data),
				Tokens:     0,
				UseContext: false,
			}
			res = tx.Create(&message)
			if res.Error != nil {
				logger.Error("error with update database: ", err)
				return
			}

			// save the job
			job.UserId = task.UserId
			job.Type = task.Type.String()
			job.MessageId = data.MessageId
			job.ReferenceId = data.ReferenceId
			job.Prompt = data.Prompt
			job.ImgURL = imgURL
			job.Progress = data.Progress
			job.Hash = data.Image.Hash
			job.CreatedAt = time.Now()
			res = tx.Create(&job)
			if res.Error != nil {
				logger.Error("error with update database: ", err)
				tx.Rollback()
				return
			}
			tx.Commit()
		}

		if wsClient == nil { // 客户端断线，则丢弃
			logger.Errorf("Client is offline: %+v", data)
			return
		}

		if data.Status == Finished {
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsEnd})
			// 本次绘画完毕，移除客户端
			s.ChatClients.Delete(task.SessionId)
		} else {
			// 使用代理临时转发图片
			if data.Image.URL != "" {
				image, err := utils.DownloadImage(data.Image.URL, s.proxyURL)
				if err == nil {
					data.Image.URL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
		}
	}

	// 更新用户剩余绘图次数
	// TODO: 放大图片是否需要消耗绘图次数？
	if data.Status == Finished {
		s.db.Model(&model.User{}).Where("id = ?", task.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls - ?", 1))
		// 解除任务锁定
		s.redis.Del(context.Background(), RunningJobKey)
	}

}
