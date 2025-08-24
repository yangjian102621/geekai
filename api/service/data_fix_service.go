package service

import (
	"geekai/store/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DataFixService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewDataFixService(db *gorm.DB, redis *redis.Client) *DataFixService {
	return &DataFixService{db: db, redis: redis}
}

func (s *DataFixService) FixData() {
	s.FixColumn()
}

// 字段修正
func (s *DataFixService) FixColumn() {
	// 订单字段整理
	if s.db.Migrator().HasColumn(&model.Order{}, "pay_type") {
		s.db.Migrator().RenameColumn(&model.Order{}, "pay_type", "channel")
	}
	if !s.db.Migrator().HasColumn(&model.Order{}, "check") {
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
}
