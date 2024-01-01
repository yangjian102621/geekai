package sd

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

// SD 绘画服务

type Service struct {
	httpClient       *req.Client
	config           types.StableDiffusionConfig
	taskQueue        *store.RedisQueue
	db               *gorm.DB
	uploadManager    *oss.UploaderManager
	name             string            // service name
	maxHandleTaskNum int32             // max task number current service can handle
	handledTaskNum   int32             // already handled task number
	taskStartTimes   map[int]time.Time // task start time, to check if the task is timeout
	taskTimeout      int64
}

func NewService(name string, maxTaskNum int32, timeout int64, config types.StableDiffusionConfig, queue *store.RedisQueue, db *gorm.DB, manager *oss.UploaderManager) *Service {
	return &Service{
		name:             name,
		config:           config,
		httpClient:       req.C(),
		taskQueue:        queue,
		db:               db,
		uploadManager:    manager,
		taskTimeout:      timeout,
		maxHandleTaskNum: maxTaskNum,
		taskStartTimes:   make(map[int]time.Time),
	}
}

func (s *Service) Run() {
	for {
		s.checkTasks()
		if !s.canHandleTask() {
			// current service is full, can not handle more task
			// waiting for running task finish
			time.Sleep(time.Second * 3)
			continue
		}

		var task types.SdTask
		err := s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}
		logger.Infof("%s handle a new Stable-Diffusion task: %+v", s.name, task)
		err = s.Txt2Img(task)
		if err != nil {
			logger.Error("绘画任务执行失败：", err)
			// update the task progress
			s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumn("progress", -1)
			// restore img_call quota
			s.db.Model(&model.User{}).Where("id = ?", task.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls + ?", 1))
			// release task num
			atomic.AddInt32(&s.handledTaskNum, -1)
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
	taskInfo.UserId = uint(task.UserId)
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
	var result = make(chan CBReq)
	go func() {
		var res struct {
			Data            []interface{} `json:"data"`
			IsGenerating    bool          `json:"is_generating"`
			Duration        float64       `json:"duration"`
			AverageDuration float64       `json:"average_duration"`
		}
		var cbReq = CBReq{UserId: taskInfo.UserId, TaskId: taskInfo.TaskId, JobId: taskInfo.JobId, SessionId: taskInfo.SessionId}
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
			var cbReq = CBReq{UserId: taskInfo.UserId, TaskId: taskInfo.TaskId, Success: true, JobId: taskInfo.JobId, SessionId: taskInfo.SessionId}
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
			s.callback(cbReq)
			time.Sleep(time.Second)
		}
	}
}

func (s *Service) callback(data CBReq) {
	// release task num
	atomic.AddInt32(&s.handledTaskNum, -1)
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
			job.ImgURL = fmt.Sprintf("%s/file=%s", s.config.ApiURL, data.ImageName)
			if data.Progress == 100 {
				imageURL, err := s.uploadManager.GetUploadHandler().PutImg(job.ImgURL, false)
				if err != nil {
					logger.Error("error with download img: ", err.Error())
					return
				}
				job.ImgURL = imageURL
			}
		}

		job.Params = utils.JsonEncode(params)
		res = s.db.Updates(&job)
		if res.Error != nil {
			logger.Error("error with update job: ", res.Error)
			return
		}

		logger.Debugf("绘图进度：%d", data.Progress)
	} else { // 任务失败
		logger.Error("任务执行失败：", data.Message)
		// update the task progress
		s.db.Model(&model.SdJob{Id: uint(data.JobId)}).UpdateColumn("progress", -1)
		// restore img_calls
		s.db.Model(&model.User{}).Where("id = ? AND img_calls > 0", data.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls + ?", 1))
	}
}
