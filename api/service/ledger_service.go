package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * 积分账本服务（Ledger）
// *
// * 设计目标：
// *  1. 原子性：余额变更与积分日志必须落在【同一个事务】里，
// *     调用方自己 Begin 的事务必须真正生效（这也是 tx 为什么是显式入参）。
// *  2. 防负余额：扣减使用 `UPDATE ... WHERE id=? AND power>=?` 原子守卫，
// *     多实例并发下不可能把余额扣成负数。
// *  3. 幂等性：每笔账都带业务关联键（RefType + RefId），
// *     库表上有唯一索引 uk_power_ref 做最终闸门。
// *     支付回调重投、失败任务重复扫描、进程重启都不可能导致二次入账。
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"geekai/core/types"
	"geekai/store/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LedgerRef 账本操作的业务关联键。
//
// Type 标识业务类型（如 "order"、"job:mj"），ID 标识业务实体主键。
// 两者共同构成全局唯一的幂等键，落库时写入 PowerLog 的 ref_type / ref_id，
// 由唯一索引 uk_power_ref 保证"同一业务事件只入账一次"。
type LedgerRef struct {
	Type string
	ID   uint
}

var (
	// ErrLedgerNilTx 调用方没有传事务。
	// 账本绝不允许在无事务情况下单独提交，否则调用方后续失败时无法回滚。
	ErrLedgerNilTx = errors.New("账本操作必须传入事务")
	// ErrDuplicateLedgerRef 同一业务关联键已经入过账，说明该事件此前已处理。
	// 调用方应把它理解为"已处理过"而不是"执行失败"。
	ErrDuplicateLedgerRef = errors.New("重复的账本操作")
)

// LedgerService 积分账本服务。
// 所有会改变用户余额的新代码都应走这里，而不是 UserService.IncreasePower/DecreasePower。
type LedgerService struct {
	db *gorm.DB
}

// NewLedgerService 构造账本服务。db 仅用于自开事务的场景（如任务退款 settle），
// Charge / Refund / Grant 的写操作一律使用调用方传入的 tx。
func NewLedgerService(db *gorm.DB) *LedgerService {
	return &LedgerService{db: db}
}

