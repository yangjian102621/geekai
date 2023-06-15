package service

import (
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"gorm.io/gorm"
)

type ChatRoleService struct {
	DB    *gorm.DB
	Model interface{}
}

func NewChatRoleService(db *gorm.DB) *ChatRoleService {
	return &ChatRoleService{DB: db, Model: &model.ChatRole{}}
}

func (s *ChatRoleService) Create(value interface{}) error {
	r, ok := value.(vo.ChatRole)

	if ok {
		var role model.ChatRole
		err := utils.CopyObject(r, &role)
		if err != nil {
			return err
		}

		result := s.DB.Create(&role)
		return result.Error
	} else {
		result := s.DB.Create(value)
		return result.Error
	}
}
