package alipay

import "net/url"

// TradeWapPay 手机网站支付接口 https://docs.open.alipay.com/api_1/alipay.trade.wap.pay/
func (this *Client) TradeWapPay(param TradeWapPay) (result *url.URL, err error) {
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

// TradeWapMergePay 无线Wap合并支付接口2.0 https://opendocs.alipay.com/open/028xra
// TODO TradeWapMergePay 接口未经测试
func (this *Client) TradeWapMergePay(param TradeWapMergePay) (result *url.URL, err error) {
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
