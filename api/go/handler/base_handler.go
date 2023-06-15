package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	logger2 "chatplus/logger"
)

var logger = logger2.GetLogger()

type BaseHandler struct {
	app    *core.AppServer
	config *types.AppConfig
}
