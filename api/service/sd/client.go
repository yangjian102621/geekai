package sd

import (
	"chatplus/core/types"
	"chatplus/utils"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
	"time"
)

type Client struct {
	httpClient *req.Client
	config     *types.StableDiffusionConfig
}

func NewSdClient(config *types.AppConfig) *Client {
	return &Client{
		config:     &config.SdConfig,
		httpClient: req.C(),
	}
}

func (c *Client) Txt2Img(params types.SdTaskParams) error {
	var data []interface{}
	err := utils.JsonDecode(Text2ImgParamTemplate, &data)
	if err != nil {
		return err
	}
	data[ParamKeys["task_id"]] = params.TaskId
	data[ParamKeys["prompt"]] = params.Prompt
	data[ParamKeys["negative_prompt"]] = params.NegativePrompt
	data[ParamKeys["steps"]] = params.Steps
	data[ParamKeys["sampler"]] = params.Sampler
	data[ParamKeys["face_fix"]] = params.FaceFix
	data[ParamKeys["cfg_scale"]] = params.CfgScale
	data[ParamKeys["seed"]] = params.Seed
	data[ParamKeys["height"]] = params.Height
	data[ParamKeys["width"]] = params.Width
	data[ParamKeys["hd_fix"]] = params.HdFix
	data[ParamKeys["hd_redraw_rate"]] = params.HdRedrawRate
	data[ParamKeys["hd_scale"]] = params.HdScale
	data[ParamKeys["hd_scale_alg"]] = params.HdScaleAlg
	data[ParamKeys["hd_sample_num"]] = params.HdSampleNum
	task := TaskInfo{
		TaskId:      params.TaskId,
		Data:        data,
		EventData:   nil,
		FnIndex:     494,
		SessionHash: "ycaxgzm9ah",
	}

	go func() {
		c.runTask(task, c.httpClient)
	}()
	return nil
}

func (c *Client) runTask(taskInfo TaskInfo, client *req.Client) {
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
		var cbReq = CBReq{TaskId: taskInfo.TaskId}
		response, err := client.R().SetBody(body).SetSuccessResult(&res).Post(c.config.ApiURL + "/run/predict")
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
			cbReq.Message = err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		//for k, v := range info {
		//	fmt.Println(k, " => ", v)
		//}
		cbReq.ImageName = images[0].Name
		cbReq.Seed = utils.InterfaceToString(info["seed"])
		cbReq.Success = true
		cbReq.Progress = 100
		result <- cbReq
		close(result)

	}()

	for {
		select {
		case value := <-result:
			if value.Success {
				logger.Infof("%s/file=%s", c.config.ApiURL, value.ImageName)
			}
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
			response, err := client.R().SetBody(progressReq).SetSuccessResult(&progressRes).Post(c.config.ApiURL + "/internal/progress")
			var cbReq = CBReq{TaskId: taskInfo.TaskId, Success: true}
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
			fmt.Println("Progress: ", progressRes.Progress)
			fmt.Println("Image: ", progressRes.LivePreview)
			time.Sleep(time.Second)
		}
	}
}
