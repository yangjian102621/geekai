package wx

import (
	"encoding/xml"
	"net/url"
	"strconv"
	"strings"
)

// Message 转账消息
type Message struct {
	Des string
	Url string
}

// Transaction 解析后的交易信息
type Transaction struct {
	TransId string  `json:"trans_id"` // 微信转账交易 ID
	Amount  float64 `json:"amount"`   // 微信转账交易金额
	Remark  string  `json:"remark"`   // 转账备注
}

// 解析微信转账消息
func parseTransactionMessage(xmlData string) *Message {
	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	message := Message{}
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			var value string
			if se.Name.Local == "des" && message.Des == "" {
				if err := decoder.DecodeElement(&value, &se); err == nil {
					message.Des = strings.TrimSpace(value)
				}
				break
			}
			if se.Name.Local == "weapp_path" || se.Name.Local == "url" {
				if err := decoder.DecodeElement(&value, &se); err == nil {
					if strings.Contains(value, "?trans_id=") || strings.Contains(value, "?id=") {
						message.Url = value
					}
				}
				break
			}
		}
	}

	// 兼容旧版消息记录
	if message.Url == "" {
		var msg struct {
			XMLName xml.Name `xml:"msg"`
			AppMsg  struct {
				Des string `xml:"des"`
				Url string `xml:"url"`
			} `xml:"appmsg"`
		}
		if err := xml.Unmarshal([]byte(xmlData), &msg); err == nil {
			message.Url = msg.AppMsg.Url
		}
	}
	return &message
}

// 导出交易信息
func extractTransaction(message *Message) Transaction {
	var tx = Transaction{}
	// 导出交易金额和备注
	lines := strings.Split(message.Des, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// 解析收款金额
		prefix := "收款金额￥"
		if strings.HasPrefix(line, prefix) {
			if value, err := strconv.ParseFloat(line[len(prefix):], 64); err == nil {
				tx.Amount = value
				continue
			}
		}
		// 解析收款备注
		prefix = "付款方备注"
		if strings.HasPrefix(line, prefix) {
			tx.Remark = line[len(prefix):]
			break
		}
	}

	// 解析交易 ID
	parse, err := url.Parse(message.Url)
	if err == nil {
		tx.TransId = parse.Query().Get("id")
		if tx.TransId == "" {
			tx.TransId = parse.Query().Get("trans_id")
		}
	}

	return tx
}
