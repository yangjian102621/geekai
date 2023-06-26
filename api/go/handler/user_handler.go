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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"gorm.io/gorm"
)

type UserHandler struct {
	BaseHandler
	db       *gorm.DB
	searcher *xdb.Searcher
}

func NewUserHandler(app *core.AppServer, db *gorm.DB, searcher *xdb.Searcher) *UserHandler {
	handler := &UserHandler{db: db, searcher: searcher}
	handler.App = app
	return handler
}

// Register user register
func (h *UserHandler) Register(c *gin.Context) {
	// parameters process
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	data.Username = strings.TrimSpace(data.Username)
	data.Password = strings.TrimSpace(data.Password)

	if len(data.Username) < 5 {
		resp.ERROR(c, "用户名长度不能少于5个字符")
		return
	}
	if len(data.Password) < 8 {
		resp.ERROR(c, "密码长度不能少于8个字符")
		return
	}

	// check if the username is exists
	var item model.User
	tx := h.db.Where("username = ?", data.Username).First(&item)
	if tx.RowsAffected > 0 {
		resp.ERROR(c, "用户名已存在")
		return
	}
	// 默认订阅所有角色
	var chatRoles []model.ChatRole
	h.db.Find(&chatRoles)
	var roleKeys = make([]string, 0)
	for _, r := range chatRoles {
		roleKeys = append(roleKeys, r.Key)
	}

	salt := utils.RandString(8)
	user := model.User{
		Username:  data.Username,
		Password:  utils.GenPassword(data.Password, salt),
		Nickname:  fmt.Sprintf("极客学长@%d", utils.RandomNumber(5)),
		Avatar:    "images/avatar/user.png",
		Salt:      salt,
		Status:    true,
		ChatRoles: utils.JsonEncode(roleKeys),
		ChatConfig: utils.JsonEncode(types.ChatConfig{
			Temperature:   h.App.ChatConfig.Temperature,
			MaxTokens:     h.App.ChatConfig.MaxTokens,
			EnableContext: h.App.ChatConfig.EnableContext,
			EnableHistory: true,
			Model:         h.App.ChatConfig.Model,
			ApiKey:        "",
		}),
	}
	// 初始化调用次数
	var cfg model.Config
	h.db.Where("marker = ?", "system").First(&cfg)
	var config types.SystemConfig
	err := utils.JsonDecode(cfg.Config, &config)
	if err != nil || config.UserInitCalls <= 0 {
		user.Calls = types.UserInitCalls
	} else {
		user.Calls = config.UserInitCalls
	}
	res := h.db.Create(&user)
	if res.Error != nil {
		resp.ERROR(c, "保存数据失败")
		logger.Error(res.Error)
		return
	}

	resp.SUCCESS(c, user)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var data struct {
		Username string
		Password string
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var user model.User
	res := h.db.Where("username = ?", data.Username).First(&user)
	if res.Error != nil {
		resp.ERROR(c, "用户名不存在")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	if password != user.Password {
		resp.ERROR(c, "用户名或密码错误")
		return
	}

	// 更新最后登录时间和IP
	user.LastLoginIp = c.ClientIP()
	user.LastLoginAt = time.Now().Unix()
	h.db.Model(&user).Updates(user)

	sessionId := utils.RandString(42)
	err := utils.SetLoginUser(c, user)
	if err != nil {
		resp.ERROR(c, "保存会话失败")
		logger.Error("Error for save session: ", err)
		return
	}

	// 记录登录信息在服务端
	h.App.ChatSession.Put(sessionId, types.ChatSession{ClientIP: c.ClientIP(), UserId: user.Id, Username: data.Username, SessionId: sessionId})

	// 加载用户订阅的聊天角色
	var roleKeys []string
	err = utils.JsonDecode(user.ChatRoles, &roleKeys)
	var chatRoles interface{}
	if err == nil {
		var roles []model.ChatRole
		res = h.db.Where("marker IN ?", roleKeys).Find(&roles)
		if res.Error == err {
			type Item struct {
				Name string
				Key  string
				Icon string
			}
			items := make([]Item, 0)
			for _, r := range roles {
				items = append(items, Item{Name: r.Name, Key: r.Key, Icon: r.Icon})
			}
			chatRoles = items
		}
	}

	h.db.Create(&model.UserLoginLog{
		UserId:       user.Id,
		Username:     user.Username,
		LoginIp:      c.ClientIP(),
		LoginAddress: utils.Ip2Region(h.searcher, c.ClientIP()),
	})
	var chatConfig types.ChatConfig
	err = utils.JsonDecode(user.ChatConfig, &chatConfig)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"session_id":     sessionId,
		"id":             user.Id,
		"nickname":       user.Nickname,
		"avatar":         user.Avatar,
		"username":       user.Username,
		"tokens":         user.Tokens,
		"calls":          user.Calls,
		"expiredTime":    user.ExpiredTime,
		"chatRoles":      chatRoles,
		"api_key":        chatConfig.ApiKey,
		"model":          chatConfig.Model,
		"temperature":    chatConfig.Temperature,
		"max_tokens":     chatConfig.MaxTokens,
		"enable_context": chatConfig.EnableContext,
		"enable_history": chatConfig.EnableHistory,
	})
}

// Logout 注 销
func (h *UserHandler) Logout(c *gin.Context) {
	sessionId := c.GetHeader(types.SessionName)
	session := sessions.Default(c)
	session.Delete(types.SessionUser)
	err := session.Save()
	if err != nil {
		logger.Error("Error for save session: ", err)
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

func (h *UserHandler) ProfileUpdate(c *gin.Context) {
	var data vo.User
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
	user.Nickname = data.Nickname
	user.Avatar = data.Avatar

	var chatConfig types.ChatConfig
	err = utils.JsonDecode(user.ChatConfig, &chatConfig)
	if err != nil {
		resp.ERROR(c, "用户配置解析失败")
		return
	}

	chatConfig.EnableHistory = data.ChatConfig.EnableHistory
	chatConfig.EnableContext = data.ChatConfig.EnableContext
	chatConfig.Model = data.ChatConfig.Model
	chatConfig.MaxTokens = data.ChatConfig.MaxTokens
	chatConfig.ApiKey = data.ChatConfig.ApiKey
	chatConfig.Temperature = data.ChatConfig.Temperature

	user.ChatConfig = utils.JsonEncode(chatConfig)
	res := h.db.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c, "更新用户信息失败")
		return
	}
	resp.SUCCESS(c)
}

func (h *UserHandler) Profile(c *gin.Context) {
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	h.db.First(&user, user.Id)
	var userVo vo.User
	err = utils.CopyObject(user, &userVo)
	if err != nil {
		logger.Error("对象拷贝失败：", err.Error())
		resp.ERROR(c, "获取用户信息失败")
		return
	}

	userVo.Id = user.Id
	userVo.CreatedAt = user.CreatedAt.Unix()
	userVo.UpdatedAt = user.UpdatedAt.Unix()
	resp.SUCCESS(c, userVo)
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
