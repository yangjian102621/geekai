package sd

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

// SD 绘画服务

type Service struct {
	httpClient    *req.Client
	taskQueue     *store.RedisQueue
	notifyQueue   *store.RedisQueue
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	leveldb       *store.LevelDB
	Clients       *types.LMap[uint, *types.WsClient] // UserId => Client
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, levelDB *store.LevelDB, redisCli *redis.Client) *Service {
	return &Service{
		httpClient:    req.C(),
		taskQueue:     store.NewRedisQueue("StableDiffusion_Task_Queue", redisCli),
		notifyQueue:   store.NewRedisQueue("StableDiffusion_Queue", redisCli),
		db:            db,
		leveldb:       levelDB,
		Clients:       types.NewLMap[uint, *types.WsClient](),
		uploadManager: manager,
	}
}

func (s *Service) Run() {
	logger.Infof("Starting Stable-Diffusion job consumer")
	go func() {
		for {
			var task types.SdTask
			err := s.taskQueue.LPop(&task)
			if err != nil {
				logger.Errorf("taking task with error: %v", err)
				continue
			}

			// translate prompt
			if utils.HasChinese(task.Params.Prompt) {
				content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.RewritePromptTemplate, task.Params.Prompt), "gpt-4o-mini")
				if err == nil {
					task.Params.Prompt = content
				} else {
					logger.Warnf("error with translate prompt: %v", err)
				}
			}

			// translate negative prompt
			if task.Params.NegPrompt != "" && utils.HasChinese(task.Params.NegPrompt) {
				content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.TranslatePromptTemplate, task.Params.NegPrompt), "gpt-4o-mini")
				if err == nil {
					task.Params.NegPrompt = content
				} else {
					logger.Warnf("error with translate prompt: %v", err)
				}
			}

			logger.Infof("handle a new Stable-Diffusion task: %+v", task)
			err = s.Txt2Img(task)
			if err != nil {
				logger.Error("绘画任务执行失败：", err.Error())
				// update the task progress
				s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumns(map[string]interface{}{
					"progress": service.FailTaskProgress,
					"err_msg":  err.Error(),
				})
				// 通知前端，任务失败
				s.notifyQueue.RPush(service.NotifyMessage{UserId: task.UserId, JobId: task.Id, Message: service.TaskStatusFailed})
				continue
			}
		}
	}()
}

