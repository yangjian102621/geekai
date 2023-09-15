package types

// BizVo 业务返回 VO
type BizVo struct {
	Code     BizCode     `json:"code"`
	Page     int         `json:"page,omitempty"`
	PageSize int         `json:"page_size,omitempty"`
	Total    int         `json:"total,omitempty"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// WsMessage Websocket message
type WsMessage struct {
	Type    WsMsgType   `json:"type"` // 消息类别，start, end, img
	Content interface{} `json:"content"`
}
type WsMsgType string

const (
	WsStart  = WsMsgType("start")
	WsMiddle = WsMsgType("middle")
	WsEnd    = WsMsgType("end")
	WsMjImg  = WsMsgType("mj")
)

type BizCode int

const (
	Success       = BizCode(0)
	Failed        = BizCode(1)
	NotAuthorized = BizCode(400) // 未授权

	OkMsg       = "Success"
	ErrorMsg    = "系统开小差了"
	InvalidArgs = "非法参数或参数解析失败"
	NoData      = "No Data"
)
