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
	"fmt"
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
		license = types.License{
			Key:       "",
			MachineId: "",
			Configs:   types.LicenseConfig{UserNum: 0, DeCopy: false},
			ExpiredAt: 0,
			IsActive:  false,
		}
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

// 迁移配置内容
func (s *MigrationService) MigrateConfigContent() error {
	// 用户协议
	if err := s.saveConfig(types.ConfigKeyPrivacy, map[string]string{
		"content": "用户协议内容",
	}); err != nil {
		return fmt.Errorf("迁移配置内容失败: %v", err)
	}
	// 隐私政策
	if err := s.saveConfig(types.ConfigKeyAgreement, map[string]string{
		"content": "隐私政策内容",
	}); err != nil {
		return fmt.Errorf("迁移配置内容失败: %v", err)
	}
	// 思维导图
	if err := s.saveConfig(types.ConfigKeyMarkMap, map[string]string{
		"content": `# GeekAI 演示站

- 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
- 基于 Websocket 实现，完美的打字机体验。
- 内置了各种预训练好的角色应用,轻松满足你的各种聊天和应用需求。
- 支持 OPenAI，Azure，文心一言，讯飞星火，清华 ChatGLM等多个大语言模型。
- 支持 MidJourney / Stable Diffusion AI 绘画集成，开箱即用。
- 支持使用个人微信二维码作为充值收费的支付渠道，无需企业支付通道。
- 已集成支付宝支付功能，微信支付，支持多种会员套餐和点卡购买功能。
- 集成插件 API 功能，可结合大语言模型的 function 功能开发各种强大的插件。`,
	}); err != nil {
		return fmt.Errorf("迁移配置内容失败: %v", err)
	}

	// 微信登录配置
	if err := s.saveConfig(types.ConfigKeyWxLogin, map[string]string{
		"api_key":    "",
		"notify_url": "",
		"enabled":    "false",
	}); err != nil {
		return fmt.Errorf("迁移配置内容失败: %v", err)
	}

	// 验证码配置
	if err := s.saveConfig(types.ConfigKeyCaptcha, map[string]string{
		"api_key": "",
		"type":    "dot",
		"enabled": "false",
	}); err != nil {
		return fmt.Errorf("迁移配置内容失败: %v", err)
	}

	// 文本审核
	if err := s.saveConfig(types.ConfigKeyModeration, map[string]any{
		"enable":       "false",
		"active":       "gitee",
		"enable_guide": "false",
		"guide_prompt": "",
		"gitee": map[string]string{
			"api_key": "",
			"model":   "Security-semantic-filtering",
		},
		"baidu": map[string]string{
			"access_key": "",
			"secret_key": "",
		},
		"tencent": map[string]string{
			"access_key": "",
			"secret_key": "",
		},
	}); err != nil {
		return fmt.Errorf("迁移配置内容失败: %v", err)
	}

	return nil
}

// 数据表迁移
func (s *MigrationService) TableMigration() {
	// 新数据表
	s.db.AutoMigrate(&model.Moderation{})
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

	// 迁移配置内容
	if err := s.MigrateConfigContent(); err != nil {
		logger.Errorf("迁移配置内容失败: %v", err)
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
