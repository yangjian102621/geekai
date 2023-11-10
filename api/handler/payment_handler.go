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
	PayWayWechat = "微信支付"
)

// PaymentHandler 支付服务回调 handler
type PaymentHandler struct {
	BaseHandler
	alipayService *payment.AlipayService
	snowflake     *service.Snowflake
	db            *gorm.DB
	fs            embed.FS
	lock          sync.Mutex
}

func NewPaymentHandler(server *core.AppServer, alipayService *payment.AlipayService, snowflake *service.Snowflake, db *gorm.DB, fs embed.FS) *PaymentHandler {
	h := PaymentHandler{lock: sync.Mutex{}}
	h.App = server
	h.alipayService = alipayService
	h.snowflake = snowflake
	h.db = db
	h.fs = fs
	return &h
}

func (h *PaymentHandler) Alipay(c *gin.Context) {
	orderNo := h.GetTrim(c, "order_no")
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
}

// OrderQuery 清单状态查询
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

// AlipayQrcode 生成支付宝支付 URL 二维码
func (h *PaymentHandler) AlipayQrcode(c *gin.Context) {
	if !h.App.SysConfig.EnabledAlipay || h.alipayService == nil {
		resp.ERROR(c, "当前支付通道已经关闭，请联系管理员开通！")
		return
	}

	var data struct {
		ProductId uint `json:"product_id"`
		UserId    int  `json:"user_id"`
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

	orderNo, err := h.snowflake.Next()
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

	// 创建订单
	remark := types.OrderRemark{
		Days:     product.Days,
		Calls:    product.Calls,
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
		PayWay:    PayWayAlipay,
		Remark:    utils.JsonEncode(remark),
	}
	res = h.db.Create(&order)
	if res.Error != nil {
		resp.ERROR(c, "error with create order: "+res.Error.Error())
		return
	}

	// 生成二维码图片
	file, err := h.fs.Open("res/img/alipay.jpg")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	parse, err := url.Parse(h.App.Config.AlipayConfig.NotifyURL)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	imageURL := fmt.Sprintf("%s://%s/api/payment/alipay?order_no=%s", parse.Scheme, parse.Host, orderNo)
	imgData, err := utils.GenQrcode(imageURL, 400, file)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	imgDataBase64 := base64.StdEncoding.EncodeToString(imgData)
	resp.SUCCESS(c, gin.H{"order_no": orderNo, "image": fmt.Sprintf("data:image/jpg;base64, %s", imgDataBase64), "url": imageURL})
}

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

	var order model.Order
	res := h.db.Where("order_no = ?", r.OutTradeNo).First(&order)
	if res.Error != nil {
		logger.Error(res.Error)
		c.String(http.StatusOK, "fail")
		return
	}
	var user model.User
	res = h.db.First(&user, order.UserId)
	if res.Error != nil {
		logger.Error(res.Error)
		c.String(http.StatusOK, "fail")
		return
	}
	var remark types.OrderRemark
	err = utils.JsonDecode(order.Remark, &remark)
	if err != nil {
		logger.Error(res.Error)
		c.String(http.StatusOK, "fail")
		return
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

	// 更新用户信息
	res = h.db.Updates(&user)
	if res.Error != nil {
		logger.Error(res.Error)
		c.String(http.StatusOK, "fail")
		return
	}

	// 更新订单状态
	order.PayTime = time.Now().Unix()
	order.Status = types.OrderPaidSuccess
	h.db.Updates(&order)

	// 更新产品销量
	h.db.Model(&model.Product{}).Where("id = ?", order.ProductId).UpdateColumn("sales", gorm.Expr("sales + ?", 1))

	c.String(http.StatusOK, "success")
}
