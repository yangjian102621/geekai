package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"io"
	"net/http"
	"net/url"
)

var logger = logger2.GetLogger()

// SendMsg 回复客户片段端消息
func SendMsg(client *types.WsClient, message types.ReplyMessage) {
	message.ClientId = client.Id
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

// SendAndFlush 回复客户端一条完整的消息
func SendAndFlush(ws *types.WsClient, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeText, Body: message})
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeEnd})
}

func SendChunkMsg(ws *types.WsClient, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeText, Body: message})
}

// SendErrMsg 向客户端发送错误消息
func SendErrMsg(ws *types.WsClient, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeErr, Body: message})
}

func SendChannelMsg(ws *types.WsClient, channel types.WsChannel, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: channel, Type: types.MsgTypeText, Body: message})
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func GetBaseURL(strURL string) string {
	u, err := url.Parse(strURL)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}
