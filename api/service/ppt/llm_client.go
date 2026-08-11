package ppt

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/log"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

var logger = log.GetLogger()

// slidePlan LLM 输出的分镜结构
type slidePlan struct {
	SlideIndex  int      `json:"slide_index"`
	Theme       string   `json:"theme"`
	Title       string   `json:"title"`
	Points      []string `json:"points"`
	ImagePrompt string   `json:"image_prompt"`
}

// systemPrompt 图文并茂幻灯片：每页有插图且画面上含与主题一致的文字，文字量与生成模式相关
const systemPrompt = `
# Role

你是一位顶级的专业演示文稿（PPT）策划专家和 AI 图像提示词（Prompt）工程师。任务是根据用户提供的「内容大纲」或「设计要求」，生成一套逻辑清晰、视觉风格高度统一的幻灯片分镜数据。目标是生成**图文并茂**的幻灯片：每页既有文字又有插图，图片上直接呈现与本页内容一致的文字，而不是留白让用户后加文字。

# Rules

1. 全局风格锚定：根据大纲推断或遵循用户要求的全局视觉风格。所有配图必须严格遵循此风格。
2. 结构化拆解：合理拆分为多张幻灯片，单页最多 3-4 个简短要点。
3. 视觉转译（图文并茂）：
   - 为每页构思的 image_prompt 既要描述**插图画面**，也要明确**画面上应出现的文字**（如本页标题、要点或短句），与当页 theme、title、points 内容一致。不要描述留白或“用于排版文字的空间”。
   - 图片中出现的所有文字必须使用与 theme、title、points **相同的语言**（即本次请求指定的输出语言）。
   - 图片上文字的量由「生成模式」决定（见用户输入中的模式说明）：
     - **详细演示文稿**：图片上的说明文字可适当多一些，如本页要点、一两句说明。
     - **演示用幻灯片**：图片上的文字尽量精简，如仅主标题或少量关键词，便于演讲时配合口述。
   - image_prompt 须包含前缀「[全局风格描述]」，并清晰描述画面中的插图与文字内容（含具体要出现的文字及其语言），不要包含“不要在图片中生成任何文字”的约束。
4. 严格输出合法的纯 JSON 数组：
[
  {"slide_index": 1, "theme": "...", "title": "...", "points": ["..."], "image_prompt": "..."}
]
禁止输出任何 Markdown 标记或多余文本。`

// notebookSystemPrompt 文档提炼：NotebookLM 风格输出 PPT 可用大纲文本（纯文本/Markdown）。
const notebookSystemPrompt = `
# Role

你是一位“NotebookLM 风格”的专业文档理解与提炼助手。任务是基于用户提供的「原始文档文本」和「设计要求」，提炼出可用于制作 PPT 的结构化大纲内容。

# Output Requirements

1. 输出必须是纯文本/Markdown（允许使用标题与列表），禁止输出任何 JSON。
2. 禁止输出代码块（不要出现代码块语法）。
3. 不要输出解释过程、不要复述提示词。
4. 大纲必须是“内容大纲/要点”，用于后续继续拆分成幻灯片，而不是直接输出最终幻灯片分镜。

# Rules

1. 文档优先：尽可能从原始文档中提取信息与措辞；若文档缺失关键点，则给出合理补全的“建议方向”，并明确标注为“（建议）”。
2. 贴合设计要求：根据设计要求调整大纲的语气、侧重点、术语风格，使内容更符合目标受众与整体风格。
3. 结构清晰：使用分层标题（例如：# 总主题、## 模块/章节、### 要点），并为每个模块给出 2-4 个要点句（可直接用于 PPT 每页标题/要点）。
4. 语言一致：所有输出语言必须与本次请求指定的语言一致。
`

// LLMClient 分镜 LLM 客户端
type LLMClient struct {
	httpClient *req.Client
}

func NewLLMClient() *LLMClient {
	return &LLMClient{
		httpClient: req.C().SetTimeout(2 * time.Minute),
	}
}

