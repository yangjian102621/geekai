package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * 失败任务积分结算（Settle）
// *
// * 背景：MJ / 图片 / Suno / 视频 / 即梦 五个业务各自复制了一份退款循环，
// * 形如「查出 status=failed 且 power>0 的任务 → IncreasePower → power=0」。
// * 这段代码有三个钱上的 bug：
// *   1. IncreasePower 自己开事务提交，跟 power=0 的更新不是原子的；
// *   2. 先加钱后置零，进程在两步之间挂掉 → 重启后重复扫描，重复退款；
// *   3. 多实例同时扫描 → 同一笔任务退两次。
// *
// * 这里把它们收敛成一个函数，规则只有一句话：
// *   「谁把任务行的 power 从 >0 改成 0，谁才有资格退款」。
// * 抢占式的条件更新 + 同一事务内的幂等账本，两者叠加让重复扫描彻底无害。
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"

	"geekai/core/types"
	"geekai/store/model"

	"gorm.io/gorm"
)

// ErrLedgerNotReady 账本服务未初始化（漏了 fx 注入或构造参数传了 nil）
var ErrLedgerNotReady = errors.New("账本服务未初始化")

// FailedJobRefund 描述一笔待退还的失败任务。
//
// Table 是任务表名（如 geekai_mj_jobs），各业务包自己填自己的表，
// 这样 settle 逻辑可以通用，且不会出现 service/mj 反向 import service/image 的循环依赖。
type FailedJobRefund struct {
	Table   string // 任务表名
	JobId   uint   // 任务主键
	UserId  uint   // 扣费用户
	Power   int    // 该任务实际扣掉的积分
	RefType string // 业务类型，如 job:mj、job:image
	Log     model.PowerLog
}

// validate 关键字段必须齐备，否则幂等键会退化成"每次都能退"
func (r FailedJobRefund) validate() error {
	if r.Table == "" {
		return errors.New("任务退款缺少表名")
	}
	if r.JobId == 0 {
		return errors.New("任务退款缺少任务ID")
	}
	if r.UserId == 0 {
		return errors.New("任务退款缺少用户ID")
	}
	if r.RefType == "" {
		return errors.New("任务退款缺少业务类型")
	}
	return nil
}

// RefundFailedJob 退还一笔失败任务的积分，天然幂等。
//
// 返回值：
//   - (true, nil)  本次真的退了（余额已增加）；
//   - (false, nil) 这笔任务已被处理过（power 已是 0，或账本里已有同一 ref 的记录），
//     重复扫描不会二次退款；
//   - (false, err) 真正的失败，任务行的 power 不会被清零，下一轮扫描会重试。
func (s *LedgerService) RefundFailedJob(refund FailedJobRefund) (bool, error) {
	if s == nil || s.db == nil {
		return false, ErrLedgerNotReady
	}
	if err := refund.validate(); err != nil {
		return false, err
	}
	// 没有扣费的任务（免费额度/失败即不扣费）无需退款
	if refund.Power <= 0 {
		return false, nil
	}

	log := refund.Log
	log.UserId = refund.UserId
	log.Amount = refund.Power
	if log.Type == 0 {
		log.Type = types.PowerRefund
	}

	var refunded bool
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 抢占：只有把 power 从 >0 改成 0 的执行者才有退款资格。
		// 并发/重复扫描的其它执行者会拿到 RowsAffected==0，直接放弃。
		res := tx.Table(refund.Table).
			Where("id = ? AND power > 0", refund.JobId).
			UpdateColumn("power", 0)
		if res.Error != nil {
			return fmt.Errorf("锁定失败任务 %d 失败：%w", refund.JobId, res.Error)
		}
		if res.RowsAffected == 0 {
			// 已被其它实例或上一轮扫描处理过
			return nil
		}

		// 积分与 power=0 在同一个事务里；Refund 内部再用 uk_power_ref 兜一次幂等。
		ok, err := s.Refund(tx, refund.JobId, LedgerRef{Type: refund.RefType, ID: refund.JobId}, log)
		if err != nil {
			return err
		}
		refunded = ok
		return nil
	})
	if err != nil {
		return false, err
	}
	return refunded, nil
}
