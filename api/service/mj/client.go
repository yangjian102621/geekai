package mj

import (
	"chatplus/core/types"
	"fmt"
	"github.com/imroc/req/v3"
	"time"
)

// MidJourney client

type Client struct {
	client *req.Client
	config *types.MidJourneyConfig
}

func NewClient(config *types.AppConfig) *Client {
	client := req.C().SetTimeout(10 * time.Second)
	// set proxy URL
	if config.ProxyURL != "" {
		client.SetProxyURL(config.ProxyURL)
	}
	return &Client{client: client, config: &config.MjConfig}
}

func (c *Client) Imagine(prompt string) error {
	interactionsReq := &InteractionsRequest{
		Type:          2,
		ApplicationID: ApplicationID,
		GuildID:       c.config.GuildId,
		ChannelID:     c.config.ChanelId,
		SessionID:     SessionID,
		Data: map[string]any{
			"version": "1166847114203123795",
			"id":      "938956540159881230",
			"name":    "imagine",
			"type":    "1",
			"options": []map[string]any{
				{
					"type":  3,
					"name":  "prompt",
					"value": prompt,
				},
			},
			"application_command": map[string]any{
				"id":                         "938956540159881230",
				"application_id":             ApplicationID,
				"version":                    "1118961510123847772",
				"default_permission":         true,
				"default_member_permissions": nil,
				"type":                       1,
				"nsfw":                       false,
				"name":                       "imagine",
				"description":                "Create images with Midjourney",
				"dm_permission":              true,
				"options": []map[string]any{
					{
						"type":        3,
						"name":        "prompt",
						"description": "The prompt to imagine",
						"required":    true,
					},
				},
				"attachments": []any{},
			},
		},
	}

	url := "https://discord.com/api/v9/interactions"
	r, err := c.client.R().SetHeader("Authorization", c.config.UserToken).
		SetHeader("Content-Type", "application/json").
		SetBody(interactionsReq).
		Post(url)

	if err != nil || r.IsErrorState() {
		return fmt.Errorf("error with http request: %w%v", err, r.Err)
	}

	return nil
}

// Upscale 放大指定的图片
func (c *Client) Upscale(index int, messageId string, hash string) error {
	flags := 0
	interactionsReq := &InteractionsRequest{
		Type:          3,
		ApplicationID: ApplicationID,
		GuildID:       c.config.GuildId,
		ChannelID:     c.config.ChanelId,
		MessageFlags:  &flags,
		MessageID:     &messageId,
		SessionID:     SessionID,
		Data: map[string]any{
			"component_type": 2,
			"custom_id":      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, hash),
		},
		Nonce: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	url := "https://discord.com/api/v9/interactions"
	var res InteractionsResult
	r, err := c.client.R().SetHeader("Authorization", c.config.UserToken).
		SetHeader("Content-Type", "application/json").
		SetBody(interactionsReq).
		SetErrorResult(&res).
		Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("error with http request: %v%v%v", err, r.Err, res.Message)
	}

	return nil
}

// Variation  以指定的图片的视角进行变换再创作，注意需要在对应的频道中关闭 Remix 变换，否则 Variation 指令将不会生效
func (c *Client) Variation(index int, messageId string, hash string) error {
	flags := 0
	interactionsReq := &InteractionsRequest{
		Type:          3,
		ApplicationID: ApplicationID,
		GuildID:       c.config.GuildId,
		ChannelID:     c.config.ChanelId,
		MessageFlags:  &flags,
		MessageID:     &messageId,
		SessionID:     SessionID,
		Data: map[string]any{
			"component_type": 2,
			"custom_id":      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, hash),
		},
		Nonce: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	url := "https://discord.com/api/v9/interactions"
	var res InteractionsResult
	r, err := c.client.R().SetHeader("Authorization", c.config.UserToken).
		SetHeader("Content-Type", "application/json").
		SetBody(interactionsReq).
		SetErrorResult(&res).
		Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("error with http request: %v%v%v", err, r.Err, res.Message)
	}

	return nil
}
