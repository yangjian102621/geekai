package vo

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// PPTSlideImageVersion 单页配图历史中的一条（图片地址 + 文生图/图生图提示）
type PPTSlideImageVersion struct {
	ImageURL string `json:"image_url"`
	Prompt   string `json:"prompt"`
}

// PPTSlideImageVersions 存库 JSON 数组；UnmarshalJSON 兼容旧版仅字符串数组
type PPTSlideImageVersions []PPTSlideImageVersion

// UnmarshalJSON 兼容 ["url1","url2"] 与 [{"image_url":"...","prompt":"..."}]
func (v *PPTSlideImageVersions) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*v = nil
		return nil
	}
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	out := make([]PPTSlideImageVersion, 0, len(raw))
	for _, item := range raw {
		item = bytes.TrimSpace(item)
		if len(item) == 0 || string(item) == "null" {
			continue
		}
		if item[0] == '"' {
			var s string
			if err := json.Unmarshal(item, &s); err != nil {
				return err
			}
			out = append(out, PPTSlideImageVersion{ImageURL: s, Prompt: ""})
			continue
		}
		var o PPTSlideImageVersion
		if err := json.Unmarshal(item, &o); err != nil {
			return err
		}
		out = append(out, o)
	}
	*v = out
	return nil
}

// PPTParams 生成参数（语言、生成模式、页数）
type PPTParams struct {
	Language string `json:"language"` // 输出语言，如 zh-CN, en
	Mode     string `json:"mode"`     // 生成模式：detailed | slides
	Pages    int    `json:"pages"`    // 目标页数
}

// Value 实现 driver.Valuer，序列化为 JSON 存储
func (p PPTParams) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现 sql.Scanner，从数据库读取 JSON
func (p *PPTParams) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal PPTParams value: %v", value)
	}
	return json.Unmarshal(bytes, p)
}

// PPTSlideData 单页幻灯片数据（与 service/ppt SlideData 一致）
type PPTSlideData struct {
	SlideIndex   int                   `json:"slide_index"`
	Theme        string                `json:"theme"`
	Title        string                `json:"title"`
	Points       []string              `json:"points"`
	ImagePrompt  string                `json:"image_prompt"`
	ImageURL     string                `json:"image_url"`
	ImageHistory PPTSlideImageVersions `json:"image_history,omitempty"`
}

// PPTSlides 幻灯片列表（JSON 数组）
type PPTSlides []PPTSlideData

// Value 实现 driver.Valuer
func (s PPTSlides) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner
func (s *PPTSlides) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal PPTSlides value: %v", value)
	}
	return json.Unmarshal(bytes, s)
}
