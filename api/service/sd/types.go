package sd

import logger2 "chatplus/logger"

var logger = logger2.GetLogger()

type TaskInfo struct {
	UserId      uint          `json:"user_id"`
	SessionId   string        `json:"session_id"`
	JobId       int           `json:"job_id"`
	TaskId      string        `json:"task_id"`
	Data        []interface{} `json:"data"`
	EventData   interface{}   `json:"event_data"`
	FnIndex     int           `json:"fn_index"`
	SessionHash string        `json:"session_hash"`
}

type CBReq struct {
	UserId    uint
	SessionId string
	JobId     int
	TaskId    string
	ImageName string
	ImageData string
	Progress  int
	Seed      int64
	Success   bool
	Message   string
}

var ParamKeys = map[string]int{
	"task_id":         0,
	"prompt":          1,
	"negative_prompt": 2,
	"steps":           4,
	"sampler":         5,
	"face_fix":        7, // 面部修复
	"cfg_scale":       8,
	"seed":            27,
	"height":          10,
	"width":           9,
	"hd_fix":          11,
	"hd_redraw_rate":  12, //高清修复重绘幅度
	"hd_scale":        13, // 高清修复放大倍数
	"hd_scale_alg":    14, // 高清修复放大算法
	"hd_sample_num":   15, // 高清修复采样次数
}
