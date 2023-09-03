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

const (
	FuncZaoBao     = "zao_bao"     // 每日早报
	FuncHeadLine   = "headline"    // 今日头条
	FuncWeibo      = "weibo_hot"   // 微博热搜
	FuncMidJourney = "mid_journey" // MJ 绘画
)

var InnerFunctions = []Function{
	{
		Name:        FuncZaoBao,
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
		Name:        FuncWeibo,
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
		Name:        FuncHeadLine,
		Description: "今日头条，给用户推荐当天的头条新闻，周榜热文",
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
		Name:        FuncMidJourney,
		Description: "AI 绘画工具，使用 MJ MidJourney API 进行 AI 绘画",
		Parameters: Parameters{
			Type: "object",
			Properties: map[string]Property{
				"prompt": {
					Type:        "string",
					Description: "绘画内容描述，提示词，如果该参数中有中文的话，则需要翻译成英文",
				},
				"ar": {
					Type:        "string",
					Description: "图片长宽比，默认值 16:9",
				},
				"niji": {
					Type:        "string",
					Description: "动漫模型版本，默认值空",
				},
				"v": {
					Type:        "string",
					Description: "模型版本，默认值: 5.2",
				},
			},
			Required: []string{},
		},
	},
}
