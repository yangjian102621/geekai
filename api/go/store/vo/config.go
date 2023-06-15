package vo

import "chatplus/core/types"

type Config struct {
	Id           uint               `json:"id"`
	Key          string             `json:"key"`
	ChatConfig   types.ChatConfig   `json:"chat_config"`
	SystemConfig types.SystemConfig `json:"system_config"`
}
