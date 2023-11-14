package sd

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
	"time"
)

// SD 绘画服务

const RunningJobKey = "StableDiffusion_Running_Job"

type Service struct {
	httpClient    *req.Client
	config        *types.StableDiffusionConfig
	taskQueue     *store.RedisQueue
	redis         *redis.Client
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	Clients       *types.LMap[string, *types.WsClient] // SD 绘画页面 websocket 连接池
}

func NewService(config *types.AppConfig, redisCli *redis.Client, db *gorm.DB, manager *oss.UploaderManager) *Service {
	return &Service{
		config:        &config.SdConfig,
		httpClient:    req.C(),
		redis:         redisCli,
		db:            db,
		uploadManager: manager,
		Clients:       types.NewLMap[string, *types.WsClient](),
		taskQueue:     store.NewRedisQueue("stable_diffusion_task_queue", redisCli),
	}
}

func (s *Service) Run() {
	logger.Info("Starting StableDiffusion job consumer.")
	ctx := context.Background()
	for {
		_, err := s.redis.Get(ctx, RunningJobKey).Result()
		if err == nil { // 队列串行执行
			time.Sleep(time.Second * 3)
			continue
		}
		var task types.SdTask
		err = s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}
		logger.Infof("Consuming Task: %+v", task)
		err = s.Txt2Img(task)
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
		s.db.Model(&model.SdJob{}).Where("id = ?", task.Id).UpdateColumn("started", true)
		// 锁定任务执行通道，直到任务超时（5分钟）
		s.redis.Set(ctx, RunningJobKey, utils.JsonEncode(task), time.Minute*5)
	}
}

// PushTask 推送任务到队列
func (s *Service) PushTask(task types.SdTask) {
	logger.Infof("add a new Stable Diffusion Task: %+v", task)
	s.taskQueue.RPush(task)
}

// Txt2Img 文生图 API
func (s *Service) Txt2Img(task types.SdTask) error {
	var taskInfo TaskInfo
	bytes, err := os.ReadFile(s.config.Txt2ImgJsonPath)
	if err != nil {
		return fmt.Errorf("error with load text2img json template file: %s", err.Error())
	}

	err = json.Unmarshal(bytes, &taskInfo)
	if err != nil {
		return fmt.Errorf("error with decode json params: %s", err.Error())
	}

	data := taskInfo.Data
	params := task.Params
	data[ParamKeys["task_id"]] = params.TaskId
	data[ParamKeys["prompt"]] = params.Prompt
	data[ParamKeys["negative_prompt"]] = params.NegativePrompt
	data[ParamKeys["steps"]] = params.Steps
	data[ParamKeys["sampler"]] = params.Sampler
	// @fix bug: 有些 stable diffusion 没有面部修复功能
	//data[ParamKeys["face_fix"]] = params.FaceFix
	data[ParamKeys["cfg_scale"]] = params.CfgScale
	data[ParamKeys["seed"]] = params.Seed
	data[ParamKeys["height"]] = params.Height
	data[ParamKeys["width"]] = params.Width
	data[ParamKeys["hd_fix"]] = params.HdFix
	data[ParamKeys["hd_redraw_rate"]] = params.HdRedrawRate
	data[ParamKeys["hd_scale"]] = params.HdScale
	data[ParamKeys["hd_scale_alg"]] = params.HdScaleAlg
	data[ParamKeys["hd_sample_num"]] = params.HdSteps

	taskInfo.SessionId = task.SessionId
	taskInfo.TaskId = params.TaskId
	taskInfo.Data = data
	taskInfo.JobId = task.Id
	go func() {
		s.runTask(taskInfo, s.httpClient)
	}()
	return nil
}

