package sms

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var logger = logger2.GetLogger()

type SmsBaoSmsService struct {
	config *types.SmsBaoSmsConfig
}

func NewSmsBaoSmsService(appConfig *types.AppConfig) *SmsBaoSmsService {
	return &SmsBaoSmsService{
		config: &appConfig.SMS.SMSBAO,
	}
}

var statusStr = map[string]string{
	"0":  "短信发送成功",
	"-1": "参数不全",
	"-2": "服务器空间不支持，请确认支持curl或者fsocket，联系您的空间商解决或者更换空间",
	"30": "密码错误",
	"40": "账号不存在",
	"41": "余额不足",
	"42": "账户已过期",
	"43": "IP地址限制",
	"50": "内容含有敏感词",
}

func (s *SmsBaoSmsService) SendVerifyCode(mobile string, code int) error {

	content := fmt.Sprintf("%s%s", s.config.Sign, s.config.CodeTemplate)
	template := replaceTemplate(content, s.config.Num, code)

	md5Hash := s.config.ApiKey
	params := url.Values{}
	params.Set("u", s.config.Account)
	params.Set("p", md5Hash)
	params.Set("m", mobile)
	params.Set("c", template)

	// 判断 s.config.Domain 是否为空
	if s.config.Domain == "" {
		// 设置默认值
		s.config.Domain = "api.smsbao.com"
		// 记录日志，提醒用户默认值被使用
		logger.Infof("SmsBao.config.Domain is empty. Using default value: %s", s.config.Domain)
	}
	real_url := fmt.Sprintf("https://%s/sms?", s.config.Domain)
	sendURL := real_url + params.Encode()
	logger.Infof("send SmsBao content: %v", template)

	response, err := http.Get(sendURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	result := string(body)
	logger.Infof("send SmsBao result: %v", statusStr[result])

	if result != "0" {
		return fmt.Errorf("failed to send SMS:%v", statusStr[result])
	}
	return nil
}

func replaceTemplate(template, num string, code int) string {
	result := strings.ReplaceAll(template, "{code}", strconv.Itoa(code))
	result = strings.ReplaceAll(result, "{num}", num)
	return result
}

var _ SmsService = &SmsBaoSmsService{}
