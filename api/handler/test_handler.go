package handler

import (
	"chatplus/service"
	"chatplus/service/payment"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

type TestHandler struct {
	db        *gorm.DB
	snowflake *service.Snowflake
	js        *payment.PayJS
}

func NewTestHandler(db *gorm.DB, snowflake *service.Snowflake, js *payment.PayJS) *TestHandler {
	return &TestHandler{db: db, snowflake: snowflake, js: js}
}

type reqBody struct {
	BotType       string        `json:"botType"`
	Prompt        string        `json:"prompt"`
	Base64Array   []interface{} `json:"base64Array,omitempty"`
	AccountFilter struct {
		InstanceId          string        `json:"instanceId"`
		Modes               []interface{} `json:"modes"`
		Remix               bool          `json:"remix"`
		RemixAutoConsidered bool          `json:"remixAutoConsidered"`
	} `json:"accountFilter,omitempty"`
	NotifyHook string `json:"notifyHook"`
	State      string `json:"state,omitempty"`
}

type resBody struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Properties  struct {
	} `json:"properties"`
	Result string `json:"result"`
}

func (h *TestHandler) Test(c *gin.Context) {
	query(c)

}

func upscale(c *gin.Context) {
	apiURL := "https://api.openai1s.cn/mj/submit/action"
	token := "sk-QpBaQn9Z5vngsjJaFdDfC9Db90C845EaB5E764578a7d292a"
	body := map[string]string{
		"customId":   "MJ::JOB::upsample::1::c80a8eb1-f2d1-4f40-8785-97eb99b7ba0a",
		"taskId":     "1704880156226095",
		"notifyHook": "http://r9it.com:6004/api/test/mj",
	}
	var res resBody
	var resErr errRes
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&resErr).
		Post(apiURL)
	if err != nil {
		resp.ERROR(c, "请求出错："+err.Error())
		return
	}

	if r.IsErrorState() {
		resp.ERROR(c, "返回错误状态："+resErr.Error.Message)
		return
	}

	resp.SUCCESS(c, res)

}

type queryRes struct {
	Action  string `json:"action"`
	Buttons []struct {
		CustomId string `json:"customId"`
		Emoji    string `json:"emoji"`
		Label    string `json:"label"`
		Style    int    `json:"style"`
		Type     int    `json:"type"`
	} `json:"buttons"`
	Description string `json:"description"`
	FailReason  string `json:"failReason"`
	FinishTime  int    `json:"finishTime"`
	Id          string `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	Progress    string `json:"progress"`
	Prompt      string `json:"prompt"`
	PromptEn    string `json:"promptEn"`
	Properties  struct {
	} `json:"properties"`
	StartTime  int    `json:"startTime"`
	State      string `json:"state"`
	Status     string `json:"status"`
	SubmitTime int    `json:"submitTime"`
}

func query(c *gin.Context) {
	apiURL := "https://api.openai1s.cn/mj/task/1704960661008372/fetch"
	token := "sk-QpBaQn9Z5vngsjJaFdDfC9Db90C845EaB5E764578a7d292a"
	var res queryRes
	r, err := req.C().R().SetHeader("Authorization", "Bearer "+token).
		SetSuccessResult(&res).
		Get(apiURL)

	if err != nil {
		resp.ERROR(c, "请求出错："+err.Error())
		return
	}

	if r.IsErrorState() {
		resp.ERROR(c, "返回错误状态："+r.Status)
		return
	}

	resp.SUCCESS(c, res)
}

type errRes struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func image(c *gin.Context) {
	apiURL := "https://api.openai1s.cn/mj-fast/mj/submit/imagine"
	token := "sk-QpBaQn9Z5vngsjJaFdDfC9Db90C845EaB5E764578a7d292a"
	body := reqBody{
		BotType:    "MID_JOURNEY",
		Prompt:     "一个中国美女，手上拿着一桶爆米花，脸上带着迷人的微笑，白色衣服 --s 750 --v 6",
		NotifyHook: "http://r9it.com:6004/api/test/mj",
	}
	var res resBody
	var resErr errRes
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&resErr).
		Post(apiURL)
	if err != nil {
		resp.ERROR(c, "请求出错："+err.Error())
		return
	}

	if r.IsErrorState() {
		resp.ERROR(c, "返回错误状态："+resErr.Error.Message)
		return
	}

	resp.SUCCESS(c, res)
}

type cbReq struct {
	Id          string      `json:"id"`
	Action      string      `json:"action"`
	Status      string      `json:"status"`
	Prompt      string      `json:"prompt"`
	PromptEn    string      `json:"promptEn"`
	Description string      `json:"description"`
	SubmitTime  int64       `json:"submitTime"`
	StartTime   int64       `json:"startTime"`
	FinishTime  int64       `json:"finishTime"`
	Progress    string      `json:"progress"`
	ImageUrl    string      `json:"imageUrl"`
	FailReason  interface{} `json:"failReason"`
	Properties  struct {
		FinalPrompt string `json:"finalPrompt"`
	} `json:"properties"`
}

func (h *TestHandler) Mj(c *gin.Context) {
	var data cbReq
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error(err)
	}
	logger.Debugf("任务ID：%s,任务进度：%s,图片地址：%s, 最终提示词：%s", data.Id, data.Progress, data.ImageUrl, data.Properties.FinalPrompt)
	apiURL := "https://api.openai1s.cn/mj/task/" + data.Id + "/fetch"
	token := "sk-QpBaQn9Z5vngsjJaFdDfC9Db90C845EaB5E764578a7d292a"
	var res queryRes
	_, _ = req.C().R().SetHeader("Authorization", "Bearer "+token).
		SetSuccessResult(&res).
		Get(apiURL)

	fmt.Println(res.State, ",", res.ImageUrl, ",", res.Progress)
}

func (h *TestHandler) initUserNickname(c *gin.Context) {
	var users []model.User
	tx := h.db.Find(&users)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	for _, u := range users {
		u.Nickname = fmt.Sprintf("极客学长@%d", utils.RandomNumber(6))
		h.db.Updates(&u)
	}

	resp.SUCCESS(c)
}

func (h *TestHandler) initMjTaskId(c *gin.Context) {
	var jobs []model.MidJourneyJob
	tx := h.db.Find(&jobs)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	for _, job := range jobs {
		id, _ := h.snowflake.Next(true)
		job.TaskId = id
		h.db.Updates(&job)
	}

	resp.SUCCESS(c)
}
