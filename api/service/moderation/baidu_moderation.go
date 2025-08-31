package moderation

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"geekai/core/types"
)

type BaiduAIModeration struct {
	config types.ModerationBaiduConfig
}

func NewBaiduAIModeration(sysConfig *types.SystemConfig) *BaiduAIModeration {
	return &BaiduAIModeration{
		config: sysConfig.Moderation.Baidu,
	}
}

func (s *BaiduAIModeration) UpdateConfig(config types.ModerationBaiduConfig) {
	s.config = config
}

func (s *BaiduAIModeration) Moderate(text string) (types.ModerationResult, error) {
	return types.ModerationResult{}, errors.New("not implemented")
}

var _ Service = (*BaiduAIModeration)(nil)
