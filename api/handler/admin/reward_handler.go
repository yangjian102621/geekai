package admin

import (
	"chatplus/core"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RewardHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewRewardHandler(app *core.AppServer, db *gorm.DB) *RewardHandler {
	h := RewardHandler{db: db}
	h.App = app
	return &h
}

func (h *RewardHandler) List(c *gin.Context) {
	var items []model.Reward
	res := h.db.Order("id DESC").Find(&items)
	var rewards = make([]vo.Reward, 0)
	if res.Error == nil {
		userIds := make([]uint, 0)
		for _, v := range items {
			userIds = append(userIds, v.UserId)
		}
		var users []model.User
		h.db.Where("id IN ?", userIds).Find(&users)
		var userMap = make(map[uint]model.User)
		for _, u := range users {
			userMap[u.Id] = u
		}

		for _, v := range items {
			var r vo.Reward
			err := utils.CopyObject(v, &r)
			if err != nil {
				continue
			}

			r.Id = v.Id
			r.Username = userMap[v.UserId].Mobile
			r.CreatedAt = v.CreatedAt.Unix()
			r.UpdatedAt = v.UpdatedAt.Unix()
			rewards = append(rewards, r)
		}
	}

	resp.SUCCESS(c, rewards)
}

func (h *RewardHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.db.Where("id = ?", id).Delete(&model.Reward{})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}
