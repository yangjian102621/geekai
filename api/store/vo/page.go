package vo

import "math"

type Page struct {
	Items     interface{} `json:"items"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	Total     int64       `json:"total"`
	TotalPage int         `json:"total_page"`
}

func NewPage(total int64, page int, pageSize int, items interface{}) Page {
	totalPage := math.Ceil(float64(total) / float64(pageSize))
	return Page{
		Items:     items,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: int(totalPage),
	}
}
