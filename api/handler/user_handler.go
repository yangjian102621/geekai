package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/store"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserHandler struct {
	BaseHandler
	searcher       *xdb.Searcher
	redis          *redis.Client
	levelDB        *store.LevelDB
	captchaService *service.CaptchaService
	ledgerService  *service.LedgerService
	wxLoginService *service.WxLoginService
	ipSearcher     *xdb.Searcher
	wechatService  *service.WxGzhService
	tokengineOAuth *service.TokengineOAuth
}

func NewUserHandler(
	app *core.AppServer,
	db *gorm.DB,
	searcher *xdb.Searcher,
	client *redis.Client,
	levelDB *store.LevelDB,
	captcha *service.CaptchaService,
	ledgerService *service.LedgerService,
	wxLoginService *service.WxLoginService,
	wechatService *service.WxGzhService,
	ipSearcher *xdb.Searcher) *UserHandler {
	return &UserHandler{
		BaseHandler:    BaseHandler{DB: db, App: app},
		searcher:       searcher,
		redis:          client,
		levelDB:        levelDB,
		captchaService: captcha,
		ledgerService:  ledgerService,
		wxLoginService: wxLoginService,
		wechatService:  wechatService,
		ipSearcher:     ipSearcher,
		tokengineOAuth: service.NewTokengineOAuth(app.Config.TokengineOAuth),
	}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(r *core.Router) {
	// 公开接口，不需要授权。
	// 注意：logout 也在公开段 —— 登出必须允许"token 已失效"的用户调用，
	// 否则用户会卡在无法登出、而服务端会话 key 仍在有效期内继续可用。
	public := r.PublicGroup("/api/user/")
	public.POST("register", h.Register)
	public.POST("login", h.Login)
	public.POST("login/mobile", h.MobileLogin)
	public.GET("login/tokengine/config", h.TokengineConfig)
	public.GET("login/tokengine/start", h.TokengineStart)
	public.GET("login/tokengine/callback", h.TokengineCallback)
	public.POST("login/tokengine/complete", h.TokengineComplete)
	public.POST("resetPass", h.ResetPass)
	public.GET("login/qrcode", h.GetWxLoginQRCode)
	public.POST("login/callback", h.WxLoginCallback)
	public.GET("login/status", h.GetWxLoginState)
	public.GET("logout", h.Logout)
	public.POST("wxAuthLogin", h.WxAuthLogin)
	checkin := r.PublicGroup("/api")
	checkin.POST("/checkin", h.TokengineCheckin)
	checkin.POST("/v1/checkin", h.TokengineCheckin)

	// 需要用户授权的接口
	group := r.UserGroup("/api/user/")
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

var registrationMobilePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)
var registrationCodePattern = regexp.MustCompile(`^\d{6}$`)

// Register 使用自定义用户名和已验证的手机号注册，密码可选。
func (h *UserHandler) Register(c *gin.Context) {
	// parameters process
	var data struct {
		Username   string `json:"username"`
		Mobile     string `json:"mobile"`
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
	if !h.App.SysConfig.Base.EnabledRegister {
		resp.ERROR(c, "注册功能已关闭")
		return
	}
	if !utils.Contains(h.App.SysConfig.Base.RegisterWays, "username") && !utils.Contains(h.App.SysConfig.Base.RegisterWays, "mobile") {
		resp.ERROR(c, "系统已禁用手机号验证注册")
		return
	}
	data.Username = strings.TrimSpace(data.Username)
	data.Mobile = strings.TrimSpace(data.Mobile)
	data.Code = strings.TrimSpace(data.Code)
	data.Password = strings.TrimSpace(data.Password)
	if data.Username == "" || len([]rune(data.Username)) > 30 {
		resp.ERROR(c, "请输入不超过30个字符的用户名")
		return
	}
	if !registrationMobilePattern.MatchString(data.Mobile) {
		resp.ERROR(c, "请输入合法的手机号")
		return
	}
	if !registrationCodePattern.MatchString(data.Code) {
		resp.ERROR(c, "请输入6位短信验证码")
		return
	}
	if data.Password != "" && (len([]rune(data.Password)) < 8 || len([]rune(data.Password)) > 16) {
		resp.ERROR(c, "密码长度应为8-16个字符")
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

	// 按手机号串行处理注册，避免并发请求重复绑定同一手机号。
	unlock, err := h.lockMobile(c, data.Mobile)
	if err != nil {
		resp.ERROR(c, "注册请求正在处理中，请稍后重试")
		return
	}
	defer unlock()

	var item model.User
	err = h.DB.Where("username = ? OR mobile = ?", data.Username, data.Mobile).First(&item).Error
	if err == nil {
		resp.ERROR(c, "用户名或手机号已被注册")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.ERROR(c, "查询用户失败，请稍后重试")
		return
	}

	if !h.consumeSmsCode(c, data.Mobile, data.Code) {
		resp.ERROR(c, "验证码错误或已过期，请重新获取")
		return
	}

	user, err := h.createNewUser(model.User{
		Username: data.Username,
		Mobile:   data.Mobile,
		Password: data.Password,
	}, data.InviteCode)
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
func (h *UserHandler) createNewUser(user model.User, code string) (model.User, error) {
	if user.OpenId != "" {
		user.Platform = "wechat"
		// 如果未设置昵称，则生成默认昵称
		if user.Nickname == "" {
			user.Nickname = fmt.Sprintf("微信用户@%d", utils.RandomNumber(6))
		}
		// 如果未设置用户名，则生成默认用户名
		if user.Username == "" {
			user.Username = fmt.Sprintf("wx@%d", utils.RandomNumber(8))
		}
		// 如果未设置密码，则生成默认密码
		if user.Password == "" {
			user.Password = "geekai123"
		}
	} else {
		if user.Nickname == "" && user.TokengineSubject == "" {
			user.Nickname = fmt.Sprintf("用户@%d", utils.RandomNumber(6))
		}
		if user.Username == "" || (user.Password == "" && user.Mobile == "" && (user.TokengineID == nil || *user.TokengineID == "")) {
			return user, fmt.Errorf("用户名和手机号不能为空")
		}
	}

	salt := utils.RandString(8)
	user.Salt = salt
	// 未设置密码时保留空值，避免空密码被哈希后用于密码登录。
	if user.Password != "" {
		user.Password = utils.GenPassword(user.Password, salt)
	}
	// 如果未设置头像，则使用默认头像
	if user.Avatar == "" {
		user.Avatar = "/images/avatar/user.png"
	}
	user.Status = true
	user.ChatConfig = "{}"
	user.ChatModels = "{}"
	user.Power = h.App.SysConfig.Base.InitPower

	// 创建用户
	tx := h.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	if code != "" {
		if err := h.applyRegistrationInvite(tx, &user, code); err != nil {
			tx.Rollback()
			return user, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return user, err
	}

	return user, nil
}

const maxInvitesPerUser = 30

func (h *UserHandler) applyRegistrationInvite(tx *gorm.DB, user *model.User, code string) error {
	var inviteCode model.InviteCode
	err := tx.Where("code = ?", code).First(&inviteCode).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("无效的邀请码")
	}
	if err != nil {
		return err
	}

	// Lock the inviter's counters so concurrent registrations share the same limit across codes.
	var counts []int
	if err := tx.Model(&model.InviteCode{}).Where("user_id = ?", inviteCode.UserId).
		Clauses(clause.Locking{Strength: "UPDATE"}).Pluck("reg_num", &counts).Error; err != nil {
		return err
	}
	total := 0
	for _, count := range counts {
		total += count
	}
	if total >= maxInvitesPerUser {
		return nil
	}

	if err := tx.Model(&model.InviteCode{}).Where("id = ?", inviteCode.Id).
		UpdateColumn("reg_num", gorm.Expr("reg_num + ?", 1)).Error; err != nil {
		return err
	}
	invitePower := h.App.SysConfig.Base.InvitePower
	if invitePower <= 0 {
		return nil
	}
	for _, reward := range []struct {
		userID  uint
		refType string
		remark  string
	}{
		{inviteCode.UserId, "invite", fmt.Sprintf("邀请用户注册奖励，邀请码：%s，新用户：%s", inviteCode.Code, user.Username)},
		{user.Id, "invite:invitee", fmt.Sprintf("受邀注册奖励，邀请码：%s", inviteCode.Code)},
	} {
		granted, err := h.ledgerService.Grant(tx, reward.userID, invitePower, service.LedgerRef{Type: reward.refType, ID: user.Id}, model.PowerLog{
			Type:   types.PowerInvite,
			Model:  "Invite",
			Remark: reward.remark,
		})
		if err != nil {
			return err
		}
		if !granted {
			return errors.New("邀请奖励已发放")
		}
	}
	if err := tx.Create(&model.InviteLog{
		InviterId:  inviteCode.UserId,
		UserId:     user.Id,
		Username:   user.Username,
		InviteCode: inviteCode.Code,
		Remark:     fmt.Sprintf("双方各奖励 %d 积分", invitePower),
	}).Error; err != nil {
		return err
	}
	user.Power += invitePower
	return nil
}

// doLogin 执行登录操作
func (h *UserHandler) doLogin(user *model.User, ip string) (string, error) {
	// 更新最后登录时间和IP
	user.LastLoginIp = ip
	user.LastLoginAt = time.Now().Unix()
	err := h.DB.Model(user).Updates(map[string]any{
		"last_login_ip": user.LastLoginIp,
		"last_login_at": user.LastLoginAt,
	}).Error
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
	GemIds      []uint `json:"gem_ids"`
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
	if user.GemIds != "" {
		var raw []interface{}
		if utils.JsonDecode(user.GemIds, &raw) == nil {
			for _, v := range raw {
				if n, ok := v.(float64); ok {
					profile.GemIds = append(profile.GemIds, uint(n))
				}
			}
		}
	}
	if profile.GemIds == nil {
		profile.GemIds = []uint{}
	}
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
	if data.GemIds != nil {
		if len(data.GemIds) > 8 {
			data.GemIds = data.GemIds[:8]
		}
		user.GemIds = utils.JsonEncode(data.GemIds)
	}
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
	switch data.Type {
	case "email":
		session = session.Where("email", data.Email)
		key = CodeStorePrefix + data.Email
	case "mobile":
		session = session.Where("mobile", data.Mobile)
		key = CodeStorePrefix + data.Mobile
	default:
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

	unlock, err := h.lockMobile(c, data.Mobile)
	if err != nil {
		resp.ERROR(c, "手机号请求正在处理中，请稍后重试")
		return
	}
	defer unlock()

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
	_, alreadyCheckedIn, err := h.checkIn(c.Request.Context(), h.GetLoginUserId(c), false)
	if err != nil {
		resp.ERROR(c, "签到失败，请稍后重试")
		return
	}
	if alreadyCheckedIn {
		resp.ERROR(c, "今日已签到，请明日再来！")
		return
	}
	resp.SUCCESS(c)
}

// 微信公众号 小程序授权登录
func (h *UserHandler) WxAuthLogin(c *gin.Context) {
	var data struct {
		Code       string `json:"code"`
		InviteCode string `json:"invite_code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 根据 code 获取 openid
	openID, accessToken, err := h.wechatService.GetOpenIDByCode(data.Code)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 获取微信用户昵称头像
	userInfo, err := h.wechatService.GetUserInfo(accessToken, openID)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	nickname := userInfo["nickname"].(string)
	headimgurl := userInfo["headimgurl"].(string)

	// 查询用户是否存在
	var user model.User
	h.DB.Where("openid = ?", openID).First(&user)
	if user.Id > 0 {
		// 用户存在，更新用户信息并登录
		user.Nickname = nickname
		user.Avatar = headimgurl
		if err := h.DB.Save(&user).Error; err != nil {
			resp.ERROR(c, "更新用户信息失败")
			return
		}

		token, err := h.doLogin(&user, c.ClientIP())
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}

		resp.SUCCESS(c, gin.H{"token": token, "user_id": user.Id, "username": user.Username})
		return
	}

	// 用户不存在，创建新用户
	user = model.User{
		OpenId:   openID,
		Nickname: nickname,
		Avatar:   headimgurl,
	}

	// 被邀请人也获得赠送积分
	if data.InviteCode != "" {
		user.Power = h.App.SysConfig.Base.InitPower * 2
	}

	user, err = h.createNewUser(user, data.InviteCode)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 自动登录
	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"token": token, "user_id": user.Id, "username": user.Username})
}
