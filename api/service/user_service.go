package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * UserService 积分增减（历史实现，保留兼容）
// *
// * ⚠️ 新代码一律走 LedgerService（见 ledger_service.go）：
// * 本文件里的 IncreasePower / DecreasePower 会自己开事务并提交，
// * 调用方哪怕外面包了一层 tx := db.Begin() 也回滚不到它们，
// * 这正是支付回调"积分已发、订单状态没改"的根因。
// *
// * 现阶段仍有 20+ 个调用点在用这两个方法，为避免一次性改动过大，
// * 先保持行为与签名不变，待 Phase 2-B 逐步迁移到 LedgerService 后删除。
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

var ErrInsufficientPower = errors.New("用户积分不足")

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// IncreasePower 增加用户积分
func (s *UserService) IncreasePower(userId uint, power int, log model.PowerLog) error {
	if power == 0 {
		return nil
	}
	if power < 0 {
		return errors.New("积分数值不能为负数")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.User{}).Where("id = ?", userId).UpdateColumn("power", gorm.Expr("power + ?", power))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("用户不存在或积分未变更")
		}
		var user model.User
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}
		return tx.Create(&model.PowerLog{
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
	})
}

// DecreasePower 减少用户积分
func (s *UserService) DecreasePower(userId uint, power int, log model.PowerLog) error {
	if power == 0 {
		return nil
	}
	if power < 0 {
		return errors.New("积分数值不能为负数")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 在数据库中原子检查余额，避免多个实例并发扣成负数。
		result := tx.Model(&model.User{}).Where("id = ? AND power >= ?", userId, power).UpdateColumn("power", gorm.Expr("power - ?", power))
		if result.Error != nil {
			return fmt.Errorf("扣减积分失败：%w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrInsufficientPower
		}
		var user model.User
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      log.Type,
			Amount:    power,
			Balance:   user.Power,
			Mark:      types.PowerSub,
			Model:     log.Model,
			Remark:    log.Remark,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("记录积分日志失败：%w", err)
		}
		return nil
	})
}
