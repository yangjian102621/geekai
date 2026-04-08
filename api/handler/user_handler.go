package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"gorm.io/gorm"
)

type UserHandler struct {
	BaseHandler
	searcher       *xdb.Searcher
	redis          *redis.Client
	levelDB        *store.LevelDB
	licenseService *service.LicenseService
	captchaService *service.CaptchaService
	userService    *service.UserService
	wxLoginService *service.WxLoginService
	ipSearcher     *xdb.Searcher
}

func NewUserHandler(
	app *core.AppServer,
	db *gorm.DB,
	searcher *xdb.Searcher,
	client *redis.Client,
	levelDB *store.LevelDB,
	captcha *service.CaptchaService,
	userService *service.UserService,
	wxLoginService *service.WxLoginService,
	ipSearcher *xdb.Searcher,
	licenseService *service.LicenseService) *UserHandler {
	return &UserHandler{
		BaseHandler:    BaseHandler{DB: db, App: app},
		searcher:       searcher,
		redis:          client,
		levelDB:        levelDB,
		captchaService: captcha,
		licenseService: licenseService,
		userService:    userService,
		wxLoginService: wxLoginService,
		ipSearcher:     ipSearcher,
	}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/user/")

	// 公开接口，不需要授权
	group.POST("register", h.Register)
	group.POST("login", h.Login)
	group.POST("resetPass", h.ResetPass)
	group.GET("login/qrcode", h.GetWxLoginQRCode)
	group.POST("login/callback", h.WxLoginCallback)
	group.GET("login/status", h.GetWxLoginState)
	group.GET("logout", h.Logout)

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.GET("session", h.Session)
		group.GET("profile", h.Profile)
		group.POST("profile/update", h.ProfileUpdate)
		group.POST("password", h.UpdatePass)
		group.POST("bind/mobile", h.BindMobile)
		group.POST("bind/email", h.BindEmail)
		group.GET("signin", h.SignIn)
	}
}

