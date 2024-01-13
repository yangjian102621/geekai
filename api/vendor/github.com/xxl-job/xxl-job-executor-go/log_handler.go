package xxl

import (
	"encoding/json"
	"net/http"
)

/**
用来日志查询，显示到xxl-job-admin后台
*/

type LogHandler func(req *LogReq) *LogRes

//默认返回
func defaultLogHandler(req *LogReq) *LogRes {
	return &LogRes{Code: SuccessCode, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这是日志默认返回，说明没有设置LogHandler",
		IsEnd:       true,
	}}
}

//请求错误
func reqErrLogHandler(w http.ResponseWriter, req *LogReq, err error) {
	res := &LogRes{Code: FailureCode, Msg: err.Error(), Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   0,
		LogContent:  err.Error(),
		IsEnd:       true,
	}}
	str, _ := json.Marshal(res)
	_, _ = w.Write(str)
}
