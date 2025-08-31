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
}
