package oss

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"
	"strings"

	"geekai/log"
)

var logger = log.GetLogger()

// 默认缩略图模板（本地存储格式）
const DefaultThumbTemplate = "?imageView2/4/w/{width}/h/{height}/q/75"

type UploaderManager struct {
	local     *LocalStorage
	aliyun    *AliYunOss
	mini      *MiniOss
	qiniu     *QiNiuOss
	tencent   *TencentOss
	active    string
	ossConfig types.OSSConfig // 保存当前OSS配置
}

func NewUploaderManager(sysConfig *types.SystemConfig, local *LocalStorage, aliyun *AliYunOss, mini *MiniOss, qiniu *QiNiuOss, tencent *TencentOss) (*UploaderManager, error) {
	if sysConfig.OSS.Active == "" {
		sysConfig.OSS.Active = Local
	}

	return &UploaderManager{
		active:    sysConfig.OSS.Active,
		local:     local,
		aliyun:    aliyun,
		mini:      mini,
		qiniu:     qiniu,
		tencent:   tencent,
		ossConfig: sysConfig.OSS,
	}, nil
}

func (m *UploaderManager) GetUploadHandler() Uploader {
	switch m.active {
	case Local:
		return m.local
	case AliYun:
		return m.aliyun
	case Minio:
		return m.mini
	case QiNiu:
		return m.qiniu
	case Tencent:
		return m.tencent
	}
	return m.local
}

func (m *UploaderManager) UpdateConfig(config types.OSSConfig) {
	switch config.Active {
	case Local:
		m.local.UpdateConfig(config.Local)
	case AliYun:
		m.aliyun.UpdateConfig(config.AliYun)
	case Minio:
		m.mini.UpdateConfig(config.Minio)
	case QiNiu:
		m.qiniu.UpdateConfig(config.QiNiu)
	case Tencent:
		m.tencent.UpdateConfig(config.Tencent)
	}
	m.active = config.Active
	m.ossConfig = config
}

// GetThumbURL 根据原始图片URL和尺寸生成缩略图URL
// 如果模板为空，使用默认模板（本地存储格式）
// 如果明确设置为空字符串（表示不支持缩略图），返回原图URL
func (m *UploaderManager) GetThumbURL(originalURL string, width, height int) string {
	var template string

	// 根据当前激活的存储引擎获取对应的模板
	switch m.active {
	case Local:
		template = m.ossConfig.Local.ThumbTemplate
	case AliYun:
		template = m.ossConfig.AliYun.ThumbTemplate
	case Minio:
		template = m.ossConfig.Minio.ThumbTemplate
	case QiNiu:
		template = m.ossConfig.QiNiu.ThumbTemplate
	case Tencent:
		template = m.ossConfig.Tencent.ThumbTemplate
	default:
		template = m.ossConfig.Local.ThumbTemplate
	}

	// 如果模板为空，使用默认模板（兼容旧配置）
	if template == "" {
		template = DefaultThumbTemplate
	}

	// 替换变量
	thumbURL := strings.ReplaceAll(template, "{width}", fmt.Sprintf("%d", width))
	thumbURL = strings.ReplaceAll(thumbURL, "{height}", fmt.Sprintf("%d", height))

	// 拼接原始URL和缩略图参数
	return originalURL + thumbURL
}

// GetThumbTemplate 获取当前存储引擎的缩略图模板
func (m *UploaderManager) GetThumbTemplate() string {
	var template string

	switch m.active {
	case Local:
		template = m.ossConfig.Local.ThumbTemplate
	case AliYun:
		template = m.ossConfig.AliYun.ThumbTemplate
	case Minio:
		template = m.ossConfig.Minio.ThumbTemplate
	case QiNiu:
		template = m.ossConfig.QiNiu.ThumbTemplate
	case Tencent:
		template = m.ossConfig.Tencent.ThumbTemplate
	default:
		template = m.ossConfig.Local.ThumbTemplate
	}

	// 如果模板为空，返回默认模板
	if template == "" {
		return DefaultThumbTemplate
	}

	return template
}
