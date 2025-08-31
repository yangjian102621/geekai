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
	"geekai/store"
	"geekai/store/model"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	// 迁移状态Redis key
	MigrationStatusKey = "config_migration:status"
	// 迁移完成标志
	MigrationCompleted = "completed"
)

// MigrationService 配置迁移服务
type MigrationService struct {
	db             *gorm.DB
	redisClient    *redis.Client
	appConfig      *types.AppConfig
	levelDB        *store.LevelDB
	licenseService *LicenseService
}

func NewMigrationService(db *gorm.DB, redisClient *redis.Client, appConfig *types.AppConfig, levelDB *store.LevelDB, licenseService *LicenseService) *MigrationService {
	return &MigrationService{
		db:             db,
		redisClient:    redisClient,
		appConfig:      appConfig,
		levelDB:        levelDB,
		licenseService: licenseService,
	}
}

func (s *MigrationService) StartMigrate() {
	go func() {
		s.MigrateConfig(s.appConfig)
		s.TableMigration()
		s.MigrateLicense()
	}()
}

// 迁移 License
func (s *MigrationService) MigrateLicense() {
	key := "migrate:license"
	if s.redisClient.Get(context.Background(), key).Val() == "1" {
		logger.Info("License 已迁移，跳过迁移")
		return
	}

	logger.Info("开始迁移 License...")
	var license types.License
	err := s.levelDB.Get(types.LicenseKey, &license)
	if err != nil {
		logger.Errorf("迁移 License 失败: %v", err)
		return
	}
	logger.Infof("迁移 License: %+v", license)
	if err := s.saveConfig(types.ConfigKeyLicense, license); err != nil {
		logger.Errorf("迁移 License 失败: %v", err)
		return
	}
	s.licenseService.SetLicense(license.Key)
	logger.Info("迁移 License 完成")
	s.redisClient.Set(context.Background(), key, "1", 0)
}

// 数据表迁移
func (s *MigrationService) TableMigration() {
	// 订单字段整理
	if s.db.Migrator().HasColumn(&model.Order{}, "pay_type") {
		s.db.Migrator().RenameColumn(&model.Order{}, "pay_type", "channel")
	}
	if !s.db.Migrator().HasColumn(&model.Order{}, "checked") {
		s.db.Migrator().AddColumn(&model.Order{}, "checked")
	}

	// 重命名 config 表字段
	if s.db.Migrator().HasColumn(&model.Config{}, "config_json") {
		s.db.Migrator().RenameColumn(&model.Config{}, "config_json", "value")
	}
	if s.db.Migrator().HasColumn(&model.Config{}, "marker") {
		s.db.Migrator().RenameColumn(&model.Config{}, "marker", "name")
	}
	if s.db.Migrator().HasIndex(&model.Config{}, "idx_chatgpt_configs_key") {
		s.db.Migrator().DropIndex(&model.Config{}, "idx_chatgpt_configs_key")
	}
	if s.db.Migrator().HasIndex(&model.Config{}, "marker") {
		s.db.Migrator().DropIndex(&model.Config{}, "marker")
	}

	// 手动删除字段
	if s.db.Migrator().HasColumn(&model.Order{}, "deleted_at") {
		s.db.Migrator().DropColumn(&model.Order{}, "deleted_at")
	}
	if s.db.Migrator().HasColumn(&model.ChatItem{}, "deleted_at") {
		s.db.Migrator().DropColumn(&model.ChatItem{}, "deleted_at")
	}
	if s.db.Migrator().HasColumn(&model.ChatMessage{}, "deleted_at") {
		s.db.Migrator().DropColumn(&model.ChatMessage{}, "deleted_at")
	}
	if s.db.Migrator().HasColumn(&model.User{}, "chat_config") {
		s.db.Migrator().DropColumn(&model.User{}, "chat_config")
	}
	if s.db.Migrator().HasColumn(&model.ChatModel{}, "category") {
		s.db.Migrator().DropColumn(&model.ChatModel{}, "category")
	}
	if s.db.Migrator().HasColumn(&model.ChatModel{}, "description") {
		s.db.Migrator().DropColumn(&model.ChatModel{}, "description")
	}
	if s.db.Migrator().HasColumn(&model.Product{}, "discount") {
		s.db.Migrator().DropColumn(&model.Product{}, "discount")
	}
	if s.db.Migrator().HasColumn(&model.Product{}, "days") {
		s.db.Migrator().DropColumn(&model.Product{}, "days")
	}
	if s.db.Migrator().HasColumn(&model.Product{}, "app_url") {
		s.db.Migrator().DropColumn(&model.Product{}, "app_url")
	}
	if s.db.Migrator().HasColumn(&model.Product{}, "url") {
		s.db.Migrator().DropColumn(&model.Product{}, "url")
	}
}

// 迁移配置数据
func (s *MigrationService) MigrateConfig(config *types.AppConfig) error {

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

	logger.Info("配置迁移完成")
	return nil
}

// 迁移支付配置
func (s *MigrationService) migratePaymentConfig(config *types.AppConfig) error {

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
func (s *MigrationService) migrateStorageConfig(config *types.AppConfig) error {

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
func (s *MigrationService) migrateCommunicationConfig(config *types.AppConfig) error {
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
		"active": strings.ToLower(config.SMS.Active),
		"aliyun": map[string]any{
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
func (s *MigrationService) saveConfig(key string, config any) error {
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
