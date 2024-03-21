package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"fmt"
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

const SuperManagerID = 1

type ManagerHandler struct {
	handler.BaseHandler
	redis *redis.Client
}

func NewAdminHandler(app *core.AppServer, db *gorm.DB, client *redis.Client) *ManagerHandler {
	return &ManagerHandler{BaseHandler: handler.BaseHandler{DB: db, App: app}, redis: client}
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
	res := h.DB.Model(&model.AdminUser{}).Where("username = ?", data.Username).First(&manager)
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
	if manager.Id != SuperManagerID && manager.Status == false {
		resp.ERROR(c, "该用户已被禁止登录，请联系超级管理员")
		return
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": manager.Id,
		"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.AdminSession.SecretKey))
	if err != nil {
		resp.ERROR(c, "Failed to generate token, "+err.Error())
		return
	}
	// 保存到 redis
	key := fmt.Sprintf("admin/%d", manager.Id)
	if _, err := h.redis.Set(context.Background(), key, tokenString, 0).Result(); err != nil {
		resp.ERROR(c, "error with save token: "+err.Error())
		return
	}

	// 更新最后登录时间和IP
	manager.LastLoginIp = c.ClientIP()
	manager.LastLoginAt = time.Now().Unix()
	h.DB.Updates(&manager)

	var result = struct {
		IsSuperAdmin bool   `json:"is_super_admin"`
		Token        string `json:"token"`
	}{
		IsSuperAdmin: manager.Id == 1,
		Token:        tokenString,
	}

	resp.SUCCESS(c, result)
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
	id := h.GetLoginUserId(c)
	key := fmt.Sprintf("admin/%d", id)
	if _, err := h.redis.Get(context.Background(), key).Result(); err != nil {
		resp.NotAuth(c)
		return
	}
	var manager model.AdminUser
	res := h.DB.Where("id", id).First(&manager)
	if res.Error != nil {
		resp.NotAuth(c)
		return
	}

	resp.SUCCESS(c, manager)
}

// List 数据列表
func (h *ManagerHandler) List(c *gin.Context) {
	var items []model.AdminUser
	res := h.DB.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	users := make([]vo.AdminUser, 0)
	for _, item := range items {
		var u vo.AdminUser
		err := utils.CopyObject(item, &u)
		if err != nil {
			continue
		}
		u.Id = item.Id
		u.CreatedAt = item.CreatedAt.Unix()
		users = append(users, u)
	}

	resp.SUCCESS(c, users)

}

func (h *ManagerHandler) Save(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Status   bool   `json:"status"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.AdminUser
	res := h.DB.Where("username", data.Username).First(&user)
	if res.Error == nil {
		resp.ERROR(c, "用户名已存在")
		return
	}

	// 生成密码
	salt := utils.RandString(8)
	password := utils.GenPassword(data.Password, salt)
	res = h.DB.Save(&model.AdminUser{
		Username: data.Username,
		Password: password,
		Salt:     salt,
		Status:   data.Status,
	})
	if res.Error != nil {
		resp.ERROR(c, "failed with update database")
		return
	}

	resp.SUCCESS(c)
}

// Remove 删除管理员
func (h *ManagerHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if id == SuperManagerID {
		resp.ERROR(c, "超级管理员不能删除")
		return
	}

	res := h.DB.Where("id", id).Delete(&model.AdminUser{})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}

// Enable 启用/禁用
func (h *ManagerHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.AdminUser{}).Where("id", data.Id).UpdateColumn("status", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}
	resp.SUCCESS(c)
}

// ResetPass 重置密码
func (h *ManagerHandler) ResetPass(c *gin.Context) {
	id := h.GetLoginUserId(c)
	if id != SuperManagerID {
		resp.ERROR(c, "只有超级管理员能够进行该操作")
		return
	}

	var data struct {
		Id       int    `json:"id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.AdminUser
	res := h.DB.Where("id", data.Id).First(&user)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	user.Password = password
	res = h.DB.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}
