package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/service/dalle"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

type FunctionHandler struct {
	BaseHandler
	config        types.ApiConfig
	uploadManager *oss.UploaderManager
	dallService   *dalle.Service
}

func NewFunctionHandler(
	server *core.AppServer,
	db *gorm.DB,
	config *types.AppConfig,
	manager *oss.UploaderManager,
	dallService *dalle.Service) *FunctionHandler {
	return &FunctionHandler{
		BaseHandler: BaseHandler{
			App: server,
			DB:  db,
		},
		config:        config.ApiConfig,
		uploadManager: manager,
		dallService:   dallService,
	}
}

type resVo struct {
	Code    types.BizCode `json:"code"`
	Message string        `json:"message"`
	Data    struct {
		Title     string     `json:"title"`
		UpdatedAt string     `json:"updated_at"`
		Items     []dataItem `json:"items"`
	} `json:"data"`
}

type dataItem struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Remark string `json:"remark"`
}

// check authorization
func (h *FunctionHandler) checkAuth(c *gin.Context) error {
	tokenString := c.GetHeader(types.UserAuthHeader)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(h.App.Config.Session.SecretKey), nil
	})

	if err != nil {
		return fmt.Errorf("error with parse auth token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("token is invalid")
	}

	expr := utils.IntValue(utils.InterfaceToString(claims["expired"]), 0)
	if expr > 0 && int64(expr) < time.Now().Unix() {
		return errors.New("token is expired")
	}

	return nil
}

// WeiBo 微博热搜
func (h *FunctionHandler) WeiBo(c *gin.Context) {
	if err := h.checkAuth(c); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if h.config.Token == "" {
		resp.ERROR(c, "无效的 API Token")
		return
	}

	url := fmt.Sprintf("%s/api/weibo/fetch", h.config.ApiURL)
	var res resVo
	r, err := req.C().R().
		SetHeader("AppId", h.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", h.config.Token)).
		SetSuccessResult(&res).Get(url)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("%v", err))
		return
	}
	if r.IsErrorState() {
		resp.ERROR(c, fmt.Sprintf("error http code status: %v", r.Status))
	}

	if res.Code != types.Success {
		resp.ERROR(c, res.Message)
		return
	}

	builder := make([]string, 0)
	builder = append(builder, fmt.Sprintf("**%s**，最新更新：%s", res.Data.Title, res.Data.UpdatedAt))
	for i, v := range res.Data.Items {
		builder = append(builder, fmt.Sprintf("%d、 [%s](%s) [热度：%s]", i+1, v.Title, v.Url, v.Remark))
	}
	resp.SUCCESS(c, strings.Join(builder, "\n\n"))
}

// ZaoBao 今日早报
func (h *FunctionHandler) ZaoBao(c *gin.Context) {
	if err := h.checkAuth(c); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if h.config.Token == "" {
		resp.ERROR(c, "无效的 API Token")
		return
	}

	url := fmt.Sprintf("%s/api/zaobao/fetch", h.config.ApiURL)
	var res resVo
	r, err := req.C().R().
		SetHeader("AppId", h.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", h.config.Token)).
		SetSuccessResult(&res).Get(url)
	if err != nil || r.IsErrorState() {
		resp.ERROR(c, fmt.Sprintf("%v%v", err, r.Err))
		return
	}

	if res.Code != types.Success {
		resp.ERROR(c, res.Message)
		return
	}

	builder := make([]string, 0)
	builder = append(builder, fmt.Sprintf("**%s 早报：**", res.Data.UpdatedAt))
	for _, v := range res.Data.Items {
		builder = append(builder, v.Title)
	}
	builder = append(builder, fmt.Sprintf("%s", res.Data.Title))
	resp.SUCCESS(c, strings.Join(builder, "\n\n"))
}

// Dall3 DallE3 AI 绘图
func (h *FunctionHandler) Dall3(c *gin.Context) {
	if err := h.checkAuth(c); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	logger.Debugf("绘画参数：%+v", params)
	var user model.User
	res := h.DB.Where("id = ?", params["user_id"]).First(&user)
	if res.Error != nil {
		resp.ERROR(c, "当前用户不存在！")
		return
	}

	if user.Power < h.App.SysConfig.DallPower {
		resp.ERROR(c, "创建 DALL-E 绘图任务失败，算力不足")
		return
	}

	// create dall task
	prompt := utils.InterfaceToString(params["prompt"])
	job := model.DallJob{
		UserId: user.Id,
		Prompt: prompt,
		Power:  h.App.SysConfig.DallPower,
	}
	res = h.DB.Create(&job)

	if res.Error != nil {
		resp.ERROR(c, "创建 DALL-E 绘图任务失败："+res.Error.Error())
		return
	}

	content, err := h.dallService.Image(types.DallTask{
		Id:      job.Id,
		UserId:  user.Id,
		Prompt:  job.Prompt,
		N:       1,
		Quality: "standard",
		Size:    "1024x1024",
		Style:   "vivid",
		Power:   job.Power,
	}, true)
	if err != nil {
		resp.ERROR(c, "任务执行失败："+err.Error())
		return
	}

	resp.SUCCESS(c, content)
}

// List 获取所有的工具函数列表
func (h *FunctionHandler) List(c *gin.Context) {
	var items []model.Function
	err := h.DB.Where("enabled", true).Find(&items).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	tools := make([]vo.Function, 0)
	for _, v := range items {
		var f vo.Function
		err = utils.CopyObject(v, &f)
		if err != nil {
			continue
		}
		f.Action = ""
		f.Token = ""
		tools = append(tools, f)
	}

	resp.SUCCESS(c, tools)
}
