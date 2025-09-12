package jimeng

import (
	"geekai/core/types"
)

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

const CodeSuccess = 10000

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Type      types.JMTaskType `json:"type"`
	Prompt    string           `json:"prompt"`
	Params    map[string]any   `json:"params"`
	ReqKey    string           `json:"req_key"`
	ImageUrls []string         `json:"image_urls,omitempty"`
	Power     int              `json:"power,omitempty"`
}

const (
	ImageEffectReqKey      = "i2i_multi_style_zx2x"
	DoubaoSeedream40ReqKey = "doubao-seedream-4-0-250828"
)