// Register user register
func (h *UserHandler) Register(c *gin.Context) {
	// parameters process
	var data struct {
		RegWay     string `json:"reg_way"`
		Username   string `json:"username"`
		Mobile     string `json:"mobile"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Code       string `json:"code"`
		InviteCode string `json:"invite_code"`
		Key        string `json:"key,omitempty"`
		Dots       string `json:"dots,omitempty"`
		X          int    `json:"x,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 人机验证
	if h.captchaService.GetConfig().Enabled {
		var check bool
		if data.X != 0 {
			check = h.captchaService.SlideCheck(data)
		} else {
			check = h.captchaService.Check(data)
		}
		if !check {
			resp.ERROR(c, "请先完人机验证")
			return
		}
	}

	data.Password = strings.TrimSpace(data.Password)
	if len(data.Password) < 8 {
		resp.ERROR(c, "密码长度不能少于8个字符")
		return
	}

	// 检测最大注册人数
	var totalUser int64
	h.DB.Model(&model.User{}).Count(&totalUser)
	if h.licenseService.GetLicense().Configs.UserNum > 0 && int(totalUser) >= h.licenseService.GetLicense().Configs.UserNum {
		resp.ERROR(c, "当前注册用户数已达上限，请请升级 License")
		return
	}

	// 检查验证码
	var key string
	if data.RegWay == "email" {
		key = CodeStorePrefix + data.Email
		code, err := h.redis.Get(c, key).Result()
		if err != nil || code != data.Code {
			resp.ERROR(c, "验证码错误")
			return
		}
	} else if data.RegWay == "mobile" {
		key = CodeStorePrefix + data.Mobile
		code, err := h.redis.Get(c, key).Result()
		if err != nil || code != data.Code {
			resp.ERROR(c, "验证码错误")
			return
		}
	}

	// check if the username is existing
	user := model.User{Username: data.Username, Password: data.Password}
	var item model.User
	session := h.DB.Session(&gorm.Session{})
	if data.Mobile != "" {
		session = session.Where("mobile = ?", data.Mobile)
		user.Username = data.Mobile
		user.Mobile = data.Mobile
	} else if data.Email != "" {
		session = session.Where("email = ?", data.Email)
		user.Username = data.Email
		user.Email = data.Email
	} else if data.Username != "" {
		session = session.Where("username = ?", data.Username)
	}
	session.First(&item)
	if item.Id > 0 {
		resp.ERROR(c, "该用户名已经被注册")
		return
	}

	user, err := h.createNewUser(user, data.InviteCode)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"token": token, "user_id": user.Id, "username": user.Username})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Key      string `json:"key,omitempty"`
		Dots     string `json:"dots,omitempty"`
		X        int    `json:"x,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if h.captchaService.GetConfig().Enabled {
		var check bool
		if data.X != 0 {
			check = h.captchaService.SlideCheck(data)
		} else {
			check = h.captchaService.Check(data)
		}
		if !check {
			resp.ERROR(c, "请先完人机验证")
			return
		}
	}

	var user model.User
	res := h.DB.Where("username = ?", data.Username).First(&user)
	if res.Error != nil {
		resp.ERROR(c, "用户名不存在")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	if password != user.Password {
		resp.ERROR(c, "用户名或密码错误")
		return
	}

	if !user.Status {
		resp.ERROR(c, "该用户已被禁止登录，请联系管理员")
		return
	}

	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"token": token, "user_id": user.Id, "username": user.Username})
}

// Logout 注 销
func (h *UserHandler) Logout(c *gin.Context) {
	key := h.GetUserKey(c)
	if _, err := h.redis.Del(c, key).Result(); err != nil {
		logger.Error("error with delete session: ", err)
	}
	resp.SUCCESS(c)
}

// GetWxLoginQRCode 获取微信登录二维码URL
func (h *UserHandler) GetWxLoginQRCode(c *gin.Context) {
	if !h.wxLoginService.GetConfig().Enabled {
		resp.ERROR(c, "微信登录功能未启用")
		return
	}

	if h.wxLoginService.GetConfig().ApiKey == "" {
		resp.ERROR(c, "微信登录服务令牌未配置")
		return
	}

	state := utils.RandString(32)
	qrCodeURL, err := h.wxLoginService.GetLoginQrCodeUrl(state)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"url":   qrCodeURL,
		"state": state,
	})
}

// 查询微信登录状态
func (h *UserHandler) GetWxLoginState(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		resp.ERROR(c, "参数错误")
		return
	}

	status, err := h.wxLoginService.GetLoginStatus(state)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if status.Status != service.LoginStatusSuccess {
		resp.SUCCESS(c, status)
		return
	}

	// 登录成功
	var user model.User
	h.DB.Where("openid = ?", status.OpenID).First(&user)
	if user.Id == 0 {
		// 创建新用户
		user, err = h.createNewUser(model.User{OpenId: status.OpenID}, "")
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	status.Status = service.LoginStatusExpired
	h.wxLoginService.SetLoginStatus(state, *status)

	status.Status = service.LoginStatusSuccess
	status.Token = token
	resp.SUCCESS(c, status)
}

// createNewUser 创建新用户
func (h *UserHandler) createNewUser(user model.User, inviteCode string) (model.User, error) {
	if user.OpenId != "" {
		user.Platform = "wechat"
		user.Nickname = fmt.Sprintf("微信用户@%d", utils.RandomNumber(6))
		user.Username = fmt.Sprintf("wx@%d", utils.RandomNumber(8))
		user.Password = "geekai123"
	} else {
		user.Nickname = fmt.Sprintf("用户@%d", utils.RandomNumber(6))
		if user.Username == "" || user.Password == "" {
			return user, fmt.Errorf("用户名或密码不能为空")
		}
	}

	salt := utils.RandString(8)
	user.Salt = salt
	user.Password = utils.GenPassword(user.Password, salt)
	user.Avatar = "/images/avatar/user.png"
	user.Status = true
	user.ChatRoles = utils.JsonEncode([]string{"gpt"})
	user.ChatConfig = "{}"
	user.ChatModels = "{}"
	user.Power = h.App.SysConfig.Base.InitPower

	// 创建用户
	tx := h.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		return user, err
	}

	// 记录邀请关系
	if inviteCode != "" {
		inviteCode := model.InviteCode{}
		err := h.DB.Where("code = ?", inviteCode).First(&inviteCode).Error
		if err != nil {
			return user, fmt.Errorf("无效的邀请码")
		}

		// 增加邀请数量
		h.DB.Model(&model.InviteCode{}).Where("code = ?", inviteCode).UpdateColumn("reg_num", gorm.Expr("reg_num + ?", 1))
		if h.App.SysConfig.Base.InvitePower > 0 {
			err := h.userService.IncreasePower(inviteCode.UserId, h.App.SysConfig.Base.InvitePower, model.PowerLog{
				Type:   types.PowerInvite,
				Model:  "Invite",
				Remark: fmt.Sprintf("邀请用户注册奖励，金额：%d，邀请码：%s，新用户：%s", h.App.SysConfig.Base.InvitePower, inviteCode.Code, user.Username),
			})
			if err != nil {
				tx.Rollback()
				return user, err
			}

			// 添加邀请记录
			err = tx.Create(&model.InviteLog{
				InviterId:  inviteCode.UserId,
				UserId:     user.Id,
				Username:   user.Username,
				InviteCode: inviteCode.Code,
				Remark:     fmt.Sprintf("奖励 %d 算力", h.App.SysConfig.Base.InvitePower),
			}).Error
			if err != nil {
				tx.Rollback()
				return user, err
			}
		}
	}

	tx.Commit()

	return user, nil
}

// doLogin 执行登录操作
func (h *UserHandler) doLogin(user *model.User, ip string) (string, error) {
	// 更新最后登录时间和IP
	user.LastLoginIp = ip
	user.LastLoginAt = time.Now().Unix()
	err := h.DB.Model(user).Updates(user).Error
	if err != nil {
		return "", fmt.Errorf("failed to update user: %v", err)
	}

	// 记录登录日志
	h.DB.Create(&model.UserLoginLog{
		UserId:       user.Id,
		Username:     user.Username,
		LoginIp:      ip,
		LoginAddress: utils.Ip2Region(h.ipSearcher, ip),
	})

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// 保存到 redis
	sessionKey := fmt.Sprintf("users/%d", user.Id)
	if _, err = h.redis.Set(context.Background(), sessionKey, tokenString, 0).Result(); err != nil {
		return "", fmt.Errorf("error with save token: %v", err)
	}

	return tokenString, nil
}

// WxLoginCallback 微信登录回调处理
func (h *UserHandler) WxLoginCallback(c *gin.Context) {
	var data struct {
		OpenID string `json:"openid"`
		State  string `json:"state"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.OpenID == "" || data.State == "" {
		resp.ERROR(c, "参数错误")
		return
	}

	// 设置登录状态
	status := service.LoginStatus{
		Status: service.LoginStatusSuccess,
		OpenID: data.OpenID,
	}
	h.wxLoginService.SetLoginStatus(data.State, status)

	resp.SUCCESS(c, status)
}

