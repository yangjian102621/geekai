package utils

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetLoginUser(c *gin.Context, userId uint) error {
	session := sessions.Default(c)
	session.Set(types.SessionUserId, userId)
	return session.Save()
}

func GetLoginUser(c *gin.Context, db *gorm.DB) (model.User, error) {
	value, exists := c.Get(types.LoginUserCache)
	if exists {
		return value.(model.User), nil
	}

	session := sessions.Default(c)
	userId := session.Get(types.SessionUserId)
	if userId == nil {
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
