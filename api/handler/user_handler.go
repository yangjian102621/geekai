package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
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
	searcher *xdb.Searcher
	redis    *redis.Client
}

func NewUserHandler(
	app *core.AppServer,
	db *gorm.DB,
	searcher *xdb.Searcher,
	client *redis.Client) *UserHandler {
	return &UserHandler{BaseHandler: BaseHandler{DB: db, App: app}, searcher: searcher, redis: client}
}

// Register user register
func (h *UserHandler) Register(c *gin.Context) {
	// parameters process
	var data struct {
		RegWay     string `json:"reg_way"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Code       string `json:"code"`
		InviteCode string `json:"invite_code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	data.Password = strings.TrimSpace(data.Password)
	if len(data.Password) < 8 {
		resp.ERROR(c, "密码长度不能少于8个字符")
		return
	}

	// 检查验证码
	var key string
	if data.RegWay == "email" || data.RegWay == "mobile" || data.Code != "" {
		key = CodeStorePrefix + data.Username
		code, err := h.redis.Get(c, key).Result()
		if err != nil || code != data.Code {
			resp.ERROR(c, "验证码错误")
			return
		}
	}

	// 验证邀请码
	inviteCode := model.InviteCode{}
	if data.InviteCode != "" {
		res := h.DB.Where("code = ?", data.InviteCode).First(&inviteCode)
		if res.Error != nil {
			resp.ERROR(c, "无效的邀请码")
			return
		}
	}

	// check if the username is exists
	var item model.User
	res := h.DB.Where("username = ?", data.Username).First(&item)
	if item.Id > 0 {
		resp.ERROR(c, "该用户名已经被注册")
		return
	}

	salt := utils.RandString(8)
	user := model.User{
		Username:   data.Username,
		Password:   utils.GenPassword(data.Password, salt),
		Nickname:   fmt.Sprintf("极客学长@%d", utils.RandomNumber(6)),
		Avatar:     "/images/avatar/user.png",
		Salt:       salt,
		Status:     true,
		ChatRoles:  utils.JsonEncode([]string{"gpt"}),               // 默认只订阅通用助手角色
		ChatModels: utils.JsonEncode(h.App.SysConfig.DefaultModels), // 默认开通的模型
		Power:      h.App.SysConfig.InitPower,
	}

	res = h.DB.Create(&user)
	if res.Error != nil {
		resp.ERROR(c, "保存数据失败")
		logger.Error(res.Error)
		return
	}

	// 记录邀请关系
	if data.InviteCode != "" {
		// 增加邀请数量
		h.DB.Model(&model.InviteCode{}).Where("code = ?", data.InviteCode).UpdateColumn("reg_num", gorm.Expr("reg_num + ?", 1))
		if h.App.SysConfig.InvitePower > 0 {
			h.DB.Model(&model.User{}).Where("id = ?", inviteCode.UserId).UpdateColumn("power", gorm.Expr("power + ?", h.App.SysConfig.InvitePower))
			// 记录邀请算力充值日志
			var inviter model.User
			h.DB.Where("id", inviteCode.UserId).First(&inviter)
			h.DB.Create(&model.PowerLog{
				UserId:    inviter.Id,
				Username:  inviter.Username,
				Type:      types.PowerInvite,
				Amount:    h.App.SysConfig.InvitePower,
				Balance:   inviter.Power,
				Mark:      types.PowerAdd,
				Model:     "",
				Remark:    fmt.Sprintf("邀请用户注册奖励，金额：%d，邀请码：%s，新用户：%s", h.App.SysConfig.InvitePower, inviteCode.Code, user.Username),
				CreatedAt: time.Now(),
			})
		}

		// 添加邀请记录
		h.DB.Create(&model.InviteLog{
			InviterId:  inviteCode.UserId,
			UserId:     user.Id,
			Username:   user.Username,
			InviteCode: inviteCode.Code,
			Remark:     fmt.Sprintf("奖励 %d 算力", h.App.SysConfig.InvitePower),
		})
	}

	_ = h.redis.Del(c, key) // 注册成功，删除短信验证码

	// 自动登录创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
	if err != nil {
		resp.ERROR(c, "Failed to generate token, "+err.Error())
		return
	}
	// 保存到 redis
	key = fmt.Sprintf("users/%d", user.Id)
	if _, err := h.redis.Set(c, key, tokenString, 0).Result(); err != nil {
		resp.ERROR(c, "error with save token: "+err.Error())
		return
	}
	resp.SUCCESS(c, tokenString)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
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

	if user.Status == false {
		resp.ERROR(c, "该用户已被禁止登录，请联系管理员")
		return
	}

	// 更新最后登录时间和IP
	user.LastLoginIp = c.ClientIP()
	user.LastLoginAt = time.Now().Unix()
	h.DB.Model(&user).Updates(user)

	h.DB.Create(&model.UserLoginLog{
		UserId:       user.Id,
		Username:     user.Username,
		LoginIp:      c.ClientIP(),
		LoginAddress: utils.Ip2Region(h.searcher, c.ClientIP()),
	})

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
	if err != nil {
		resp.ERROR(c, "Failed to generate token, "+err.Error())
		return
	}
	// 保存到 redis
	key := fmt.Sprintf("users/%d", user.Id)
	if _, err := h.redis.Set(c, key, tokenString, 0).Result(); err != nil {
		resp.ERROR(c, "error with save token: "+err.Error())
		return
	}
	resp.SUCCESS(c, tokenString)
}

// Logout 注 销
func (h *UserHandler) Logout(c *gin.Context) {
	sessionId := c.GetHeader(types.ChatTokenHeader)
	key := h.GetUserKey(c)
	if _, err := h.redis.Del(c, key).Result(); err != nil {
		logger.Error("error with delete session: ", err)
	}
	// 删除 websocket 会话列表
	h.App.ChatSession.Delete(sessionId)
	// 关闭 socket 连接
	client := h.App.ChatClients.Get(sessionId)
	if client != nil {
		client.Close()
	}
	resp.SUCCESS(c)
}

// Session 获取/验证会话
func (h *UserHandler) Session(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err == nil {
		var userVo vo.User
		err := utils.CopyObject(user, &userVo)
		if err != nil {
			resp.ERROR(c)
		}
		userVo.Id = user.Id
		resp.SUCCESS(c, userVo)
	} else {
		resp.NotAuth(c)
	}

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
	res := h.DB.Model(&user).UpdateColumn("password", newPass)
	if res.Error != nil {
		logger.Error("更新数据库失败: ", res.Error)
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c)
}

// ResetPass 重置密码
func (h *UserHandler) ResetPass(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Code     string `json:"code"`     // 验证码
		Password string `json:"password"` // 新密码
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.User
	res := h.DB.Where("username", data.Username).First(&user)
	if res.Error != nil {
		resp.ERROR(c, "用户不存在！")
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Username
	code, err := h.redis.Get(c, key).Result()
	if err != nil || code != data.Code {
		resp.ERROR(c, "短信验证码错误")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	user.Password = password
	res = h.DB.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c)
	} else {
		h.redis.Del(c, key)
		resp.SUCCESS(c)
	}
}

// BindUsername 重置账号
func (h *UserHandler) BindUsername(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Username
	code, err := h.redis.Get(c, key).Result()
	if err != nil || code != data.Code {
		resp.ERROR(c, "验证码错误")
		return
	}

	// 检查手机号是否被其他账号绑定
	var item model.User
	res := h.DB.Where("username = ?", data.Username).First(&item)
	if res.Error == nil {
		resp.ERROR(c, "该账号已经被其他账号绑定")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	res = h.DB.Model(&user).UpdateColumn("username", data.Username)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败")
		return
	}

	_ = h.redis.Del(c, key) // 删除短信验证码
	resp.SUCCESS(c)
}
