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
	"github.com/shopspring/decimal"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	PayWayAlipay = "支付宝"
	PayWayXunHu  = "虎皮椒"
	PayWayJs     = "PayJS"
)

// PaymentHandler 支付服务回调 handler
type PaymentHandler struct {
	BaseHandler
	alipayService  *payment.AlipayService
	huPiPayService *payment.HuPiPayService
	js             *payment.PayJS
	snowflake      *service.Snowflake
	fs             embed.FS
	lock           sync.Mutex
}

func NewPaymentHandler(
	server *core.AppServer,
	alipayService *payment.AlipayService,
	huPiPayService *payment.HuPiPayService,
	js *payment.PayJS,
	db *gorm.DB,
	snowflake *service.Snowflake,
	fs embed.FS) *PaymentHandler {
	return &PaymentHandler{
		alipayService:  alipayService,
		huPiPayService: huPiPayService,
		js:             js,
		snowflake:      snowflake,
		fs:             fs,
		lock:           sync.Mutex{},
		BaseHandler: BaseHandler{
			App: server,
			DB:  db,
		},
	}
}

func (h *PaymentHandler) DoPay(c *gin.Context) {
	orderNo := h.GetTrim(c, "order_no")
	payWay := h.GetTrim(c, "pay_way")

	if orderNo == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var order model.Order
	res := h.DB.Where("order_no = ?", orderNo).First(&order)
	if res.Error != nil {
		resp.ERROR(c, "Order not found")
		return
	}

	// fix: 这里先检查一下订单状态，如果已经支付了，就直接返回
	if order.Status == types.OrderPaidSuccess {
		resp.ERROR(c, "This order had been paid, please do not pay twice")
		return
	}

	// 更新扫码状态
	h.DB.Model(&order).UpdateColumn("status", types.OrderScanned)
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
		params := payment.HuPiPayReq{
			Version:      "1.1",
			TradeOrderId: orderNo,
			TotalFee:     fmt.Sprintf("%f", order.Amount),
			Title:        order.Subject,
			NotifyURL:    h.App.Config.HuPiPayConfig.NotifyURL,
			WapName:      "极客学长",
		}
		r, err := h.huPiPayService.Pay(params)
		if err != nil {
			resp.ERROR(c, err.Error())
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
	res := h.DB.Where("order_no = ?", data.OrderNo).First(&order)
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
		h.DB.Where("order_no = ?", data.OrderNo).First(&item)
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
	res := h.DB.First(&product, data.ProductId)
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
	res = h.DB.First(&user, data.UserId)
	if res.Error != nil {
		resp.ERROR(c, "Invalid user ID")
		return
	}

	var payWay string
	var notifyURL string
	switch data.PayWay {
	case "hupi":
		payWay = PayWayXunHu
		notifyURL = h.App.Config.HuPiPayConfig.NotifyURL
	case "payjs":
		payWay = PayWayJs
		notifyURL = h.App.Config.JPayConfig.NotifyURL
	default:
		payWay = PayWayAlipay
		notifyURL = h.App.Config.AlipayConfig.NotifyURL
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
		Remark:    utils.JsonEncode(remark),
	}
	res = h.DB.Create(&order)
	if res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "error with create order: "+res.Error.Error())
		return
	}

	// PayJs 单独处理，只能用官方生成的二维码
	if data.PayWay == "payjs" {
		params := payment.JPayReq{
			TotalFee:   int(math.Ceil(order.Amount * 100)),
			OutTradeNo: order.OrderNo,
			Subject:    product.Name,
		}
		r := h.js.Pay(params)
		if r.IsOK() {
			resp.SUCCESS(c, gin.H{"order_no": order.OrderNo, "image": r.Qrcode})
			return
		} else {
			resp.ERROR(c, "error with generating payment qrcode: "+r.ReturnMsg)
			return
		}
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

	parse, err := url.Parse(notifyURL)
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

// Mobile 移动端支付
func (h *PaymentHandler) Mobile(c *gin.Context) {
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
	res := h.DB.First(&product, data.ProductId)
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
	res = h.DB.First(&user, data.UserId)
	if res.Error != nil {
		resp.ERROR(c, "Invalid user ID")
		return
	}

	amount, _ := decimal.NewFromFloat(product.Price).Sub(decimal.NewFromFloat(product.Discount)).Float64()
	var payWay string
	var notifyURL, returnURL string
	var payURL string
	switch data.PayWay {
	case "hupi":
		payWay = PayWayXunHu
		notifyURL = h.App.Config.HuPiPayConfig.NotifyURL
		returnURL = h.App.Config.HuPiPayConfig.ReturnURL
		params := payment.HuPiPayReq{
			Version:      "1.1",
			TradeOrderId: orderNo,
			TotalFee:     fmt.Sprintf("%f", amount),
			Title:        product.Name,
			NotifyURL:    notifyURL,
			ReturnURL:    returnURL,
			CallbackURL:  returnURL,
			WapName:      "极客学长",
		}
		r, err := h.huPiPayService.Pay(params)
		if err != nil {
			logger.Error("error with generating Pay URL: ", err.Error())
			resp.ERROR(c, "error with generating Pay URL: "+err.Error())
			return
		}
		payURL = r.URL
	case "payjs":
		payWay = PayWayJs
		notifyURL = h.App.Config.JPayConfig.NotifyURL
		returnURL = h.App.Config.JPayConfig.ReturnURL
		totalFee := decimal.NewFromFloat(product.Price).Sub(decimal.NewFromFloat(product.Discount)).Mul(decimal.NewFromInt(100)).IntPart()
		params := url.Values{}
		params.Add("total_fee", fmt.Sprintf("%d", totalFee))
		params.Add("out_trade_no", orderNo)
		params.Add("body", product.Name)
		params.Add("notify_url", notifyURL)
		params.Add("auto", "0")
		payURL = h.js.PayH5(params)
	case "alipay":
		payWay = PayWayAlipay
		notifyURL = h.App.Config.AlipayConfig.NotifyURL
		returnURL = h.App.Config.AlipayConfig.ReturnURL
		payURL, err = h.alipayService.PayUrlMobile(orderNo, notifyURL, returnURL, fmt.Sprintf("%.2f", amount), product.Name)
		if err != nil {
			resp.ERROR(c, "error with generating Pay URL: "+err.Error())
			return
		}
	default:
		resp.ERROR(c, "Unsupported pay way: "+data.PayWay)
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

	order := model.Order{
		UserId:    user.Id,
		Username:  user.Username,
		ProductId: product.Id,
		OrderNo:   orderNo,
		Subject:   product.Name,
		Amount:    amount,
		Status:    types.OrderNotPaid,
		PayWay:    payWay,
		Remark:    utils.JsonEncode(remark),
	}
	res = h.DB.Create(&order)
	if res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "error with create order: "+res.Error.Error())
		return
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
	if opt != "" {
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
	data := gin.H{}
	if h.App.Config.AlipayConfig.Enabled {
		data["alipay"] = gin.H{"name": "alipay"}
	}
	if h.App.Config.HuPiPayConfig.Enabled {
		data["hupi"] = gin.H{"name": h.App.Config.HuPiPayConfig.Name}
	}
	if h.App.Config.JPayConfig.Enabled {
		data["payjs"] = gin.H{"name": h.App.Config.JPayConfig.Name}
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
	res := h.alipayService.TradeVerify(c.Request.Form)
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

// PayJsNotify PayJs 支付异步回调
func (h *PaymentHandler) PayJsNotify(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	orderNo := c.Request.Form.Get("out_trade_no")
	returnCode := c.Request.Form.Get("return_code")
	logger.Infof("收到订单支付回调，订单 NO：%s，支付结果代码：%v", orderNo, returnCode)
	// 支付失败
	if returnCode != "1" {
		return
	}

	// 校验订单支付状态
	tradeNo := c.Request.Form.Get("payjs_order_id")
	err = h.js.Check(tradeNo)
	if err != nil {
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
