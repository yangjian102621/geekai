package core

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/store/model"
	"geekai/utils"
	"os"

	"github.com/BurntSushi/toml"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

func NewDefaultConfig() *types.AppConfig {
	return &types.AppConfig{
		Listen:    "0.0.0.0:5678",
		ProxyURL:  "",
		StaticDir: "./static",
		StaticUrl: "http://localhost/5678/static",
		Redis:     types.RedisConfig{Host: "localhost", Port: 6379, Password: ""},
		Session: types.Session{
			SecretKey: utils.RandString(64),
			MaxAge:    86400,
		},
		OSS: types.OSSConfig{
			Active: "local",
			Local: types.LocalStorageConfig{
				BaseURL:  "http://localhost/5678/static/upload",
				BasePath: "./static/upload",
			},
		},
	}
}

func LoadConfig(configFile string) (*types.AppConfig, error) {
	var config *types.AppConfig
	_, err := os.Stat(configFile)
	if err != nil {
		logger.Info("creating new config file: ", configFile)
		config = NewDefaultConfig()
		config.Path = configFile
		// save config
		err := SaveConfig(config)
		if err != nil {
			return nil, err
		}

		return config, nil
	}
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}

	return config, err
}

func SaveConfig(config *types.AppConfig) error {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(&config); err != nil {
		return err
	}

	return os.WriteFile(config.Path, buf.Bytes(), 0644)
}

func LoadSystemConfig(db *gorm.DB) *types.SystemConfig {
	// 加载系统配置
	var sysConfig model.Config
	var baseConfig types.BaseConfig
	db.Where("name", "system").First(&sysConfig)
	err := utils.JsonDecode(sysConfig.Value, &baseConfig)
	if err != nil {
		logger.Error("load system config error: ", err)
	}

	// 加载许可证配置
	var license types.License
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeyLicense).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &license)
	if err != nil {
		logger.Error("load license config error: ", err)
	}

	// 加载验证码配置
	var captchaConfig types.CaptchaConfig
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeyCaptcha).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &captchaConfig)
	if err != nil {
		logger.Error("load geek service config error: ", err)
	}

	// 加载微信登录配置
	var wxLoginConfig types.WxLoginConfig
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeyWxLogin).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &wxLoginConfig)
	if err != nil {
		logger.Error("load wx login config error: ", err)
	}

	// 加载短信配置
	var smsConfig types.SMSConfig
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeySms).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &smsConfig)
	if err != nil {
		logger.Error("load sms config error: ", err)
	}

	// 加载 OSS 配置
	var ossConfig types.OSSConfig
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeyOss).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &ossConfig)
	if err != nil {
		logger.Error("load oss config error: ", err)
	}

	// 加载 SMTP 配置
	var smtpConfig types.SmtpConfig
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeySmtp).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &smtpConfig)
	if err != nil {
		logger.Error("load smtp config error: ", err)
	}

	// 加载支付配置
	var paymentConfig types.PaymentConfig
	sysConfig.Id = 0
	db.Where("name", types.ConfigKeyPayment).First(&sysConfig)
	err = utils.JsonDecode(sysConfig.Value, &paymentConfig)
	if err != nil {
		logger.Error("load payment config error: ", err)
	}

	return &types.SystemConfig{
		Base:    baseConfig,
		License: license,
		SMS:     smsConfig,
		OSS:     ossConfig,
		SMTP:    smtpConfig,
		Payment: paymentConfig,
		Captcha: captchaConfig,
		WxLogin: wxLoginConfig,
	}
}
