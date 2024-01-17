package mj

const (
	ApplicationID string = "936929561302675456"
	SessionID     string = "ea8816d857ba9ae2f74c59ae1a953afe"
)

type InteractionsRequest struct {
	Type          int            `json:"type"`
	ApplicationID string         `json:"application_id"`
	MessageFlags  *int           `json:"message_flags,omitempty"`
	MessageID     *string        `json:"message_id,omitempty"`
	GuildID       string         `json:"guild_id"`
	ChannelID     string         `json:"channel_id"`
	SessionID     string         `json:"session_id"`
	Data          map[string]any `json:"data"`
	Nonce         string         `json:"nonce,omitempty"`
}

type InteractionsResult struct {
	Code    int `json:"code"`
	Message string
	Error   map[string]any
}

type CBReq struct {
	ChannelId   string     `json:"channel_id"`
	MessageId   string     `json:"message_id"`
	ReferenceId string     `json:"reference_id"`
	Image       Image      `json:"image"`
	Content     string     `json:"content"`
	Prompt      string     `json:"prompt"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"`
}

type MjAPIRequest struct {
	// 结果回调地址，[请参考回调使用指南](https://chuzhanai.apifox.cn/doc-3556414)
	Callback string `json:"callback,omitempty"`
	// 图片链接，base64图片或者图片url地址,用于图生图 （如果是图片url请确保此url公网可以访问）
	ImageURL string `json:"imageUrl,omitempty"`
	// 请求回调标识，任意字符串，具体[请参考回调使用指南](https://chuzhanai.apifox.cn/doc-3556414)
	Nonce string `json:"nonce,omitempty"`
	// 提示词，提示词，支持中英文，如果img字段为空则必填
	Prompt string `json:"prompt,omitempty"`
	// 生成通道，生成通道，fast：快速生成，消耗40积分，relax：普速，需排队等待，消耗25积分，默认是relax
	Type string `json:"type,omitempty"`
}

//	{
//	    "data": {
//	        "paintingSign": "74cda5279c6a4097980079fe87b486b9",
//	        "used": 4,
//	        "balance": 9996,
//	    },
//	    "code": 0,
//	    "subCode": 1000,
//	    "success": true,
//	    "msg" : "词条中涉及敏感词，请调整后再生成图片"
//	}

type MjAPIResponse struct {
	Data struct {
		PaintingSign string `json:"paintingSign"`
		Used         int    `json:"used"`
		Balance      int    `json:"balance"`
	} `json:"data"`
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// {
//     "data": {
//         "progress": "100%",
//         "progressValue": 0.1,
//         "current_image": "/9j/4AAQSkZJRgABAQAAAQABAAD/2w",
//         "state": "running",
//         "imgUrl": "https://imagesource.huashi6.com/test.jpg",
//         "audit": 1
//     },
//     "code": 0,
//     "success": true
// }

type MjTaskStatusData struct {
	Progress      string     `json:"progress"`
	ProgressValue float32    `json:"progressValue"`
	CurrentImage  string     `json:"current_image"`
	State         TaskStatus `json:"state"`
	ImgURL        string     `json:"imgUrl"`
	Audit         int        `json:"audit"`
}

type MjAPICheckStatuResponse struct {
	Data    *MjTaskStatusData `json:"data"`
	Code    int               `json:"code"`
	Success bool              `json:"success"`
}

type CheckStatusRequest struct {
	TaskID  string `json:"taskId"`
	Preview bool   `json:"preview"`
}

func NewCheckStatusRequest(taskID string) *CheckStatusRequest {
	return &CheckStatusRequest{
		TaskID:  taskID,
		Preview: false,
	}
}
