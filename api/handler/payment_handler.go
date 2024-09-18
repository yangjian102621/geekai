package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"embed"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/payment"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/shopspring/decimal"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PayWay struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PaymentHandler 支付服务回调 handler
type PaymentHandler struct {
	BaseHandler
	alipayService    *payment.AlipayService
	huPiPayService   *payment.HuPiPayService
	geekPayService   *payment.GeekPayService
	wechatPayService *payment.WechatPayService
	snowflake        *service.Snowflake
	fs               embed.FS
	lock             sync.Mutex
	signKey          string // 用来签名的随机秘钥
}

func NewPaymentHandler(
	server *core.AppServer,
	alipayService *payment.AlipayService,
	huPiPayService *payment.HuPiPayService,
	geekPayService *payment.GeekPayService,
	wechatPayService *payment.WechatPayService,
	db *gorm.DB,
	snowflake *service.Snowflake,
	fs embed.FS) *PaymentHandler {
	return &PaymentHandler{
		alipayService:    alipayService,
		huPiPayService:   huPiPayService,
		geekPayService:   geekPayService,
		wechatPayService: wechatPayService,
		snowflake:        snowflake,
		fs:               fs,
		lock:             sync.Mutex{},
		BaseHandler: BaseHandler{
			App: server,
			DB:  db,
		},
		signKey: utils.RandString(32),
	}
}

func (h *PaymentHandler) Pay(c *gin.Context) {
	payWay := c.Query("pay_way")
	payType := c.Query("pay_type")
	productId := c.Query("pid")
	device := c.Query("device")
	userId := c.Query("user_id")

	var product model.Product
	err := h.DB.First(&product, productId).Error
	if err != nil {
		resp.ERROR(c, "Product not found")
		return
	}

	orderNo, err := h.snowflake.Next(false)
	if err != nil {
		resp.ERROR(c, "error with generate trade no: "+err.Error())
		return
	}
	var user model.User
	err = h.DB.Where("id", userId).First(&user).Error
	if err != nil {
		resp.NotAuth(c)
		return
	}
	// 创建订单
	remark := types.OrderRemark{
		Days:     product.Days,
		Power:    product.Power,
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
	}

	amount, _ := decimal.NewFromFloat(product.Price).Sub(decimal.NewFromFloat(product.Discount)).Float64()

	order := model.Order{
		UserId:    user.Id,
		Username:  user.Username,
		ProductId: product.Id,
		OrderNo:   orderNo,
		Subject:   product.Name,
		Amount:    amount,
		Status:    types.OrderNotPaid,
		PayWay:    payWay,
		PayType:   payType,
		Remark:    utils.JsonEncode(remark),
	}
	err = h.DB.Create(&order).Error
	if err != nil {
		resp.ERROR(c, "error with create order: "+err.Error())
		return
	}

	var payURL string
	if payWay == "alipay" { // 支付宝
		money := fmt.Sprintf("%.2f", order.Amount)
		notifyURL := fmt.Sprintf("%s/api/payment/notify/alipay", h.App.Config.AlipayConfig.NotifyHost)
		returnURL := fmt.Sprintf("%s/member", h.App.Config.AlipayConfig.NotifyHost)
		if device == "mobile" {
			payURL, err = h.alipayService.PayMobile(payment.AlipayParams{
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   money,
				ReturnURL:  returnURL,
				NotifyURL:  notifyURL,
			})
		} else {
			payURL, err = h.alipayService.PayPC(payment.AlipayParams{
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   money,
				ReturnURL:  returnURL,
				NotifyURL:  notifyURL,
			})
		}
		if err != nil {
			resp.ERROR(c, "error with generate pay url: "+err.Error())
			return
		}
	} else if order.PayWay == "hupi" { // 虎皮椒支付
		params := payment.HuPiPayReq{
			Version:      "1.1",
			TradeOrderId: orderNo,
			TotalFee:     fmt.Sprintf("%f", order.Amount),
			Title:        order.Subject,
			NotifyURL:    fmt.Sprintf("%s/api/payment/notify/hupi", h.App.Config.HuPiPayConfig.NotifyHost),
			WapName:      "GeekAI助手",
		}
		r, err := h.huPiPayService.Pay(params)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}

		c.Redirect(302, r.URL)
	} else if order.PayWay == "wechat" {
		//uri, err := h.wechatPayService.PayUrlNative(order.OrderNo, int(order.Amount*100), order.Subject)
		uri, err := h.wechatPayService.PayUrlNative(payment.WechatPayParams{})
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}

		c.Redirect(302, uri)
	} else if order.PayWay == "geek" {
		notifyURL := fmt.Sprintf("%s/api/payment/notify/geek", h.App.Config.GeekPayConfig.NotifyHost)
		returnURL := fmt.Sprintf("%s/member", h.App.Config.GeekPayConfig.NotifyHost)
		params := payment.GeekPayParams{
			OutTradeNo: orderNo,
			Method:     "web",
			Name:       order.Subject,
			Money:      fmt.Sprintf("%f", order.Amount),
			ClientIP:   c.ClientIP(),
			Device:     device,
			Type:       payType,
			ReturnURL:  returnURL,
			NotifyURL:  notifyURL,
		}

		res, err := h.geekPayService.Pay(params)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
		payURL = res.PayURL
	}
	resp.SUCCESS(c, payURL)
}

