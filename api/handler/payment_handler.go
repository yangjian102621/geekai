package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/payment"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"embed"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	PayWayAlipay = "支付宝"
	PayWayXunHu  = "虎皮椒"
)

// PaymentHandler 支付服务回调 handler
type PaymentHandler struct {
	BaseHandler
	alipayService  *payment.AlipayService
	huPiPayService *payment.HuPiPayService
	snowflake      *service.Snowflake
	db             *gorm.DB
	fs             embed.FS
	lock           sync.Mutex
}

func NewPaymentHandler(server *core.AppServer, alipayService *payment.AlipayService, huPiPayService *payment.HuPiPayService, snowflake *service.Snowflake, db *gorm.DB, fs embed.FS) *PaymentHandler {
	h := PaymentHandler{
		alipayService:  alipayService,
		huPiPayService: huPiPayService,
		snowflake:      snowflake,
		fs:             fs,
		db:             db,
		lock:           sync.Mutex{},
	}
	h.App = server
	return &h
}

func (h *PaymentHandler) DoPay(c *gin.Context) {
	orderNo := h.GetTrim(c, "order_no")
	payWay := h.GetTrim(c, "pay_way")

	if orderNo == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var order model.Order
	res := h.db.Where("order_no = ?", orderNo).First(&order)
	if res.Error != nil {
		resp.ERROR(c, "Order not found")
		return
	}

	// 更新扫码状态
	h.db.Model(&order).UpdateColumn("status", types.OrderScanned)
	if payWay == "alipay" { // 支付宝
		// 生成支付链接
		notifyURL := h.App.Config.AlipayConfig.NotifyURL
		returnURL := "" // 关闭同步回跳
		amount := fmt.Sprintf("%.2f", order.Amount)

		uri, err := h.alipayService.PayUrlMobile(order.OrderNo, notifyURL, returnURL, amount, order.Subject)
		if err != nil {
			resp.ERROR(c, "error with generate pay url: "+err.Error())
			return
		}

		c.Redirect(302, uri)
		return
	} else if payWay == "hupi" { // 虎皮椒支付
		params := map[string]string{
			"version":        "1.1",
			"trade_order_id": orderNo,
			"total_fee":      fmt.Sprintf("%f", order.Amount),
			"title":          order.Subject,
			"notify_url":     h.App.Config.HuPiPayConfig.NotifyURL,
			"return_url":     "",
			"wap_name":       "极客学长",
			"callback_url":   "",
		}

		res, err := h.huPiPayService.Pay(params)
		if err != nil {
			resp.ERROR(c, "error with generate pay url: "+err.Error())
			return
		}

		var r struct {
			Openid    interface{} `json:"openid"`
			UrlQrcode string      `json:"url_qrcode"`
			URL       string      `json:"url"`
			ErrCode   int         `json:"errcode"`
			ErrMsg    string      `json:"errmsg,omitempty"`
		}
		err = utils.JsonDecode(res, &r)
		if err != nil {
			logger.Debugf(res)
			resp.ERROR(c, "error with decode payment result: "+err.Error())
			return
		}

		if r.ErrCode != 0 {
			resp.ERROR(c, "error with generate pay url: "+r.ErrMsg)
			return
		}
		c.Redirect(302, r.URL)
	}
	resp.ERROR(c, "Invalid operations")
}

