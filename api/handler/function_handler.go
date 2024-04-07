package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
	"strings"
	"time"
)

type FunctionHandler struct {
	BaseHandler
	config        types.ChatPlusApiConfig
	uploadManager *oss.UploaderManager
}

func NewFunctionHandler(server *core.AppServer, db *gorm.DB, config *types.AppConfig, manager *oss.UploaderManager) *FunctionHandler {
	return &FunctionHandler{
		BaseHandler: BaseHandler{
			App: server,
			DB:  db,
		},
		config:        config.ApiConfig,
		uploadManager: manager,
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
	if err != nil || r.IsErrorState() {
		resp.ERROR(c, fmt.Sprintf("%v%v", err, r.Err))
		return
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

type imgReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type imgRes struct {
	Created int64 `json:"created"`
	Data    []struct {
		RevisedPrompt string `json:"revised_prompt"`
		Url           string `json:"url"`
	} `json:"data"`
}

type ErrRes struct {
	Error struct {
		Code    interface{} `json:"code"`
		Message string      `json:"message"`
		Param   interface{} `json:"param"`
		Type    string      `json:"type"`
	} `json:"error"`
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
	tx := h.DB.Where("id = ?", params["user_id"]).First(&user)
	if tx.Error != nil {
		resp.ERROR(c, "当前用户不存在！")
		return
	}

	if user.Power < h.App.SysConfig.DallPower {
		resp.ERROR(c, "当前用户剩余算力不足以完成本次绘画！")
		return
	}

	prompt := utils.InterfaceToString(params["prompt"])
	// get image generation API KEY
	var apiKey model.ApiKey
	tx = h.DB.Where("platform = ?", types.OpenAI).Where("type = ?", "img").Where("enabled = ?", true).Order("last_used_at ASC").First(&apiKey)
	if tx.Error != nil {
		resp.ERROR(c, "获取绘图 API KEY 失败: "+tx.Error.Error())
		return
	}

	// translate prompt
	const translatePromptTemplate = "Translate the following painting prompt words into English keyword phrases. Without any explanation, directly output the keyword phrases separated by commas. The content to be translated is: [%s]"
	pt, err := utils.OpenAIRequest(h.DB, fmt.Sprintf(translatePromptTemplate, params["prompt"]))
	if err == nil {
		logger.Debugf("翻译绘画提示词，原文：%s，译文：%s", prompt, pt)
		prompt = pt
	}
	var res imgRes
	var errRes ErrRes
	var request *req.Request
	if len(apiKey.ProxyURL) > 5 {
		request = req.C().SetProxyURL(apiKey.ProxyURL).R()
	} else {
		request = req.C().R()
	}
	logger.Debugf("Sending %s request, ApiURL:%s, API KEY:%s, PROXY: %s", apiKey.Platform, apiKey.ApiURL, apiKey.Value, apiKey.ProxyURL)
	r, err := request.SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(imgReq{
			Model:  "dall-e-3",
			Prompt: prompt,
			N:      1,
			Size:   "1024x1024",
		}).
		SetErrorResult(&errRes).
		SetSuccessResult(&res).Post(apiKey.ApiURL)
	if r.IsErrorState() {
		resp.ERROR(c, "请求 OpenAI API 失败: "+errRes.Error.Message)
		return
	}
	// 更新 API KEY 的最后使用时间
	h.DB.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())
	logger.Debugf("%+v", res)
	// 存储图片
	imgURL, err := h.uploadManager.GetUploadHandler().PutImg(res.Data[0].Url, false)
	if err != nil {
		resp.ERROR(c, "下载图片失败: "+err.Error())
		return
	}

	content := fmt.Sprintf("下面是根据您的描述创作的图片，它描绘了 【%s】 的场景。 \n\n![](%s)\n", prompt, imgURL)
	// 更新用户算力
	tx = h.DB.Model(&model.User{}).Where("id", user.Id).UpdateColumn("power", gorm.Expr("power - ?", h.App.SysConfig.DallPower))
	// 记录算力变化日志
	if tx.Error == nil && tx.RowsAffected > 0 {
		var u model.User
		h.DB.Where("id", user.Id).First(&u)
		h.DB.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      types.PowerConsume,
			Amount:    h.App.SysConfig.DallPower,
			Balance:   u.Power,
			Mark:      types.PowerSub,
			Model:     "dall-e-3",
			Remark:    fmt.Sprintf("绘画提示词：%s", utils.CutWords(prompt, 10)),
			CreatedAt: time.Now(),
		})
	}

	resp.SUCCESS(c, content)
}
