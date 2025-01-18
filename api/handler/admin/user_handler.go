package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	handler.BaseHandler
	licenseService *service.LicenseService
	redis          *redis.Client
}

func NewUserHandler(app *core.AppServer, db *gorm.DB, licenseService *service.LicenseService, redisCli *redis.Client) *UserHandler {
	return &UserHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}, licenseService: licenseService, redis: redisCli}
}

// List 用户列表
func (h *UserHandler) List(c *gin.Context) {
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
	res := session.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
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
		Mobile      string   `json:"mobile"`
		Email       string   `json:"email"`
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
	// 检测最大注册人数
	var totalUser int64
	h.DB.Model(&model.User{}).Count(&totalUser)
	if h.licenseService.GetLicense().Configs.UserNum > 0 && int(totalUser) >= h.licenseService.GetLicense().Configs.UserNum {
		resp.ERROR(c, "当前注册用户数已达上限，请请升级 License")
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
		user.Email = data.Email
		user.Mobile = data.Mobile
		user.Status = data.Status
		user.Vip = data.Vip
		user.Power = data.Power
		user.ChatRoles = utils.JsonEncode(data.ChatRoles)
		user.ChatModels = utils.JsonEncode(data.ChatModels)
		user.ExpiredTime = utils.Str2stamp(data.ExpiredTime)

		res = h.DB.Select("username", "mobile", "email", "status", "vip", "power", "chat_roles_json", "chat_models_json", "expired_time").Updates(&user)

		if res.Error != nil {
			logger.Error("error with update database：", res.Error)
			resp.ERROR(c, res.Error.Error())
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
		// 如果禁用了用户，则将用户踢下线
		if user.Status == false {
			key := fmt.Sprintf("users/%v", user.Id)
			if _, err := h.redis.Del(c, key).Result(); err != nil {
				logger.Error("error with delete session: ", err)
			}
		}
	} else {
		// 检查用户是否已经存在
		h.DB.Where("username", data.Username).First(&user)
		if user.Id > 0 {
			resp.ERROR(c, "用户名已存在")
			return
		}

		salt := utils.RandString(8)
		u := model.User{
			Username:    data.Username,
			Password:    utils.GenPassword(data.Password, salt),
			Mobile:      data.Mobile,
			Email:       data.Email,
			Avatar:      "/images/avatar/user.png",
			Salt:        salt,
			Power:       data.Power,
			Status:      true,
			ChatRoles:   utils.JsonEncode(data.ChatRoles),
			ChatModels:  utils.JsonEncode(data.ChatModels),
			ExpiredTime: utils.Str2stamp(data.ExpiredTime),
		}
		if h.licenseService.GetLicense().Configs.DeCopy {
			u.Nickname = fmt.Sprintf("用户@%d", utils.RandomNumber(6))
		} else {
			u.Nickname = fmt.Sprintf("极客学长@%d", utils.RandomNumber(6))
		}
		res = h.DB.Create(&u)
		_ = utils.CopyObject(u, &userVo)
		userVo.Id = u.Id
		userVo.CreatedAt = u.CreatedAt.Unix()
		userVo.UpdatedAt = u.UpdatedAt.Unix()
	}

	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
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
	id := c.Query("id")
	ids := c.QueryArray("ids[]")
	if id != "" {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	tx := h.DB.Begin()
	var err error
	for _, id = range ids {
		// 删除用户
		if err = tx.Where("id", id).Delete(&model.User{}).Error; err != nil {
			break
		}
		// 删除聊天记录
		if err = tx.Unscoped().Where("user_id = ?", id).Delete(&model.ChatItem{}).Error; err != nil {
			break
		}
		// 删除聊天历史记录
		if err = tx.Unscoped().Where("user_id = ?", id).Delete(&model.ChatMessage{}).Error; err != nil {
			break
		}
		// 删除登录日志
		if err = tx.Where("user_id = ?", id).Delete(&model.UserLoginLog{}).Error; err != nil {
			break
		}
		// 删除算力日志
		if err = tx.Where("user_id = ?", id).Delete(&model.PowerLog{}).Error; err != nil {
			break
		}
		if err = tx.Where("user_id = ?", id).Delete(&model.InviteLog{}).Error; err != nil {
			break
		}
		// 删除众筹日志
		if err = tx.Where("user_id = ?", id).Delete(&model.Redeem{}).Error; err != nil {
			break
		}
		// 删除绘图任务
		if err = tx.Where("user_id = ?", id).Delete(&model.MidJourneyJob{}).Error; err != nil {
			break
		}
		if err = tx.Where("user_id = ?", id).Delete(&model.SdJob{}).Error; err != nil {
			break
		}
		if err = tx.Where("user_id = ?", id).Delete(&model.DallJob{}).Error; err != nil {
			break
		}
		if err = tx.Where("user_id = ?", id).Delete(&model.SunoJob{}).Error; err != nil {
			break
		}
		if err = tx.Where("user_id = ?", id).Delete(&model.VideoJob{}).Error; err != nil {
			break
		}
	}
	if err != nil {
		resp.ERROR(c, err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
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
