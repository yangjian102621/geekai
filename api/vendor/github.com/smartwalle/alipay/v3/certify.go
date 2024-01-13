package alipay

import "net/url"

// UserCertifyOpenInitialize 身份认证初始化服务接口 https://docs.open.alipay.com/api_2/alipay.user.certify.open.initialize
func (this *Client) UserCertifyOpenInitialize(param UserCertifyOpenInitialize) (result *UserCertifyOpenInitializeRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// UserCertifyOpenCertify 身份认证开始认证接口 https://docs.open.alipay.com/api_2/alipay.user.certify.open.certify
func (this *Client) UserCertifyOpenCertify(param UserCertifyOpenCertify) (result *url.URL, err error) {
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

// UserCertifyOpenQuery 身份认证记录查询接口 https://docs.open.alipay.com/api_2/alipay.user.certify.open.query/
func (this *Client) UserCertifyOpenQuery(param UserCertifyOpenQuery) (result *UserCertifyOpenQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// UserCertDocCertVerifyPreConsult 实名证件信息比对验证预咨询 https://opendocs.alipay.com/apis/api_2/alipay.user.certdoc.certverify.preconsult
func (this *Client) UserCertDocCertVerifyPreConsult(param UserCertDocCertVerifyPreConsult) (result *UserCertDocCertVerifyPreConsultRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// UserCertDocCertVerifyConsult 实名证件信息比对验证咨询 https://opendocs.alipay.com/apis/api_2/alipay.user.certdoc.certverify.consult
func (this *Client) UserCertDocCertVerifyConsult(param UserCertDocCertVerifyConsult) (result *UserCertDocCertVerifyConsultRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}
