package service

import (
	"encoding/json"
	"fmt"
	"geekai/store/model"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

type WechatService struct {
	DB *gorm.DB
}

func NewWechatService(db *gorm.DB) *WechatService {
	return &WechatService{DB: db}
}

type WechatConfig struct {
	WechatAppID    string `json:"wechat_app_id"`
	WechatSecret   string `json:"wechat_secret"`
	WechatToken    string `json:"wechat_token"`
	WechatAesKey   string `json:"wechat_aes_key"`
	WechatCallback string `json:"wechat_callback"`
}

// GetWechatConfig 从数据库中获取微信配置
func (s *WechatService) GetWechatConfig() (WechatConfig, error) {
	var config model.Config
	res := s.DB.Where("marker = ?", "wechat").First(&config)
	if res.Error != nil {
		return WechatConfig{}, res.Error
	}

	var wechatConfig WechatConfig
	err := json.Unmarshal([]byte(config.Config), &wechatConfig)
	if err != nil {
		return WechatConfig{}, err
	}

	return wechatConfig, nil
}

// GetOpenIDByCode 根据 code 获取 openid 和 access_token
func (s *WechatService) GetOpenIDByCode(code string) (string, string, error) {
	wechatConfig, err := s.GetWechatConfig()
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wechatConfig.WechatAppID, wechatConfig.WechatSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
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
func (s *WechatService) GetUserInfo(accessToken string, openID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if errMsg, ok := result["errmsg"].(string); ok {
		return nil, fmt.Errorf("微信 API 错误: %s", errMsg)
	}

	return result, nil
}