// GenerateSlides 调用大模型生成分镜列表。language 约束输出语言，mode 约束图中文字量，maxPages 约束恰好生成 N 页。
func (c *LLMClient) GenerateSlides(ctx context.Context, cfg types.PPTConfig, content, prompt, language, mode string, maxPages int) ([]slidePlan, error) {
	if cfg.OutlineLLMApiURL == "" {
		return nil, fmt.Errorf("outline LLM api url is empty")
	}
	if cfg.OutlineLLMApiKey == "" {
		return nil, fmt.Errorf("outline LLM api key is empty")
	}
	if maxPages <= 0 {
		maxPages = 10
	}
	if mode != "detailed" && mode != "slides" {
		mode = "slides"
	}

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// 动态 system：加入页数约束
	systemContent := systemPrompt + fmt.Sprintf("\n\n# 页数约束\n请将内容拆分为恰好 %d 页的幻灯片分镜，保证逻辑完整、故事线连贯，不要多也不要少。输出 JSON 数组长度必须为 %d。", maxPages, maxPages)

	// 组装用户输入
	userContent := fmt.Sprintf("下面是用户提供的演示文稿大纲内容：\n\n%s", content)
	if prompt != "" {
		userContent = fmt.Sprintf("%s\n\n额外的设计要求：%s", userContent, prompt)
	}
	if language != "" {
		langHint := "中文"
		if language == "en" || language == "en-US" {
			langHint = "英文"
		} else if language == "zh-CN" || language == "zh" {
			langHint = "中文"
		} else {
			langHint = "语言代码 " + language + " 对应的语言"
		}
		userContent = fmt.Sprintf("%s\n\n请用%s输出所有分镜内容（theme、title、points、image_prompt 等均使用该语言；图片中出现的文字也必须是%s）。", userContent, langHint, langHint)
	}
	modeHint := "演示用幻灯片"
	if mode == "detailed" {
		modeHint = "详细演示文稿"
	}
	userContent = fmt.Sprintf("%s\n\n本次生成模式为：%s。请按上述规则控制每页插图中文字的量。", userContent, modeHint)

	modelName := cfg.OutlineLLMModel
	if modelName == "" {
		modelName = "gpt-5.2"
	}
	body := map[string]any{
		"model": modelName,
		"messages": []message{
			{Role: "user", Content: systemContent + "\n\n" + userContent},
		},
		"temperature": 0.8,
	}

	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	logger.Infof("generate PPT slides with outline LLM, api: %s", cfg.OutlineLLMApiURL)
	r, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+cfg.OutlineLLMApiKey).
		SetBody(body).
		SetSuccessResult(&respBody).
		Post(cfg.OutlineLLMApiURL)
	if err != nil {
		return nil, fmt.Errorf("request outline LLM failed: %v", err)
	}
	if r.IsErrorState() {
		return nil, fmt.Errorf("outline LLM returned error status: %s", r.Status)
	}

	if len(respBody.Choices) == 0 {
		return nil, fmt.Errorf("outline LLM returned empty choices")
	}

	contentStr := respBody.Choices[0].Message.Content
	var plans []slidePlan
	if err := json.Unmarshal([]byte(contentStr), &plans); err != nil {
		return nil, fmt.Errorf("parse outline LLM json failed: %v, raw: %s", err, contentStr)
	}

	return plans, nil
}

// GenerateNotebookContent 调用文档提炼 LLM，把 rawDocText -> PPT 可用的 content（大纲/结构化要点）。
func (c *LLMClient) GenerateNotebookContent(ctx context.Context, cfg types.PPTConfig, rawDocText, designPrompt, language string) (string, error) {
	if cfg.OutlineLLMApiURL == "" {
		return "", fmt.Errorf("outline LLM api url is empty")
	}
	if cfg.OutlineLLMApiKey == "" {
		return "", fmt.Errorf("outline LLM api key is empty")
	}

	rawDocText = strings.TrimSpace(rawDocText)
	if rawDocText == "" {
		return "", fmt.Errorf("rawDocText is empty")
	}

	// 对超长输入做保守截断，避免请求体过大或上下文溢出。
	// 这里按“字符数”截断，真实 token 仍可能超出，但作为兜底足够。
	const maxChars = 25000
	runes := []rune(rawDocText)
	if len(runes) > maxChars {
		rawDocText = string(runes[:maxChars])
	}

	langHint := "中文"
	if language == "en" || language == "en-US" {
		langHint = "英文"
	} else if language == "zh-CN" || language == "zh" {
		langHint = "中文"
	} else if language != "" {
		langHint = "语言代码 " + language + " 对应的语言"
	}

	maxSlides := cfg.MaxSlidesPerTask
	if maxSlides <= 0 {
		maxSlides = 10
	}

	userContent := fmt.Sprintf("原始文档文本如下（可能很长）：\n\n%s", rawDocText)
	if strings.TrimSpace(designPrompt) != "" {
		userContent = fmt.Sprintf("%s\n\n设计要求（风格/受众/侧重点等）：\n%s", userContent, designPrompt)
	}
	userContent = fmt.Sprintf(
		"%s\n\n请用%s输出 PPT 大纲内容。该大纲应便于拆分为不超过 %d 页的 PPT。",
		userContent,
		langHint,
		maxSlides,
	)

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	systemContent := notebookSystemPrompt + fmt.Sprintf("\n\n# 语言约束\n输出语言：%s。", langHint)

	modelName := cfg.OutlineLLMModel
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

	body := map[string]any{
		"model": modelName,
		"messages": []message{
			{Role: "user", Content: systemContent + "\n\n" + userContent},
		},
		"temperature": 0.4,
	}

	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	logger.Infof("generate PPT content with outline LLM, api: %s", cfg.OutlineLLMApiURL)
	r, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+cfg.OutlineLLMApiKey).
		SetBody(body).
		SetSuccessResult(&respBody).
		Post(cfg.OutlineLLMApiURL)
	if err != nil {
		return "", fmt.Errorf("request outline LLM failed: %v", err)
	}
	if r.IsErrorState() {
		return "", fmt.Errorf("outline LLM returned error status: %s", r.Status)
	}
	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("outline LLM returned empty choices")
	}

	return strings.TrimSpace(respBody.Choices[0].Message.Content), nil
}

