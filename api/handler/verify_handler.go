package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/store"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 短信验证码控制器

type VerifyHandler struct {
	BaseHandler
	sms *service.AliYunSmsService
	db  *store.LevelDB
}

const TokenStorePrefix = "/verify/tokens/"
const CodeStorePrefix = "/verify/codes/"
const MobileStatPrefix = "/verify/stats/"

func NewVerifyHandler(app *core.AppServer, sms *service.AliYunSmsService, db *store.LevelDB) *VerifyHandler {
	handler := &VerifyHandler{sms: sms, db: db}
	handler.App = app
	return handler
}

type VerifyToken struct {
	Token     string
	Timestamp int64
}

// CodeStats 验证码发送统计
type CodeStats struct {
	Mobile string
	Count  uint
	Time   int64
}

// Token 生成自验证 token
func (h *VerifyHandler) Token(c *gin.Context) {
	// 如果不是通过浏览器访问，则返回错误的 token
	// TODO: 引入验证码机制防刷机制
	if c.GetHeader("Sec-Fetch-Mode") != "cors" {
		token := fmt.Sprintf("%s:%d", utils.RandString(32), time.Now().Unix())
		encrypt, err := utils.AesEncrypt(h.App.Config.AesEncryptKey, []byte(token))
		if err != nil {
			resp.ERROR(c, "Token 加密出错")
			return
		}
		resp.SUCCESS(c, encrypt)
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

	if time.Now().Unix()-token.Timestamp > 30 {
		resp.ERROR(c, "Token 已过期，请刷新页面重试")
		return
	}

	// 验证当前手机号发送次数，24 小时内相同手机号只允许发送 2 次
	var stat CodeStats
	err = h.db.Get(MobileStatPrefix+data.Mobile, &stat)
	if err != nil {
		stat = CodeStats{
			Mobile: data.Mobile,
			Count:  0,
			Time:   time.Now().Unix(),
		}
	} else if stat.Count == 2 {
		if time.Now().Unix()-stat.Time > 86400 {
			stat.Count = 0
			stat.Time = time.Now().Unix()
		} else {
			resp.ERROR(c, "触发流量预警，请 24 小时后再操作！")
			return
		}
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

	// 更新发送次数
	stat.Count = stat.Count + 1
	_ = h.db.Put(MobileStatPrefix+data.Mobile, stat)
	logger.Infof("%+v", stat)

	resp.SUCCESS(c)
}
