package sd

import (
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"fmt"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
	"strings"
	"time"
)

// SD 绘画服务

type Service struct {
	httpClient    *req.Client
	config        types.StableDiffusionConfig
	taskQueue     *store.RedisQueue
	notifyQueue   *store.RedisQueue
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	name          string // service name
	leveldb       *store.LevelDB
}

func NewService(name string, config types.StableDiffusionConfig, taskQueue *store.RedisQueue, notifyQueue *store.RedisQueue, db *gorm.DB, manager *oss.UploaderManager, levelDB *store.LevelDB) *Service {
	config.ApiURL = strings.TrimRight(config.ApiURL, "/")
	return &Service{
		name:          name,
		config:        config,
		httpClient:    req.C(),
		taskQueue:     taskQueue,
		notifyQueue:   notifyQueue,
		db:            db,
		leveldb:       levelDB,
		uploadManager: manager,
	}
}

func (s *Service) Run() {
	for {
		var task types.SdTask
		err := s.taskQueue.LPop(&task)
		if err != nil {
			logger.Errorf("taking task with error: %v", err)
			continue
		}

		// translate prompt
		if utils.HasChinese(task.Params.Prompt) {
			content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.RewritePromptTemplate, task.Params.Prompt))
			if err == nil {
				task.Params.Prompt = content
			} else {
				logger.Warnf("error with translate prompt: %v", err)
			}
		}

		// translate negative prompt
		if task.Params.NegPrompt != "" && utils.HasChinese(task.Params.NegPrompt) {
			content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.TranslatePromptTemplate, task.Params.NegPrompt))
			if err == nil {
				task.Params.NegPrompt = content
			} else {
				logger.Warnf("error with translate prompt: %v", err)
			}
		}

		logger.Infof("%s handle a new Stable-Diffusion task: %+v", s.name, task)
		err = s.Txt2Img(task)
		if err != nil {
			logger.Error("绘画任务执行失败：", err.Error())
			// update the task progress
			s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumns(map[string]interface{}{
				"progress": -1,
				"err_msg":  err.Error(),
			})
			// 通知前端，任务失败
			s.notifyQueue.RPush(task.UserId)
			continue
		}
	}
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
	apiURL := fmt.Sprintf("%s/sdapi/v1/txt2img", s.config.ApiURL)
	logger.Debugf("send image request to %s", apiURL)
	go func() {
		response, err := s.httpClient.R().SetBody(body).SetSuccessResult(&res).Post(apiURL)
		if err != nil {
			errChan <- err
			return
		}
		if response.IsErrorState() {
			errChan <- fmt.Errorf("error http code status: %v", response.Status)
			return
		}

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
		s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumns(model.SdJob{ImgURL: imgURL, Params: utils.JsonEncode(task.Params)})
		errChan <- nil
	}()

	for {
		select {
		case err := <-errChan: // 任务完成
			if err != nil {
				return err
			}
			s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumn("progress", 100)
			s.notifyQueue.RPush(task.UserId)
			// 从 leveldb 中删除预览图片数据
			_ = s.leveldb.Delete(task.Params.TaskId)
			return nil
		default:
			err, resp := s.checkTaskProgress()
			// 更新任务进度
			if err == nil && resp.Progress > 0 {
				s.db.Model(&model.SdJob{Id: uint(task.Id)}).UpdateColumn("progress", int(resp.Progress*100))
				// 发送更新状态信号
				s.notifyQueue.RPush(task.UserId)
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
func (s *Service) checkTaskProgress() (error, *TaskProgressResp) {
	apiURL := fmt.Sprintf("%s/sdapi/v1/progress?skip_current_image=false", s.config.ApiURL)
	var res TaskProgressResp
	response, err := s.httpClient.R().SetSuccessResult(&res).Get(apiURL)
	if err != nil {
		return err, nil
	}
	if response.IsErrorState() {
		return fmt.Errorf("error http code status: %v", response.Status), nil
	}

	return nil, &res
}
