package vo

type ChatApp struct {
	BaseVo
	Tid          uint   `json:"tid"`
	Name         string `json:"name"`          // 角色名称
	UserId       uint   `json:"user_id"`       // 所属用户 ID，0 表示系统内置
	SystemPrompt string `json:"system_prompt,omitempty"` // 系统提示词（列表接口对内置智能体不下发）
	HelloMsg     string `json:"hello_msg"`     // 打招呼的消息
	Icon         string `json:"icon"`          // 角色聊天图标
	Enable       bool   `json:"enable"`        // 是否启用
	SortNum      int    `json:"sort"`          // 排序
	ModelId      uint   `json:"model_id"`      // 绑定模型 ID
	ModelName    string `json:"model_name"`    // 模型名称
	TypeName     string `json:"type_name"`     // 分类名称
}
