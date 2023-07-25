package function

import "chatplus/core/types"

type Function interface {
	Invoke(...interface{}) (string, error)
	Name() string
}

type resVo struct {
	Code    types.BizCode `json:"code"`
	Message string        `json:"message"`
	Data    struct {
		Title     string     `json:"title"`
		UpdatedAt string     `json:"updated_at"`
		Items     []dataItem `json:"items"`
	} `json:"data"`
}

type dataItem struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Remark string `json:"remark"`
}
