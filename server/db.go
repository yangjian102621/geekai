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

func GetToken(name string) (types.Token, error) {
	key := TokenPrefix + name
	token, err := db.Get(key)
	if err != nil {
		return types.Token{}, err
	}

	return token.(types.Token), nil
}

func RemoveToken(token string) error {
	key := TokenPrefix + token
	return db.Delete(key)
}

// GetChatRoles 获取聊天角色
// chat/roles
func GetChatRoles() map[string]types.ChatRole {
	items := db.Search(ChatRolePrefix)
	var roles = make(map[string]types.ChatRole)
	for _, v := range items {
		var role types.ChatRole
		err := json.Unmarshal([]byte(v), &role)
		if err != nil {
			continue
		}
		roles[role.Key] = role
	}
	return roles
}

func PutChatRole(role types.ChatRole) error {
	key := ChatRolePrefix + role.Key
	return db.Put(key, role)
}

// GetChatHistory 获取聊天历史记录
// chat/history/{token}/{role}
func GetChatHistory() []types.Message {
	return nil
}
