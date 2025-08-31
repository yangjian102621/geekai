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

type TencentAIModeration struct {
	config types.ModerationTencentConfig
}

func NewTencentAIModeration(sysConfig *types.SystemConfig) *TencentAIModeration {
	return &TencentAIModeration{
		config: sysConfig.Moderation.Tencent,
	}
}

func (s *TencentAIModeration) UpdateConfig(config types.ModerationTencentConfig) {
	s.config = config
}

func (s *TencentAIModeration) Moderate(text string) (types.ModerationResult, error) {
	return types.ModerationResult{}, errors.New("not implemented")
}

var _ Service = (*TencentAIModeration)(nil)
