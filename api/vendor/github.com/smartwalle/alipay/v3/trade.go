package alipay

import (
	"net/url"
)

// TradePagePay 统一收单下单并支付页面接口 https://docs.open.alipay.com/api_1/alipay.trade.page.pay
func (this *Client) TradePagePay(param TradePagePay) (result *url.URL, err error) {
	p, err := this.URLValues(param)
	if err != nil {
		return nil, err
	}

	result, err = url.Parse(this.host + "?" + p.Encode())
	if err != nil {
		return nil, err
	}
	return result, err
}

// TradeAppPay App支付接口 https://docs.open.alipay.com/api_1/alipay.trade.app.pay
func (this *Client) TradeAppPay(param TradeAppPay) (result string, err error) {
	p, err := this.URLValues(param)
	if err != nil {
		return "", err
	}
	return p.Encode(), err
}

// TradeFastPayRefundQuery 统一收单交易退款查询接口 https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
func (this *Client) TradeFastPayRefundQuery(param TradeFastPayRefundQuery) (result *TradeFastPayRefundQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeOrderSettle 统一收单交易结算接口 https://docs.open.alipay.com/api_1/alipay.trade.order.settle
func (this *Client) TradeOrderSettle(param TradeOrderSettle) (result *TradeOrderSettleRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeClose 统一收单交易关闭接口 https://docs.open.alipay.com/api_1/alipay.trade.close/
func (this *Client) TradeClose(param TradeClose) (result *TradeCloseRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeCancel 统一收单交易撤销接口 https://docs.open.alipay.com/api_1/alipay.trade.cancel/
func (this *Client) TradeCancel(param TradeCancel) (result *TradeCancelRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeRefund 统一收单交易退款接口 https://docs.open.alipay.com/api_1/alipay.trade.refund/
func (this *Client) TradeRefund(param TradeRefund) (result *TradeRefundRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradePreCreate 统一收单线下交易预创建接口 https://docs.open.alipay.com/api_1/alipay.trade.precreate/
func (this *Client) TradePreCreate(param TradePreCreate) (result *TradePreCreateRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeQuery 统一收单线下交易查询接口 https://docs.open.alipay.com/api_1/alipay.trade.query/
func (this *Client) TradeQuery(param TradeQuery) (result *TradeQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeCreate 统一收单交易创建接口 https://docs.open.alipay.com/api_1/alipay.trade.create/
func (this *Client) TradeCreate(param TradeCreate) (result *TradeCreateRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradePay 统一收单交易支付接口 https://docs.open.alipay.com/api_1/alipay.trade.pay/
func (this *Client) TradePay(param TradePay) (result *TradePayRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeOrderInfoSync 支付宝订单信息同步接口 https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync/
func (this *Client) TradeOrderInfoSync(param TradeOrderInfoSync) (result *TradeOrderInfoSyncRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeRefundAsync 统一收单交易退款(异步)接口 https://opendocs.alipay.com/pre-apis/api_pre/alipay.trade.refund.apply
func (this *Client) TradeRefundAsync(param TradeRefundAsync) (result *TradeRefundAsyncRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeMergePreCreate 统一收单合并支付预创建接口请求参数 https://opendocs.alipay.com/open/028xr9
// TODO TradeMergePreCreate 接口未经测试
func (this *Client) TradeMergePreCreate(param TradeMergePreCreate) (result *TradeMergePreCreateRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// TradeAppMergePay App合并支付接口 https://opendocs.alipay.com/open/028py8
// TODO TradeAppMergePay 接口未经测试
func (this *Client) TradeAppMergePay(param TradeAppPay) (result string, err error) {
	p, err := this.URLValues(param)
	if err != nil {
		return "", err
	}
	return p.Encode(), err
}
