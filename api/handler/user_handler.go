package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"gorm.io/gorm"
)

type UserHandler struct {
	BaseHandler
	db       *gorm.DB
	searcher *xdb.Searcher
	leveldb  *store.LevelDB
	redis    *redis.Client
}

func NewUserHandler(
	app *core.AppServer,
	db *gorm.DB,
	searcher *xdb.Searcher,
	levelDB *store.LevelDB,
	client *redis.Client) *UserHandler {
	handler := &UserHandler{db: db, searcher: searcher, leveldb: levelDB, redis: client}
	handler.App = app
	return handler
}

// Register user register
func (h *UserHandler) Register(c *gin.Context) {
	// parameters process
	var data struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
		Code     int    `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	data.Password = strings.TrimSpace(data.Password)

	if len(data.Mobile) < 10 {
		resp.ERROR(c, "请输入合法的手机号")
		return
	}
	if len(data.Password) < 8 {
		resp.ERROR(c, "密码长度不能少于8个字符")
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Mobile
	if h.App.SysConfig.EnabledMsg {
		var code int
		err := h.leveldb.Get(key, &code)
		if err != nil || code != data.Code {
			resp.ERROR(c, "短信验证码错误")
			return
		}
	}

	// check if the username is exists
	var item model.User
	res := h.db.Where("mobile = ?", data.Mobile).First(&item)
	if res.RowsAffected > 0 {
		resp.ERROR(c, "该手机号码已经被注册，请更换其他手机号")
		return
	}

	salt := utils.RandString(8)
	user := model.User{
		Password:   utils.GenPassword(data.Password, salt),
		Avatar:     "/images/avatar/user.png",
		Salt:       salt,
		Status:     true,
		Mobile:     data.Mobile,
		ChatRoles:  utils.JsonEncode([]string{"gpt"}),               // 默认只订阅通用助手角色
		ChatModels: utils.JsonEncode(h.App.SysConfig.DefaultModels), // 默认开通的模型
		ChatConfig: utils.JsonEncode(types.UserChatConfig{
			ApiKeys: map[types.Platform]string{
				types.OpenAI:  "",
				types.Azure:   "",
				types.ChatGLM: "",
			},
		}),
		Calls:    h.App.SysConfig.UserInitCalls,
		ImgCalls: h.App.SysConfig.InitImgCalls,
	}
	res = h.db.Create(&user)
	if res.Error != nil {
		resp.ERROR(c, "保存数据失败")
		logger.Error(res.Error)
		return
	}

	if h.App.SysConfig.EnabledMsg {
		_ = h.leveldb.Delete(key) // 注册成功，删除短信验证码
	}

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
		Mobile   string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var user model.User
	res := h.db.Where("mobile = ?", data.Mobile).First(&user)
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
	h.db.Model(&user).Updates(user)

	h.db.Create(&model.UserLoginLog{
		UserId:       user.Id,
		Username:     user.Mobile,
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
	user, err := utils.GetLoginUser(c, h.db)
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
	Id          uint                 `json:"id"`
	Mobile      string               `json:"mobile"`
	Avatar      string               `json:"avatar"`
	ChatConfig  types.UserChatConfig `json:"chat_config"`
	Calls       int                  `json:"calls"`
	ImgCalls    int                  `json:"img_calls"`
	TotalTokens int64                `json:"total_tokens"`
}

func (h *UserHandler) Profile(c *gin.Context) {
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	h.db.First(&user, user.Id)
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

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}
	h.db.First(&user, user.Id)
	user.Avatar = data.Avatar
	user.ChatConfig = utils.JsonEncode(data.ChatConfig)
	res := h.db.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c, "更新用户信息失败")
		return
	}

	resp.SUCCESS(c)
}

// Password 更新密码
func (h *UserHandler) Password(c *gin.Context) {
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

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	password := utils.GenPassword(data.OldPass, user.Salt)
	logger.Info(user.Salt, ",", user.Password, ",", password, ",", data.OldPass)
	if password != user.Password {
		resp.ERROR(c, "原密码错误")
		return
	}

	newPass := utils.GenPassword(data.Password, user.Salt)
	res := h.db.Model(&user).UpdateColumn("password", newPass)
	if res.Error != nil {
		logger.Error("更新数据库失败: ", res.Error)
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c)
}

// BindMobile 绑定手机号
func (h *UserHandler) BindMobile(c *gin.Context) {
	var data struct {
		Mobile string `json:"mobile"`
		Code   int    `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 检查手机号是否被其他账号绑定
	var item model.User
	res := h.db.Where("mobile = ?", data.Mobile).First(&item)
	if res.Error == nil {
		resp.ERROR(c, "该手机号已经被其他账号绑定")
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Mobile
	var code int
	err := h.leveldb.Get(key, &code)
	if err != nil || code != data.Code {
		resp.ERROR(c, "短信验证码错误")
		return
	}

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	res = h.db.Model(&user).UpdateColumn("mobile", data.Mobile)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败")
		return
	}

	_ = h.leveldb.Delete(key) // 删除短信验证码
	resp.SUCCESS(c)
}
