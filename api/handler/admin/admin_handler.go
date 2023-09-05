package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"strings"
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
			"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)),
		})
		tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
		if err != nil {
			resp.ERROR(c, "Failed to generate token, "+err.Error())
			return
		}
		// 保存到 redis
		if _, err := h.redis.Set(context.Background(), "users/"+manager.Username, tokenString, 0).Result(); err != nil {
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
	token := c.GetHeader(types.AdminAuthHeader)
	if _, err := h.redis.Del(c, token).Result(); err != nil {
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

// Migrate 数据修正
func (h *ManagerHandler) Migrate(c *gin.Context) {
	opt := c.Query("opt")
	switch opt {
	case "user":
		// 将用户订阅角色的数据结构从 map 改成数组
		var users []model.User
		h.db.Find(&users)
		for _, u := range users {
			var m map[string]int
			var roleKeys = make([]string, 0)
			err := utils.JsonDecode(u.ChatRoles, &m)
			if err != nil {
				continue
			}

			for k := range m {
				roleKeys = append(roleKeys, k)
			}
			u.ChatRoles = utils.JsonEncode(roleKeys)
			h.db.Updates(&u)

		}
		break
	case "role":
		// 修改角色图片，改成绝对路径
		var roles []model.ChatRole
		h.db.Find(&roles)
		for _, r := range roles {
			if !strings.HasPrefix(r.Icon, "/") {
				r.Icon = "/" + r.Icon
				h.db.Updates(&r)
			}
		}
		break
	case "history":
		// 修改角色图片，改成绝对路径
		var message []model.HistoryMessage
		h.db.Find(&message)
		for _, r := range message {
			if !strings.HasPrefix(r.Icon, "/") {
				r.Icon = "/" + r.Icon
				h.db.Updates(&r)
			}

		}
		break

	case "avatar":
		// 更新用户的头像地址
		var users []model.User
		h.db.Find(&users)
		for _, u := range users {
			if !strings.HasPrefix(u.Avatar, "/") {
				u.Avatar = "/" + u.Avatar
				h.db.Updates(&u)
			}
		}
		break
	}

	resp.SUCCESS(c, "SUCCESS")
}