// 异步通知回调公共逻辑
func (h *PaymentHandler) notify(orderNo string, tradeNo string) error {
	var order model.Order
	res := h.DB.Where("order_no = ?", orderNo).First(&order)
	if res.Error != nil {
		err := fmt.Errorf("error with fetch order: %v", res.Error)
		logger.Error(err)
		return err
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	// 已支付订单，直接返回
	if order.Status == types.OrderPaidSuccess {
		return nil
	}

	var user model.User
	res = h.DB.First(&user, order.UserId)
	if res.Error != nil {
		err := fmt.Errorf("error with fetch user info: %v", res.Error)
		logger.Error(err)
		return err
	}

	var remark types.OrderRemark
	err := utils.JsonDecode(order.Remark, &remark)
	if err != nil {
		err := fmt.Errorf("error with decode order remark: %v", err)
		logger.Error(err)
		return err
	}

	var opt string
	var power int
	if remark.Days > 0 { // VIP 充值
		if user.ExpiredTime >= time.Now().Unix() {
			user.ExpiredTime = time.Unix(user.ExpiredTime, 0).AddDate(0, 0, remark.Days).Unix()
			opt = "VIP充值，VIP 没到期，只延期不增加算力"
		} else {
			user.ExpiredTime = time.Now().AddDate(0, 0, remark.Days).Unix()
			user.Power += h.App.SysConfig.VipMonthPower
			power = h.App.SysConfig.VipMonthPower
			opt = "VIP充值"
		}
		user.Vip = true
	} else { // 充值点卡，直接增加次数即可
		user.Power += remark.Power
		opt = "点卡充值"
		power = remark.Power
	}

	// 更新用户信息
	res = h.DB.Updates(&user)
	if res.Error != nil {
		err := fmt.Errorf("error with update user info: %v", res.Error)
		logger.Error(err)
		return err
	}

	// 更新订单状态
	order.PayTime = time.Now().Unix()
	order.Status = types.OrderPaidSuccess
	order.TradeNo = tradeNo
	res = h.DB.Updates(&order)
	if res.Error != nil {
		err := fmt.Errorf("error with update order info: %v", res.Error)
		logger.Error(err)
		return err
	}

	// 更新产品销量
	h.DB.Model(&model.Product{}).Where("id = ?", order.ProductId).UpdateColumn("sales", gorm.Expr("sales + ?", 1))

	// 记录算力充值日志
	if power > 0 {
		h.DB.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      types.PowerRecharge,
			Amount:    power,
			Balance:   user.Power,
			Mark:      types.PowerAdd,
			Model:     order.PayWay,
			Remark:    fmt.Sprintf("%s，金额：%f，订单号：%s", opt, order.Amount, order.OrderNo),
			CreatedAt: time.Now(),
		})
	}

	return nil
}

