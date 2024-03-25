package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	handler.BaseHandler
}

func NewUserHandler(app *core.AppServer, db *gorm.DB) *UserHandler {
	return &UserHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// List 用户列表
func (h *UserHandler) List(c *gin.Context) {
	if err := utils.CheckPermission(c, h.DB); err != nil {
		resp.NotPermission(c)
		return
	}

	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	username := h.GetTrim(c, "username")

	offset := (page - 1) * pageSize
	var items []model.User
	var users = make([]vo.User, 0)
	var total int64

	session := h.DB.Session(&gorm.Session{})
	if username != "" {
		session = session.Where("username LIKE ?", "%"+username+"%")
	}

	session.Model(&model.User{}).Count(&total)
	res := session.Offset(offset).Limit(pageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var user vo.User
			err := utils.CopyObject(item, &user)
			if err == nil {
				user.Id = item.Id
				user.CreatedAt = item.CreatedAt.Unix()
				user.UpdatedAt = item.UpdatedAt.Unix()
				users = append(users, user)
			} else {
				logger.Error(err)
			}
		}
	}
	pageVo := vo.NewPage(total, page, pageSize, users)
	resp.SUCCESS(c, pageVo)
}

func (h *UserHandler) Save(c *gin.Context) {
	var data struct {
		Id          uint     `json:"id"`
		Password    string   `json:"password"`
		Username    string   `json:"username"`
		ChatRoles   []string `json:"chat_roles"`
		ChatModels  []int    `json:"chat_models"`
		ExpiredTime string   `json:"expired_time"`
		Status      bool     `json:"status"`
		Vip         bool     `json:"vip"`
		Power       int      `json:"power"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var user = model.User{}
	var res *gorm.DB
	var userVo vo.User
	if data.Id > 0 { // 更新
		res = h.DB.Where("id", data.Id).First(&user)
		if res.Error != nil {
			resp.ERROR(c, "user not found")
			return
		}
		var oldPower = user.Power
		user.Username = data.Username
		user.Status = data.Status
		user.Vip = data.Vip
		user.Power = data.Power
		user.ChatRoles = utils.JsonEncode(data.ChatRoles)
		user.ChatModels = utils.JsonEncode(data.ChatModels)
		user.ExpiredTime = utils.Str2stamp(data.ExpiredTime)

		res = h.DB.Select("username", "status", "vip", "power", "chat_roles_json", "chat_models_json", "expired_time").Updates(&user)
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
		// 记录算力日志
		if oldPower != user.Power {
			mark := types.PowerAdd
			amount := user.Power - oldPower
			if oldPower > user.Power {
				mark = types.PowerSub
				amount = oldPower - user.Power
			}
			h.DB.Create(&model.PowerLog{
				UserId:    user.Id,
				Username:  user.Username,
				Type:      types.PowerGift,
				Amount:    amount,
				Balance:   user.Power,
				Mark:      mark,
				Model:     "管理员",
				Remark:    fmt.Sprintf("后台管理员强制修改用户算力，修改前：%d,修改后:%d, 管理员ID：%d", oldPower, user.Power, h.GetLoginUserId(c)),
				CreatedAt: time.Now(),
			})
		}
	} else {
		salt := utils.RandString(8)
		u := model.User{
			Username:    data.Username,
			Nickname:    fmt.Sprintf("极客学长@%d", utils.RandomNumber(6)),
			Password:    utils.GenPassword(data.Password, salt),
			Avatar:      "/images/avatar/user.png",
			Salt:        salt,
			Power:       data.Power,
			Status:      true,
			ChatRoles:   utils.JsonEncode(data.ChatRoles),
			ChatModels:  utils.JsonEncode(data.ChatModels),
			ExpiredTime: utils.Str2stamp(data.ExpiredTime),
		}
		res = h.DB.Create(&u)
		_ = utils.CopyObject(u, &userVo)
		userVo.Id = u.Id
		userVo.CreatedAt = u.CreatedAt.Unix()
		userVo.UpdatedAt = u.UpdatedAt.Unix()
	}

	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c, userVo)
}

// ResetPass 重置密码
func (h *UserHandler) ResetPass(c *gin.Context) {
	var data struct {
		Id       uint
		Password string
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.User
	res := h.DB.First(&user, data.Id)
	if res.Error != nil {
		resp.ERROR(c, "No user found")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	user.Password = password
	res = h.DB.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c)
	} else {
		resp.SUCCESS(c)
	}
}

func (h *UserHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 删除用户
	res := h.DB.Where("id = ?", id).Delete(&model.User{})
	if res.Error != nil {
		resp.ERROR(c, "删除失败")
		return
	}

	// 删除聊天记录
	h.DB.Where("user_id = ?", id).Delete(&model.ChatItem{})
	// 删除聊天历史记录
	h.DB.Where("user_id = ?", id).Delete(&model.ChatMessage{})
	// 删除登录日志
	h.DB.Where("user_id = ?", id).Delete(&model.UserLoginLog{})
	// 删除算力日志
	h.DB.Where("user_id = ?", id).Delete(&model.PowerLog{})
	// 删除众筹日志
	h.DB.Where("user_id = ?", id).Delete(&model.Reward{})
	// 删除绘图任务
	h.DB.Where("user_id = ?", id).Delete(&model.MidJourneyJob{})
	h.DB.Where("user_id = ?", id).Delete(&model.SdJob{})
	//  删除订单
	h.DB.Where("user_id = ?", id).Delete(&model.Order{})
	resp.SUCCESS(c)
}

func (h *UserHandler) LoginLog(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	var total int64
	h.DB.Model(&model.UserLoginLog{}).Count(&total)
	offset := (page - 1) * pageSize
	var items []model.UserLoginLog
	res := h.DB.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "获取数据失败")
		return
	}
	var logs []vo.UserLoginLog
	for _, v := range items {
		var log vo.UserLoginLog
		err := utils.CopyObject(v, &log)
		if err == nil {
			log.Id = v.Id
			log.CreatedAt = v.CreatedAt.Unix()
			logs = append(logs, log)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, logs))
}
