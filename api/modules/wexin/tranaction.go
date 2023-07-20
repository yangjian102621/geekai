package wexin

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Message 转账消息
type Message struct {
	XMLName xml.Name `xml:"msg"`
	AppMsg  struct {
		Des string `xml:"des"`
		Url string `xml:"url"`
	} `xml:"appmsg"`
}

// Transaction 解析后的交易信息
type Transaction struct {
	TransId string  `json:"trans_id"` // 微信转账交易 ID
	Amount  float64 `json:"amount"`   // 微信转账交易金额
	Remark  string  `json:"remark"`   // 转账备注
}

// 解析微信转账消息
func parseTransactionMessage(xmlData string) (*Message, error) {
	var msg Message
	if err := xml.Unmarshal([]byte(xmlData), &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// 导出交易信息
func extractTransaction(message *Message) Transaction {
	var tx = Transaction{}
	// 导出交易金额和备注
	lines := strings.Split(message.AppMsg.Des, "\n")
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
	index := strings.Index(message.AppMsg.Url, "trans_id=")
	if index != -1 {
		end := strings.LastIndex(message.AppMsg.Url, "&")
		tx.TransId = strings.TrimSpace(message.AppMsg.Url[index+9 : end])
	}
	return tx
}
