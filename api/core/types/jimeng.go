package types

// JimengConfig 即梦AI配置
type JimengConfig struct {
	AccessKey string      `json:"access_key"`
	SecretKey string      `json:"secret_key"`
	Power     JimengPower `json:"power"`
}

// JimengPower 即梦AI算力配置
type JimengPower struct {
	Image          int `json:"image"`           // 图片生成算力，单位：积分/张
	Video          int `json:"video"`           // 视频生成算力，单位：积分/秒
	VirtualHuman   int `json:"virtual_human"`   // 数字人视频生成算力，单位：积分/秒
	ActionTransfer int `json:"action_transfer"` // 视频动作迁移算力，单位：积分/秒
}
