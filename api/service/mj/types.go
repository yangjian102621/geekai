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
