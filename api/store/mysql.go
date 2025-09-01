package store

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var log = logger2.GetLogger()

func NewGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "geekai_", // 设置表前缀
			SingularTable: false,     // 使用单数表名形式
		},
	}
}

func NewMysql(config *gorm.Config, appConfig *types.AppConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(appConfig.MysqlDns), config)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(32)
	sqlDB.SetMaxOpenConns(512)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("开始重命名数据表...")

	// 重命名数据表
	tableRenames := map[string]string{
		"chatgpt_chat_roles":      "geekai_chat_roles",
		"chatgpt_admin_users":     "geekai_admin_users",
		"chatgpt_api_keys":        "geekai_api_keys",
		"chatgpt_app_types":       "geekai_app_types",
		"chatgpt_chat_history":    "geekai_chat_history",
		"chatgpt_chat_items":      "geekai_chat_items",
		"chatgpt_chat_models":     "geekai_chat_models",
		"chatgpt_users":           "geekai_users",
		"chatgpt_orders":          "geekai_orders",
		"chatgpt_products":        "geekai_products",
		"chatgpt_configs":         "geekai_configs",
		"chatgpt_sd_jobs":         "geekai_sd_jobs",
		"chatgpt_mj_jobs":         "geekai_mj_jobs",
		"chatgpt_suno_jobs":       "geekai_suno_jobs",
		"chatgpt_dall_jobs":       "geekai_dall_jobs",
		"chatgpt_video_jobs":      "geekai_video_jobs",
		"chatgpt_jimeng_jobs":     "geekai_jimeng_jobs",
		"chatgpt_files":           "geekai_files",
		"chatgpt_menus":           "geekai_menus",
		"chatgpt_functions":       "geekai_functions",
		"chatgpt_invite_codes":    "geekai_invite_codes",
		"chatgpt_invite_logs":     "geekai_invite_logs",
		"chatgpt_redeems":         "geekai_redeems",
		"chatgpt_power_logs":      "geekai_power_logs",
		"chatgpt_user_login_logs": "geekai_user_login_logs",
	}

	// 执行重命名操作
	for oldTableName, newTableName := range tableRenames {
		// 检查新表是否已存在
		if !db.Migrator().HasTable(newTableName) {
			err := db.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldTableName, newTableName)).Error
			if err != nil {
				log.Errorf("重命名数据表 %s 到 %s 失败: %v", oldTableName, newTableName, err)
			} else {
				log.Infof("成功重命名数据表: %s -> %s", oldTableName, newTableName)
			}
		}
	}
	return db, nil
}
