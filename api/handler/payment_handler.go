package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"embed"
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/payment"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"math"
	"net/http"
	"strconv"
	"strings"
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
	alipayService *payment.AlipayService
	epayService   *payment.EPayService
	wxpayService  *payment.WxPayService
	stripeService *payment.StripeService
	snowflake     *service.Snowflake
	userService   *service.UserService
	ledgerService *service.LedgerService
	fs            embed.FS
	lock          sync.Mutex
	config        *types.PaymentConfig
}

func NewPaymentHandler(
	server *core.AppServer,
	alipayService *payment.AlipayService,
	geekPayService *payment.EPayService,
	wxpayService *payment.WxPayService,
	stripeService *payment.StripeService,
	db *gorm.DB,
	userService *service.UserService,
	ledgerService *service.LedgerService,
	snowflake *service.Snowflake,
	fs embed.FS,
	sysConfig *types.SystemConfig) *PaymentHandler {
	return &PaymentHandler{
		alipayService: alipayService,
		epayService:   geekPayService,
		wxpayService:  wxpayService,
		stripeService: stripeService,
		snowflake:     snowflake,
		userService:   userService,
		ledgerService: ledgerService,
		fs:            fs,
		lock:          sync.Mutex{},
		BaseHandler: BaseHandler{
			App: server,
			DB:  db,
		},
		config: &sysConfig.Payment,
	}
}

// RegisterRoutes 注册路由
func (h *PaymentHandler) RegisterRoutes(r *core.Router) {
	// 支付回调接口（公开）：由第三方网关服务器发起，天然不带用户 token。
	// 误收进 UserGroup 的直接后果是"用户已付款、回调被 401、积分永不到账"。
	public := r.PublicGroup("/api/payment/")
	public.POST("notify/alipay", h.AlipayNotify)
	public.GET("notify/epay", h.EPayNotify)
	public.POST("notify/wxpay", h.WxpayNotify)
	public.POST("notify/stripe", h.StripeNotify)

	// 需要用户登录的接口
	group := r.UserGroup("/api/payment/")
	{
		group.POST("create", h.CreateOrder)
	}
}

