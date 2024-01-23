package sms

import (
	"chatplus/core/types"
	"chatplus/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type BaoSmsService struct {
	config *types.SmsConfigBao
}

func NewSmsBaoSmsService(appConfig *types.AppConfig) *BaoSmsService {
	config := appConfig.SMS.Bao
	if config.Domain == "" { // use default domain
		config.Domain = "api.smsbao.com"
		logger.Infof("Using default domain for SMS-BAO: %s", config.Domain)
	}
	return &BaoSmsService{
		config: &config,
	}
}

var errMsg = map[string]string{
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

func (s *BaoSmsService) SendVerifyCode(mobile string, code int) error {

	content := fmt.Sprintf("%s%s", s.config.Sign, s.config.CodeTemplate)
	content = strings.ReplaceAll(content, "{code}", strconv.Itoa(code))
	password := utils.Md5(s.config.Password)
	params := url.Values{}
	params.Set("u", s.config.Username)
	params.Set("p", password)
	params.Set("m", mobile)
	params.Set("c", content)

	apiURL := fmt.Sprintf("https://%s/sms?%s", s.config.Domain, params.Encode())
	response, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	result := string(body)
	logger.Debugf("send SmsBao result: %v", errMsg[result])

	if result != "0" {
		return fmt.Errorf("failed to send SMS:%v", errMsg[result])
	}
	return nil
}

var _ Service = &BaoSmsService{}
