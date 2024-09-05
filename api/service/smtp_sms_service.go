package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"geekai/core/types"
	"mime"
	"net/smtp"
	"net/textproto"
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
	subject := fmt.Sprintf("%s 注册验证码", s.config.AppName)
	body := fmt.Sprintf("您正在注册 %s 账户，注册验证码为 %d，请不要告诉他人。如非本人操作，请忽略此邮件。", s.config.AppName, code)

	auth := smtp.PlainAuth("", s.config.From, s.config.Password, s.config.Host)
	if s.config.UseTls {
		return s.sendTLS(auth, to, subject, body)
	} else {
		return s.send(auth, to, subject, body)
	}
}

func (s *SmtpService) send(auth smtp.Auth, to string, subject string, body string) error {
	// 对主题进行MIME编码
	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)
	// 组装邮件
	message := bytes.NewBuffer(nil)
	message.WriteString(fmt.Sprintf("From: \"%s\" <%s>\r\n", s.config.AppName, s.config.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	message.WriteString("\r\n" + body)

	// 发送邮件
	err := smtp.SendMail(s.config.Host+":"+fmt.Sprint(s.config.Port), auth, s.config.From, []string{to}, message.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return err

}

func (s *SmtpService) sendTLS(auth smtp.Auth, to string, subject string, body string) error {
	// TLS配置
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	// 建立TLS连接
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("error connecting to SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("error creating SMTP client: %v", err)
	}
	defer client.Quit()

	// 身份验证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("error authenticating: %v", err)
	}

	// 设置寄件人
	if err = client.Mail(s.config.From); err != nil {
		return fmt.Errorf("error setting sender: %v", err)
	}

	// 设置收件人
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("error setting recipient: %v", err)
	}

	// 发送邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("error getting data writer: %v", err)
	}
	defer wc.Close()

	header := make(textproto.MIMEHeader)
	header.Set("From", s.config.From)
	header.Set("To", to)
	header.Set("Subject", subject)

	// 将邮件头写入
	for key, values := range header {
		for _, value := range values {
			_, err = fmt.Fprintf(wc, "%s: %s\r\n", key, value)
			if err != nil {
				return fmt.Errorf("error sending email header: %v", err)
			}
		}
	}
	_, _ = fmt.Fprintln(wc)
	// 将邮件内容写入
	_, err = fmt.Fprintf(wc, body)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	// 发送完毕
	err = wc.Close()
	if err != nil {
		return fmt.Errorf("error closing data writer: %v", err)
	}

	return nil
}
