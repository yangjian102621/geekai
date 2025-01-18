package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OrderStatus int

const (
	OrderNotPaid     = OrderStatus(0)
	OrderScanned     = OrderStatus(1) // 已扫码
	OrderPaidSuccess = OrderStatus(2)
)

type OrderRemark struct {
	Days     int     `json:"days"`  // 有效期
	Power    int     `json:"power"` // 增加算力点数
	Name     string  `json:"name"`  // 产品名称
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
}

var PayMethods = map[string]string{
	"alipay": "支付宝商号",
	"wechat": "微信商号",
	"hupi":   "虎皮椒",
	"geek":   "易支付",
}
var PayNames = map[string]string{
	"alipay": "支付宝",
	"wxpay":  "微信支付",
	"qqpay":  "QQ钱包",
	"jdpay":  "京东支付",
	"douyin": "抖音支付",
	"paypal": "PayPal支付",
}