// 执行任务
func (s *Service) runTask(taskInfo TaskInfo, client *req.Client) {
	body := map[string]any{
		"data":         taskInfo.Data,
		"event_data":   taskInfo.EventData,
		"fn_index":     taskInfo.FnIndex,
		"session_hash": taskInfo.SessionHash,
	}
	logger.Debug(utils.JsonEncode(body))
	var result = make(chan CBReq)
	go func() {
		var res struct {
			Data            []interface{} `json:"data"`
			IsGenerating    bool          `json:"is_generating"`
			Duration        float64       `json:"duration"`
			AverageDuration float64       `json:"average_duration"`
		}
		var cbReq = CBReq{TaskId: taskInfo.TaskId, JobId: taskInfo.JobId, SessionId: taskInfo.SessionId}
		response, err := client.R().SetBody(body).SetSuccessResult(&res).Post(s.config.ApiURL + "/run/predict")
		if err != nil {
			cbReq.Message = "error with send request: " + err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		if response.IsErrorState() {
			bytes, _ := io.ReadAll(response.Body)
			cbReq.Message = "error http status code: " + string(bytes)
			cbReq.Success = false
			result <- cbReq
			return
		}

		var images []struct {
			Name   string      `json:"name"`
			Data   interface{} `json:"data"`
			IsFile bool        `json:"is_file"`
		}
		err = utils.ForceCovert(res.Data[0], &images)
		if err != nil {
			cbReq.Message = "error with decode image:" + err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		var info map[string]any
		err = utils.JsonDecode(utils.InterfaceToString(res.Data[1]), &info)
		if err != nil {
			logger.Error(res.Data)
			cbReq.Message = "error with decode image url:" + err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		// 获取真实的 seed 值
		cbReq.ImageName = images[0].Name
		seed, _ := strconv.ParseInt(utils.InterfaceToString(info["seed"]), 10, 64)
		cbReq.Seed = seed
		cbReq.Success = true
		cbReq.Progress = 100
		result <- cbReq
		close(result)

	}()

	for {
		select {
		case value := <-result:
			s.callback(value)
			return
		default:
			var progressReq = map[string]any{
				"id_task":         taskInfo.TaskId,
				"id_live_preview": 1,
			}

			var progressRes struct {
				Active        bool        `json:"active"`
				Queued        bool        `json:"queued"`
				Completed     bool        `json:"completed"`
				Progress      float64     `json:"progress"`
				Eta           float64     `json:"eta"`
				LivePreview   string      `json:"live_preview"`
				IDLivePreview int         `json:"id_live_preview"`
				TextInfo      interface{} `json:"textinfo"`
			}
			response, err := client.R().SetBody(progressReq).SetSuccessResult(&progressRes).Post(s.config.ApiURL + "/internal/progress")
			var cbReq = CBReq{TaskId: taskInfo.TaskId, Success: true, JobId: taskInfo.JobId, SessionId: taskInfo.SessionId}
			if err != nil { // TODO: 这里可以考虑设置失败重试次数
				logger.Error(err)
				return
			}

			if response.IsErrorState() {
				bytes, _ := io.ReadAll(response.Body)
				logger.Error(string(bytes))
				return
			}

			cbReq.ImageData = progressRes.LivePreview
			cbReq.Progress = int(progressRes.Progress * 100)
			logger.Debug(cbReq)
			s.callback(cbReq)
			time.Sleep(time.Second)
		}
	}
}

func (s *Service) callback(data CBReq) {
	// 释放任务锁
	s.redis.Del(context.Background(), RunningJobKey)
	client := s.Clients.Get(data.SessionId)
	if data.Success { // 任务成功
		var job model.SdJob
		res := s.db.Where("id = ?", data.JobId).First(&job)
		if res.Error != nil {
			logger.Warn("非法任务：", res.Error)
			return
		}
		// 更新任务进度
		job.Progress = data.Progress
		// 更新任务 seed
		var params types.SdTaskParams
		err := utils.JsonDecode(job.Params, &params)
		if err != nil {
			logger.Error("任务解析失败：", err)
			return
		}

		params.Seed = data.Seed
		if data.ImageName != "" { // 下载图片
			imageURL := fmt.Sprintf("%s/file=%s", s.config.ApiURL, data.ImageName)
			imageURL, err := s.uploadManager.GetUploadHandler().PutImg(imageURL, false)
			if err != nil {
				logger.Error("error with download img: ", err.Error())
				return
			}
			job.ImgURL = imageURL
		}

		job.Params = utils.JsonEncode(params)
		res = s.db.Updates(&job)
		if res.Error != nil {
			logger.Error("error with update job: ", res.Error)
			return
		}

		var jobVo vo.SdJob
		err = utils.CopyObject(job, &jobVo)
		if err != nil {
			logger.Error("error with copy object: ", err)
			return
		}

		if data.Progress < 100 && data.ImageData != "" {
			jobVo.ImgURL = data.ImageData
		}
		// 扣减绘图次数
		s.db.Model(&model.User{}).Where("id = ?", jobVo.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls - ?", 1))
		// 推送任务到前端
		if client != nil {
			utils.ReplyChunkMessage(client, jobVo)
		}
	} else { // 任务失败
		logger.Error("任务执行失败：", data.Message)
		// 删除任务
		s.db.Delete(&model.SdJob{Id: uint(data.JobId)})
		// 推送消息到前端
		if client != nil {
			utils.ReplyChunkMessage(client, vo.SdJob{
				Id:       uint(data.JobId),
				Progress: -1,
				TaskId:   data.TaskId,
			})
		}
	}
}
