package utils

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLoginUser(c *gin.Context, db *gorm.DB) (model.User, error) {
	value, exists := c.Get(types.LoginUserCache)
	if exists {
		return value.(model.User), nil
	}

	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return model.User{}, errors.New("user not login")
	}

	var user model.User
	res := db.First(&user, userId)
	// 更新缓存
	if res.Error == nil {
		c.Set(types.LoginUserCache, user)
	}
	return user, res.Error
}
