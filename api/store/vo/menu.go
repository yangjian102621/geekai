package vo

// Menu 系统菜单
type Menu struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	URL     string `json:"url"`
	SortNum int    `json:"sort_num"`
	Enabled bool   `json:"enabled"`
}