// Txt2ImgReq 文生图请求实体
type Txt2ImgReq struct {
	Prompt            string  `json:"prompt"`
	NegativePrompt    string  `json:"negative_prompt"`
	Seed              int64   `json:"seed,omitempty"`
	Steps             int     `json:"steps"`
	CfgScale          float32 `json:"cfg_scale"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	SamplerName       string  `json:"sampler_name"`
	Scheduler         string  `json:"scheduler"`
	EnableHr          bool    `json:"enable_hr,omitempty"`
	HrScale           int     `json:"hr_scale,omitempty"`
	HrUpscaler        string  `json:"hr_upscaler,omitempty"`
	HrSecondPassSteps int     `json:"hr_second_pass_steps,omitempty"`
	DenoisingStrength float32 `json:"denoising_strength,omitempty"`
	ForceTaskId       string  `json:"force_task_id,omitempty"`
}

// Txt2ImgResp 文生图响应实体
type Txt2ImgResp struct {
	Images     []string `json:"images"`
	Parameters struct {
	} `json:"parameters"`
	Info string `json:"info"`
}

// TaskProgressResp 任务进度响应实体
type TaskProgressResp struct {
	Progress     float64 `json:"progress"`
	EtaRelative  float64 `json:"eta_relative"`
	CurrentImage string  `json:"current_image"`
}

// Txt2Img 文生图 API
func (s *Service) Txt2Img(task types.SdTask) error {
	body := Txt2ImgReq{
		Prompt:         task.Params.Prompt,
		NegativePrompt: task.Params.NegPrompt,
		Steps:          task.Params.Steps,
		CfgScale:       task.Params.CfgScale,
		Width:          task.Params.Width,
		Height:         task.Params.Height,
		SamplerName:    task.Params.Sampler,
		Scheduler:      task.Params.Scheduler,
		ForceTaskId:    task.Params.TaskId,
	}
	if task.Params.Seed > 0 {
		body.Seed = task.Params.Seed
	}
	if task.Params.HdFix {
		body.EnableHr = true
		body.HrScale = task.Params.HdScale
		body.HrUpscaler = task.Params.HdScaleAlg
		body.HrSecondPassSteps = task.Params.HdSteps
		body.DenoisingStrength = task.Params.HdRedrawRate
	}
	var res Txt2ImgResp
	var errChan = make(chan error)

	var apiKey model.ApiKey
	err := s.db.Where("type", "sd").Where("enabled", true).Order("last_used_at ASC").First(&apiKey).Error
	if err != nil {
		return fmt.Errorf("no available Stable-Diffusion api key: %v", err)
	}

	apiURL := fmt.Sprintf("%s/sdapi/v1/txt2img", apiKey.ApiURL)
	logger.Debugf("send image request to %s", apiURL)
	// send a request to sd api endpoint
	go func() {
		response, err := s.httpClient.R().
			SetHeader("Authorization", apiKey.Value).
			SetBody(body).
			SetSuccessResult(&res).
			Post(apiURL)
		if err != nil {
			errChan <- err
			return
		}
		if response.IsErrorState() {
			errChan <- fmt.Errorf("error http code status: %v", response.Status)
			return
		}

		// update the last used time
		apiKey.LastUsedAt = time.Now().Unix()
		s.db.Updates(&apiKey)

		// 保存 Base64 图片
		imgURL, err := s.uploadManager.GetUploadHandler().PutBase64(res.Images[0])
		if err != nil {
			errChan <- fmt.Errorf("error with upload image: %v", err)
			return
		}
		// 获取绘画真实的 seed
		var info map[string]interface{}
		err = utils.JsonDecode(res.Info, &info)
		if err != nil {
			errChan <- fmt.Errorf("error with decode task response: %v", err)
			return
		}
		task.Params.Seed = int64(utils.IntValue(utils.InterfaceToString(info["seed"]), -1))
		s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumns(model.SdJob{ImgURL: imgURL, Params: utils.JsonEncode(task.Params), Prompt: task.Params.Prompt})
		errChan <- nil
	}()

	// waiting for task finish
	for {
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}

			// task finished
			s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumn("progress", 100)
			s.notifyQueue.RPush(service.NotifyMessage{UserId: task.UserId, JobId: task.Id, Message: service.TaskStatusFinished})
			// 从 leveldb 中删除预览图片数据
			_ = s.leveldb.Delete(task.Params.TaskId)
			return nil
		default:
			err, resp := s.checkTaskProgress(apiKey)
			// 更新任务进度
			if err == nil && resp.Progress > 0 {
				s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumn("progress", int(resp.Progress*100))
				// 发送更新状态信号
				s.notifyQueue.RPush(service.NotifyMessage{UserId: task.UserId, JobId: task.Id, Message: service.TaskStatusRunning})
				// 保存预览图片数据
				if resp.CurrentImage != "" {
					_ = s.leveldb.Put(task.Params.TaskId, resp.CurrentImage)
				}
			}
			time.Sleep(time.Second)
		}
	}

}

// 执行任务
func (s *Service) checkTaskProgress(apiKey model.ApiKey) (error, *TaskProgressResp) {
	apiURL := fmt.Sprintf("%s/sdapi/v1/progress?skip_current_image=false", apiKey.ApiURL)
	var res TaskProgressResp
	response, err := s.httpClient.R().
		SetHeader("Authorization", apiKey.Value).
		SetSuccessResult(&res).
		Get(apiURL)
	if err != nil {
		return err, nil
	}
	if response.IsErrorState() {
		return fmt.Errorf("error http code status: %v", response.Status), nil
	}

	return nil, &res
}

func (s *Service) PushTask(task types.SdTask) {
	logger.Debugf("add a new MidJourney task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) CheckTaskNotify() {
	go func() {
		logger.Info("Running Stable-Diffusion task notify checking ...")
		for {
			var message service.NotifyMessage
			err := s.notifyQueue.LPop(&message)
			if err != nil {
				continue
			}
			client := s.Clients.Get(uint(message.UserId))
			if client == nil {
				continue
			}
			err = client.Send([]byte(message.Message))
			if err != nil {
				continue
			}
		}
	}()
}

// CheckTaskStatus 检查任务状态，自动删除过期或者失败的任务
func (s *Service) CheckTaskStatus() {
	go func() {
		logger.Info("Running Stable-Diffusion task status checking ...")
		for {
			var jobs []model.SdJob
			res := s.db.Where("progress < ?", 100).Find(&jobs)
			if res.Error != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, job := range jobs {
				// 5 分钟还没完成的任务标记为失败
				if time.Now().Sub(job.CreatedAt) > time.Minute*5 {
					job.Progress = service.FailTaskProgress
					job.ErrMsg = "任务超时"
					s.db.Updates(&job)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
