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

// WsMessage Websocket message
type WsMessage struct {
	Type    WsMsgType `json:"type"` // 消息类别，start, end
	Content string    `json:"content"`
}
type WsMsgType string

const (
	WsStart  = WsMsgType("start")
	WsMiddle = WsMsgType("middle")
	WsEnd    = WsMsgType("end")
)

type BizCode int

const (
	Success       = BizCode(0)
	Failed        = BizCode(1)
	InvalidParams = BizCode(101) // 非法参数
	NotAuthorized = BizCode(400) // 未授权

	OkMsg    = "Success"
	ErrorMsg = "系统开小差了"
)

const TokenName = "ChatGPT-Token"
const SessionKey = "WEB_SSH_SESSION"
