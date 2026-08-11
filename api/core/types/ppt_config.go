package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// PPTImageProvider PPT 图片生成提供方
type PPTImageProvider string

const (
	PPTImageProviderNanoBanana PPTImageProvider = "nano_banana"
	PPTImageProviderSeedream   PPTImageProvider = "seedream"
)

// PPTConfig PPT 生成配置（存储在 config 表，name = 'ppt'）
type PPTConfig struct {
	// 分镜 LLM 配置
	OutlineLLMApiURL string `json:"outline_llm_api_url"` // 分镜 LLM API 地址
	OutlineLLMApiKey string `json:"outline_llm_api_key"` // 分镜 LLM API Key
	OutlineLLMModel  string `json:"outline_llm_model"`   // 分镜 LLM 模型名称，如 gpt-4o-mini

	// 图片生成通用配置
	ActiveImageProvider   PPTImageProvider `json:"active_image_provider"`   // 当前启用的图片模型提供方
	MaxSlidesPerTask      int              `json:"max_slides_per_task"`     // 单个任务最多生成的 PPT 页数
	PowerCostPerSlide     int              `json:"power_cost_per_slide"`    // 每张 PPT 图片消耗的算力
	MaxConcurrentRequests int              `json:"max_concurrent_requests"` // 图片生成最大并发数
	QPSLimit              int              `json:"qps_limit"`               // 外部图片 API 的 QPS 限制

	// Nano Banana 配置
	NanoBananaApiURL         string `json:"nano_banana_api_url"`         // Nano Banana API 地址，如 https://xxx/v1/images/generations
	NanoBananaApiKey         string `json:"nano_banana_api_key"`         // Nano Banana API Key
	NanoBananaModel          string `json:"nano_banana_model"`           // 模型名称，如 nano-banana、nano-banana-hd
	NanoBananaResponseFormat string `json:"nano_banana_response_format"` // 响应格式：url 或 b64_json
	NanoBananaAspectRatio    string `json:"nano_banana_aspect_ratio"`    // 宽高比，如 1:1、4:3、16:9；空则默认 16:9

	// Doubao Seedream 配置
	SeedreamBaseURL      string `json:"seedream_base_url"`        // Seedream base url，例如：https://ark.cn-beijing.volces.com/api/v3
	SeedreamApiKey       string `json:"seedream_api_key"`         // Seedream API Key（ARK_API_KEY）
	SeedreamModel        string `json:"seedream_model"`           // Seedream 模型 ID，例如：doubao-seedream-5-0-260128
	SeedreamSize         string `json:"seedream_size"`            // 图片尺寸 WxH（如 1920x1080）或 2K/4K 等；空则默认 1920x1080（16:9）
	SeedreamOutputFormat string `json:"seedream_output_format"`   // 输出格式，例如：png
	SeedreamResponseType string `json:"seedream_response_format"` // 响应格式，例如：url
	SeedreamWatermark    bool   `json:"seedream_watermark"`       // 是否开启水印
}
