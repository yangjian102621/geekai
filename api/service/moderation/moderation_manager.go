package moderation

import (
	"geekai/core/types"

	logger2 "geekai/logger"
)

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

var logger = logger2.GetLogger()

type Service interface {
	Moderate(text string) (types.ModerationResult, error)
}

type ServiceManager struct {
	gitee   *GiteeAIModeration
	baidu   *BaiduAIModeration
	tencent *TencentAIModeration
	active  string
}

func NewServiceManager(gitee *GiteeAIModeration, baidu *BaiduAIModeration, tencent *TencentAIModeration) *ServiceManager {
	return &ServiceManager{
		gitee:   gitee,
		baidu:   baidu,
		tencent: tencent,
	}
}

func (s *ServiceManager) GetService() Service {
	switch s.active {
	case types.ModerationBaidu:
		return s.baidu
	case types.ModerationTencent:
		return s.tencent
	default:
		return s.gitee
	}
}

func (s *ServiceManager) UpdateConfig(config types.ModerationConfig) {
	switch config.Active {
	case types.ModerationGitee:
		s.gitee.UpdateConfig(config.Gitee)
	case types.ModerationBaidu:
		s.baidu.UpdateConfig(config.Baidu)
	case types.ModerationTencent:
		s.tencent.UpdateConfig(config.Tencent)
	}
	s.active = config.Active
}
