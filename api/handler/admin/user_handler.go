package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewUserHandler(app *core.AppServer, db *gorm.DB) *UserHandler {
	h := UserHandler{db: db}
	h.App = app
	return &h
}

// List 用户列表
func (h *UserHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	mobile := h.GetTrim(c, "mobile")

	offset := (page - 1) * pageSize
	var items []model.User
	var users = make([]vo.User, 0)
	var total int64

	session := h.db.Session(&gorm.Session{})
	if mobile != "" {
		session = session.Where("mobile LIKE ?", "%"+mobile+"%")
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
		Mobile      string   `json:"mobile"`
		Calls       int      `json:"calls"`
		ImgCalls    int      `json:"img_calls"`
		ChatRoles   []string `json:"chat_roles"`
		ChatModels  []string `json:"chat_models"`
		ExpiredTime string   `json:"expired_time"`
		Status      bool     `json:"status"`
		Vip         bool     `json:"vip"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var user = model.User{}
	var res *gorm.DB
	var userVo vo.User
	if data.Id > 0 { // 更新
		user.Id = data.Id
		// 此处需要用 map 更新，用结构体无法更新 0 值
		res = h.db.Model(&user).Updates(map[string]interface{}{
			"mobile":           data.Mobile,
			"calls":            data.Calls,
			"img_calls":        data.ImgCalls,
			"status":           data.Status,
			"vip":              data.Vip,
			"chat_roles_json":  utils.JsonEncode(data.ChatRoles),
			"chat_models_json": utils.JsonEncode(data.ChatModels),
			"expired_time":     utils.Str2stamp(data.ExpiredTime),
		})
	} else {
		salt := utils.RandString(8)
		u := model.User{
			Mobile:      data.Mobile,
			Password:    utils.GenPassword(data.Password, salt),
			Avatar:      "/images/avatar/user.png",
			Salt:        salt,
			Status:      true,
			ChatRoles:   utils.JsonEncode(data.ChatRoles),
			ChatModels:  utils.JsonEncode(data.ChatModels),
			ExpiredTime: utils.Str2stamp(data.ExpiredTime),
			ChatConfig: utils.JsonEncode(types.UserChatConfig{
				ApiKeys: map[types.Platform]string{
					types.OpenAI:  "",
					types.Azure:   "",
					types.ChatGLM: "",
				},
			}),
			Calls:    data.Calls,
			ImgCalls: data.ImgCalls,
		}
		res = h.db.Create(&u)
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
	res := h.db.First(&user, data.Id)
	if res.Error != nil {
		resp.ERROR(c, "No user found")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	user.Password = password
	res = h.db.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c)
	} else {
		resp.SUCCESS(c)
	}
}

func (h *UserHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id > 0 {
		tx := h.db.Begin()
		res := h.db.Where("id = ?", id).Delete(&model.User{})
		if res.Error != nil {
			resp.ERROR(c, "删除失败")
			return
		}
		// 删除聊天记录
		res = h.db.Where("user_id = ?", id).Delete(&model.ChatItem{})
		if res.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		// 删除聊天历史记录
		res = h.db.Where("user_id = ?", id).Delete(&model.HistoryMessage{})
		if res.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		// 删除登录日志
		res = h.db.Where("user_id = ?", id).Delete(&model.UserLoginLog{})
		if res.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		tx.Commit()
	}
	resp.SUCCESS(c)
}

func (h *UserHandler) LoginLog(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	var total int64
	h.db.Model(&model.UserLoginLog{}).Count(&total)
	offset := (page - 1) * pageSize
	var items []model.UserLoginLog
	res := h.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
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
