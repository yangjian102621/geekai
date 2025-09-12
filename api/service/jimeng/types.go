package jimeng

import (
	"geekai/core/types"
)

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	ReqKey string `json:"req_key"`
	// 文生图参数
	Prompt    string  `json:"prompt,omitempty"`
	Seed      int64   `json:"seed,omitempty"`
	Scale     float64 `json:"scale,omitempty"`
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	UsePreLLM bool    `json:"use_pre_llm,omitempty"`
	// 图生图参数
	ImageInput       string   `json:"image_input,omitempty"`
	ImageUrls        []string `json:"image_urls,omitempty"`
	BinaryDataBase64 []string `json:"binary_data_base64,omitempty"`
	Gpen             float64  `json:"gpen,omitempty"`
	Skin             float64  `json:"skin,omitempty"`
	SkinUnifi        float64  `json:"skin_unifi,omitempty"`
	GenMode          string   `json:"gen_mode,omitempty"`
	// 图像编辑参数
	// 图像特效参数
	ImageInput1 string `json:"image_input1,omitempty"`
	TemplateId  string `json:"template_id,omitempty"`
	// 视频生成参数
	AspectRatio string `json:"aspect_ratio,omitempty"`
}

// SubmitTaskResponse 提交任务响应
type SubmitTaskResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RequestId   string `json:"request_id"`
	Status      int    `json:"status"`
	TimeElapsed string `json:"time_elapsed"`
	Data        struct {
		TaskId string `json:"task_id"`
	} `json:"data"`
}

// QueryTaskRequest 查询任务请求
type QueryTaskRequest struct {
	ReqKey  string `json:"req_key"`
	TaskId  string `json:"task_id"`
	ReqJson string `json:"req_json,omitempty"`
}

// QueryTaskResponse 查询任务响应
type QueryTaskResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RequestId   string `json:"request_id"`
	Status      int    `json:"status"`
	TimeElapsed string `json:"time_elapsed"`
	Data        struct {
		AlgorithmBaseResp struct {
			StatusCode    int    `json:"status_code"`
			StatusMessage string `json:"status_message"`
		} `json:"algorithm_base_resp"`
		BinaryDataBase64  []string           `json:"binary_data_base64"`
		ImageUrls         []string           `json:"image_urls"`
		VideoUrl          string             `json:"video_url"`
		RespData          string             `json:"resp_data"`
		Status            types.JMTaskStatus `json:"status"`
		LlmResult         string             `json:"llm_result"`
		PeResult          string             `json:"pe_result"`
		PredictTagsResult string             `json:"predict_tags_result"`
		RephraserResult   string             `json:"rephraser_result"`
		VlmResult         string             `json:"vlm_result"`
		InferCtx          any                `json:"infer_ctx"`
	} `json:"data"`
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Type      types.JMTaskType `json:"type"`
	Prompt    string           `json:"prompt"`
	Params    map[string]any   `json:"params"`
	ReqKey    string           `json:"req_key"`
	ImageUrls []string         `json:"image_urls,omitempty"`
	Power     int              `json:"power,omitempty"`
}
