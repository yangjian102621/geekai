package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	handler.BaseHandler
	redis *redis.Client
}

func NewUserHandler(app *core.AppServer, db *gorm.DB, redisCli *redis.Client) *UserHandler {
	return &UserHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}, redis: redisCli}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/admin/user/")

	// 需要管理员授权的接口
	group.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		group.GET("list", h.List)
		group.POST("save", h.Save)
		group.GET("remove", h.Remove)
		group.GET("loginLog", h.LoginLog)
		group.GET("genLoginLink", h.GenLoginLink)
		group.POST("resetPass", h.ResetPass)
		group.GET("import/template", h.ImportTemplate)
		group.POST("import", h.ImportUsers)
	}
}

// ImportTemplate 下载用户导入模板
func (h *UserHandler) ImportTemplate(c *gin.Context) {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	// 表头
	headers := []string{"用户名", "密码", "手机", "邮箱", "剩余算力", "启用状态"}
	for i, title := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheetName, cell, title)
	}
	// 示例数据
	sample := []interface{}{"user001", "Passw0rd!", "13800000000", "user001@example.com", 100, 1}
	for i, v := range sample {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		_ = f.SetCellValue(sheetName, cell, v)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		logger.Error("failed to generate user import template: ", err)
		resp.ERROR(c, "生成模板失败")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"user_import_template.xlsx\"")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ImportUsers 批量导入用户
func (h *UserHandler) ImportUsers(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		resp.ERROR(c, "文件上传失败: "+err.Error())
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" {
		resp.ERROR(c, "只支持 .xlsx 格式的 Excel 文件")
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		resp.ERROR(c, "无法读取上传文件: "+err.Error())
		return
	}
	defer src.Close()

	// 读取到内存，避免多次读取问题
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(src); err != nil {
		resp.ERROR(c, "读取文件内容失败: "+err.Error())
		return
	}

	excel, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		resp.ERROR(c, "解析 Excel 失败: "+err.Error())
		return
	}
	defer func() {
		_ = excel.Close()
	}()

	rows, err := excel.GetRows("Sheet1")
	if err != nil {
		resp.ERROR(c, "读取工作表失败: "+err.Error())
		return
	}
	if len(rows) < 2 {
		resp.ERROR(c, "Excel 中没有可导入的数据")
		return
	}

	type rowError struct {
		Row   int    `json:"row"`
		Error string `json:"error"`
	}

	var (
		successCount int
		failedCount  int
		errorsList   []rowError
	)

	usernameSet := make(map[string]struct{})

	// 从第二行开始读取
	for index, row := range rows[1:] {
		line := index + 2 // Excel 行号
		if len(row) == 0 {
			continue
		}

		get := func(i int) string {
			if i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}

		username := get(0)
		password := get(1)
		mobile := get(2)
		email := get(3)
		powerStr := get(4)
		statusStr := get(5)

		// 基础校验
		if username == "" {
			failedCount++
			errorsList = append(errorsList, rowError{Row: line, Error: "用户名不能为空"})
			continue
		}
		if _, ok := usernameSet[username]; ok {
			failedCount++
			errorsList = append(errorsList, rowError{Row: line, Error: "同一文件中用户名重复"})
			continue
		}
		usernameSet[username] = struct{}{}

		if len(password) < 8 || len(password) > 16 {
			failedCount++
			errorsList = append(errorsList, rowError{Row: line, Error: "密码必须为 8-16 位"})
			continue
		}

		if mobile != "" && len(mobile) != 11 {
			failedCount++
			errorsList = append(errorsList, rowError{Row: line, Error: "手机号必须为 11 位"})
			continue
		}

		// 解析算力
		power := 0
		if powerStr != "" {
			p, err := strconv.Atoi(powerStr)
			if err != nil {
				failedCount++
				errorsList = append(errorsList, rowError{Row: line, Error: "剩余算力必须为数字"})
				continue
			}
			if p < 0 {
				failedCount++
				errorsList = append(errorsList, rowError{Row: line, Error: "剩余算力不能为负数"})
				continue
			}
			power = p
		}

		// 解析启用状态
		status := true
		if statusStr != "" {
			switch strings.TrimSpace(statusStr) {
			case "0", "否", "false", "停用":
				status = false
			case "1", "是", "true", "启用":
				status = true
			default:
				failedCount++
				errorsList = append(errorsList, rowError{Row: line, Error: "启用状态只支持 1/是 或 0/否"})
				continue
			}
		}

		// 检查用户名是否已存在
		var exist model.User
		if err = h.DB.Where("username = ?", username).First(&exist).Error; err == nil && exist.Id > 0 {
			failedCount++
			errorsList = append(errorsList, rowError{Row: line, Error: "用户名已存在"})
			continue
		}

		salt := utils.RandString(8)
		u := model.User{
			Username:    username,
			Password:    utils.GenPassword(password, salt),
			Mobile:      mobile,
			Email:       email,
			Avatar:      "/images/avatar/user.png",
			Salt:        salt,
			Power:       power,
			Status:      status,
			ChatRoles:   utils.JsonEncode([]string{}),
			ChatConfig:  "{}",
			ChatModels:  utils.JsonEncode([]int{}),
			ExpiredTime: 0, // 长期有效
			Vip:         false,
		}
		u.Nickname = fmt.Sprintf("用户@%d", utils.RandomNumber(6))

		if err = h.DB.Create(&u).Error; err != nil {
			failedCount++
			errorsList = append(errorsList, rowError{Row: line, Error: "写入数据库失败: " + err.Error()})
			continue
		}

		successCount++
	}

	resp.SUCCESS(c, gin.H{
		"success": successCount,
		"failed":  failedCount,
		"errors":  errorsList,
	})
}