// Session 获取/验证会话
func (h *UserHandler) Session(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c, err.Error())
		return
	}

	var userVo vo.User
	err = utils.CopyObject(user, &userVo)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 用户 VIP 到期
	if user.ExpiredTime > 0 && user.ExpiredTime < time.Now().Unix() {
		h.DB.Model(&user).UpdateColumn("vip", false)
	}
	userVo.Id = user.Id
	resp.SUCCESS(c, userVo)

}

type userProfile struct {
	Id          uint   `json:"id"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Power       int    `json:"power"`
	ExpiredTime int64  `json:"expired_time"`
	Vip         bool   `json:"vip"`
}

func (h *UserHandler) Profile(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	h.DB.First(&user, user.Id)
	var profile userProfile
	err = utils.CopyObject(user, &profile)
	if err != nil {
		logger.Error("对象拷贝失败：", err.Error())
		resp.ERROR(c, "获取用户信息失败")
		return
	}

	profile.Id = user.Id
	resp.SUCCESS(c, profile)
}

func (h *UserHandler) ProfileUpdate(c *gin.Context) {
	var data userProfile
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}
	h.DB.First(&user, user.Id)
	user.Avatar = data.Avatar
	user.Nickname = data.Nickname
	res := h.DB.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c, "更新用户信息失败")
		return
	}

	resp.SUCCESS(c)
}

// UpdatePass 更新密码
func (h *UserHandler) UpdatePass(c *gin.Context) {
	var data struct {
		OldPass  string `json:"old_pass"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Password) < 8 {
		resp.ERROR(c, "密码长度不能少于8个字符")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	password := utils.GenPassword(data.OldPass, user.Salt)
	logger.Debugf(user.Salt, ",", user.Password, ",", password, ",", data.OldPass)
	if password != user.Password {
		resp.ERROR(c, "原密码错误")
		return
	}

	newPass := utils.GenPassword(data.Password, user.Salt)
	err = h.DB.Model(&user).UpdateColumn("password", newPass).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// ResetPass 找回密码
func (h *UserHandler) ResetPass(c *gin.Context) {
	var data struct {
		Type     string `json:"type"`     // 验证类别：mobile, email
		Mobile   string `json:"mobile"`   // 手机号
		Email    string `json:"email"`    // 邮箱地址
		Code     string `json:"code"`     // 验证码
		Password string `json:"password"` // 新密码
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	var key string
	if data.Type == "email" {
		session = session.Where("email", data.Email)
		key = CodeStorePrefix + data.Email
	} else if data.Type == "mobile" {
		session = session.Where("mobile", data.Mobile)
		key = CodeStorePrefix + data.Mobile
	} else {
		resp.ERROR(c, "验证类别错误")
		return
	}
	var user model.User
	err := session.First(&user).Error
	if err != nil {
		resp.ERROR(c, "用户不存在！")
		return
	}

	// 检查验证码
	code, err := h.redis.Get(c, key).Result()
	if err != nil || code != data.Code {
		resp.ERROR(c, "验证码错误")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	err = h.DB.Model(&user).UpdateColumn("password", password).Error
	if err != nil {
		resp.ERROR(c, err.Error())
	} else {
		h.redis.Del(c, key)
		resp.SUCCESS(c)
	}
}

// BindMobile 绑定手机号
func (h *UserHandler) BindMobile(c *gin.Context) {
	var data struct {
		Mobile string `json:"mobile"`
		Code   string `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Mobile
	code, err := h.redis.Get(c, key).Result()
	if err != nil || code != data.Code {
		resp.ERROR(c, "验证码错误")
		return
	}

	// 检查手机号是否被其他账号绑定
	var item model.User
	res := h.DB.Where("mobile", data.Mobile).First(&item)
	if res.Error == nil {
		resp.ERROR(c, "该手机号已经绑定了其他账号，请更换手机号")
		return
	}

	userId := h.GetLoginUserId(c)

	err = h.DB.Model(&item).Where("id", userId).UpdateColumn("mobile", data.Mobile).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	_ = h.redis.Del(c, key) // 删除短信验证码
	resp.SUCCESS(c)
}

// BindEmail 绑定邮箱
func (h *UserHandler) BindEmail(c *gin.Context) {
	var data struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Email
	code, err := h.redis.Get(c, key).Result()
	if err != nil || code != data.Code {
		resp.ERROR(c, "验证码错误")
		return
	}

	// 检查手机号是否被其他账号绑定
	var item model.User
	res := h.DB.Where("email", data.Email).First(&item)
	if res.Error == nil {
		resp.ERROR(c, "该邮箱地址已经绑定了其他账号，请更邮箱地址")
		return
	}

	userId := h.GetLoginUserId(c)

	err = h.DB.Model(&item).Where("id", userId).UpdateColumn("email", data.Email).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	_ = h.redis.Del(c, key) // 删除短信验证码
	resp.SUCCESS(c)
}

// SignIn 每日签到
func (h *UserHandler) SignIn(c *gin.Context) {
	// 获取当前日期
	date := time.Now().Format("2006-01-02")

	// 检查是否已经签到
	userId := h.GetLoginUserId(c)
	key := fmt.Sprintf("signin/%d/%s", userId, date)
	var signIn bool
	err := h.levelDB.Get(key, &signIn)
	if err == nil && signIn {
		resp.ERROR(c, "今日已签到，请明日再来！")
		return
	}

	// 签到
	h.levelDB.Put(key, true)
	if h.App.SysConfig.Base.DailyPower > 0 {
		h.userService.IncreasePower(userId, h.App.SysConfig.Base.DailyPower, model.PowerLog{
			Type:   types.PowerSignIn,
			Model:  "SignIn",
			Remark: fmt.Sprintf("每日签到奖励，金额：%d", h.App.SysConfig.Base.DailyPower),
		})
	}
	resp.SUCCESS(c)
}
