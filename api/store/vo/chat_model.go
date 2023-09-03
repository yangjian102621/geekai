package vo

type ChatModel struct {
	BaseVo
	Platform string `json:"platform"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Enabled  bool   `json:"enabled"`
}
