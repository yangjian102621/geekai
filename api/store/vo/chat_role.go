package vo

import "chatplus/core/types"

type ChatRole struct {
	BaseVo
	Key      string          `json:"key"`       // 角色唯一标识
	Name     string          `json:"name"`      // 角色名称
	Context  []types.Message `json:"context"`   // 角色语料信息
	HelloMsg string          `json:"hello_msg"` // 打招呼的消息
	Icon     string          `json:"icon"`      // 角色聊天图标
	Enable   bool            `json:"enable"`    // 是否启用被启用
	SortNum  int             `json:"sort"`      // 排序
}
