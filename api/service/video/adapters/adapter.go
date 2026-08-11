package adapters

import "geekai/log"

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var logger = log.GetLogger()

// CreateTaskResponse 创建任务响应
type CreateTaskResponse struct {
	TaskId    string `json:"task_id"`    // 任务ID
	Channel   string `json:"channel"`    // 渠道标识
	Prompt    string `json:"prompt"`     // 优化后的提示词（如果有）
	State     string `json:"state"`      // 任务状态
	CreatedAt string `json:"created_at"` // 创建时间
}

// QueryTaskResponse 查询任务响应
type QueryTaskResponse struct {
	TaskId    string `json:"task_id"`    // 任务ID
	Status    string `json:"status"`     // 任务状态
	Progress  int    `json:"progress"`   // 进度（0-100）
	VideoURL  string `json:"video_url"`  // 视频URL
	Prompt    string `json:"prompt"`     // 提示词
	ErrMsg    string `json:"err_msg"`    // 错误信息
	StatusMsg string `json:"status_msg"` // 状态消息
	Output    string `json:"output"`     // 任务输出的原始信息（JSON字符串）
}
