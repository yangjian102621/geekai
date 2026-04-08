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

	"github.com/imroc/req/v3"
)

type GiteeAIModeration struct {
	config types.ModerationGiteeConfig
	apiURL string
}

func NewGiteeAIModeration(sysConfig *types.SystemConfig) *GiteeAIModeration {
	return &GiteeAIModeration{
		config: sysConfig.Moderation.Gitee,
		apiURL: "https://ai.gitee.com/v1/moderations",
	}
}

func (s *GiteeAIModeration) UpdateConfig(config types.ModerationGiteeConfig) {
	s.config = config
}

type GiteeAIModerationResult struct {
	ID      string                   `json:"id"`
	Model   string                   `json:"model"`
	Results []types.ModerationResult `json:"results"`
}

func (s *GiteeAIModeration) Moderate(text string) (types.ModerationResult, error) {

	body := map[string]any{
		"input": text,
		"model": s.config.Model,
	}
	var res GiteeAIModerationResult
	r, err := req.C().R().SetHeader("Authorization", "Bearer "+s.config.ApiKey).SetBody(body).SetSuccessResult(&res).Post(s.apiURL)
	if err != nil {
		return types.ModerationResult{}, err
	}

	if r.IsErrorState() {
		return types.ModerationResult{}, errors.New(r.String())
	}

	return res.Results[0], nil
}

var _ Service = (*GiteeAIModeration)(nil)
