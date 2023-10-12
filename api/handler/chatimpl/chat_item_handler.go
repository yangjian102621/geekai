package chatimpl

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// List 获取会话列表
func (h *ChatHandler) List(c *gin.Context) {
	userId := h.GetInt(c, "user_id", 0)
	if userId == 0 {
		resp.ERROR(c, "The parameter 'user_id' is needed.")
		return
	}
	var items = make([]vo.ChatItem, 0)
	var chats []model.ChatItem
	res := h.db.Where("user_id = ?", userId).Order("id DESC").Find(&chats)
	if res.Error == nil {
		var roleIds = make([]uint, 0)
		for _, chat := range chats {
			roleIds = append(roleIds, chat.RoleId)
		}
		var roles []model.ChatRole
		res = h.db.Find(&roles, roleIds)
		if res.Error == nil {
			roleMap := make(map[uint]model.ChatRole)
			for _, role := range roles {
				roleMap[role.Id] = role
			}

			for _, chat := range chats {
				var item vo.ChatItem
				err := utils.CopyObject(chat, &item)
				if err == nil {
					item.Id = chat.Id
					item.Icon = roleMap[chat.RoleId].Icon
					items = append(items, item)
				}
			}
		}

	}
	resp.SUCCESS(c, items)
}

// Update 更新会话标题
func (h *ChatHandler) Update(c *gin.Context) {
	var data struct {
		ChatId string `json:"chat_id"`
		Title  string `json:"title"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	res := h.db.Model(&model.ChatItem{}).Where("chat_id = ?", data.ChatId).UpdateColumn("title", data.Title)
	if res.Error != nil {
		resp.ERROR(c, "Failed to update database")
		return
	}

	resp.SUCCESS(c, types.OkMsg)
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

// History 获取聊天历史记录
func (h *ChatHandler) History(c *gin.Context) {
	chatId := c.Query("chat_id") // 会话 ID
	var items []model.HistoryMessage
	var messages = make([]vo.HistoryMessage, 0)
	res := h.db.Where("chat_id = ?", chatId).Find(&items)
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

// Remove 删除会话
func (h *ChatHandler) Remove(c *gin.Context) {
	chatId := h.GetTrim(c, "chat_id")
	if chatId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	res := h.db.Where("user_id = ? AND chat_id = ?", user.Id, chatId).Delete(&model.ChatItem{})
	if res.Error != nil {
		resp.ERROR(c, "Failed to update database")
		return
	}

	// 删除当前会话的聊天记录
	res = h.db.Where("user_id = ? AND chat_id =?", user.Id, chatId).Delete(&model.ChatItem{})
	if res.Error != nil {
		resp.ERROR(c, "Failed to remove chat from database.")
		return
	}

	// TODO: 是否要删除 MidJourney 绘画记录和图片文件？

	// 清空会话上下文
	h.App.ChatContexts.Delete(chatId)
	resp.SUCCESS(c, types.OkMsg)
}

// Detail 对话详情，用户导出对话
func (h *ChatHandler) Detail(c *gin.Context) {
	chatId := h.GetTrim(c, "chat_id")
	if utils.IsEmptyValue(chatId) {
		resp.ERROR(c, "Invalid chatId")
		return
	}

	var chatItem model.ChatItem
	res := h.db.Where("chat_id = ?", chatId).First(&chatItem)
	if res.Error != nil {
		resp.ERROR(c, "No chat found")
		return
	}

	var chatItemVo vo.ChatItem
	err := utils.CopyObject(chatItem, &chatItemVo)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, chatItemVo)
}
