package types

type ToolCall struct {
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
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
	FuncZaoBao   = "zao_bao"    // 每日早报
	FuncHeadLine = "headline"   // 今日头条
	FuncWeibo    = "weibo_hot"  // 微博热搜
	FuncImage    = "draw_image" // AI 绘画
)

var InnerFunctions = []Function{
	{
		Name:        FuncZaoBao,
		Description: "每日早报，获取当天新闻事件列表",
		Parameters: Parameters{

			Type:       "object",
			Properties: map[string]Property{},
			Required:   []string{},
		},
	},
	{
		Name:        FuncWeibo,
		Description: "新浪微博热搜榜，微博当日热搜榜单",
		Parameters: Parameters{
			Type:       "object",
			Properties: map[string]Property{},
			Required:   []string{},
		},
	},

	{
		Name:        FuncImage,
		Description: "AI 绘画工具，根据输入的绘图描述用 AI 工具进行绘画",
		Parameters: Parameters{
			Type: "object",
			Properties: map[string]Property{
				"prompt": {
					Type:        "string",
					Description: "提示词，请自动将该参数翻译成英文。",
				},
			},
			Required: []string{},
		},
	},
}
