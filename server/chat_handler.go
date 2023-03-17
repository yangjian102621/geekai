package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"openai/types"
	"time"
)

func (s *Server) Chat(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Infof("New websocket connected, IP: %s", c.Request.RemoteAddr)
	client := NewWsClient(ws)
	go func() {
		for {
			_, message, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				return
			}

			logger.Info(string(message))
			for {
				err = client.Send([]byte("H"))
				time.Sleep(time.Second)
			}
			// TODO: 根据会话请求，传入不同的用户 ID
			//err = s.sendMessage("test", string(message), client)
			//if err != nil {
			//	logger.Error(err)
			//}
		}
	}()
}

func (s *Server) sendMessage(userId string, text string, ws Client) error {
	var r = types.ApiRequest{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.9,
		MaxTokens:   1024,
		Stream:      true,
	}
	var history []types.Message
	if v, ok := s.History[userId]; ok {
		history = v
	} else {
		history = make([]types.Message, 0)
	}
	r.Messages = append(history, types.Message{
		Role:    "user",
		Content: text,
	})

	logger.Info("上下文历史消息:%+v", s.History[userId])
	requestBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, s.Config.OpenAi.ApiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// TODO: API KEY 负载均衡
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.Config.OpenAi.ApiKey[0]))

	uri := url.URL{}
	proxy, _ := uri.Parse(s.Config.ProxyURL)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	response, err := client.Do(request)
	var retryCount = 3
	for err != nil {
		if retryCount <= 0 {
			return err
		}
		response, err = client.Do(request)
		retryCount--
	}

	var message = types.Message{}
	var contents = make([]string, 0)
	var responseBody = types.ApiResponse{}

	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		}

		if line == "" {
			break
		} else if len(line) < 20 {
			continue
		}

		err = json.Unmarshal([]byte(line[6:]), &responseBody)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 初始化 role
		if responseBody.Choices[0].Delta.Role != "" && message.Role == "" {
			message.Role = responseBody.Choices[0].Delta.Role
			continue
		} else {
			contents = append(contents, responseBody.Choices[0].Delta.Content)
		}

		err = ws.(*WsClient).Send([]byte(responseBody.Choices[0].Delta.Content))
		if err != nil {
			logger.Error(err)
		}
		fmt.Print(responseBody.Choices[0].Delta.Content)
		if responseBody.Choices[0].FinishReason != "" {
			fmt.Println()
			break
		}
	}

	// 追加历史消息
	history = append(history, message)
	return nil
}
