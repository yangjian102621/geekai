package service

import (
	"bytes"
	"chatplus/core/types"
	"fmt"
	"mime"
	"net/smtp"
)

type SmtpService struct {
	config *types.SmtpConfig
}

func NewSmtpService(appConfig *types.AppConfig) *SmtpService {
	return &SmtpService{
		config: &appConfig.SmtpConfig,
	}
}

func (s *SmtpService) SendVerifyCode(to string, code int) error {
	subject := "ChatPlus注册验证码"
	body := fmt.Sprintf("您正在注册 ChatPlus AI 助手账户，注册验证码为 %d，请不要告诉他人。如非本人操作，请忽略此邮件。", code)

	// 设置SMTP客户端配置
	auth := smtp.PlainAuth("", s.config.From, s.config.Password, s.config.Host)

	// 对主题进行MIME编码
	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)
	// 组装邮件
	message := bytes.NewBuffer(nil)
	message.WriteString(fmt.Sprintf("From: \"%s\" <%s>\r\n", s.config.AppName, s.config.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	message.WriteString("\r\n" + body)

	// 发送邮件
	// 发送邮件
	err := smtp.SendMail(s.config.Host+":"+fmt.Sprint(s.config.Port), auth, s.config.From, []string{to}, message.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
