package service

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/store/model"
	"chatplus/utils"
	"context"
	"fmt"
	"github.com/xxl-job/xxl-job-executor-go"
	"gorm.io/gorm"
	"time"
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
	e.executor.RegTask("ResetVipPower", e.ResetVipPower)
	return e.executor.Run()
}

// ClearOrders 清理未支付的订单，如果没有抛出异常则表示执行成功
func (e *XXLJobExecutor) ClearOrders(cxt context.Context, param *xxl.RunReq) (msg string) {
	logger.Debug("执行清理未支付订单...")
	var sysConfig model.Config
	res := e.db.Where("marker", "system").First(&sysConfig)
	if res.Error != nil {
		return "error with get system config: " + res.Error.Error()
	}

	var config types.SystemConfig
	err := utils.JsonDecode(sysConfig.Config, &config)
	if err != nil {
		return "error with decode system config: " + err.Error()
	}

	if config.OrderPayTimeout == 0 { // 默认未支付订单的生命周期为 30 分钟
		config.OrderPayTimeout = 1800
	}
	timeout := time.Now().Unix() - int64(config.OrderPayTimeout)
	start := utils.Stamp2str(timeout)
	// 这里不是用软删除，而是永久删除订单
	res = e.db.Unscoped().Where("status != ? AND created_at < ?", types.OrderPaidSuccess, start).Delete(&model.Order{})
	return fmt.Sprintf("Clear order successfully, affect rows: %d", res.RowsAffected)
}

// ResetVipPower 清理过期的 VIP 会员
func (e *XXLJobExecutor) ResetVipPower(cxt context.Context, param *xxl.RunReq) (msg string) {
	logger.Info("开始进行月底账号盘点...")
	var users []model.User
	res := e.db.Where("vip = ?", 1).Find(&users)
	if res.Error != nil {
		return "No vip users found"
	}

	var sysConfig model.Config
	res = e.db.Where("marker", "system").First(&sysConfig)
	if res.Error != nil {
		return "error with get system config: " + res.Error.Error()
	}

	var config types.SystemConfig
	err := utils.JsonDecode(sysConfig.Config, &config)
	if err != nil {
		return "error with decode system config: " + err.Error()
	}

	for _, u := range users {
		if u.Power <= 0 {
			u.Power = 0
		}
		u.Power += config.VipMonthPower
		// update user
		tx := e.db.Updates(&u)
		// 记录算力充值日志
		if tx.Error == nil {
			e.db.Create(&model.PowerLog{
				UserId:    u.Id,
				Username:  u.Username,
				Type:      types.PowerRecharge,
				Amount:    config.VipMonthPower,
				Mark:      types.PowerAdd,
				Balance:   u.Power + config.VipMonthPower,
				Model:     "",
				Remark:    fmt.Sprintf("月底盘点，会员每月赠送算力：%d", config.VipMonthPower),
				CreatedAt: time.Now(),
			})
		}
	}
	logger.Info("月底盘点完成！")
	return "success"
}

type customLogger struct{}

func (l *customLogger) Info(format string, a ...interface{}) {
	logger.Debugf(format, a...)
}

func (l *customLogger) Error(format string, a ...interface{}) {
	logger.Errorf(format, a...)
}
