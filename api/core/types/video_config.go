package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// VideoConfig 视频生成配置（存储在 config 表，name = 'video'）
type VideoConfig struct {
	ApiURL      string                     `json:"api_url"`      // API 地址
	ApiKey      string                     `json:"api_key"`      // API 密钥
	VideoPowers map[string]VideoModelPower `json:"video_powers"` // 模型算力配置
}

// VideoModelPower 单个模型的算力配置
// PowerConfig 说明：
// 新的价格配置方式：根据模型的 priceParams 生成笛卡尔乘积，每个组合对应一个价格
// 固定价格示例：{"fixed": 20}
// 多参数组合示例：{"5_720P": 10, "5_1080P": 20, "10_720P": 10, "10_1080P": 40}
// 复杂组合示例：{"std_5_sound": 10, "std_5_silent": 5, "pro_10_sound": 20, "pro_10_silent": 10}
type VideoModelPower struct {
	Provider    string         `json:"provider"`     // 服务提供商（不带版本号：veo, sora, luma）
	Model       string         `json:"model"`        // 模型名称（带版本号：veo-2.0, sora-2.0）
	PowerConfig map[string]int `json:"power_config"` // 算力配置（基于 priceParams 的笛卡尔乘积）
	ApiKeyType  string         `json:"api_key_type"` // ApiKey 表的 type 字段（可选，用于多 API Key 场景）
}

// 视频任务状态常量
const (
	VideoStatusPending     = "pending"     // 等待处理
	VideoStatusInProgress  = "in_progress" // 处理中
	VideoStatusDownloading = "downloading" // 视频下载中
	VideoStatusSuccess     = "success"     // 成功
	VideoStatusFailed      = "failed"      // 失败
)
