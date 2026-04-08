package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 文本审查
type ModerationConfig struct {
	Enable      bool                    `json:"enable"` // 是否启用文本审查
	Active      string                  `json:"active"`
	EnableGuide bool                    `json:"enable_guide"` // 是否启用模型引导提示词
	GuidePrompt string                  `json:"guide_prompt"` // 模型引导提示词
	Gitee       ModerationGiteeConfig   `json:"gitee"`
	Baidu       ModerationBaiduConfig   `json:"baidu"`
	Tencent     ModerationTencentConfig `json:"tencent"`
}

const (
	ModerationGitee   = "gitee"
	ModerationBaidu   = "baidu"
	ModerationTencent = "tencent"
)

// GiteeAI 文本审查配置
type ModerationGiteeConfig struct {
	ApiKey string `json:"api_key"`
	Model  string `json:"model"` // 文本审核模型
}

// 百度文本审查配置
type ModerationBaiduConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// 腾讯云文本审查配置
type ModerationTencentConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type ModerationResult struct {
	Flagged        bool               `json:"flagged"`
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
}

var ModerationCategories = map[string]string{
	"politic":  "内容涉及人物、事件或敏感的政治观点",
	"porn":     "明确的色情内容",
	"insult":   "具有侮辱、攻击性语言、人身攻击或冒犯性表达",
	"violence": "包含暴力、血腥、攻击行为或煽动暴力的言论",
	"illegal":  "涉及违法活动的内容，如诈骗、赌博等",
	"terror":   "宣扬恐怖主义、极端暴力或煽动恐怖行为的内容",
	"ad":       "垃圾广告或未经许可的推广内容",
	"spam":     "无意义重复内容或诱导性信息",
	"abuse":    "人身攻击、恶意辱骂或侮辱性言论",
	"polity":   "涉及国家政治、领导人或政策的违规讨论内容",
}

// 敏感词来源
const (
	ModerationSourceChat   = "chat"
	ModerationSourceMJ     = "mj"
	ModerationSourceDalle  = "dalle"
	ModerationSourceSD     = "sd"
	ModerationSourceSuno   = "suno"
	ModerationSourceVideo  = "video"
	ModerationSourceJiMeng = "jimeng"
)
