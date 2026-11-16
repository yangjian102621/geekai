package handler

import (
	"context"
	"errors"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MobileLogin 手机号验证码登录，未注册时自动创建无密码账号。
func (h *UserHandler) MobileLogin(c *gin.Context) {
	var data struct {
		Mobile     string `json:"mobile"`
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
	data.Mobile = strings.TrimSpace(data.Mobile)
	data.Code = strings.TrimSpace(data.Code)
	if !registrationMobilePattern.MatchString(data.Mobile) {
		resp.ERROR(c, "请输入合法的手机号")
		return
	}
	if !registrationCodePattern.MatchString(data.Code) {
		resp.ERROR(c, "请输入6位短信验证码")
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
			resp.ERROR(c, "请先完成人机验证")
			return
		}
	}

	unlock, err := h.lockMobile(c, data.Mobile)
	if err != nil {
		resp.ERROR(c, "登录请求正在处理中，请稍后重试")
		return
	}
	defer unlock()

	var user model.User
	err = h.DB.Where("mobile = ?", data.Mobile).First(&user).Error
	newUser := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !newUser {
		resp.ERROR(c, "查询用户失败，请稍后重试")
		return
	}
	if newUser {
		if !h.App.SysConfig.Base.EnabledRegister {
			resp.ERROR(c, "注册功能已关闭，请使用已注册的手机号登录")
			return
		}
		if !utils.Contains(h.App.SysConfig.Base.RegisterWays, "username") && !utils.Contains(h.App.SysConfig.Base.RegisterWays, "mobile") {
			resp.ERROR(c, "系统已禁用手机号验证注册")
			return
		}
	} else if !user.Status {
		resp.ERROR(c, "该用户已被禁止登录，请联系管理员")
		return
	}

	if !h.consumeSmsCode(c, data.Mobile, data.Code) {
		resp.ERROR(c, "验证码错误或已过期，请重新获取")
		return
	}
	if newUser {
		username := data.Mobile
		var existing model.User
		err = h.DB.Where("username = ?", username).First(&existing).Error
		if err == nil {
			// 同名但未绑定该手机号的账号不属于本次登录用户。
			username += "_" + utils.RandString(8)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			resp.ERROR(c, "查询用户失败，请稍后重试")
			return
		}
		user, err = h.createNewUser(model.User{Username: username, Mobile: data.Mobile}, strings.TrimSpace(data.InviteCode))
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
	resp.SUCCESS(c, gin.H{"token": token, "user_id": user.Id, "username": user.Username})
}

// 注册和快捷注册共用手机号锁，避免并发创建重复绑定。
func (h *UserHandler) lockMobile(ctx context.Context, mobile string) (func(), error) {
	key := "register/mobile/" + mobile
	token := utils.RandString(32)
	locked, err := h.redis.SetNX(ctx, key, token, time.Minute).Result()
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, errors.New("手机号请求正在处理中")
	}
	return func() {
		h.redis.Eval(ctx, `if redis.call('GET', KEYS[1]) == ARGV[1] then return redis.call('DEL', KEYS[1]) end return 0`, []string{key}, token)
	}, nil
}

// 原子消费验证码；五次错误后需重新获取，拒绝无过期时间的旧验证码。
func (h *UserHandler) consumeSmsCode(ctx context.Context, mobile, code string) bool {
	key := CodeStorePrefix + mobile
	verified, err := h.redis.Eval(ctx, `
local code = redis.call('GET', KEYS[1])
if not code or redis.call('TTL', KEYS[1]) <= 0 then return 0 end
if code ~= ARGV[1] then
    local attempts = redis.call('INCR', KEYS[2])
    if attempts == 1 then redis.call('EXPIRE', KEYS[2], 300) end
    if attempts >= 5 then redis.call('DEL', KEYS[1]) end
    return 0
end
redis.call('DEL', KEYS[1], KEYS[2])
return 1
`, []string{key, key + "/attempts"}, code).Int()
	return err == nil && verified == 1
}
