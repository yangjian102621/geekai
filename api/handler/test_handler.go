package handler

import (
	"chatplus/service"
	"chatplus/service/payment"
	"gorm.io/gorm"
)

type TestHandler struct {
	db        *gorm.DB
	snowflake *service.Snowflake
	js        *payment.PayJS
}

func NewTestHandler(db *gorm.DB, snowflake *service.Snowflake, js *payment.PayJS) *TestHandler {
	return &TestHandler{db: db, snowflake: snowflake, js: js}
}
