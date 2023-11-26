package fun

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/service/oss"
	"gorm.io/gorm"
)

type Function interface {
	Invoke(map[string]interface{}) (string, error)
	Name() string
}

var logger = logger2.GetLogger()

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

func NewFunctions(config *types.AppConfig, db *gorm.DB, manager *oss.UploaderManager) map[string]Function {
	return map[string]Function{
		types.FuncZaoBao:   NewZaoBao(config.ApiConfig),
		types.FuncWeibo:    NewWeiboHot(config.ApiConfig),
		types.FuncHeadLine: NewHeadLines(config.ApiConfig),
		types.FuncImage:    NewImageFunc(db, manager, config),
	}
}
