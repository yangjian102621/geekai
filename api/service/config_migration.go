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

	paymentConfig := types.PaymentConfig{
		Alipay: config.AlipayConfig,
		Epay:   config.GeekPayConfig,
		WxPay:  config.WechatPayConfig,
	}
	if err := s.saveConfig(types.ConfigKeyPayment, paymentConfig); err != nil {
		return err
	}

	return nil
}

// 迁移存储配置
func (s *ConfigMigrationService) migrateStorageConfig(config *types.AppConfig) error {

	ossConfig := types.OSSConfig{
		Active: config.OSS.Active,
		Local:  config.OSS.Local,
		Minio:  config.OSS.Minio,
		QiNiu:  config.OSS.QiNiu,
		AliYun: config.OSS.AliYun,
	}
	return s.saveConfig(types.ConfigKeyOss, ossConfig)
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
	if err := s.saveConfig(types.ConfigKeySmtp, smtpConfig); err != nil {
		return err
	}

	// 短信配置
	smsConfig := map[string]any{
		"active": config.SMS.Active,
		"ali": map[string]any{
			"access_key":    config.SMS.Ali.AccessKey,
			"access_secret": config.SMS.Ali.AccessSecret,
			"sign":          config.SMS.Ali.Sign,
			"code_temp_id":  config.SMS.Ali.CodeTempId,
		},
		"bao": map[string]any{
			"username":      config.SMS.Bao.Username,
			"password":      config.SMS.Bao.Password,
			"sign":          config.SMS.Bao.Sign,
			"code_template": config.SMS.Bao.CodeTemplate,
		},
	}
	return s.saveConfig(types.ConfigKeySms, smsConfig)
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
