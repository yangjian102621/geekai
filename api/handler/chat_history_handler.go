package handler

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Update 更新会话标题
func (h *ChatHandler) Update(c *gin.Context) {
	var data struct {
		Id    uint   `json:"id"`
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var m = model.ChatItem{}
	m.Id = data.Id
	res := h.db.Model(&m).UpdateColumn("title", data.Title)
	if res.Error != nil {
		resp.ERROR(c, "Failed to update database")
		return
	}

	resp.SUCCESS(c, types.OkMsg)
}

// History 获取聊天历史记录
func (h *ChatHandler) History(c *gin.Context) {
	chatId := c.Query("chat_id") // 会话 ID
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}
	var items []model.HistoryMessage
	var messages = make([]vo.HistoryMessage, 0)
	res := h.db.Where("chat_id = ? AND user_id = ?", chatId, user.Id).Find(&items)
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

// Clear 清空所有聊天记录
func (h *ChatHandler) Clear(c *gin.Context) {
	// 获取当前登录用户所有的聊天会话
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	var chats []model.ChatItem
	res := h.db.Where("user_id = ?", user.Id).Find(&chats)
	if res.Error != nil {
		resp.ERROR(c, "No chats found")
		return
	}

	var chatIds = make([]string, 0)
	for _, chat := range chats {
		chatIds = append(chatIds, chat.ChatId)
		// 清空会话上下文
		h.App.ChatContexts.Delete(chat.ChatId)
	}
	err = h.db.Transaction(func(tx *gorm.DB) error {
		res := h.db.Where("user_id =?", user.Id).Delete(&model.ChatItem{})
		if res.Error != nil {
			return res.Error
		}

		res = h.db.Where("user_id = ? AND chat_id IN ?", user.Id, chatIds).Delete(&model.HistoryMessage{})
		if res.Error != nil {
			return res.Error
		}

		// TODO: 是否要删除 MidJourney 绘画记录和图片文件？
		return nil
	})

	if err != nil {
		logger.Errorf("Error with delete chats: %+v", err)
		resp.ERROR(c, "Failed to remove chat from database.")
		return
	}

	resp.SUCCESS(c, types.OkMsg)
}
