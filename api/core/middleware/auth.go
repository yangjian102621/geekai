package middleware

import (
	"context"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

// 前端用户授权验证
func UserAuthMiddleware(secretKey string, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(types.UserAuthHeader)
		if tokenString == "" {
			resp.NotAuth(c, "无效的授权令牌")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支持的令牌签名方法: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			resp.NotAuth(c, fmt.Sprintf("解析授权令牌失败: %v", err))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			resp.NotAuth(c, "令牌无效")
			c.Abort()
			return
		}

		expr := utils.IntValue(utils.InterfaceToString(claims["expired"]), 0)
		if expr > 0 && int64(expr) < time.Now().Unix() {
			resp.NotAuth(c, "令牌过期")
			c.Abort()
			return
		}

		key := fmt.Sprintf("users/%v", claims["user_id"])
		if _, err := redis.Get(context.Background(), key).Result(); err != nil {
			resp.NotAuth(c, "当前用户已退出登录")
			c.Abort()
			return
		}
		c.Set(types.LoginUserID, claims["user_id"])
	}
}

// 管理后台用户授权验证
func AdminAuthMiddleware(secretKey string, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(types.AdminAuthHeader)
		if tokenString == "" {
			resp.NotAuth(c, "无效的授权令牌")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支持的令牌签名方法: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			resp.NotAuth(c, fmt.Sprintf("解析授权令牌失败: %v", err))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			resp.NotAuth(c, "令牌无效")
			c.Abort()
			return
		}

		expr := utils.IntValue(utils.InterfaceToString(claims["expired"]), 0)
		if expr > 0 && int64(expr) < time.Now().Unix() {
			resp.NotAuth(c, "令牌过期")
			c.Abort()
			return
		}

		key := fmt.Sprintf("admin/%v", claims["user_id"])
		if _, err := redis.Get(context.Background(), key).Result(); err != nil {
			resp.NotAuth(c, "当前用户已退出登录")
			c.Abort()
			return
		}
		c.Set(types.AdminUserID, claims["user_id"])
	}
}
