package sms

const Ali = "ALI"
const Bao = "BAO"

type Service interface {
	SendVerifyCode(mobile string, code int) error
}
