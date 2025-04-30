package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"geekai/core/types"
	logger2 "geekai/logger"

	"github.com/xxl-job/xxl-job-executor-go"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

type XXLJobExecutor struct {
	executor xxl.Executor
	db       *gorm.DB
}

func NewXXLJobExecutor(config *types.AppConfig, db *gorm.DB) *XXLJobExecutor {
	if !config.XXLConfig.Enabled {
		logger.Info("XXL-JOB service is disabled")
		return nil
	}

	exec := xxl.NewExecutor(
		xxl.ServerAddr(config.XXLConfig.ServerAddr),
		xxl.AccessToken(config.XXLConfig.AccessToken),   //请求令牌(默认为空)
		xxl.ExecutorIp(config.XXLConfig.ExecutorIp),     //可自动获取
		xxl.ExecutorPort(config.XXLConfig.ExecutorPort), //默认9999（非必填）
		xxl.RegistryKey(config.XXLConfig.RegistryKey),   //执行器名称
		xxl.SetLogger(&customLogger{}),                  //自定义日志
	)
	exec.Init()
	return &XXLJobExecutor{executor: exec, db: db}
}

func (e *XXLJobExecutor) Run() error {
	e.executor.RegTask("ClearOrders", e.ClearOrders)
	return e.executor.Run()
}

// ClearOrders 清理未支付的订单，如果没有抛出异常则表示执行成功
func (e *XXLJobExecutor) ClearOrders(cxt context.Context, param *xxl.RunReq) (msg string) {
	logger.Info("执行清理未支付订单...")
	
	return "success"
}

type customLogger struct{}

func (l *customLogger) Info(format string, a ...interface{}) {
	logger.Debugf(format, a...)
}

func (l *customLogger) Error(format string, a ...interface{}) {
	logger.Errorf(format, a...)
}
