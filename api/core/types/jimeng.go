package types

// JimengConfig 即梦AI配置
type JimengConfig struct {
	AccessKey string      `json:"access_key"`
	SecretKey string      `json:"secret_key"`
	Power     JimengPower `json:"power"`
}

// JimengPower 即梦AI算力配置
type JimengPower struct {
	TextToImage  int `json:"text_to_image"`
	ImageToImage int `json:"image_to_image"`
	ImageEdit    int `json:"image_edit"`
	ImageEffects int `json:"image_effects"`
	TextToVideo  int `json:"text_to_video"`
	ImageToVideo int `json:"image_to_video"`
}
