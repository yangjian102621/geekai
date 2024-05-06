package vo

import "geekai/core/types"

type Config struct {
	Id           uint               `json:"id"`
	Key          string             `json:"key"`
	SystemConfig types.SystemConfig `json:"system_config"`
}
