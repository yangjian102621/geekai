package sd

import logger2 "chatplus/logger"

var logger = logger2.GetLogger()

type TaskInfo struct {
	SessionId   string
	JobId       int
	TaskId      string
	Data        []interface{}
	EventData   interface{}
	FnIndex     int
	SessionHash string
}

type CBReq struct {
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
	"face_fix":        6,
	"cfg_scale":       10,
	"seed":            11,
	"height":          17,
	"width":           18,
	"hd_fix":          19,
	"hd_redraw_rate":  20, //高清修复重绘幅度
	"hd_scale":        21, // 高清修复放大倍数
	"hd_scale_alg":    22, // 高清修复放大算法
	"hd_sample_num":   23, // 高清修复采样次数
}

const Text2ImgParamTemplate = `[
  "task(6sm0b3j17tag2gd)",
  "A beautiful  Chinese girl wearing a cheongsam walks on the bluestone street",
  "",
  [],
  50,
  "Euler a",
  false,
  false,
  1,
  1,
  15,
  null,
  -1,
  0,
  0,
  0,
  false,
  768,
  512,
  false,
  0.7,
  2,
  "ESRGAN_4x",
  10,
  0,
  0,
  "Use same sampler",
  "",
  "",
  [],
  "None",
  null,
  false,
  false,
  "positive",
  "comma",
  0,
  false,
  false,
  "",
  "Seed",
  "",
  [],
  "Nothing",
  "",
  [],
  "Nothing",
  "",
  [],
  true,
  false,
  false,
  false,
  0,
  "Not set",
  true,
  true,
  "",
  "",
  "",
  "",
  "",
  1.3,
  "Not set",
  "Not set",
  1.3,
  "Not set",
  1.3,
  "Not set",
  1.3,
  1.3,
  "Not set",
  1.3,
  "Not set",
  1.3,
  "Not set",
  1.3,
  "Not set",
  1.3,
  "Not set",
  1.3,
  "Not set",
  false,
  "None",
  null,
  false,
  50,
  [],
  "",
  "",
  ""
]`
