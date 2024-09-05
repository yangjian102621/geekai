package handler

import (
	"geekai/service"
	"geekai/service/payment"
	"gorm.io/gorm"
)

type TestHandler struct {
	db        *gorm.DB
	snowflake *service.Snowflake
	js        *payment.JPayService
}

func NewTestHandler(db *gorm.DB, snowflake *service.Snowflake, js *payment.JPayService) *TestHandler {
	return &TestHandler{db: db, snowflake: snowflake, js: js}
}
