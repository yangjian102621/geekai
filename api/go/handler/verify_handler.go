package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/store"
	"chatplus/utils"
	"chatplus/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
)

// 生成验证的控制器

type VerifyHandler struct {
	BaseHandler
	sms *service.AliYunSmsService
	db  *store.LevelDB
}

const TokenStorePrefix = "/tokens/"
const CodeStorePrefix = "/codes/"

func NewVerifyHandler(app *core.AppServer, sms *service.AliYunSmsService, db *store.LevelDB) *VerifyHandler {
	handler := &VerifyHandler{sms: sms, db: db}
	handler.App = app
	return handler
}

type VerifyToken struct {
	Token     string
	Timestamp int64
}

// Token 生成自验证 token
func (h *VerifyHandler) Token(c *gin.Context) {
	// 确保是通过浏览器访问
	if c.GetHeader("Sec-Fetch-Mode") != "cors" {
		resp.HACKER(c)
		return
	}

	token := VerifyToken{
		Token:     utils.RandString(32),
		Timestamp: time.Now().Unix(),
	}
	json := utils.JsonEncode(token)
	encrypt, err := utils.AesEncrypt(h.App.Config.AesEncryptKey, []byte(json))
	if err != nil {
		resp.ERROR(c, "Token 加密出错")
		return
	}
	err = h.db.Put(TokenStorePrefix+token.Token, token)
	if err != nil {
		resp.ERROR(c, "Token 存储失败")
		return
	}

	resp.SUCCESS(c, encrypt)
}

// SendMsg 发送验证码短信
func (h *VerifyHandler) SendMsg(c *gin.Context) {
	var data struct {
		Mobile string `json:"mobile"`
		Token  string `json:"token"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	decrypt, err := utils.AesDecrypt(h.App.Config.AesEncryptKey, data.Token)
	if err != nil {
		resp.ERROR(c, "Token 解密失败")
		return
	}

	var token VerifyToken
	err = utils.JsonDecode(string(decrypt), &token)
	if err != nil {
		resp.ERROR(c, "Token 解码失败")
		return
	}

	_, err = h.db.Get(TokenStorePrefix + token.Token)
	if err != nil {
		resp.HACKER(c)
		return
	}

	if time.Now().Unix()-token.Timestamp > 30 {
		resp.ERROR(c, "Token 已过期，请刷新页面重试")
		return
	}

	code := utils.RandomNumber(6)
	err = h.sms.SendVerifyCode(data.Mobile, code)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 每个 token 用完一次立即失效
	_ = h.db.Delete(TokenStorePrefix + token.Token)
	// 存储验证码，等待后面注册验证
	err = h.db.Put(CodeStorePrefix+data.Mobile, code)
	if err != nil {
		resp.ERROR(c, "验证码保存失败")
		return
	}

	resp.SUCCESS(c)
}