// List 用户列表
func (h *UserHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	username := h.GetTrim(c, "username")
	mobile := h.GetTrim(c, "mobile")
	email := h.GetTrim(c, "email")

	offset := (page - 1) * pageSize
	var items []model.User
	var users = make([]vo.User, 0)
	var total int64

	session := h.DB.Session(&gorm.Session{})
	if username != "" {
		session = session.Where("username LIKE ?", "%"+username+"%")
	}
	if mobile != "" {
		session = session.Where("mobile LIKE ?", "%"+mobile+"%")
	}
	if email != "" {
		session = session.Where("email LIKE ?", "%"+email+"%")
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
		Id          uint   `json:"id"`
		Password    string `json:"password"`
		Username    string `json:"username"`
		Mobile      string `json:"mobile"`
		Email       string `json:"email"`
		ChatModels  []int  `json:"chat_models"`
		ExpiredTime string `json:"expired_time"`
		Status      bool   `json:"status"`
		Vip         bool   `json:"vip"`
		Power       int    `json:"power"`
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
		user.Email = data.Email
		user.Mobile = data.Mobile
		user.Status = data.Status
		user.Vip = data.Vip
		user.Power = data.Power
		user.ChatModels = utils.JsonEncode(data.ChatModels)
		user.ExpiredTime = utils.Str2stamp(data.ExpiredTime)

		res = h.DB.Select("username", "mobile", "email", "status", "vip", "power", "chat_models_json", "expired_time").Updates(&user)

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
			ChatConfig:  "{}",
			ChatModels:  utils.JsonEncode(data.ChatModels),
			ExpiredTime: utils.Str2stamp(data.ExpiredTime),
		}
		u.Nickname = fmt.Sprintf("用户@%d", utils.RandomNumber(6))
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
		if err = tx.Where("user_id = ?", id).Delete(&model.ImageJob{}).Error; err != nil {
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

// GenLoginLink 生成登录链接
func (h *UserHandler) GenLoginLink(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var user model.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		resp.ERROR(c, "用户不存在")
		return
	}

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
	sessionKey := fmt.Sprintf("users/%d", user.Id)
	if _, err = h.redis.Set(c, sessionKey, tokenString, 0).Result(); err != nil {
		resp.ERROR(c, "error with save token: "+err.Error())
		return
	}

	resp.SUCCESS(c, tokenString)
}
