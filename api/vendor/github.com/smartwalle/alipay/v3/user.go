package alipay

import (
	"net/url"
)

// AgreementPageSign 支付宝个人协议页面签约接口 https://docs.open.alipay.com/api_2/alipay.user.agreement.page.sign
func (this *Client) AgreementPageSign(param AgreementPageSign) (result *url.URL, err error) {
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

// AgreementQuery 支付宝个人代扣协议查询接口 https://opendocs.alipay.com/open/02fkao?scene=8837b4183390497f84bb53783b488ecc
func (this *Client) AgreementQuery(param AgreementQuery) (result *AgreementQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// AgreementUnsign 支付宝个人代扣协议解约接口 https://docs.open.alipay.com/api_2/alipay.user.agreement.unsign
func (this *Client) AgreementUnsign(param AgreementUnsign) (result *AgreementUnsignRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// AgreementExecutionPlanModify 周期性扣款协议执行计划修改接口 https://docs.open.alipay.com/api_2/alipay.user.agreement.executionplan.modify
func (this *Client) AgreementExecutionPlanModify(param AgreementExecutionPlanModify) (result *AgreementExecutionPlanModifyRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// DecodePhoneNumber 小程序获取会员手机号  https://opendocs.alipay.com/mini/api/getphonenumber
//
// 本方法用于解码小程序端 my.getPhoneNumber 获取的数据
func (this *Client) DecodePhoneNumber(data string) (result *MobileNumber, err error) {
	if err = this.decode([]byte(data), "response", true, &result); err != nil {
		return nil, err
	}
	return result, nil
}
