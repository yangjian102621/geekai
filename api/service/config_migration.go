package service

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Copyright 2023 The Geek-AI Authors. All rights reserved.
// Use of this source code is governed by a Apache-2.0 license
// that can be found in the LICENSE file.
// @Author yangjian102621@163.com
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"encoding/json"
	"geekai/core/types"
	"geekai/store/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	// 迁移状态Redis key
	MigrationStatusKey = "config_migration:status"
	// 迁移完成标志
	MigrationCompleted = "completed"
)

// ConfigMigrationService 配置迁移服务
type ConfigMigrationService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewConfigMigrationService(db *gorm.DB, redisClient *redis.Client) *ConfigMigrationService {
	return &ConfigMigrationService{
		db:          db,
		redisClient: redisClient,
	}
}

// MigrateFromConfig 从 config.toml 迁移配置到数据库（仅首次启动时执行）
func (s *ConfigMigrationService) MigrateFromConfig(config *types.AppConfig) error {
	// 检查是否已经迁移过
	if s.isMigrationCompleted() {
		logger.Info("配置迁移已完成，跳过迁移")
		return nil
	}

	logger.Info("开始迁移配置到数据库...")

	// 迁移支付配置
	if err := s.migratePaymentConfig(config); err != nil {
		logger.Errorf("迁移支付配置失败: %v", err)
		return err
	}

	// 迁移存储配置
	if err := s.migrateStorageConfig(config); err != nil {
		logger.Errorf("迁移存储配置失败: %v", err)
		return err
	}

	// 迁移通信配置
	if err := s.migrateCommunicationConfig(config); err != nil {
		logger.Errorf("迁移通信配置失败: %v", err)
		return err
	}

	// 迁移API配置
	if err := s.migrateApiConfig(config); err != nil {
		logger.Errorf("迁移API配置失败: %v", err)
		return err
	}

	// 标记迁移完成
	if err := s.markMigrationCompleted(); err != nil {
		logger.Errorf("标记迁移完成失败: %v", err)
		return err
	}

	logger.Info("配置迁移完成")
	return nil
}

// 检查是否已经迁移完成
func (s *ConfigMigrationService) isMigrationCompleted() bool {
	ctx := context.Background()
	status, err := s.redisClient.Get(ctx, MigrationStatusKey).Result()
	if err != nil {
		// Redis中没有找到标志，说明未迁移过
		return false
	}
	return status == MigrationCompleted
}

// 标记迁移完成
func (s *ConfigMigrationService) markMigrationCompleted() error {
	ctx := context.Background()
	// 设置迁移完成标志，永不过期
	return s.redisClient.Set(ctx, MigrationStatusKey, MigrationCompleted, 0).Err()
}

// 迁移支付配置
func (s *ConfigMigrationService) migratePaymentConfig(config *types.AppConfig) error {
	// 支付宝配置
	alipayConfig := map[string]any{
		"enabled":           config.AlipayConfig.Enabled,
		"sand_box":          config.AlipayConfig.SandBox,
		"app_id":            config.AlipayConfig.AppId,
		"private_key":       config.AlipayConfig.PrivateKey,
		"alipay_public_key": config.AlipayConfig.AlipayPublicKey,
		"notify_url":        config.AlipayConfig.NotifyURL,
		"return_url":        config.AlipayConfig.ReturnURL,
	}
	if err := s.saveConfig("alipay", alipayConfig); err != nil {
		return err
	}

	// 微信支付配置
	wechatConfig := map[string]any{
		"enabled":     config.WechatPayConfig.Enabled,
		"app_id":      config.WechatPayConfig.AppId,
		"mch_id":      config.WechatPayConfig.MchId,
		"serial_no":   config.WechatPayConfig.SerialNo,
		"private_key": config.WechatPayConfig.PrivateKey,
		"api_v3_key":  config.WechatPayConfig.ApiV3Key,
		"notify_url":  config.WechatPayConfig.NotifyURL,
	}
	if err := s.saveConfig("wechat", wechatConfig); err != nil {
		return err
	}

	// 虎皮椒配置
	hupiConfig := map[string]any{
		"enabled":    config.HuPiPayConfig.Enabled,
		"app_id":     config.HuPiPayConfig.AppId,
		"app_secret": config.HuPiPayConfig.AppSecret,
		"api_url":    config.HuPiPayConfig.ApiURL,
		"notify_url": config.HuPiPayConfig.NotifyURL,
		"return_url": config.HuPiPayConfig.ReturnURL,
	}
	if err := s.saveConfig("hupi", hupiConfig); err != nil {
		return err
	}

	// GeekPay配置
	geekpayConfig := map[string]any{
		"enabled":     config.GeekPayConfig.Enabled,
		"app_id":      config.GeekPayConfig.AppId,
		"private_key": config.GeekPayConfig.PrivateKey,
		"api_url":     config.GeekPayConfig.ApiURL,
		"notify_url":  config.GeekPayConfig.NotifyURL,
		"return_url":  config.GeekPayConfig.ReturnURL,
		"methods":     config.GeekPayConfig.Methods,
	}
	if err := s.saveConfig("geekpay", geekpayConfig); err != nil {
		return err
	}

	return nil
}

