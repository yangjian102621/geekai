package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * LedgerService 单元测试
// *
// * 重点验证三件事：
// *  1. 事务是强制的：nil tx 一律拒绝（防止"看似在事务里、其实自己提交了"的老 bug）。
// *  2. 余额不会被扣成负数：条件更新 RowsAffected==0 即视为余额不足。
// *  3. 幂等：同一 RefType+RefId 第二次调用时，余额完全不发生变更。
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"geekai/core/types"
	"geekai/store/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// newLedgerTestDB 构造一个由 sqlmock 代理的 gorm.DB
func newLedgerTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建 sqlmock 失败：%v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	db, err := gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}),
		&gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)},
	)
	if err != nil {
		t.Fatalf("打开 gorm 失败：%v", err)
	}
	return db, mock
}

// dupKeyErr 模拟 MySQL 唯一键冲突（1062）
func dupKeyErr(ref string) error {
	return errors.New("Error 1062: Duplicate entry '" + ref + "' for key 'uk_power_ref'")
}

const (
	userSelectSQL = "SELECT .*geekai_users"
	deductSQL     = "UPDATE `geekai_users` SET `power`=power - \\? WHERE id = \\? AND power >= \\?"
	creditSQL     = "UPDATE `geekai_users` SET `power`=power \\+ \\? WHERE id = \\?"
	logInsertSQL  = "INSERT INTO `geekai_power_logs`"
)

// TestLedgerRejectsNilTx 账本绝不允许脱离事务执行
func TestLedgerRejectsNilTx(t *testing.T) {
	ledger := NewLedgerService(nil)

	if err := ledger.Charge(nil, 7, 2, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{}); !errors.Is(err, ErrLedgerNilTx) {
		t.Fatalf("Charge 应拒绝 nil tx，实际返回 %v", err)
	}
	if _, err := ledger.Refund(nil, 5, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{Amount: 2}); !errors.Is(err, ErrLedgerNilTx) {
		t.Fatalf("Refund 应拒绝 nil tx，实际返回 %v", err)
	}
	if _, err := ledger.Grant(nil, 7, 2, LedgerRef{Type: "order", ID: 9}, model.PowerLog{}); !errors.Is(err, ErrLedgerNilTx) {
		t.Fatalf("Grant 应拒绝 nil tx，实际返回 %v", err)
	}
}