// GetPayWays 获取支付方式
func (h *PaymentHandler) GetPayWays(c *gin.Context) {
	payWays := make([]gin.H, 0)
	if h.App.Config.AlipayConfig.Enabled {
		payWays = append(payWays, gin.H{"pay_way": "alipay", "pay_type": "alipay"})
	}
	if h.App.Config.HuPiPayConfig.Enabled {
		payWays = append(payWays, gin.H{"pay_way": "hupi", "pay_type": h.App.Config.HuPiPayConfig.Name})
	}
	if h.App.Config.GeekPayConfig.Enabled {
		payWays = append(payWays, gin.H{"pay_way": "geek", "pay_type": "alipay"})
		payWays = append(payWays, gin.H{"pay_way": "geek", "pay_type": "wxpay"})
		payWays = append(payWays, gin.H{"pay_way": "geek", "pay_type": "qqpay"})
		payWays = append(payWays, gin.H{"pay_way": "geek", "pay_type": "jdpay"})
		payWays = append(payWays, gin.H{"pay_way": "geek", "pay_type": "douyin"})
		payWays = append(payWays, gin.H{"pay_way": "geek", "pay_type": "paypal"})
	}
	if h.App.Config.WechatPayConfig.Enabled {
		payWays = append(payWays, gin.H{"pay_way": "wechat", "pay_type": "wxpay"})
	}
	resp.SUCCESS(c, payWays)
}

// HuPiPayNotify 虎皮椒支付异步回调
func (h *PaymentHandler) HuPiPayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	orderNo := c.Request.Form.Get("trade_order_id")
	tradeNo := c.Request.Form.Get("open_order_id")
	logger.Infof("收到虎皮椒订单支付回调，订单 NO：%s，交易流水号：%s", orderNo, tradeNo)

	if err = h.huPiPayService.Check(tradeNo); err != nil {
		logger.Error("订单校验失败：", err)
		c.String(http.StatusOK, "fail")
		return
	}
	err = h.notify(orderNo, tradeNo)
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// AlipayNotify 支付宝支付回调
func (h *PaymentHandler) AlipayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	// TODO：验证交易签名
	res := h.alipayService.TradeVerify(c.Request)
	logger.Infof("验证支付结果：%+v", res)
	if !res.Success() {
		logger.Error("订单校验失败：", res.Message)
		c.String(http.StatusOK, "fail")
		return
	}

	tradeNo := c.Request.Form.Get("trade_no")
	err = h.notify(res.OutTradeNo, tradeNo)
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// GeekPayNotify 支付异步回调
func (h *PaymentHandler) GeekPayNotify(c *gin.Context) {
	var params = make(map[string]string)
	for k := range c.Request.URL.Query() {
		params[k] = c.Query(k)
	}

	logger.Infof("收到GeekPay订单支付回调：%+v", params)
	// 检查支付状态
	if params["trade_status"] != "TRADE_SUCCESS" {
		c.String(http.StatusOK, "success")
		return
	}

	sign := h.geekPayService.Sign(params)
	if sign != c.Query("sign") {
		logger.Errorf("签名验证失败, %s, %s", sign, c.Query("sign"))
		c.String(http.StatusOK, "fail")
		return
	}

	err := h.notify(params["out_trade_no"], params["trade_no"])
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// WechatPayNotify 微信商户支付异步回调
func (h *PaymentHandler) WechatPayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	result := h.wechatPayService.TradeVerify(c.Request)
	if !result.Success() {
		logger.Error("订单校验失败：", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FAIL",
			"message": err.Error(),
		})
		return
	}

	err = h.notify(result.OutTradeNo, result.TradeId)
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}
