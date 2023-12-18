package vo

type Product struct {
	BaseVo
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
	Days     int     `json:"days"`
	Calls    int     `json:"calls"`
	ImgCalls int     `json:"img_calls"`
	Enabled  bool    `json:"enabled"`
	Sales    int     `json:"sales"`
	SortNum  int     `json:"sort_num"`
}
