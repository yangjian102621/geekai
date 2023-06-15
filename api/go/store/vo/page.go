package vo

type Page struct {
	Items     interface{} `json:"items"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	Total     int64       `json:"total"`
	TotalPage int         `json:"total_page"`
}

func NewPage(total int64, page int, pageSize int, items interface{}) Page {
	totalPage := int(total / int64(pageSize))
	return Page{
		Items:     items,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}
