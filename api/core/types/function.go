package types

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

var InnerFunctions = []Function{
	{
		Name:        "zao_bao",
		Description: "每日早报，获取当天全球的热门新闻事件列表",
		Parameters: Parameters{

			Type: "object",
			Properties: map[string]Property{
				"text": {
					Type:        "string",
					Description: "",
				},
			},
			Required: []string{},
		},
	},
	{
		Name:        "weibo_hot",
		Description: "新浪微博热搜榜，微博当日热搜榜单",
		Parameters: Parameters{
			Type: "object",
			Properties: map[string]Property{
				"text": {
					Type:        "string",
					Description: "",
				},
			},
			Required: []string{},
		},
	},

	{
		Name:        "zhihu_top",
		Description: "知乎热榜,知乎当日话题讨论榜单",
		Parameters: Parameters{
			Type: "object",
			Properties: map[string]Property{
				"text": {
					Type:        "string",
					Description: "",
				},
			},
			Required: []string{},
		},
	},
}

var FunctionNameMap = map[string]string{
	"zao_bao":   "每日早报",
	"weibo_hot": "微博热搜",
	"zhihu_top": "知乎热榜",
}
