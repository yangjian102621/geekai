package alipay

// FundTransToAccountTransfer 单笔转账到支付宝账户接口 https://docs.open.alipay.com/api_28/alipay.fund.trans.toaccount.transfer
func (this *Client) FundTransToAccountTransfer(param FundTransToAccountTransfer) (result *FundTransToAccountTransferRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundTransOrderQuery 查询转账订单接口 https://docs.open.alipay.com/api_28/alipay.fund.trans.order.query/
func (this *Client) FundTransOrderQuery(param FundTransOrderQuery) (result *FundTransOrderQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAuthOrderVoucherCreate 资金授权发码接口 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.voucher.create/
func (this *Client) FundAuthOrderVoucherCreate(param FundAuthOrderVoucherCreate) (result *FundAuthOrderVoucherCreateRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAuthOrderFreeze 资金授权冻结接口 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.freeze/
func (this *Client) FundAuthOrderFreeze(param FundAuthOrderFreeze) (result *FundAuthOrderFreezeRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAuthOrderUnfreeze 资金授权解冻接口 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.unfreeze/
func (this *Client) FundAuthOrderUnfreeze(param FundAuthOrderUnfreeze) (result *FundAuthOrderUnfreezeRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAuthOperationCancel 资金授权撤销接口 https://docs.open.alipay.com/api_28/alipay.fund.auth.operation.cancel/
func (this *Client) FundAuthOperationCancel(param FundAuthOperationCancel) (result *FundAuthOperationCancelRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAuthOperationDetailQuery 资金授权操作查询接口 https://docs.open.alipay.com/api_28/alipay.fund.auth.operation.detail.query/
func (this *Client) FundAuthOperationDetailQuery(param FundAuthOperationDetailQuery) (result *FundAuthOperationDetailQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAuthOrderAppFreeze 线上资金授权冻结接口 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.app.freeze
func (this *Client) FundAuthOrderAppFreeze(param FundAuthOrderAppFreeze) (result string, err error) {
	p, err := this.URLValues(param)
	if err != nil {
		return "", err
	}
	return p.Encode(), err
}

// FundTransUniTransfer 单笔转账接口 https://docs.open.alipay.com/api_28/alipay.fund.trans.uni.transfer/
func (this *Client) FundTransUniTransfer(param FundTransUniTransfer) (result *FundTransUniTransferRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundTransCommonQuery 转账业务单据查询接口 https://docs.open.alipay.com/api_28/alipay.fund.trans.common.query/
func (this *Client) FundTransCommonQuery(param FundTransCommonQuery) (result *FundTransCommonQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundAccountQuery 支付宝资金账户资产查询接口  https://docs.open.alipay.com/api_28/alipay.fund.account.query
func (this *Client) FundAccountQuery(param FundAccountQuery) (result *FundAccountQueryRsp, err error) {
	err = this.doRequest("POST", param, &result)
	return result, err
}

// FundTransAppPay 现金红包无线支付接口 https://opendocs.alipay.com/open/03rbyf
func (this *Client) FundTransAppPay(param FundTransAppPay) (result string, err error) {
	p, err := this.URLValues(param)
	if err != nil {
		return "", err
	}
	return p.Encode(), err
}
