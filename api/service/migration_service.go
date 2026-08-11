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
	db          *gorm.DB
	redisClient *redis.Client
	appConfig   *types.AppConfig
}

func NewMigrationService(db *gorm.DB, redisClient *redis.Client, appConfig *types.AppConfig) *MigrationService {
	return &MigrationService{
		db:          db,
		redisClient: redisClient,
		appConfig:   appConfig,
	}
}

func (s *MigrationService) StartMigrate() {
	// 表结构迁移必须在业务服务启动前完成，避免新表和新列尚未创建就被后台任务查询。
	s.TableMigration()
	go func() {
		s.MigrateConfig(s.appConfig)
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

// fullTableMigration 第一步：全量表迁移，同步所有表结构（新增表、新增字段、字段类型）
// 适用于首次安装或导入旧版数据库后同步 schema，AutoMigrate 会补齐缺失的表和列
func (s *MigrationService) fullTableMigration() {
	logger.Info("执行全量表迁移（同步 schema）...")
	models := []any{
		&model.Config{},
		&model.AdminUser{},
		&model.ChatApp{},
		&model.ApiKey{},
		&model.AppType{},
		&model.ChatModel{},
		&model.User{},
		&model.ChatItem{},
		&model.ChatMessage{},
		&model.Order{},
		&model.Product{},
		&model.Function{},
		&model.Menu{},
		&model.InviteCode{},
		&model.InviteLog{},
		&model.Redeem{},
		&model.PowerLog{},
		&model.File{},
		&model.UserLoginLog{},
		&model.MidJourneyJob{},
		&model.SunoJob{},
		&model.VideoJob{},
		&model.JimengJob{},
		&model.PPTJob{},
		&model.Moderation{},
		&model.ImageJob{},
	}
	if err := s.db.AutoMigrate(models...); err != nil {
		logger.Errorf("全量表迁移失败: %v", err)
		return
	}
	logger.Info("全量表迁移完成")
}

// fixTableConstraints 第一步之后：根据模型定义修复各表的主键、自增和关键索引
// 主要解决初始化 SQL 中缺少 AUTO_INCREMENT 或 PRIMARY KEY 导致插入失败的问题
func (s *MigrationService) fixTableConstraints() {
	logger.Info("开始修复各表的主键、自增属性和索引...")

	// 当前数据库名
	var dbName string
	if err := s.db.Raw("SELECT DATABASE()").Scan(&dbName).Error; err != nil {
		logger.Errorf("获取当前数据库名失败: %v", err)
		return
	}
	if dbName == "" {
		logger.Warn("当前连接未选择数据库，跳过约束修复")
		return
	}

	// === 修复所有包含 id 字段的表的主键 + 自增 ===
	type columnInfo struct {
		TableName  string `gorm:"column:TABLE_NAME"`
		ColumnName string `gorm:"column:COLUMN_NAME"`
		ColumnKey  string `gorm:"column:COLUMN_KEY"`
		Extra      string `gorm:"column:EXTRA"`
		DataType   string `gorm:"column:DATA_TYPE"`
	}

	var idColumns []columnInfo
	if err := s.db.Raw(`
		SELECT TABLE_NAME, COLUMN_NAME, COLUMN_KEY, EXTRA, DATA_TYPE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = ? AND COLUMN_NAME = 'id'
	`, dbName).Scan(&idColumns).Error; err != nil {
		logger.Errorf("查询各表 id 字段信息失败: %v", err)
	} else {
		for _, col := range idColumns {
			if col.ColumnName == "" {
				continue
			}

			// 检查当前表是否已经存在主键
			type pkInfo struct {
				ColumnName string `gorm:"column:COLUMN_NAME"`
			}
			var pkColumns []pkInfo
			if err := s.db.Raw(`
				SELECT COLUMN_NAME
				FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND CONSTRAINT_NAME = 'PRIMARY'
			`, dbName, col.TableName).Scan(&pkColumns).Error; err != nil {
				logger.Errorf("查询表 %s 的主键信息失败: %v", col.TableName, err)
				continue
			}

			hasPK := len(pkColumns) > 0
			pkOnIdOnly := hasPK && len(pkColumns) == 1 && pkColumns[0].ColumnName == "id"

			needAlter := false
			switch {
			case !hasPK:
				// 没有任何主键：允许将 id 设置为自增主键
				if col.ColumnKey != "PRI" || col.Extra == "" || !strings.Contains(col.Extra, "auto_increment") {
					needAlter = true
				}
			case pkOnIdOnly:
				// 只有 id 作为主键：只补充自增属性
				if col.Extra == "" || !strings.Contains(col.Extra, "auto_increment") {
					needAlter = true
				}
			default:
				// 已存在非 id 或组合主键：避免破坏原有主键，直接跳过
				logger.Infof("表 %s 已存在非 id 主键，跳过 id 自增主键修复", col.TableName)
			}
			if !needAlter {
				continue
			}

			// 维持原来的数据类型，避免与历史 SQL 冲突
			dataType := col.DataType
			if dataType == "" {
				dataType = "int"
			}

			alterSQL := fmt.Sprintf(
				"ALTER TABLE `%s` MODIFY COLUMN id %s NOT NULL AUTO_INCREMENT PRIMARY KEY",
				col.TableName,
				dataType,
			)
			if err := s.db.Exec(alterSQL).Error; err != nil {
				logger.Errorf("修复表 %s 的 id 自增主键失败: %v", col.TableName, err)
			} else {
				logger.Infof("已修复表 %s 的 id 为 AUTO_INCREMENT PRIMARY KEY", col.TableName)
			}
		}
	}

	// === 为 geekai_users.username 补充唯一索引（根据模型 uniqueIndex 定义） ===
	var usernameUniqueCount int64
	err := s.db.Raw(`
		SELECT COUNT(1)
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = ?
		  AND TABLE_NAME = 'geekai_users'
		  AND COLUMN_NAME = 'username'
		  AND NON_UNIQUE = 0
	`, dbName).Scan(&usernameUniqueCount).Error
	if err != nil {
		logger.Errorf("检查 geekai_users.username 唯一索引失败: %v", err)
	} else if usernameUniqueCount == 0 {
		// 索引名尽量固定，避免重复创建
		if err := s.db.Exec("ALTER TABLE geekai_users ADD UNIQUE KEY idx_geekai_users_username (username)").Error; err != nil {
			logger.Errorf("创建 geekai_users.username 唯一索引失败: %v", err)
		} else {
			logger.Info("已为 geekai_users.username 创建唯一索引 idx_geekai_users_username")
		}
	} else {
		logger.Info("geekai_users.username 唯一索引已存在，跳过创建")
	}

	logger.Info("关键表主键、自增和索引修复完成")
}

// incrementalTableMigration 第二步：增量迁移，仅处理删除字段与数据迁移
// AutoMigrate 不会删除列，故需在此显式 DropColumn；字段重命名需先拷贝数据再删除旧列
func (s *MigrationService) incrementalTableMigration() {
	logger.Info("执行增量迁移（删除字段 + 数据迁移）...")

	// ========== 字段重命名：全量迁移已添加新列，需将旧列数据拷贝到新列后删除旧列 ==========

	// ChatApp(geekai_chat_roles): context_json -> system_prompt 历史数据迁移
	if s.db.Migrator().HasColumn(&model.ChatApp{}, "context_json") {
		// 将旧列 context_json 的值拷贝到 system_prompt（NULL 转为空字符串，保证 NOT NULL 约束）
		if err := s.db.Exec(`
			UPDATE geekai_chat_roles 
			SET system_prompt = IFNULL(NULLIF(TRIM(COALESCE(context_json, '')), ''), '')
		`).Error; err != nil {
			logger.Errorf("迁移 geekai_chat_roles.context_json -> system_prompt 失败: %v", err)
		} else {
			if err := s.db.Migrator().DropColumn(&model.ChatApp{}, "context_json"); err != nil {
				logger.Errorf("删除 geekai_chat_roles.context_json 失败: %v", err)
			} else {
				logger.Info("geekai_chat_roles: context_json 已迁移至 system_prompt 并删除旧列")
			}
		}
	}

	// ChatApp: 将 user_id 为 NULL 的历史记录置为 0（系统内置）
	if s.db.Migrator().HasColumn(&model.ChatApp{}, "user_id") {
		if err := s.db.Exec(`UPDATE geekai_chat_roles SET user_id = 0 WHERE user_id IS NULL`).Error; err != nil {
			logger.Errorf("初始化 geekai_chat_roles.user_id 失败: %v", err)
		}
	}

	// ChatApp(geekai_chat_roles): 删除 marker 列（应用仅通过 id 区分）
	var hasMarker int
	if s.db.Raw("SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'geekai_chat_roles' AND COLUMN_NAME = 'marker'").Scan(&hasMarker).Error == nil && hasMarker > 0 {
		// 先删除可能存在的唯一索引（不同版本 SQL 索引名不同）
		for _, idxName := range []string{"marker", "idx_chatgpt_chat_roles_marker", "idx_chatgpt_chat_roles_key", "idx_geekai_chat_roles_marker"} {
			_ = s.db.Exec(fmt.Sprintf("ALTER TABLE geekai_chat_roles DROP INDEX `%s`", idxName)).Error
		}
		if err := s.db.Exec("ALTER TABLE geekai_chat_roles DROP COLUMN marker").Error; err != nil {
			logger.Errorf("删除 geekai_chat_roles.marker 失败: %v", err)
		} else {
			logger.Info("geekai_chat_roles: 已删除 marker 列")
		}
	}

	// Config: config_json -> value, marker -> name
	if s.db.Migrator().HasColumn(&model.Config{}, "config_json") {
		s.db.Exec("UPDATE geekai_configs SET `value` = config_json WHERE config_json IS NOT NULL AND config_json != ''")
		s.db.Migrator().DropColumn(&model.Config{}, "config_json")
	}
	if s.db.Migrator().HasColumn(&model.Config{}, "marker") {
		s.db.Exec("UPDATE geekai_configs SET `name` = marker WHERE marker IS NOT NULL AND marker != ''")
		s.db.Migrator().DropColumn(&model.Config{}, "marker")
	}
	if s.db.Migrator().HasIndex(&model.Config{}, "idx_chatgpt_configs_key") {
		s.db.Migrator().DropIndex(&model.Config{}, "idx_chatgpt_configs_key")
	}
	if s.db.Migrator().HasIndex(&model.Config{}, "marker") {
		s.db.Migrator().DropIndex(&model.Config{}, "marker")
	}

	// Order: pay_type -> channel
	if s.db.Migrator().HasColumn(&model.Order{}, "pay_type") {
		s.db.Exec("UPDATE geekai_orders SET channel = pay_type WHERE pay_type IS NOT NULL AND pay_type != ''")
		s.db.Migrator().DropColumn(&model.Order{}, "pay_type")
	}

	// JimengJob: task_params -> params
	if s.db.Migrator().HasColumn(&model.JimengJob{}, "task_params") {
		s.db.Exec("UPDATE geekai_jimeng_jobs SET params = task_params WHERE task_params IS NOT NULL AND task_params != ''")
		s.db.Migrator().DropColumn(&model.JimengJob{}, "task_params")
	}

	// VideoJob: task_info -> params, raw_data -> output
	if s.db.Migrator().HasColumn(&model.VideoJob{}, "task_info") {
		s.db.Exec("UPDATE geekai_video_jobs SET params = task_info WHERE task_info IS NOT NULL AND task_info != ''")
		s.db.Migrator().DropColumn(&model.VideoJob{}, "task_info")
	}
	if s.db.Migrator().HasColumn(&model.VideoJob{}, "raw_data") {
		s.db.Exec("UPDATE geekai_video_jobs SET `output` = raw_data WHERE raw_data IS NOT NULL AND raw_data != ''")
		s.db.Migrator().DropColumn(&model.VideoJob{}, "raw_data")
	}

	// SunoJob: task_info -> params, raw_data -> output
	if s.db.Migrator().HasColumn(&model.SunoJob{}, "task_info") {
		s.db.Exec("UPDATE geekai_suno_jobs SET params = task_info WHERE task_info IS NOT NULL AND task_info != ''")
		s.db.Migrator().DropColumn(&model.SunoJob{}, "task_info")
	}
	if s.db.Migrator().HasColumn(&model.SunoJob{}, "raw_data") {
		s.db.Exec("UPDATE geekai_suno_jobs SET `output` = raw_data WHERE raw_data IS NOT NULL AND raw_data != ''")
		s.db.Migrator().DropColumn(&model.SunoJob{}, "raw_data")
	}

	// ========== 删除不再使用的字段 ==========

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
	if s.db.Migrator().HasColumn(&model.VideoJob{}, "water_url") {
		s.db.Migrator().DropColumn(&model.VideoJob{}, "water_url")
	}
	if s.db.Migrator().HasColumn(&model.VideoJob{}, "cover_url") {
		s.db.Migrator().DropColumn(&model.VideoJob{}, "cover_url")
	}
	if s.db.Migrator().HasColumn(&model.VideoJob{}, "prompt_ext") {
		s.db.Migrator().DropColumn(&model.VideoJob{}, "prompt_ext")
	}
	if s.db.Migrator().HasColumn(&model.SunoJob{}, "instrumental") {
		s.db.Migrator().DropColumn(&model.SunoJob{}, "instrumental")
	}
	if s.db.Migrator().HasColumn(&model.SunoJob{}, "tags") {
		s.db.Migrator().DropColumn(&model.SunoJob{}, "tags")
	}
	if s.db.Migrator().HasColumn(&model.SunoJob{}, "extend_secs") {
		s.db.Migrator().DropColumn(&model.SunoJob{}, "extend_secs")
	}
	if s.db.Migrator().HasColumn(&model.SunoJob{}, "model_name") {
		s.db.Migrator().DropColumn(&model.SunoJob{}, "model_name")
	}

	// ========== 数据迁移：根据业务逻辑更新现有数据 ==========

	// video_job: 根据 progress 填充 status
	if s.db.Migrator().HasColumn(&model.VideoJob{}, "status") {
		s.db.Exec(`UPDATE geekai_video_jobs SET status = CASE 
			WHEN progress < 100 THEN 'in_progress'
			WHEN progress = 100 THEN 'success'
			WHEN progress = 101 THEN 'failed'
			WHEN progress = 102 THEN 'downloading'
			ELSE 'pending'
		END WHERE status = '' OR status IS NULL`)
	}

	// suno_job: 从 output 提取 tags/model_name 填入 params
	s.migrateSunoJobData()

	logger.Info("增量迁移完成")
}

// TableMigration 数据表迁移入口：先全量同步 schema，再增量删除字段并迁移数据
func (s *MigrationService) TableMigration() {
	s.fullTableMigration()
	s.fixTableConstraints()
	s.incrementalTableMigration()
	s.migrateChatAppSystemPromptFromJSON()
}

// migrateChatAppSystemPromptFromJSON 将智能体 system_prompt 字段中历史 JSON 数组
// 解析后取出 role 为 system 的 content，覆盖回 system_prompt（纯文本）
func (s *MigrationService) migrateChatAppSystemPromptFromJSON() {
	key := "migrate:chat_app_system_prompt_json"
	if s.redisClient.Get(context.Background(), key).Val() == "1" {
		logger.Info("ChatApp system_prompt JSON 已迁移，跳过")
		return
	}
	logger.Info("开始迁移智能体 system_prompt 历史 JSON 数据...")

	var apps []model.ChatApp
	if err := s.db.Find(&apps).Error; err != nil {
		logger.Errorf("查询 ChatApp 失败: %v", err)
		return
	}

	updated := 0
	for i := range apps {
		raw := strings.TrimSpace(apps[i].SystemPrompt)
		if raw == "" {
			continue
		}
		if len(raw) < 2 || raw[0] != '[' {
			continue
		}
		var messages []types.Message
		if err := json.Unmarshal([]byte(raw), &messages); err != nil {
			continue
		}
		var systemContent string
		for _, m := range messages {
			if strings.ToLower(strings.TrimSpace(m.Role)) == "system" && m.Content != "" {
				systemContent = m.Content
				break
			}
		}
		if err := s.db.Model(&model.ChatApp{}).Where("id = ?", apps[i].Id).Update("system_prompt", systemContent).Error; err != nil {
			logger.Warnf("更新 ChatApp id=%d system_prompt 失败: %v", apps[i].Id, err)
			continue
		}
		updated++
	}

	logger.Infof("智能体 system_prompt JSON 迁移完成，共更新 %d 条", updated)
	s.redisClient.Set(context.Background(), key, "1", 0)
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
		"tencent": map[string]any{
			"secret_id":      config.SMS.Tencent.SecretId,
			"secret_key":     config.SMS.Tencent.SecretKey,
			"sms_sdk_app_id": config.SMS.Tencent.SmsSdkAppId,
			"sign":           config.SMS.Tencent.Sign,
			"code_temp_id":   config.SMS.Tencent.CodeTempId,
			"region":         config.SMS.Tencent.Region,
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

// migrateSunoJobData 合并 suno_job 数据：从 Output 原始数据中解析出 tags 和 model_name 填入 params 字段
func (s *MigrationService) migrateSunoJobData() {
	key := "migrate:suno_job_data"
	if s.redisClient.Get(context.Background(), key).Val() == "1" {
		logger.Info("SunoJob 数据已合并，跳过迁移")
		return
	}

	logger.Info("开始合并 SunoJob 数据...")

	// 查询所有有 output 数据的记录
	var jobs []model.SunoJob
	if err := s.db.Where("output != ? AND output != ''", "").Find(&jobs).Error; err != nil {
		logger.Errorf("查询 SunoJob 数据失败: %v", err)
		return
	}

	updatedCount := 0
	for _, job := range jobs {
		if job.Output == "" {
			continue
		}

		// 解析 Output JSON 数据
		var outputData struct {
			Metadata struct {
				Tags string `json:"tags"`
			} `json:"metadata"`
			ModelName string `json:"model_name"`
		}

		if err := json.Unmarshal([]byte(job.Output), &outputData); err != nil {
			logger.Warnf("解析 Output 数据失败 (ID: %d): %v", job.Id, err)
			continue
		}

		// 检查是否需要更新 params
		needUpdate := false
		params := job.Params

		// 如果 params 中的 tags 为空，但 output 中有 tags，则更新
		if params.Tags == "" && outputData.Metadata.Tags != "" {
			params.Tags = outputData.Metadata.Tags
			// 修复 tags 字段过长导致更新失败
			if len(params.Tags) > 255 {
				params.Tags = params.Tags[:255]
			}
			needUpdate = true
		}

		// 如果 params 中的 model 为空，但 output 中有 model_name，则更新
		if params.Model == "" && outputData.ModelName != "" {
			params.Model = outputData.ModelName
			needUpdate = true
		}

		// 如果需要更新，则保存
		if needUpdate {
			if err := s.db.Model(&model.SunoJob{}).Where("id = ?", job.Id).Update("params", params).Error; err != nil {
				logger.Errorf("更新 SunoJob 数据失败 (ID: %d): %v", job.Id, err)
				continue
			}
			updatedCount++
		}
	}

	logger.Infof("SunoJob 数据合并完成，共更新 %d 条记录", updatedCount)
	s.redisClient.Set(context.Background(), key, "1", 0)
}
