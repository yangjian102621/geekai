package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatHandler struct {
	handler.BaseHandler
}

func NewChatHandler(app *core.AppServer, db *gorm.DB) *ChatHandler {
	return &ChatHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

type chatItemVo struct {
	Username  string      `json:"username"`
	UserId    uint        `json:"user_id"`
	ChatId    string      `json:"chat_id"`
	Title     string      `json:"title"`
	Role      vo.ChatRole `json:"role"`
	Model     string      `json:"model"`
	Token     int         `json:"token"`
	CreatedAt int64       `json:"created_at"`
	MsgNum    int         `json:"msg_num"` // 消息数量
}

func (h *ChatHandler) List(c *gin.Context) {
	if err := utils.CheckPermission(c, h.DB); err != nil {
		resp.NotPermission(c)
		return
	}

	var data struct {
		Title    string   `json:"title"`
		UserId   uint     `json:"user_id"`
		Model    string   `json:"model"`
		CreateAt []string `json:"created_time"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Title != "" {
		session = session.Where("title LIKE ?", "%"+data.Title+"%")
	}
	if data.UserId > 0 {
		session = session.Where("user_id = ?", data.UserId)
	}
	if data.Model != "" {
		session = session.Where("model = ?", data.Model)
	}
	if len(data.CreateAt) == 2 {
		start := utils.Str2stamp(data.CreateAt[0] + " 00:00:00")
		end := utils.Str2stamp(data.CreateAt[1] + " 00:00:00")
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	var total int64
	session.Model(&model.ChatItem{}).Count(&total)
	var items []model.ChatItem
	var list = make([]chatItemVo, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		userIds := make([]uint, 0)
		chatIds := make([]string, 0)
		roleIds := make([]uint, 0)
		for _, item := range items {
			userIds = append(userIds, item.UserId)
			chatIds = append(chatIds, item.ChatId)
			roleIds = append(roleIds, item.RoleId)
		}
		var messages []model.ChatMessage
		var users []model.User
		var roles []model.ChatRole
		h.DB.Where("chat_id IN ?", chatIds).Find(&messages)
		h.DB.Where("id IN ?", userIds).Find(&users)
		h.DB.Where("id IN ?", roleIds).Find(&roles)

		tokenMap := make(map[string]int)
		userMap := make(map[uint]string)
		msgMap := make(map[string]int)
		roleMap := make(map[uint]vo.ChatRole)
		for _, msg := range messages {
			tokenMap[msg.ChatId] += msg.Tokens
			msgMap[msg.ChatId] += 1
		}
		for _, user := range users {
			userMap[user.Id] = user.Username
		}
		for _, r := range roles {
			var roleVo vo.ChatRole
			err := utils.CopyObject(r, &roleVo)
			if err != nil {
				continue
			}
			roleMap[r.Id] = roleVo
		}
		for _, item := range items {
			list = append(list, chatItemVo{
				UserId:    item.UserId,
				Username:  userMap[item.UserId],
				ChatId:    item.ChatId,
				Title:     item.Title,
				Model:     item.Model,
				Token:     tokenMap[item.ChatId],
				MsgNum:    msgMap[item.ChatId],
				Role:      roleMap[item.RoleId],
				CreatedAt: item.CreatedAt.Unix(),
			})
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}

type chatMessageVo struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Model     string `json:"model"`
	Token     int    `json:"token"`
	Icon      string `json:"icon"`
	CreatedAt int64  `json:"created_at"`
}

// Messages 读取聊天记录列表
func (h *ChatHandler) Messages(c *gin.Context) {
	var data struct {
		UserId   uint     `json:"user_id"`
		Content  string   `json:"content"`
		Model    string   `json:"model"`
		CreateAt []string `json:"created_time"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Content != "" {
		session = session.Where("content LIKE ?", "%"+data.Content+"%")
	}
	if data.UserId > 0 {
		session = session.Where("user_id = ?", data.UserId)
	}
	if data.Model != "" {
		session = session.Where("model = ?", data.Model)
	}
	if len(data.CreateAt) == 2 {
		start := utils.Str2stamp(data.CreateAt[0] + " 00:00:00")
		end := utils.Str2stamp(data.CreateAt[1] + " 00:00:00")
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	var total int64
	session.Model(&model.ChatMessage{}).Count(&total)
	var items []model.ChatMessage
	var list = make([]chatMessageVo, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		userIds := make([]uint, 0)
		for _, item := range items {
			userIds = append(userIds, item.UserId)
		}
		var users []model.User
		h.DB.Where("id IN ?", userIds).Find(&users)
		userMap := make(map[uint]string)
		for _, user := range users {
			userMap[user.Id] = user.Username
		}
		for _, item := range items {
			list = append(list, chatMessageVo{
				Id:        item.Id,
				UserId:    item.UserId,
				Username:  userMap[item.UserId],
				Content:   item.Content,
				Model:     item.Model,
				Token:     item.Tokens,
				Icon:      item.Icon,
				Type:      item.Type,
				CreatedAt: item.CreatedAt.Unix(),
			})
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}

// History 获取聊天历史记录
func (h *ChatHandler) History(c *gin.Context) {
	chatId := c.Query("chat_id") // 会话 ID
	var items []model.ChatMessage
	var messages = make([]vo.HistoryMessage, 0)
	res := h.DB.Where("chat_id = ?", chatId).Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "No history message")
		return
	} else {
		for _, item := range items {
			var v vo.HistoryMessage
			err := utils.CopyObject(item, &v)
			v.CreatedAt = item.CreatedAt.Unix()
			v.UpdatedAt = item.UpdatedAt.Unix()
			if err == nil {
				messages = append(messages, v)
			}
		}
	}

	resp.SUCCESS(c, messages)
}

// RemoveChat 删除对话
func (h *ChatHandler) RemoveChat(c *gin.Context) {
	chatId := h.GetTrim(c, "chat_id")
	if chatId == "" {
		resp.ERROR(c, "请传入 ChatId")
		return
	}

	tx := h.DB.Begin()
	// 删除聊天记录
	res := tx.Unscoped().Debug().Where("chat_id = ?", chatId).Delete(&model.ChatMessage{})
	if res.Error != nil {
		resp.ERROR(c, "failed to remove chat message")
		return
	}

	// 删除对话
	res = tx.Unscoped().Where("chat_id = ?", chatId).Delete(model.ChatItem{})
	if res.Error != nil {
		tx.Rollback() // 回滚
		resp.ERROR(c, "failed to remove chat")
		return
	}

	tx.Commit()
	resp.SUCCESS(c)
}

// RemoveMessage 删除聊天记录
func (h *ChatHandler) RemoveMessage(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	tx := h.DB.Unscoped().Where("id = ?", id).Delete(&model.ChatMessage{})
	if tx.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}
