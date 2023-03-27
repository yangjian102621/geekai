package server

import (
	"encoding/json"
	"openai/types"
	"openai/utils"
)

const (
	TokenPrefix       = "chat/tokens/"
	ChatRolePrefix    = "chat/roles/"
	ChatHistoryPrefix = "chat/history/"
)

var db *utils.LevelDB

func init() {
	leveldb, err := utils.NewLevelDB("data")
	if err != nil {
		panic(err)
	}
	db = leveldb
}

// GetTokens 获取 token 信息
// chat/tokens
func GetTokens() []types.Token {
	items := db.Search(TokenPrefix)
	var tokens = make([]types.Token, 0)
	for _, v := range items {
		var token types.Token
		err := json.Unmarshal([]byte(v), &token)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func PutToken(token types.Token) error {
	key := TokenPrefix + token.Name
	return db.Put(key, token)
}

func RemoveToken(token string) error {
	key := TokenPrefix + token
	return db.Delete(key)
}

// GetChatRoles 获取聊天角色
// chat/roles
func GetChatRoles() map[string]types.ChatRole {
	return nil
}

// GetChatHistory 获取聊天历史记录
// chat/history/{token}/{role}
func GetChatHistory() []types.Message {
	return nil
}