func (h *PaymentHandler) StartSyncOrders() {
	go func() {
		for {
			err := h.SyncOrders()
			if err != nil {
				logger.Error(err)
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

// SyncOrders 同步订单状态
func (h *PaymentHandler) SyncOrders() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("同步订单状态发生异常: %v", err)
		}
	}()
	var orders []model.Order
	err := h.DB.Where("status", types.OrderNotPaid).Where("checked", false).Find(&orders).Error
	if err != nil {
		return err
	}

	for _, order := range orders {
		time.Sleep(time.Second * 1)
		//超时15分钟的订单，直接标记为已关闭
		if time.Now().After(order.CreatedAt.Add(time.Minute * 5)) {
			h.DB.Model(&model.Order{}).Where("id", order.Id).Update("checked", true)
			logger.Errorf("订单超时：%v", order)
			continue
		}
		// 查询订单状态
		var res payment.OrderInfo
		switch order.Channel {
		case payment.PayChannelEpay:
			res, err = h.epayService.Query(order.OrderNo)
			if err != nil {
				logger.Errorf("error with query order info: %v", err)
				continue
			}
			// 微信支付
		case payment.PayChannelWX:
			res, err = h.wxpayService.Query(order.OrderNo)
			logger.Debugf("微信支付订单状态：%+v", res)
			if err != nil {
				logger.Errorf("error with query order info: %v", err)
				continue
			}
		case payment.PayChannelAL:
			res, err = h.alipayService.Query(order.OrderNo)
			logger.Debugf("支付宝支付订单状态：%+v", res)
			if err != nil {
				logger.Errorf("error with query order info: %v", err)
				continue
			}
		case payment.PayChannelStripe:
			continue
		}

		// 订单已关闭
		if res.Closed() {
			h.DB.Model(&model.Order{}).Where("id", order.Id).Updates(map[string]any{
				"checked": true,
				"status":  types.OrderPaidFailed,
			})
			logger.Errorf("订单已关闭：%v", order)
			continue
		}

		// 订单未支付，不处理，继续轮询
		if !res.Success() {
			continue
		}

		// 订单支付成功
		err = h.paySuccess(res)
		if err != nil {
			logger.Errorf("error with deal order: %v", err)
			continue
		}
	}
	return nil
}

func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	var data struct {
		PayWay  string `json:"pay_way,omitempty"` // 支付方式：支付宝，微信
		Pid     int    `json:"pid,omitempty"`
		Device  string `json:"device,omitempty"`
		Domain  string `json:"domain,omitempty"` // 支付回调域名
		Channel string `json:"channel,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var product model.Product
	err := h.DB.Where("id", data.Pid).First(&product).Error
	if err != nil {
		resp.ERROR(c, "Product not found")
		return
	}

	orderNo, err := h.snowflake.Next(false)
	if err != nil {
		resp.ERROR(c, "error with generate trade no: "+err.Error())
		return
	}
	userId := h.GetLoginUserId(c)
	var user model.User
	err = h.DB.Where("id", userId).First(&user).Error
	if err != nil {
		resp.NotAuth(c)
		return
	}

	amount := product.Price
	stripeAmount := product.StripePrice
	var payURL, notifyURL string
	switch data.PayWay {
	case "wxpay":
		logger.Debugf("微信支付，%+v", data)
		data.Channel = payment.PayChannelWX
		// 优先使用微信官方支付
		if h.config.WxPay.Enabled {
			data.Channel = "wxpay"
			if h.config.WxPay.Domain != "" {
				data.Domain = h.config.WxPay.Domain
			}
			notifyURL = fmt.Sprintf("%s/api/payment/notify/wxpay", data.Domain)
			params := payment.PayRequest{
				OutTradeNo: orderNo,
				TotalFee:   fmt.Sprintf("%d", int(amount*100)),
				Subject:    product.Name,
				NotifyURL:  notifyURL,
				ClientIP:   c.ClientIP(),
				Device:     data.Device,
				PayWay:     payment.PayWayWX,
			}
			if data.Device == "mobile" {
				params.OpenID = user.OpenId
			}
			payURL, err = h.wxpayService.Pay(params)
			if err != nil {
				resp.ERROR(c, err.Error())
				return
			}
		} else if h.config.Epay.Enabled { // 聚合支付
			logger.Debugf("聚合支付%+v", data)
			data.Channel = payment.PayChannelEpay
			if h.config.Epay.Domain != "" {
				data.Domain = h.config.Epay.Domain
			}
			notifyURL = fmt.Sprintf("%s/api/payment/notify/epay", data.Domain)
			params := payment.PayRequest{
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   fmt.Sprintf("%f", amount),
				ClientIP:   c.ClientIP(),
				Device:     data.Device,
				PayWay:     payment.PayWayWX,
				NotifyURL:  notifyURL,
			}

			r, err := h.epayService.Pay(params)
			logger.Debugf("请求支付结果，%+v", r)
			if err != nil {
				resp.ERROR(c, err.Error())
				return
			} else {
				payURL = r
			}
		} else {
			resp.ERROR(c, "系统没有配置可用的支付渠道！")
			return
		}
	case "alipay":
		if h.config.Alipay.Enabled {
			logger.Debugf("支付宝，%+v", data)
			data.Channel = payment.PayChannelAL
			if h.config.Alipay.Domain != "" { // 用于本地调试支付
				data.Domain = h.config.Alipay.Domain
			}
			notifyURL = fmt.Sprintf("%s/api/payment/notify/alipay", data.Domain)
			money := fmt.Sprintf("%.2f", amount)
			payURL, err = h.alipayService.Pay(payment.PayRequest{
				Device:     data.Device,
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   money,
				NotifyURL:  notifyURL,
			})

			if err != nil {
				resp.ERROR(c, "error with generate pay url: "+err.Error())
				return
			}
		} else if h.config.Epay.Enabled { // 聚合支付
			logger.Debugf("聚合支付，%+v", data)
			data.Channel = payment.PayChannelEpay
			if h.config.Epay.Domain != "" {
				data.Domain = h.config.Epay.Domain
			}
			notifyURL = fmt.Sprintf("%s/api/payment/notify/epay", data.Domain)
			params := payment.PayRequest{
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   fmt.Sprintf("%f", amount),
				ClientIP:   c.ClientIP(),
				Device:     data.Device,
				PayWay:     data.PayWay,
				NotifyURL:  notifyURL,
			}

			r, err := h.epayService.Pay(params)
			if err != nil {
				resp.ERROR(c, err.Error())
				return
			} else {
				payURL = r
			}
		} else {
			resp.ERROR(c, "系统没有配置可用的支付渠道！")
			return
		}
	case "stripe":
		if h.config.Stripe.Enabled {
			logger.Debugf("Stripe，%+v", data)
			data.Channel = payment.PayChannelStripe
			if h.config.Stripe.Domain != "" {
				data.Domain = h.config.Stripe.Domain
			}
			if stripeAmount <= 0 {
				resp.ERROR(c, "当前产品未配置 Stripe 价格")
				return
			}
			successURL := fmt.Sprintf("%s/member", data.Domain)
			params := payment.PayRequest{
				Device:     data.Device,
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   fmt.Sprintf("%.2f", stripeAmount),
				Currency:   h.config.Stripe.Currency,
				ReturnURL:  successURL,
			}
			payURL, err = h.stripeService.Pay(params)
			if err != nil {
				resp.ERROR(c, "error with generate stripe checkout session: "+err.Error())
				return
			}
		} else if h.config.Epay.Enabled {
			logger.Debugf("Stripe fall back to epay，%+v", data)
			data.Channel = payment.PayChannelEpay
			if h.config.Epay.Domain != "" {
				data.Domain = h.config.Epay.Domain
			}
			notifyURL = fmt.Sprintf("%s/api/payment/notify/epay", data.Domain)
			params := payment.PayRequest{
				OutTradeNo: orderNo,
				Subject:    product.Name,
				TotalFee:   fmt.Sprintf("%f", amount),
				ClientIP:   c.ClientIP(),
				Device:     data.Device,
				PayWay:     data.PayWay,
				NotifyURL:  notifyURL,
			}
			r, err := h.epayService.Pay(params)
			if err != nil {
				resp.ERROR(c, err.Error())
				return
			}
			payURL = r
		} else {
			resp.ERROR(c, "系统没有配置可用的支付渠道！")
			return
		}
	default:
		resp.ERROR(c, "不支持的支付渠道")
		return
	}

	// 创建订单
	remark := types.OrderRemark{
		Power: product.Power,
		Name:  product.Name,
		Price: product.Price,
	}
	order := model.Order{
		UserId:   user.Id,
		Username: user.Username,
		OrderNo:  orderNo,
		Subject:  product.Name,
		Amount:   amount,
		Status:   types.OrderNotPaid,
		PayWay:   data.PayWay,
		Channel:  data.Channel,
		Remark:   utils.JsonEncode(remark),
	}
	err = h.DB.Create(&order).Error
	if err != nil {
		resp.ERROR(c, "error with create order: "+err.Error())
		return
	}
	resp.SUCCESS(c, gin.H{"pay_url": payURL, "order_no": orderNo})
}

// 支付成功处理
//
// 资金安全要点（重构原因）：
//  1. 整个处理必须是【一个事务】。旧实现里积分由 userService.IncreasePower 自己开事务
//     提交，后续更新订单/销量失败时根本回滚不了，结果是"积分已发、订单还停在待支付"，
//     渠道一重试就再发一次积分 —— 钱的最大口子。
//  2. 幂等闸门用条件更新实现：只有把订单从"非成功"改成"成功"的那个执行者才有资格发货。
//     并发回调 / 轮询补偿 / 渠道重投都只会有一个赢家，其余 RowsAffected==0 直接返回。
//  3. 发货走 LedgerService.Grant，带业务幂等键 (order, order.Id)，即使闸门被绕过，
//     数据库唯一索引 uk_power_ref 仍然兜住第二次入账。
//  4. 必须校验金额：渠道回传金额与订单金额不一致时拒绝发货并返回错误，绝不按错的金额发货。
func (h *PaymentHandler) paySuccess(info payment.OrderInfo) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	var order model.Order
	err := h.DB.Where("order_no", info.OutTradeNo).First(&order).Error
	if err != nil {
		return fmt.Errorf("error with fetch order: %v", err)
	}

	// 已支付订单，直接返回（快速路径；真正的幂等闸门在下面的事务里）
	if order.Status == types.OrderPaidSuccess {
		return nil
	}

	// 金额校验：宁可失败让渠道重试/人工介入，也不能按错误金额发货
	if err := verifyPaidAmount(order, info); err != nil {
		return err
	}

	var remark types.OrderRemark
	err = utils.JsonDecode(order.Remark, &remark)
	if err != nil {
		return fmt.Errorf("error with decode order remark: %v", err)
	}

	if h.ledgerService == nil {
		return errors.New("积分账本服务未初始化，无法发放充值积分")
	}

	payTime := utils.Str2stamp(info.PayTime)
	return h.DB.Transaction(func(tx *gorm.DB) error {
		// 幂等闸门：状态不等于"已支付成功"时才允许改，改成功者才拥有发货权
		res := tx.Model(&model.Order{}).
			Where("id = ? AND status <> ?", order.Id, types.OrderPaidSuccess).
			Updates(map[string]any{
				"status":   types.OrderPaidSuccess,
				"pay_time": payTime,
				"trade_no": info.TradeId,
				"checked":  true,
			})
		if res.Error != nil {
			return fmt.Errorf("error with update order info: %v", res.Error)
		}
		if res.RowsAffected == 0 {
			// 已被其它回调/轮询处理过，本次不发货（不是错误）
			logger.Infof("订单 %s 已处理过，跳过重复发货", order.OrderNo)
			return nil
		}

		// 发放积分：与订单状态在同一事务内
		if remark.Power > 0 {
			granted, err := h.ledgerService.Grant(
				tx,
				order.UserId,
				remark.Power,
				service.LedgerRef{Type: "order", ID: order.Id},
				model.PowerLog{
					UserId: order.UserId,
					Type:   types.PowerRecharge,
					Model:  order.Subject,
					Remark: fmt.Sprintf("充值积分，金额：%f，订单号：%s", order.Amount, order.OrderNo),
				},
			)
			if err != nil {
				return err
			}
			if !granted {
				// 订单已改成功但积分此前发过：说明历史数据已入账，不重复发放
				logger.Warnf("订单 %s 的积分此前已入账，本次跳过发放", order.OrderNo)
			}
		}

		// 更新产品销量：同一事务，失败则整单回滚（渠道会重试）
		if err = tx.Model(&model.Product{}).Where("id = ?", order.ProductId).
			UpdateColumn("sales", gorm.Expr("sales + ?", 1)).Error; err != nil {
			return fmt.Errorf("error with update product sales: %v", err)
		}
		return nil
	})
}

// verifyPaidAmount 校验渠道回传的实付金额与订单金额是否一致。
// 允许 1 分钱误差（渠道侧分/元换算与浮点表示）。
// 渠道未回传金额时只告警不拦截，避免因个别渠道字段缺失导致已付款订单永远不到账。
func verifyPaidAmount(order model.Order, info payment.OrderInfo) error {
	raw := strings.TrimSpace(info.Amount)
	if raw == "" {
		logger.Warnf("支付渠道未回传金额，跳过金额校验，订单号：%s", order.OrderNo)
		return nil
	}
	paid, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("无法解析渠道回传金额 %q，订单号：%s", raw, order.OrderNo)
	}
	if math.Abs(paid-order.Amount) > 0.01 {
		logger.Errorf("支付金额与订单金额不一致，订单号：%s，应付：%.2f，实付：%.2f", order.OrderNo, order.Amount, paid)
		return fmt.Errorf("支付金额不匹配：订单 %s 应付 %.2f，实付 %.2f", order.OrderNo, order.Amount, paid)
	}
	return nil
}

// AlipayNotify 支付宝支付回调
func (h *PaymentHandler) AlipayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	orderInfo, err := h.alipayService.Query(c.Request.Form.Get("out_trade_no"))
	logger.Infof("收到支付宝商号订单支付回调：%+v", orderInfo)
	if !orderInfo.Success() {
		logger.Errorf("订单校验失败：%v", err)
		c.String(http.StatusOK, "fail")
		return
	}

	err = h.paySuccess(orderInfo)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// EPayNotify 易支付支付异步回调
func (h *PaymentHandler) EPayNotify(c *gin.Context) {
	var params = make(map[string]string)
	for k := range c.Request.URL.Query() {
		params[k] = c.Query(k)
	}

	logger.Infof("收到易支付订单支付回调：%+v", params)
	// 检查支付状态, 如果未支付，则返回成功
	if params["trade_status"] != "TRADE_SUCCESS" {
		c.String(http.StatusOK, "success")
		return
	}

	sign := h.epayService.Sign(params)
	if sign != c.Query("sign") {
		logger.Errorf("签名验证失败, %s, %s", sign, c.Query("sign"))
		c.String(http.StatusOK, "fail")
		return
	}
	// 查询订单状态
	order, err := h.epayService.Query(params["out_trade_no"])
	if err != nil {
		logger.Error(err)
		c.String(http.StatusOK, "fail")
		return
	}

	err = h.paySuccess(order)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// WxpayNotify 微信商户支付异步回调
func (h *PaymentHandler) WxpayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	orderInfo, err := h.wxpayService.TradeVerify(c.Request)
	logger.Infof("收到微信商号订单支付回调：%+v", orderInfo)
	if err != nil {
		logger.Errorf("订单校验失败：%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL"})
		return
	}

	err = h.paySuccess(orderInfo)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// StripeNotify Stripe 支付异步回调
func (h *PaymentHandler) StripeNotify(c *gin.Context) {
	orderInfo, err := h.stripeService.TradeVerify(c.Request)
	logger.Infof("收到 Stripe 订单支付回调：%+v", orderInfo)
	if err != nil {
		logger.Errorf("订单校验失败：%v", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	err = h.paySuccess(orderInfo)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}
