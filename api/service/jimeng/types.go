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

// 即梦AI错误代码常量
const (
	// 成功
	ECSuccess = 10000

	// 请求参数错误 (50200-50215)
	ECReqInvalidArgs     = 50200 // 参数错误
	ECReqMissingArgs     = 50201 // 缺少参数
	ECParseArgs          = 50204 // 参数类型错误/参数缺失
	ECImageSizeLimited   = 50205 // 图像尺寸超过限制
	ECImageEmpty         = 50206 // 请求参数中没有获取到图像
	ECImageDecodeError   = 50207 // 图像解码错误
	ECVideoEmpty         = 50209 // 请求参数中没有获取到视频
	ECVideoDecodeError   = 50210 // 视频解码错误
	ECVideoSizeLimited   = 50211 // 视频尺寸超过限制
	ECReqBodySizeLimited = 50213 // 请求Body过大
	ECVideoTimeTooLong   = 50214 // 输入视频时长过大
	ECRPCProcess         = 50215 // 请求处理失败

	// 算法服务错误 (60102-60208)
	ECJPFaceDetect      = 60102 // 算法服务需要输入人脸图，但未检测到
	ECFSLeaderRiskError = 60208 // 输入图片中包含敏感信息，未通过审核

	// 权限和系统错误 (50400-50501)
	ECAuth        = 50400 // 权限校验失败
	ECReqMethod   = 50402 // 访问的接口不存在
	ECReqLimit    = 50429 // 超过调用QPS限制
	ECInternal    = 50500 // 服务器内部错误
	ECRPCInternal = 50501 // 服务器内部RPC错误
)

// 错误代码到错误信息的映射
var errorCodeMessages = map[int]string{
	// 成功
	ECSuccess: "请求成功",

	// 请求参数错误
	ECReqInvalidArgs:     "参数错误，检查入参及MIME类型",
	ECReqMissingArgs:     "缺少参数，检查入参及MIME类型",
	ECParseArgs:          "参数类型错误/参数缺失，检查入参及MIME类型",
	ECImageSizeLimited:   "图像尺寸超过限制，参考接口文档入参要求部分",
	ECImageEmpty:         "请求参数中没有获取到图像，检查入参",
	ECImageDecodeError:   "图像解码错误：没有获取到图像或者通过image_base64参数传递图像是base64解码错误，检查输出图片或检查base64是否错误携带前缀",
	ECVideoEmpty:         "请求参数中没有获取到视频。输入为视频时可能返回此错误，检查入参",
	ECVideoDecodeError:   "视频解码错误。输入为视频时可能返回此错误，检查输入视频是否不正确",
	ECVideoSizeLimited:   "视频尺寸超过限制。输入为视频时可能返回此错误，检查输入视频大小",
	ECReqBodySizeLimited: "请求Body过大，超出接口限制，检查请求Body大小",
	ECVideoTimeTooLong:   "输入视频时长过大，检查输入视频时长",
	ECRPCProcess:         "由于输入的图片、视频、参数等不满足要求，导致请求处理失败。若接口文档中有具体说明，优先参考其具体含义，按照具体服务说明进行检查",

	// 算法服务错误
	ECJPFaceDetect:      "算法服务需要输入人脸图，但未检测到，检查输入图片是否包含人脸",
	ECFSLeaderRiskError: "输入图片中包含敏感信息，未通过审核",

	// 权限和系统错误
	ECAuth:        "权限校验失败，请检查是否已创建应用并开通服务或签名，参考接入指南及快速接入",
	ECReqMethod:   "访问的接口不存在，检查入参",
	ECReqLimit:    "超过调用QPS限制，购买QPS增项包",
	ECInternal:    "服务器内部错误，提工单",
	ECRPCInternal: "服务器内部RPC错误，提工单",
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

const (
	ImageEffectReqKey      = "i2i_multi_style_zx2x"
	DoubaoSeedream40ReqKey = "doubao-seedream-4-0-250828"
)

const (
	ASyncActionSubmit    = "CVSync2AsyncSubmitTask" // 异步提交任务
	SyncActionSubmit     = "CVSubmitTask"           // 同步提交任务
	ASyncActionGetResult = "CVSync2AsyncGetResult"  // 异步获取结果
	SyncActionGetResult  = "CVGetResult"            // 同步获取结果
)