// TestLedgerRejectsInvalidAmount 非正数金额与缺失幂等键必须在写库前被拒绝
func TestLedgerRejectsInvalidAmount(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	// 一条 SQL 都不该发出，所以这里不设置任何 Expect*
	mock.ExpectBegin()
	mock.ExpectRollback()
	err := db.Transaction(func(tx *gorm.DB) error {
		return ledger.Charge(tx, 7, 0, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{})
	})
	if err == nil {
		t.Fatal("amount=0 应被拒绝")
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := ledger.Refund(tx, 5, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{Amount: -3})
		return err
	})
	if err == nil {
		t.Fatal("负数退款金额应被拒绝")
	}

	// 幂等键缺失 => 拒绝（否则唯一约束形同虚设）
	mock.ExpectBegin()
	mock.ExpectRollback()
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := ledger.Grant(tx, 7, 2, LedgerRef{Type: "", ID: 0}, model.PowerLog{})
		return err
	})
	if err == nil {
		t.Fatal("缺失幂等键的发放应被拒绝")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerChargeInsufficientPower 条件更新命中 0 行 => 余额不足，不能扣成负数
func TestLedgerChargeInsufficientPower(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	mock.ExpectBegin()
	mock.ExpectExec(deductSQL).WithArgs(2, 7, 2).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := db.Transaction(func(tx *gorm.DB) error {
		return ledger.Charge(tx, 7, 2, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{Type: types.PowerConsume})
	})
	if !errors.Is(err, ErrInsufficientPower) {
		t.Fatalf("应返回 ErrInsufficientPower，实际 %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerChargeWritesPostDeductionBalance 扣款成功时日志里的 Balance 必须是扣减后的真实余额
func TestLedgerChargeWritesPostDeductionBalance(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	mock.ExpectBegin()
	mock.ExpectExec(deductSQL).WithArgs(2, 7, 2).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 8))
	// INSERT 的取值顺序与 model.PowerLog 字段声明顺序一致：
	// user_id, username, type, amount, balance, model, remark, mark, ref_type, ref_id, created_at
	mock.ExpectExec(logInsertSQL).WithArgs(
		7, "alice", types.PowerConsume, 2, 8, "gpt-4o", "测试扣费", types.PowerSub, "job:mj", 5, sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := db.Transaction(func(tx *gorm.DB) error {
		return ledger.Charge(tx, 7, 2, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{
			Type:   types.PowerConsume,
			Model:  "gpt-4o",
			Remark: "测试扣费",
		})
	})
	if err != nil {
		t.Fatalf("扣款应当成功，实际 %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerChargeDuplicateRefReturnsDuplicate 重复业务键入账必须能被调用方识别为"已处理过"
func TestLedgerChargeDuplicateRefReturnsDuplicate(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	mock.ExpectBegin()
	mock.ExpectExec(deductSQL).WithArgs(2, 7, 2).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 8))
	mock.ExpectExec(logInsertSQL).WillReturnError(dupKeyErr("job:mj-5"))
	// 重复入账时调用方必须回滚，扣款不能被提交
	mock.ExpectRollback()

	err := db.Transaction(func(tx *gorm.DB) error {
		return ledger.Charge(tx, 7, 2, LedgerRef{Type: "job:mj", ID: 5}, model.PowerLog{Type: types.PowerConsume})
	})
	if !errors.Is(err, ErrDuplicateLedgerRef) {
		t.Fatalf("应返回 ErrDuplicateLedgerRef，实际 %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerRefundIdempotent 幂等性证明：
// 第一笔退款 (true, nil) 且发出一条 UPDATE power=power+；
// 同一 ref 第二次退款 (false, nil)，事务里【没有任何 UPDATE】——
// 也就是余额完全没变，不可能重复退款。
func TestLedgerRefundIdempotent(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)
	ref := LedgerRef{Type: "job:mj", ID: 5}

	// 第一次：日志插入成功 -> 加余额
	mock.ExpectBegin()
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 10))
	mock.ExpectExec(logInsertSQL).WithArgs(
		7, "alice", types.PowerRefund, 3, 13, "mid-journey", "任务失败，退回积分", types.PowerAdd, "job:mj", 5, sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(creditSQL).WithArgs(3, 7).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	var refunded bool
	err := db.Transaction(func(tx *gorm.DB) error {
		ok, err := ledger.Refund(tx, 5, ref, model.PowerLog{
			UserId: 7,
			Type:   types.PowerRefund,
			Amount: 3,
			Model:  "mid-journey",
			Remark: "任务失败，退回积分",
		})
		refunded = ok
		return err
	})
	if err != nil || !refunded {
		t.Fatalf("首次退款应成功，refunded=%v err=%v", refunded, err)
	}

	// 第二次：同一 ref，唯一键冲突 -> (false, nil)，且不下发 UPDATE
	mock.ExpectBegin()
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 13))
	mock.ExpectExec(logInsertSQL).WillReturnError(dupKeyErr("job:mj-5"))
	mock.ExpectCommit() // 冲突不是错误，事务正常提交（本次没有任何写入）

	refunded = true
	err = db.Transaction(func(tx *gorm.DB) error {
		ok, err := ledger.Refund(tx, 5, ref, model.PowerLog{
			UserId: 7,
			Type:   types.PowerRefund,
			Amount: 3,
			Model:  "mid-journey",
			Remark: "任务失败，退回积分",
		})
		refunded = ok
		return err
	})
	if err != nil {
		t.Fatalf("重复退款不应返回错误，实际 %v", err)
	}
	if refunded {
		t.Fatal("重复退款必须返回 false（说明本次跳过，余额不变）")
	}
	// ExpectationsWereMet 同时证明了第二次事务里没有出现过任何 UPDATE geekai_users
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerGrantIdempotent 充值发放同样具备幂等性：同一订单号只会到账一次
func TestLedgerGrantIdempotent(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)
	ref := LedgerRef{Type: "order", ID: 1024}

	// 首次到账
	mock.ExpectBegin()
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 100))
	mock.ExpectExec(logInsertSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(creditSQL).WithArgs(500, 7).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	var granted bool
	err := db.Transaction(func(tx *gorm.DB) error {
		ok, err := ledger.Grant(tx, 7, 500, ref, model.PowerLog{Type: types.PowerRecharge})
		granted = ok
		return err
	})
	if err != nil || !granted {
		t.Fatalf("首次充值应到账，granted=%v err=%v", granted, err)
	}

	// 回调重投：订单 1024 已入账 -> 跳过，不加余额
	mock.ExpectBegin()
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 600))
	mock.ExpectExec(logInsertSQL).WillReturnError(dupKeyErr("order-1024"))
	mock.ExpectCommit()

	granted = true
	err = db.Transaction(func(tx *gorm.DB) error {
		ok, err := ledger.Grant(tx, 7, 500, ref, model.PowerLog{Type: types.PowerRecharge})
		granted = ok
		return err
	})
	if err != nil {
		t.Fatalf("重复回调不应报错，实际 %v", err)
	}
	if granted {
		t.Fatal("同一订单第二次回调必须跳过发放（余额保持不变）")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
