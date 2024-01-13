package auth

import (
	"context"
)

// MacContextKey 是用户的密钥信息
// context.Context中的键值不应该使用普通的字符串， 有可能导致命名冲突
type macContextKey struct{}

// tokenTypeKey 是签名算法类型key
type tokenTypeKey struct{}

// WithCredentials 返回一个包含密钥信息的context
func WithCredentials(ctx context.Context, cred *Credentials) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, macContextKey{}, cred)
}

// WithCredentialsType 返回一个context, 保存了密钥信息和token类型
func WithCredentialsType(ctx context.Context, cred *Credentials, t TokenType) context.Context {
	ctx = WithCredentials(ctx, cred)
	return context.WithValue(ctx, tokenTypeKey{}, t)
}

// CredentialsFromContext 从context获取密钥信息
func CredentialsFromContext(ctx context.Context) (cred *Credentials, t TokenType, ok bool) {
	cred, ok = ctx.Value(macContextKey{}).(*Credentials)
	t, yes := ctx.Value(tokenTypeKey{}).(TokenType)
	if !yes {
		t = TokenQBox
	}
	return
}
