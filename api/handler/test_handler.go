package handler

import (
	"chatplus/service"
	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	snowflake *service.Snowflake
}

func NewTestHandler(snowflake *service.Snowflake) *TestHandler {
	return &TestHandler{snowflake: snowflake}
}

func (h *TestHandler) TestPay(c *gin.Context) {
	//appId := ""                                           //Appid
	//appSecret := ""                                       //密钥
	//var host = "https://api.xunhupay.com/payment/do.html" //跳转支付页接口URL
	//client := payment.NewXunHuPay(appId, appSecret)     //初始化调用
	//
	////支付参数，appid、time、nonce_str和hash这四个参数不用传，调用的时候执行方法内部已经处理
	//orderNo, _ := h.snowflake.Next()
	//params := map[string]string{
	//	"version":        "1.1",
	//	"trade_order_id": orderNo,
	//	"total_fee":      "0.1",
	//	"title":          "测试支付",
	//	"notify_url":     "http://xxxxxxx.com",
	//	"return_url":     "http://localhost:8888",
	//	"wap_name":       "极客学长",
	//	"callback_url":   "",
	//}
	//
	//execute, err := client.Execute(host, params) //执行支付操作
	//if err != nil {
	//	logger.Error(err)
	//}
	//resp.SUCCESS(c, execute)
}
