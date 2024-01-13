package auth

// 七牛签名算法的类型：
// QBoxToken, QiniuToken, BearToken, QiniuMacToken
type TokenType int

const (
	TokenQiniu TokenType = iota
	TokenQBox
)
