package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"crypto/sha1"
	"encoding/hex"
	"geekai/core"
	"geekai/utils/resp"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type WxGzhHandler struct {
	BaseHandler
}

func NewWxGzhHandler(server *core.AppServer) *WxGzhHandler {
	return &WxGzhHandler{
		BaseHandler: BaseHandler{
			App: server,
		},
	}
}

func (h *WxGzhHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/wx/")
	group.GET("verify", h.WechatVerify)
}

// 处理微信服务器验证请求
func (h *WxGzhHandler) WechatVerify(c *gin.Context) {
	logger.Info("WechatVerify")
	// 只处理 GET 请求
	if c.Request.Method != "GET" {
		resp.ERROR(c, "Method Not Allowed")
		return
	}

	// 解析 URL 参数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 验证参数完整性
	if signature == "" || timestamp == "" || nonce == "" || echostr == "" {
		log.Println("Missing parameters")
		resp.ERROR(c, "Missing parameters")
		return
	}

	// 验证签名
	if validateSignature(signature, h.App.SysConfig.WxGzh.Token, timestamp, nonce) {
		// 验证成功，返回 echostr（必须是纯文本）
		c.String(http.StatusOK, echostr)
		log.Println("Token verification success")
	} else {
		// 验证失败
		resp.ERROR(c, "Forbidden: Invalid signature")
		log.Println("Token verification failed")
	}
}

func validateSignature(signature, token, timestamp, nonce string) bool {
	// 1. 将 token、timestamp、nonce 按字典序排序
	strs := []string{token, timestamp, nonce}
	sort.Strings(strs)

	// 2. 拼接字符串
	joined := strings.Join(strs, "")

	// 3. 计算 SHA1 哈希
	hash := sha1.New()
	hash.Write([]byte(joined))
	hashed := hex.EncodeToString(hash.Sum(nil))

	// 4. 与 signature 比对
	return hashed == signature
}

// 创建微信菜单
func (h *WxGzhHandler) CreateMenu(c *gin.Context) {

	resp.SUCCESS(c, "创建菜单成功")
}
