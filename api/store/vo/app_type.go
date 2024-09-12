package vo

type AppType struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	SortNum   int    `json:"sort_num"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at"`
}
