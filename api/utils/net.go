package utils

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

var logger = logger2.GetLogger()

// ReplyChunkMessage 回复客户片段端消息
func ReplyChunkMessage(client *types.WsClient, message interface{}) {
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}

// ReplyMessage 回复客户端一条完整的消息
func ReplyMessage(ws *types.WsClient, message interface{}) {
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: message})
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})
}

func DownloadImage(imageURL string, proxy string) ([]byte, error) {
	var client *http.Client
	if proxy == "" {
		client = http.DefaultClient
	} else {
		proxyURL, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}
	request, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}
