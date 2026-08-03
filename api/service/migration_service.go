package service

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Copyright 2023 The Geek-AI Authors. All rights reserved.
// Use of this source code is governed by a Apache-2.0 license
// that can be found in the LICENSE file.
// @Author yangjian102621@163.com
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/store"
	"geekai/store/model"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// 迁移状态Redis key
	MigrationStatusKey = "config_migration:status"
	// 迁移完成标志
	MigrationCompleted = "completed"
)

// MigrationService 配置迁移服务
type MigrationService struct {
	db          *gorm.DB
	redisClient *redis.Client
	appConfig   *types.AppConfig
	levelDB     *store.LevelDB
}

func NewMigrationService(db *gorm.DB, redisClient *redis.Client, appConfig *types.AppConfig, levelDB *store.LevelDB) *MigrationService {
	return &MigrationService{
		db:          db,
		redisClient: redisClient,
		appConfig:   appConfig,
		levelDB:     levelDB,
	}
}

func (s *MigrationService) StartMigrate() {
	// 表结构同步必须在对外服务前完成，避免缺列导致业务报错
	s.TableMigration()
	go func() {
		_ = s.MigrateConfig(s.appConfig)
	}()
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

// 永不删除的保护列（大小写不敏感）
var protectedColumns = map[string]struct{}{
	"id":         {},
	"created_at": {},
	"updated_at": {},
}

// allModels 全部需要同步的数据表 model
func allModels() []any {
	return []any{
		&model.AdminUser{},
		&model.ApiKey{},
		&model.AppType{},
		&model.ChatApp{},
		&model.ChatItem{},
		&model.ChatMessage{},
		&model.ChatModel{},
		&model.Config{},
		&model.DallJob{},
		&model.File{},
		&model.Function{},
		&model.InviteCode{},
		&model.InviteLog{},
		&model.JimengJob{},
		&model.Menu{},
		&model.MidJourneyJob{},
		&model.Moderation{},
		&model.Order{},
		&model.PowerLog{},
		&model.Product{},
		&model.Redeem{},
		&model.SdJob{},
		&model.SunoJob{},
		&model.User{},
		&model.UserLoginLog{},
		&model.VideoJob{},
	}
}

// 数据表迁移：先处理字段重命名（保数据），再全量同步 schema
func (s *MigrationService) TableMigration() {
	logger.Info("开始数据表迁移...")
	s.renameColumns()
	if err := s.SyncAllModels(); err != nil {
		logger.Errorf("同步数据表字段失败: %v", err)
		return
	}
	logger.Info("数据表迁移完成")
}

// renameColumns 只处理「改名」场景：删旧加新会丢数据，必须先 Rename
func (s *MigrationService) renameColumns() {
	m := s.db.Migrator()

	if m.HasColumn(&model.JimengJob{}, "task_params") {
		_ = m.RenameColumn(&model.JimengJob{}, "task_params", "params")
	}
	if m.HasColumn(&model.Order{}, "pay_type") {
		_ = m.RenameColumn(&model.Order{}, "pay_type", "channel")
	}
	if m.HasColumn(&model.Config{}, "config_json") {
		_ = m.RenameColumn(&model.Config{}, "config_json", "value")
	}
	if m.HasColumn(&model.Config{}, "marker") {
		_ = m.RenameColumn(&model.Config{}, "marker", "name")
	}
	if m.HasIndex(&model.Config{}, "idx_chatgpt_configs_key") {
		_ = m.DropIndex(&model.Config{}, "idx_chatgpt_configs_key")
	}
	if m.HasIndex(&model.Config{}, "marker") {
		_ = m.DropIndex(&model.Config{}, "marker")
	}
}

// SyncAllModels 按 model 定义同步所有数据表：缺列新建，多余列删除
func (s *MigrationService) SyncAllModels() error {
	var firstErr error
	for _, m := range allModels() {
		if err := s.syncModel(m); err != nil {
			logger.Errorf("同步 model %T 失败: %v", m, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (s *MigrationService) syncModel(dst any) error {
	tableName := s.tableName(dst)
	if err := s.db.AutoMigrate(dst); err != nil {
		return fmt.Errorf("AutoMigrate %s: %w", tableName, err)
	}
	if err := s.dropUnusedColumns(dst); err != nil {
		return fmt.Errorf("drop unused columns %s: %w", tableName, err)
	}
	logger.Infof("已同步数据表: %s", tableName)
	return nil
}

func (s *MigrationService) dropUnusedColumns(dst any) error {
	dbCols, err := s.db.Migrator().ColumnTypes(dst)
	if err != nil {
		return err
	}
	modelCols, err := s.modelColumnNames(dst)
	if err != nil {
		return err
	}

	for _, col := range dbCols {
		name := col.Name()
		if s.isProtectedColumn(name) {
			continue
		}
		if _, ok := modelCols[strings.ToLower(name)]; ok {
			continue
		}
		logger.Infof("删除多余字段: %s.%s", s.tableName(dst), name)
		if err := s.db.Migrator().DropColumn(dst, name); err != nil {
			return fmt.Errorf("DropColumn %s: %w", name, err)
		}
	}
	return nil
}

func (s *MigrationService) modelColumnNames(dst any) (map[string]struct{}, error) {
	parsed, err := schema.Parse(dst, &schemaCache, s.db.Config.NamingStrategy)
	if err != nil {
		return nil, err
	}
	cols := make(map[string]struct{}, len(parsed.Fields))
	for _, field := range parsed.Fields {
		if field.DBName == "" || field.IgnoreMigration {
			continue
		}
		cols[strings.ToLower(field.DBName)] = struct{}{}
	}
	return cols, nil
}

func (s *MigrationService) isProtectedColumn(name string) bool {
	_, ok := protectedColumns[strings.ToLower(name)]
	return ok
}

func (s *MigrationService) tableName(dst any) string {
	stmt := &gorm.Statement{DB: s.db}
	if err := stmt.Parse(dst); err != nil {
		return fmt.Sprintf("%T", dst)
	}
	return stmt.Schema.Table
}

// schema.Parse 进程内复用的 schema cache
var schemaCache sync.Map

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
