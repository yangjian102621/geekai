package jimeng

import "time"

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	ReqKey string `json:"req_key"`
	// 文生图参数
	Prompt     string  `json:"prompt,omitempty"`
	Seed       int64   `json:"seed,omitempty"`
	Scale      float64 `json:"scale,omitempty"`
	Width      int     `json:"width,omitempty"`
	Height     int     `json:"height,omitempty"`
	UsePreLLM  bool    `json:"use_pre_llm,omitempty"`
	// 图生图参数
	ImageInput        string  `json:"image_input,omitempty"`
	ImageUrls         []string `json:"image_urls,omitempty"`
	BinaryDataBase64  []string `json:"binary_data_base64,omitempty"`
	Gpen              float64 `json:"gpen,omitempty"`
	Skin              float64 `json:"skin,omitempty"`
	SkinUnifi         float64 `json:"skin_unifi,omitempty"`
	GenMode           string  `json:"gen_mode,omitempty"`
	// 图像编辑参数
	// 图像特效参数
	ImageInput1 string `json:"image_input1,omitempty"`
	TemplateId  string `json:"template_id,omitempty"`
	// 视频生成参数
	AspectRatio string `json:"aspect_ratio,omitempty"`
}

// SubmitTaskResponse 提交任务响应
type SubmitTaskResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RequestId   string `json:"request_id"`
	Status      int    `json:"status"`
	TimeElapsed string `json:"time_elapsed"`
	Data        struct {
		TaskId string `json:"task_id"`
	} `json:"data"`
}

// QueryTaskRequest 查询任务请求
type QueryTaskRequest struct {
	ReqKey  string `json:"req_key"`
	TaskId  string `json:"task_id"`
	ReqJson string `json:"req_json,omitempty"`
}

// QueryTaskResponse 查询任务响应
type QueryTaskResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RequestId   string `json:"request_id"`
	Status      int    `json:"status"`
	TimeElapsed string `json:"time_elapsed"`
	Data        struct {
		AlgorithmBaseResp struct {
			StatusCode    int    `json:"status_code"`
			StatusMessage string `json:"status_message"`
		} `json:"algorithm_base_resp"`
		BinaryDataBase64 []string `json:"binary_data_base64"`
		ImageUrls        []string `json:"image_urls"`
		VideoUrl         string   `json:"video_url"`
		RespData         string   `json:"resp_data"`
		Status           string   `json:"status"`
		LlmResult        string   `json:"llm_result"`
		PeResult         string   `json:"pe_result"`
		PredictTagsResult string  `json:"predict_tags_result"`
		RephraserResult  string   `json:"rephraser_result"`
		VlmResult        string   `json:"vlm_result"`
		InferCtx         interface{} `json:"infer_ctx"`
	} `json:"data"`
}

// TaskStatus 任务状态
const (
	TaskStatusInQueue    = "in_queue"    // 任务已提交
	TaskStatusGenerating = "generating"  // 任务处理中
	TaskStatusDone       = "done"        // 处理完成
	TaskStatusNotFound   = "not_found"   // 任务未找到
	TaskStatusExpired    = "expired"     // 任务已过期
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Type       string                 `json:"type"`
	Prompt     string                 `json:"prompt"`
	Params     map[string]interface{} `json:"params"`
	ReqKey     string                 `json:"req_key"`
	ImageUrls  []string               `json:"image_urls,omitempty"`
	Power      int                    `json:"power,omitempty"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	Id         uint      `json:"id"`
	UserId     uint      `json:"user_id"`
	TaskId     string    `json:"task_id"`
	Type       string    `json:"type"`
	ReqKey     string    `json:"req_key"`
	Prompt     string    `json:"prompt"`
	TaskParams string    `json:"task_params"`
	ImgURL     string    `json:"img_url"`
	VideoURL   string    `json:"video_url"`
	Progress   int       `json:"progress"`
	Status     string    `json:"status"`
	ErrMsg     string    `json:"err_msg"`
	Power      int       `json:"power"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// LogoInfo 水印信息
type LogoInfo struct {
	AddLogo         bool    `json:"add_logo"`
	Position        int     `json:"position"`
	Language        int     `json:"language"`
	Opacity         float64 `json:"opacity"`
	LogoTextContent string  `json:"logo_text_content"`
}

// ReqJsonConfig 查询配置
type ReqJsonConfig struct {
	ReturnUrl bool      `json:"return_url"`
	LogoInfo  *LogoInfo `json:"logo_info,omitempty"`
}

// ImageEffectTemplate 图像特效模板
const (
	TemplateIdFelt3DPolaroid             = "felt_3d_polaroid"                      // 毛毡3d拍立得风格
	TemplateIdMyWorld                    = "my_world"                              // 像素世界风
	TemplateIdMyWorldUniversal           = "my_world_universal"                    // 像素世界-万物通用版
	TemplateIdPlasticBubbleFigure        = "plastic_bubble_figure"                 // 盲盒玩偶风
	TemplateIdPlasticBubbleFigureCartoon = "plastic_bubble_figure_cartoon_text"    // 塑料泡罩人偶-文字卡头版
	TemplateIdFurryDreamDoll             = "furry_dream_doll"                      // 毛绒玩偶风
	TemplateIdMicroLandscapeMiniWorld    = "micro_landscape_mini_world"            // 迷你世界玩偶风
	TemplateIdMicroLandscapeProfessional = "micro_landscape_mini_world_professional" // 微型景观小世界-职业版
	TemplateIdAcrylicOrnaments           = "acrylic_ornaments"                     // 亚克力挂饰
	TemplateIdFeltKeychain               = "felt_keychain"                         // 毛毡钥匙扣
	TemplateIdLofiPixelCharacter         = "lofi_pixel_character_mini_card"        // Lofi像素人物小卡
	TemplateIdAngelFigurine              = "angel_figurine"                        // 天使形象手办
	TemplateIdLyingInFluffyBelly         = "lying_in_fluffy_belly"                 // 躺在毛茸茸肚皮里
	TemplateIdGlassBall                  = "glass_ball"                           // 玻璃球
)

// AspectRatio 视频宽高比
const (
	AspectRatio16_9 = "16:9"  // 1280×720
	AspectRatio9_16 = "9:16"  // 720×1280
	AspectRatio1_1  = "1:1"   // 960×960
	AspectRatio4_3  = "4:3"   // 960×720
	AspectRatio3_4  = "3:4"   // 720×960
	AspectRatio21_9 = "21:9"  // 1680×720
	AspectRatio9_21 = "9:21"  // 720×1680
)

// GenMode 生成模式
const (
	GenModeCreative       = "creative"       // 提示词模式
	GenModeReference      = "reference"      // 全参考模式
	GenModeReferenceChar  = "reference_char" // 人物参考模式
)