// 迁移存储配置
func (s *ConfigMigrationService) migrateStorageConfig(config *types.AppConfig) error {
	ossConfig := map[string]any{
		"active": config.OSS.Active,
		"local": map[string]any{
			"base_path": config.OSS.Local.BasePath,
			"base_url":  config.OSS.Local.BaseURL,
		},
		"minio": map[string]any{
			"endpoint":      config.OSS.Minio.Endpoint,
			"access_key":    config.OSS.Minio.AccessKey,
			"access_secret": config.OSS.Minio.AccessSecret,
			"bucket":        config.OSS.Minio.Bucket,
			"use_ssl":       config.OSS.Minio.UseSSL,
			"domain":        config.OSS.Minio.Domain,
		},
		"qiniu": map[string]any{
			"zone":          config.OSS.QiNiu.Zone,
			"access_key":    config.OSS.QiNiu.AccessKey,
			"access_secret": config.OSS.QiNiu.AccessSecret,
			"bucket":        config.OSS.QiNiu.Bucket,
			"domain":        config.OSS.QiNiu.Domain,
		},
		"aliyun": map[string]any{
			"endpoint":      config.OSS.AliYun.Endpoint,
			"access_key":    config.OSS.AliYun.AccessKey,
			"access_secret": config.OSS.AliYun.AccessSecret,
			"bucket":        config.OSS.AliYun.Bucket,
			"sub_dir":       config.OSS.AliYun.SubDir,
			"domain":        config.OSS.AliYun.Domain,
		},
	}
	return s.saveConfig("oss", ossConfig)
}

// 迁移通信配置
func (s *ConfigMigrationService) migrateCommunicationConfig(config *types.AppConfig) error {
	// SMTP配置
	smtpConfig := map[string]any{
		"use_tls":  config.SmtpConfig.UseTls,
		"host":     config.SmtpConfig.Host,
		"port":     config.SmtpConfig.Port,
		"app_name": config.SmtpConfig.AppName,
		"from":     config.SmtpConfig.From,
		"password": config.SmtpConfig.Password,
	}
	if err := s.saveConfig("smtp", smtpConfig); err != nil {
		return err
	}

	// 短信配置
	smsConfig := map[string]any{
		"active": config.SMS.Active,
		"ali": map[string]any{
			"access_key":    config.SMS.Ali.AccessKey,
			"access_secret": config.SMS.Ali.AccessSecret,
			"product":       config.SMS.Ali.Product,
			"domain":        config.SMS.Ali.Domain,
			"sign":          config.SMS.Ali.Sign,
			"code_temp_id":  config.SMS.Ali.CodeTempId,
		},
		"bao": map[string]any{
			"username":      config.SMS.Bao.Username,
			"password":      config.SMS.Bao.Password,
			"domain":        config.SMS.Bao.Domain,
			"sign":          config.SMS.Bao.Sign,
			"code_template": config.SMS.Bao.CodeTemplate,
		},
	}
	return s.saveConfig("sms", smsConfig)
}

// 迁移API配置
func (s *ConfigMigrationService) migrateApiConfig(config *types.AppConfig) error {
	apiConfig := map[string]any{
		"api_url": config.ApiConfig.ApiURL,
		"app_id":  config.ApiConfig.AppId,
		"token":   config.ApiConfig.Token,
		"jimeng_config": map[string]any{
			"access_key": config.ApiConfig.JimengConfig.AccessKey,
			"secret_key": config.ApiConfig.JimengConfig.SecretKey,
		},
	}
	return s.saveConfig("api", apiConfig)
}

// 保存配置到数据库
func (s *ConfigMigrationService) saveConfig(key string, config any) error {
	// 检查是否已存在
	var existingConfig model.Config
	if err := s.db.Where("name", key).First(&existingConfig).Error; err == nil {
		// 配置已存在，跳过
		logger.Infof("配置 %s 已存在，跳过迁移", key)
		return nil
	}

	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 保存到数据库
	newConfig := model.Config{
		Name:  key,
		Value: string(configJSON),
	}
	if err := s.db.Create(&newConfig).Error; err != nil {
		return err
	}

	logger.Infof("成功迁移配置 %s", key)
	return nil
}
