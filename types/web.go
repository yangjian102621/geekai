package types

// BizVo 业务返回 VO
type BizVo struct {
	Code     BizCode     `json:"code"`
	Page     int         `json:"page,omitempty"`
	PageSize int         `json:"page_size,omitempty"`
	Total    int         `json:"total,omitempty"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
}

// WsVo Websocket 信息 VO
type WsVo struct {
	Stop    bool
	Content string
}

type BizCode int

const (
	Success       = BizCode(0)
	Failed        = BizCode(1)
	InvalidParams = BizCode(101) // 非法参数
	NotAuthorized = BizCode(400) // 未授权

	OkMsg = "Success"
)
