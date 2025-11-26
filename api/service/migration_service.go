package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/store/model"

	"gorm.io/gorm"
)

type MigrationService struct {
	db *gorm.DB
}

func NewMigrationService(db *gorm.DB) *MigrationService {
	return &MigrationService{db: db}
}

func (s *MigrationService) Migrate() error {
	err := s.db.AutoMigrate(
		&model.AdminUser{},
		&model.ApiKey{},
		&model.AppType{},
		&model.ChatItem{},
		&model.ChatMessage{},
		&model.ChatModel{},
		&model.ChatRole{},
		&model.Config{},
		&model.DallJob{},
		&model.File{},
		&model.Function{},
		&model.InviteCode{},
		&model.InviteLog{},
		&model.Menu{},
		&model.MidJourneyJob{},
		&model.Order{},
		&model.PowerLog{},
		&model.Product{},
		&model.Redeem{},
		&model.SdJob{},
		&model.SunoJob{},
		&model.User{},
		&model.UserLoginLog{},
		&model.VideoJob{},
	)
	return err
}
