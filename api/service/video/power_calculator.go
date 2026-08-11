package video

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/utils"
	"gorm.io/gorm"
)

// GetVideoConfig 从数据库获取视频配置
func GetVideoConfig(db *gorm.DB) (*types.VideoConfig, error) {
	var config model.Config
	err := db.Where("name", types.ConfigKeyVideo).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("视频配置不存在，请在管理后台配置")
		}
		return nil, fmt.Errorf("获取视频配置失败: %v", err)
	}

	var videoConfig types.VideoConfig
	err = utils.JsonDecode(config.Value, &videoConfig)
	if err != nil {
		return nil, fmt.Errorf("解析视频配置失败: %v", err)
	}

	return &videoConfig, nil
}

// GetModelPowerConfig 获取指定模型的算力配置
func GetModelPowerConfig(db *gorm.DB, modelKey string) (*types.VideoModelPower, error) {
	config, err := GetVideoConfig(db)
	if err != nil {
		return nil, err
	}

	modelPower, ok := config.VideoPowers[modelKey]
	if !ok {
		return nil, fmt.Errorf("模型 %s 的算力配置不存在", modelKey)
	}

	return &modelPower, nil
}

// CalculatePower 根据 modelKey 和 priceKey 计算算力
// modelKey: 模型标识（如 "veo-2.0", "sora-2.0"）
// priceKey: 价格键（如 "fixed", "5_720P", "std_5_sound" 等）
func CalculatePower(db *gorm.DB, modelKey string, priceKey string) (int, error) {
	if priceKey == "" {
		return 0, errors.New("priceKey 不能为空")
	}

	modelPower, err := GetModelPowerConfig(db, modelKey)
	if err != nil {
		return 0, err
	}

	power, ok := modelPower.PowerConfig[priceKey]
	if !ok {
		return 0, fmt.Errorf("模型 %s 的价格配置 %s 不存在", modelKey, priceKey)
	}

	if power <= 0 {
		return 0, fmt.Errorf("模型 %s 的价格配置 %s 的值无效", modelKey, priceKey)
	}

	return power, nil
}
