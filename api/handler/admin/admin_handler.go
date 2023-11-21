package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/utils/resp"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

type ManagerHandler struct {
	handler.BaseHandler
	db    *gorm.DB
	redis *redis.Client
}

func NewAdminHandler(app *core.AppServer, db *gorm.DB, client *redis.Client) *ManagerHandler {
	h := ManagerHandler{db: db, redis: client}
	h.App = app
	return &h
}

// Login 登录
func (h *ManagerHandler) Login(c *gin.Context) {
	var data types.Manager
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	manager := h.App.Config.Manager
	if data.Username == manager.Username && data.Password == manager.Password {
		// 创建 token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": manager.Username,
			"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)).Unix(),
		})
		tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
		if err != nil {
			resp.ERROR(c, "Failed to generate token, "+err.Error())
			return
		}
		// 保存到 redis
		key := "users/" + manager.Username
		if _, err := h.redis.Set(context.Background(), key, tokenString, 0).Result(); err != nil {
			resp.ERROR(c, "error with save token: "+err.Error())
			return
		}
		resp.SUCCESS(c, tokenString)
	} else {
		resp.ERROR(c, "用户名或者密码错误")
	}
}

// Logout 注销
func (h *ManagerHandler) Logout(c *gin.Context) {
	key := h.GetUserKey(c)
	if _, err := h.redis.Del(c, key).Result(); err != nil {
		logger.Error("error with delete session: ", err)
	} else {
		resp.SUCCESS(c)
	}
}

// Session 会话检测
func (h *ManagerHandler) Session(c *gin.Context) {
	token := c.GetHeader(types.AdminAuthHeader)
	if token == "" {
		resp.NotAuth(c)
	} else {
		resp.SUCCESS(c)
	}
}