// Charge 扣减用户积分。
//
// 语义：
//   - amount<=0 直接报错（不做任何写入）；
//   - 先以 `power >= amount` 为条件原子扣减，RowsAffected==0 说明余额不足，
//     返回 ErrInsufficientPower，不给用户扣出负数；
//   - 再在同一事务内读取用户行（读到的是扣减后的余额）并写入积分日志；
//   - 日志因唯一键冲突插入失败时返回 ErrDuplicateLedgerRef，
//     此时调用方【必须回滚事务】（本次扣款不能被提交）。
func (s *LedgerService) Charge(tx *gorm.DB, userID uint, amount int, ref LedgerRef, log model.PowerLog) error {
	if tx == nil {
		return ErrLedgerNilTx
	}
	if amount <= 0 {
		return fmt.Errorf("扣减积分数值必须为正数：%d", amount)
	}

	// 原子守卫：余额不足时 RowsAffected 为 0，不会扣成负数。
	result := tx.Model(&model.User{}).
		Where("id = ? AND power >= ?", userID, amount).
		UpdateColumn("power", gorm.Expr("power - ?", amount))
	if result.Error != nil {
		return fmt.Errorf("扣减积分失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrInsufficientPower
	}

	user, err := s.lockUser(tx, userID)
	if err != nil {
		return err
	}
	return s.appendLog(tx, user, amount, types.PowerSub, ref, log)
}

// Refund 退还积分（任务失败退款 / 补偿）。
//
// 幂等策略与 Charge 相反：【先插日志再改余额】。
// 因为日志上的唯一键就是"这笔退款是否已经发生过"的唯一真相，
// 先插入才能用数据库的唯一约束当闸门，避免两个实例同时读到"还没退过"。
//
// 返回值：
//   - (true, nil)  本次退款成功，余额已增加 amount；
//   - (false, nil) 该 ref 此前已退过，本次跳过（既不是错误，也不要再改余额）；
//   - (false, err) 真正的失败，调用方必须回滚事务。
//
// jobID 为业务任务 ID：当 ref.ID 为 0 时用它兜底，同时用于校验幂等键指向同一实体。
// 退还金额取自 log.Amount，调用方应从任务行读取真实扣费值后填入。
func (s *LedgerService) Refund(tx *gorm.DB, jobID uint, ref LedgerRef, log model.PowerLog) (bool, error) {
	if tx == nil {
		return false, ErrLedgerNilTx
	}
	amount := log.Amount
	if amount <= 0 {
		return false, fmt.Errorf("退还积分数值必须为正数：%d", amount)
	}
	ref = ref.withFallbackID(jobID)
	if err := ref.validate(); err != nil {
		return false, err
	}

	// 行锁读用户，保证下面的 Balance 是加上 amount 之后的真实值。
	user, err := s.lockUser(tx, log.UserId)
	if err != nil {
		return false, err
	}

	// 闸门：日志插入成功才代表"这笔退款归我处理"。
	log.UserId = user.Id
	log.Username = user.Username
	log.Amount = amount
	log.Mark = types.PowerAdd
	log.Balance = user.Power + amount
	log.RefType = stringPtr(ref.Type)
	log.RefId = uintPtr(ref.ID)
	log.CreatedAt = time.Now()
	// 先插日志：故意不用 ON CONFLICT DO NOTHING，
	// 让数据库把唯一键冲突以错误形式抛回来，我们才能区分
	// "已退过（跳过）" 和 "真的写失败（要重试）"。
	if err := tx.Create(&log).Error; err != nil {
		if isDuplicateKeyError(err) {
			return false, nil
		}
		return false, fmt.Errorf("记录退款日志失败：%w", err)
	}

	result := tx.Model(&model.User{}).
		Where("id = ?", user.Id).
		UpdateColumn("power", gorm.Expr("power + ?", amount))
	if result.Error != nil {
		return false, fmt.Errorf("退还积分失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("退还积分失败：用户 %d 不存在", user.Id)
	}
	return true, nil
}

// Grant 发放积分（充值到账、系统赠送等"凭空增加余额"的场景）。
//
// 语义与 Refund 完全一致（先插幂等日志再改余额），只是语义命名更贴合充值，
// 且不绑定任务 ID：(true, nil) 表示本次真的发了积分，
// (false, nil) 表示该业务事件此前已发放过，本次跳过。
func (s *LedgerService) Grant(tx *gorm.DB, userID uint, amount int, ref LedgerRef, log model.PowerLog) (bool, error) {
	if tx == nil {
		return false, ErrLedgerNilTx
	}
	if amount <= 0 {
		return false, fmt.Errorf("发放积分数值必须为正数：%d", amount)
	}
	if err := ref.validate(); err != nil {
		return false, err
	}
	log.UserId = userID
	log.Amount = amount
	return s.Refund(tx, 0, ref, log)
}

// lockUser 在同一事务内加行锁读取用户，避免并发读到的余额值过期。
func (s *LedgerService) lockUser(tx *gorm.DB, userID uint) (*model.User, error) {
	var user model.User
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("读取用户 %d 失败：%w", userID, err)
	}
	return &user, nil
}

// appendLog 在当前事务内写入积分日志（支出方向）。
// 唯一键冲突返回 ErrDuplicateLedgerRef，其他错误原样包装。
func (s *LedgerService) appendLog(tx *gorm.DB, user *model.User, amount int, mark types.PowerMark, ref LedgerRef, log model.PowerLog) error {
	log.UserId = user.Id
	log.Username = user.Username
	log.Amount = amount
	log.Mark = mark
	log.Balance = user.Power
	log.RefType = stringPtr(ref.Type)
	log.RefId = uintPtr(ref.ID)
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	if err := tx.Create(&log).Error; err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicateLedgerRef
		}
		return fmt.Errorf("记录积分日志失败：%w", err)
	}
	return nil
}

// validate 幂等键必须完整，否则唯一约束形同虚设（NULL 之间互不冲突）。
func (ref LedgerRef) validate() error {
	if strings.TrimSpace(ref.Type) == "" {
		return errors.New("账本幂等键缺少业务类型 RefType")
	}
	if ref.ID == 0 {
		return errors.New("账本幂等键缺少业务ID RefId")
	}
	return nil
}

// withFallbackID 未显式指定 RefId 时用 jobID 兜底。
func (ref LedgerRef) withFallbackID(jobID uint) LedgerRef {
	if ref.ID == 0 {
		ref.ID = jobID
	}
	return ref
}

// isDuplicateKeyError 判定"唯一键冲突"。
// GORM 只有在开启 TranslateError 时才会返回 gorm.ErrDuplicatedKey，
// 所以这里同时兜住驱动原始错误（MySQL 1062 / Duplicate entry）。
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "1062") ||
		strings.Contains(msg, "Duplicate entry") ||
		strings.Contains(msg, "uk_power_ref")
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	val := v
	return &val
}

func uintPtr(v uint) *uint {
	if v == 0 {
		return nil
	}
	val := v
	return &val
}
