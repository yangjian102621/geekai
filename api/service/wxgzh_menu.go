package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"net/url"
	"strings"
	"time"
)

const wxGzhTokenSkew = 300 * time.Second

// GetClientCredentialToken 获取公众号 client_credential access_token（带进程内缓存）。
func (s *WxGzhService) GetClientCredentialToken(ctx context.Context, appID, secret string) (string, error) {
	appID, secret = strings.TrimSpace(appID), strings.TrimSpace(secret)
	if appID == "" || secret == "" {
		return "", fmt.Errorf("微信公众号 AppID 或 AppSecret 未配置")
	}

	now := time.Now()
	s.ccMu.Lock()
	if s.ccForApp == appID && s.ccToken != "" && now.Before(s.ccExpire) {
		tok := s.ccToken
		s.ccMu.Unlock()
		return tok, nil
	}
	s.ccMu.Unlock()

	raw := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		url.QueryEscape(appID), url.QueryEscape(secret),
	)
	body, status, err := utils.FetchURLBytes(ctx, raw, "", 30*time.Second, 2, 2<<20)
	if err != nil {
		return "", fmt.Errorf("获取 access_token 失败: status=%d: %w", status, err)
	}

	var parsed struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("解析 token 响应失败: %w", err)
	}
	if parsed.ErrCode != 0 {
		return "", fmt.Errorf("微信 API 错误(%d): %s", parsed.ErrCode, parsed.ErrMsg)
	}
	if parsed.AccessToken == "" {
		return "", fmt.Errorf("微信未返回 access_token: %s", string(body))
	}

	expireAt := now.Add(time.Duration(parsed.ExpiresIn) * time.Second).Add(-wxGzhTokenSkew)
	if parsed.ExpiresIn <= 0 {
		expireAt = now.Add(7000 * time.Second)
	}

	s.ccMu.Lock()
	s.ccForApp = appID
	s.ccToken = parsed.AccessToken
	s.ccExpire = expireAt
	s.ccMu.Unlock()

	return parsed.AccessToken, nil
}

// PublishCustomMenu 调用 menu/create 全量覆盖公众号菜单。
func (s *WxGzhService) PublishCustomMenu(ctx context.Context, appID, secret string, menu types.WxGzhMenuConfig) error {
	if err := ValidateWxGzhMenu(menu); err != nil {
		return err
	}
	token, err := s.GetClientCredentialToken(ctx, appID, secret)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(menu)
	if err != nil {
		return err
	}
	api := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", url.QueryEscape(token))
	body, status, err := utils.PostURLBytes(ctx, api, "application/json; charset=utf-8", payload, "", 45*time.Second, 2, 2<<20)
	if err != nil {
		return fmt.Errorf("创建菜单请求失败: status=%d: %w", status, err)
	}
	var res struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("解析创建菜单响应失败: %w", err)
	}
	if res.ErrCode != 0 {
		return fmt.Errorf("微信创建菜单失败(%d): %s", res.ErrCode, res.ErrMsg)
	}
	return nil
}

// GetCustomMenuFromWechat 调用 menu/get 拉取线上菜单。
func (s *WxGzhService) GetCustomMenuFromWechat(ctx context.Context, appID, secret string) (types.WxGzhMenuConfig, error) {
	var empty types.WxGzhMenuConfig
	token, err := s.GetClientCredentialToken(ctx, appID, secret)
	if err != nil {
		return empty, err
	}
	api := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s", url.QueryEscape(token))
	body, status, err := utils.FetchURLBytes(ctx, api, "", 30*time.Second, 2, 2<<20)
	if err != nil {
		return empty, fmt.Errorf("获取菜单失败: status=%d: %w", status, err)
	}

	var envelope struct {
		Menu *struct {
			Button []types.WxGzhMenuButton `json:"button"`
		} `json:"menu"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return empty, fmt.Errorf("解析菜单响应失败: %w", err)
	}
	if envelope.ErrCode != 0 {
		return empty, fmt.Errorf("微信获取菜单失败(%d): %s", envelope.ErrCode, envelope.ErrMsg)
	}
	if envelope.Menu == nil || len(envelope.Menu.Button) == 0 {
		return types.WxGzhMenuConfig{Button: nil}, nil
	}
	return types.WxGzhMenuConfig{Button: envelope.Menu.Button}, nil
}

// ValidateWxGzhMenu 校验自定义菜单是否符合微信公众平台规则（view / click / 子菜单）。
func ValidateWxGzhMenu(m types.WxGzhMenuConfig) error {
	if len(m.Button) == 0 {
		return fmt.Errorf("请至少配置一个一级菜单")
	}
	if len(m.Button) > 3 {
		return fmt.Errorf("一级菜单最多 3 个")
	}
	for _, b := range m.Button {
		if err := validateTopButton(b); err != nil {
			return err
		}
	}
	return nil
}

func validateTopButton(b types.WxGzhMenuButton) error {
	if strings.TrimSpace(b.Name) == "" {
		return fmt.Errorf("菜单名称不能为空")
	}
	if len(b.SubButton) > 0 {
		if strings.TrimSpace(b.Type) != "" {
			return fmt.Errorf("含子菜单的一级项不应设置类型")
		}
		if len(b.SubButton) > 5 {
			return fmt.Errorf("子菜单最多 5 项")
		}
		for _, sb := range b.SubButton {
			if err := validateLeafButton(sb); err != nil {
				return err
			}
		}
		return nil
	}
	return validateLeafButton(b)
}

func validateLeafButton(b types.WxGzhMenuButton) error {
	if strings.TrimSpace(b.Name) == "" {
		return fmt.Errorf("菜单名称不能为空")
	}
	if len(b.SubButton) > 0 {
		return fmt.Errorf("仅支持两级菜单，子菜单下不能再嵌套")
	}
	t := strings.ToLower(strings.TrimSpace(b.Type))
	switch t {
	case "view":
		if strings.TrimSpace(b.URL) == "" {
			return fmt.Errorf("跳转链接类型需要填写 URL")
		}
	case "click":
		if strings.TrimSpace(b.Key) == "" {
			return fmt.Errorf("点击事件类型需要填写 Key")
		}
	case "miniprogram":
		if strings.TrimSpace(b.AppID) == "" || strings.TrimSpace(b.PagePath) == "" || strings.TrimSpace(b.URL) == "" {
			return fmt.Errorf("小程序类型需要填写 appid、pagepath 与备用 url")
		}
	default:
		if t == "" {
			return fmt.Errorf("请选择菜单类型")
		}
		return fmt.Errorf("暂不支持的菜单类型: %s", b.Type)
	}
	return nil
}