// OrderQuery 查询订单状态
func (h *PaymentHandler) OrderQuery(c *gin.Context) {
	var data struct {
		OrderNo string `json:"order_no"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var order model.Order
	res := h.db.Where("order_no = ?", data.OrderNo).First(&order)
	if res.Error != nil {
		resp.ERROR(c, "Order not found")
		return
	}

	if order.Status == types.OrderPaidSuccess {
		resp.SUCCESS(c, gin.H{"status": order.Status})
		return
	}

	counter := 0
	for {
		time.Sleep(time.Second)
		var item model.Order
		h.db.Where("order_no = ?", data.OrderNo).First(&item)
		if counter >= 15 || item.Status == types.OrderPaidSuccess || item.Status != order.Status {
			order.Status = item.Status
			break
		}
		counter++
	}

	resp.SUCCESS(c, gin.H{"status": order.Status})
}

// PayQrcode 生成支付 URL 二维码
func (h *PaymentHandler) PayQrcode(c *gin.Context) {
	var data struct {
		PayWay    string `json:"pay_way"` // 支付方式
		ProductId uint   `json:"product_id"`
		UserId    int    `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var product model.Product
	res := h.db.First(&product, data.ProductId)
	if res.Error != nil {
		resp.ERROR(c, "Product not found")
		return
	}

	orderNo, err := h.snowflake.Next(false)
	if err != nil {
		resp.ERROR(c, "error with generate trade no: "+err.Error())
		return
	}
	var user model.User
	res = h.db.First(&user, data.UserId)
	if res.Error != nil {
		resp.ERROR(c, "Invalid user ID")
		return
	}

	payWay := PayWayAlipay
	if data.PayWay == "hupi" {
		payWay = PayWayXunHu
	}
	// 创建订单
	remark := types.OrderRemark{
		Days:     product.Days,
		Calls:    product.Calls,
		ImgCalls: product.ImgCalls,
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
	}
	order := model.Order{
		UserId:    user.Id,
		Mobile:    user.Mobile,
		ProductId: product.Id,
		OrderNo:   orderNo,
		Subject:   product.Name,
		Amount:    product.Price - product.Discount,
		Status:    types.OrderNotPaid,
		PayWay:    payWay,
		Remark:    utils.JsonEncode(remark),
	}
	res = h.db.Create(&order)
	if res.Error != nil {
		resp.ERROR(c, "error with create order: "+res.Error.Error())
		return
	}

	var logo string
	if data.PayWay == "alipay" {
		logo = "res/img/alipay.jpg"
	} else if data.PayWay == "hupi" {
		if h.App.Config.HuPiPayConfig.Name == "wechat" {
			logo = "res/img/wechat-pay.jpg"
		} else {
			logo = "res/img/alipay.jpg"
		}
	}

	file, err := h.fs.Open(logo)
	if err != nil {
		resp.ERROR(c, "error with open qrcode log file: "+err.Error())
		return
	}

	parse, err := url.Parse(h.App.Config.AlipayConfig.NotifyURL)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	imageURL := fmt.Sprintf("%s://%s/api/payment/doPay?order_no=%s&pay_way=%s", parse.Scheme, parse.Host, orderNo, data.PayWay)
	imgData, err := utils.GenQrcode(imageURL, 400, file)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	imgDataBase64 := base64.StdEncoding.EncodeToString(imgData)
	resp.SUCCESS(c, gin.H{"order_no": orderNo, "image": fmt.Sprintf("data:image/jpg;base64, %s", imgDataBase64), "url": imageURL})
}

// AlipayNotify 支付宝支付回调
func (h *PaymentHandler) AlipayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	// TODO：这里最好用支付宝的公钥签名签证一下交易真假
	//res := h.alipayService.TradeVerify(c.Request.Form)
	r := h.alipayService.TradeQuery(c.Request.Form.Get("out_trade_no"))
	logger.Infof("验证支付结果：%+v", r)
	if !r.Success() {
		c.String(http.StatusOK, "fail")
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	err = h.notify(r.OutTradeNo)
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// 异步通知回调公共逻辑
func (h *PaymentHandler) notify(orderNo string) error {
	var order model.Order
	res := h.db.Where("order_no = ?", orderNo).First(&order)
	if res.Error != nil {
		err := fmt.Errorf("error with fetch order: %v", res.Error)
		logger.Error(err)
		return err
	}

	// 已支付订单，直接返回
	if order.Status == types.OrderPaidSuccess {
		return nil
	}

	var user model.User
	res = h.db.First(&user, order.UserId)
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

	// 1. 点卡：days == 0, calls > 0
	// 2. vip 套餐：days > 0, calls == 0
	if remark.Days > 0 {
		if user.ExpiredTime > time.Now().Unix() {
			user.ExpiredTime = time.Unix(user.ExpiredTime, 0).AddDate(0, 0, remark.Days).Unix()
		} else {
			user.ExpiredTime = time.Now().AddDate(0, 0, remark.Days).Unix()
		}
		user.Vip = true

	} else if !user.Vip { // 充值点卡的非 VIP 用户
		user.ExpiredTime = time.Now().AddDate(0, 0, 30).Unix()
	}

	if remark.Calls > 0 { // 充值点卡
		user.Calls += remark.Calls
	} else {
		user.Calls += h.App.SysConfig.VipMonthCalls
	}

	if remark.ImgCalls > 0 {
		user.ImgCalls += remark.ImgCalls
	} else {
		user.ImgCalls += h.App.SysConfig.VipMonthImgCalls
	}

	// 更新用户信息
	res = h.db.Updates(&user)
	if res.Error != nil {
		err := fmt.Errorf("error with update user info: %v", res.Error)
		logger.Error(err)
		return err
	}

	// 更新订单状态
	order.PayTime = time.Now().Unix()
	order.Status = types.OrderPaidSuccess
	res = h.db.Updates(&order)
	if res.Error != nil {
		err := fmt.Errorf("error with update order info: %v", res.Error)
		logger.Error(err)
		return err
	}

	// 更新产品销量
	h.db.Model(&model.Product{}).Where("id = ?", order.ProductId).UpdateColumn("sales", gorm.Expr("sales + ?", 1))
	return nil
}

// GetPayWays 获取支付方式
func (h *PaymentHandler) GetPayWays(c *gin.Context) {
	data := gin.H{}
	if h.App.Config.AlipayConfig.Enabled {
		data["alipay"] = gin.H{"name": "alipay"}
	}
	if h.App.Config.HuPiPayConfig.Enabled {
		data["hupi"] = gin.H{"name": h.App.Config.HuPiPayConfig.Name}
	}
	resp.SUCCESS(c, data)
}

// HuPiPayNotify 虎皮椒支付异步回调
func (h *PaymentHandler) HuPiPayNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	orderNo := c.Request.Form.Get("trade_order_id")
	logger.Infof("收到订单支付回调，订单 NO：%s", orderNo)
	// TODO 是否要保存订单交易流水号
	h.lock.Lock()
	defer h.lock.Unlock()

	err = h.notify(orderNo)
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}
