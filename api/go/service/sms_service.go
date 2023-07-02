package service

type SmsService interface {
	SendVerifyCode(mobile string, code int) error
}
