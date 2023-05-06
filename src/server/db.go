package server

import (
	"chatplus/types"
	"chatplus/utils"
	"encoding/json"
)

const (
	UserPrefix        = "chat/users/"
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

// GetUsers 获取 user 信息
// chat/users
func GetUsers() []types.User {
	items := db.Search(UserPrefix)
	var users = make([]types.User, 0)
	for _, v := range items {
		var user types.User
		err := json.Unmarshal([]byte(v), &user)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users
}

func PutUser(user types.User) error {
	key := UserPrefix + user.Name
	return db.Put(key, user)
}

func GetUser(username string) (*types.User, error) {
	key := UserPrefix + username
	bytes, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	var user types.User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func RemoveUser(username string) error {
	key := UserPrefix + username
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

func GetChatRole(key string) (*types.ChatRole, error) {
	key = ChatRolePrefix + key
	bytes, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	var role types.ChatRole
	err = json.Unmarshal(bytes, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetChatHistory 获取聊天历史记录
// chat/history/{user}/{role}
func GetChatHistory(user string, role string) ([]types.Message, error) {
	key := ChatHistoryPrefix + user + "/" + role
	bytes, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	var message []types.Message
	err = json.Unmarshal(bytes, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// AppendChatHistory 追加聊天记录
func AppendChatHistory(user string, role string, message types.Message) error {
	messages, err := GetChatHistory(user, role)
	if err != nil {
		messages = make([]types.Message, 0)
	}

	messages = append(messages, message)
	key := ChatHistoryPrefix + user + "/" + role
	return db.Put(key, messages)
}

// ClearChatHistory 清空某个角色下的聊天记录
func ClearChatHistory(user string, role string) error {
	key := ChatHistoryPrefix + user + "/" + role
	return db.Delete(key)
}