// titleSystemPrompt 根据大纲生成用于任务列表的短标题（单行纯文本）
const titleSystemPrompt = `
# Role

你是「演示文稿命名助手」。用户会提供一份 PPT 内容大纲（可能含 Markdown）。请根据大纲主题与受众，生成一个**适合出现在任务列表中的短标题**。

# Output Rules

1. 只输出**一行**纯文本，不要换行、不要编号、不要引号包裹。
2. 长度建议 **20 个字以内**（中文）或 **8 个英文单词以内**；若大纲极长，仍只给概括性标题。
3. 输出语言必须与本次指定的「输出语言」一致。
4. 不要输出“标题：”“Title:”等前缀，不要复述本说明。
`

// GeneratePPTTitle 调用大模型根据 content 生成列表用短标题；失败时由调用方降级。
func (c *LLMClient) GeneratePPTTitle(ctx context.Context, cfg types.PPTConfig, content, language string) (string, error) {
	if cfg.OutlineLLMApiURL == "" {
		return "", fmt.Errorf("outline LLM api url is empty")
	}
	if cfg.OutlineLLMApiKey == "" {
		return "", fmt.Errorf("outline LLM api key is empty")
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return "", fmt.Errorf("content is empty")
	}

	const maxChars = 10000
	runes := []rune(content)
	if len(runes) > maxChars {
		content = string(runes[:maxChars])
	}

	langHint := "中文"
	if language == "en" || language == "en-US" {
		langHint = "英文"
	} else if language == "zh-CN" || language == "zh" {
		langHint = "中文"
	} else if language != "" {
		langHint = "语言代码 " + language + " 对应的语言"
	}

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	userContent := fmt.Sprintf("输出语言：%s。\n\n下面是用户提供的 PPT 大纲内容，请只返回列表标题：\n\n%s", langHint, content)
	systemContent := titleSystemPrompt + fmt.Sprintf("\n\n# 语言约束\n请用%s撰写标题。", langHint)

	modelName := cfg.OutlineLLMModel
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

	body := map[string]any{
		"model": modelName,
		"messages": []message{
			{Role: "user", Content: systemContent + "\n\n" + userContent},
		},
		"temperature": 0.5,
	}

	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	logger.Infof("generate PPT list title with outline LLM, api: %s", cfg.OutlineLLMApiURL)
	r, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+cfg.OutlineLLMApiKey).
		SetBody(body).
		SetSuccessResult(&respBody).
		Post(cfg.OutlineLLMApiURL)
	if err != nil {
		return "", fmt.Errorf("request outline LLM failed: %v", err)
	}
	if r.IsErrorState() {
		return "", fmt.Errorf("outline LLM returned error status: %s", r.Status)
	}
	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("outline LLM returned empty choices")
	}

	raw := strings.TrimSpace(respBody.Choices[0].Message.Content)
	if idx := strings.IndexAny(raw, "\r\n"); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	raw = strings.Trim(raw, `"'「」`)
	return raw, nil
}
