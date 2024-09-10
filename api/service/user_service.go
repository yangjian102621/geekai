package service

import (
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"gorm.io/gorm"
	"sync"
	"time"
)

type UserService struct {
	db   *gorm.DB
	lock sync.Mutex
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db, lock: sync.Mutex{}}
}

// IncreasePower 增加用户算力
func (s *UserService) IncreasePower(userId int, power int, log model.PowerLog) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx := s.db.Begin()
	err := tx.Model(&model.User{}).Where("id", userId).UpdateColumn("power", gorm.Expr("power + ?", power)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var user model.User
	tx.Where("id", userId).First(&user)
	err = tx.Create(&model.PowerLog{
		UserId:    user.Id,
		Username:  user.Username,
		Type:      log.Type,
		Amount:    power,
		Balance:   user.Power,
		Mark:      types.PowerAdd,
		Model:     log.Model,
		Remark:    log.Remark,
		CreatedAt: time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// DecreasePower 减少用户算力
func (s *UserService) DecreasePower(userId int, power int, log model.PowerLog) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx := s.db.Begin()
	err := tx.Model(&model.User{}).Where("id", userId).UpdateColumn("power", gorm.Expr("power - ?", power)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("扣减算力失败：%v", err)
	}
	var user model.User
	tx.Where("id", userId).First(&user)
	err = tx.Create(&model.PowerLog{
		UserId:    user.Id,
		Username:  user.Username,
		Type:      log.Type,
		Amount:    power,
		Balance:   user.Power,
		Mark:      types.PowerSub,
		Model:     log.Model,
		Remark:    log.Remark,
		CreatedAt: time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("记录算力日志失败：%v", err)
	}
	tx.Commit()
	return nil
}
