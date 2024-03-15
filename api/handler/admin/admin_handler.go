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
	"github.com/mojocn/base64Captcha"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

// Manager 管理员
type Manager struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`    // 验证码
	CaptchaId string `json:"captcha_id"` // 验证码id
}

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

	var data Manager
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// add captcha
	if !base64Captcha.DefaultMemStore.Verify(data.CaptchaId, data.Captcha, true) {
		resp.ERROR(c, "验证码错误!")
		return
	}

	var manager model.AdminUser
	res := h.db.Model(&model.AdminUser{}).Where("username = ?", data.Username).First(&manager)
	if res.Error != nil {
		resp.ERROR(c, "请检查用户名或者密码是否填写正确")
		return
	}
	password := utils.GenPassword(data.Password, manager.Salt)
	if password != manager.Password {
		resp.ERROR(c, "用户名或密码错误")
		return
	}

	// 超级管理员默认是ID:1
	if manager.Id != 1 && manager.Status == false {
		resp.ERROR(c, "该用户已被禁止登录，请联系超级管理员")
		return
	}

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

	// 更新最后登录时间和IP
	manager.LastLoginIp = c.ClientIP()
	manager.LastLoginAt = time.Now().Unix()
	h.db.Model(&manager).Updates(manager)

	permissions := h.GetAdminSlugs(manager.Id)
	var result = struct {
		IsSuperAdmin bool     `json:"is_super_admin"`
		Token        string   `json:"token"`
		Permissions  []string `json:"permissions"`
	}{
		IsSuperAdmin: manager.Id == 1,
		Token:        tokenString,
		Permissions:  permissions,
	}

	resp.SUCCESS(c, result)
}

func (h *ManagerHandler) GetAdminSlugs(userId uint) []string {
	var permissions []string
	err := h.db.Raw("SELECT distinct p.slug "+
		"FROM chatgpt_admin_user_roles as ur "+
		"LEFT JOIN chatgpt_admin_role_permissions as rp ON ur.role_id = rp.role_id "+
		"LEFT JOIN chatgpt_admin_permissions as p ON rp.permission_id = p.id "+
		"WHERE ur.admin_id = ?", userId).Scan(&permissions)
	if err.Error == nil {
		return []string{}
	}
	return permissions
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
