package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/utils"
	"time"

	"gorm.io/gorm"
)

// WxGzhService 微信公众号服务
type WxGzhService struct {
	config types.WxGzhConfig
	DB     *gorm.DB
}

func (s *WxGzhService) UpdateConfig(config types.WxGzhConfig) {
	s.config = config
}

func (s *WxGzhService) GetConfig() types.WxGzhConfig {
	return s.config
}

func (s *WxGzhService) SetConfig(config types.WxGzhConfig) {
	s.config = config
}

func NewWxGzhService(config *types.SystemConfig, db *gorm.DB) *WxGzhService {
	return &WxGzhService{config: config.WxGzh, DB: db}
}

// GetOpenIDByCode 根据 code 获取 openid 和 access_token
func (s *WxGzhService) GetOpenIDByCode(code string) (string, string, error) {

	var config model.Config
	s.DB.Where("name", types.ConfigKeyWxGzh).First(&config)

	var value map[string]any
	err := utils.JsonDecode(config.Value, &value)

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		value["app_id"], value["secret"], code)

	body, status, err := utils.FetchURLBytes(context.Background(), url, "", 30*time.Second, 2, 2<<20)
	if err != nil {
		return "", "", fmt.Errorf("wx get openid failed: status=%d: %w", status, err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", err
	}

	if openID, ok := result["openid"].(string); ok {
		if accessToken, ok := result["access_token"].(string); ok {
			return openID, accessToken, nil
		}
	}

	if errMsg, ok := result["errmsg"].(string); ok {
		return "", "", fmt.Errorf("微信 API 错误: %s", errMsg)
	}

	return "", "", fmt.Errorf("获取 openid 和 access_token 失败: %s", string(body))
}

// GetUserInfo 获取微信用户昵称和头像
func (s *WxGzhService) GetUserInfo(accessToken string, openID string) (map[string]any, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID)

	body, status, err := utils.FetchURLBytes(context.Background(), url, "", 30*time.Second, 2, 2<<20)
	if err != nil {
		return nil, fmt.Errorf("wx get userinfo failed: status=%d: %w", status, err)
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if errMsg, ok := result["errmsg"].(string); ok {
		return nil, fmt.Errorf("微信 API 错误: %s", errMsg)
	}

	return result, nil
}
