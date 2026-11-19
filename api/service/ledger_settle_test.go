package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * 失败任务退款的幂等回归
// *
// * 断言的核心事实：同一笔任务被扫描第二次时，
// * 抢占式 UPDATE 命中 0 行 -> 直接返回 (false, nil)，
// * 事务里不会出现任何 INSERT geekai_power_logs / UPDATE geekai_users，
// * 也就是用户余额一分钱都没变。
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	"geekai/store/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

const jobClaimSQL = "UPDATE `geekai_mj_jobs` SET `power`=\\? WHERE id = \\? AND power > 0"

func mjRefund() FailedJobRefund {
	return FailedJobRefund{
		Table:   "geekai_mj_jobs",
		JobId:   42,
		UserId:  7,
		Power:   3,
		RefType: "job:mj",
		Log: model.PowerLog{
			Type:   types.PowerRefund,
			Model:  "mid-journey",
			Remark: "任务失败，退回积分。任务ID：42",
		},
	}
}

// TestLedgerRefundFailedJobFirstTime 首次扫描：抢占成功 -> 同事务退款
func TestLedgerRefundFailedJobFirstTime(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	mock.ExpectBegin()
	mock.ExpectExec(jobClaimSQL).WithArgs(0, 42).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(userSelectSQL).WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 10))
	mock.ExpectExec(logInsertSQL).WithArgs(
		7, "alice", types.PowerRefund, 3, 13, "mid-journey", "任务失败，退回积分。任务ID：42", types.PowerAdd, "job:mj", 42, sqlmock.AnyArg(),
		0, "", "", "", "", nil,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(creditSQL).WithArgs(3, 7).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	refunded, err := ledger.RefundFailedJob(mjRefund())
	if err != nil {
		t.Fatalf("首次退款不应失败：%v", err)
	}
	if !refunded {
		t.Fatal("首次退款应返回 refunded=true")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerRefundFailedJobSecondScanSkips 第二次扫描（进程重启/多实例/定时器重入）：
// 抢占失败 -> 什么都不做，余额不变。
func TestLedgerRefundFailedJobSecondScanSkips(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	mock.ExpectBegin()
	// power 已经是 0，抢占更新命中 0 行
	mock.ExpectExec(jobClaimSQL).WithArgs(0, 42).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	refunded, err := ledger.RefundFailedJob(mjRefund())
	if err != nil {
		t.Fatalf("重复扫描不应报错：%v", err)
	}
	if refunded {
		t.Fatal("重复扫描必须返回 refunded=false")
	}
	// 若实现里还去插日志/加余额，这里会因为"unexpected call"而失败
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerRefundFailedJobZeroPowerSkips 免费任务（power=0）无需退款，也不该动库
func TestLedgerRefundFailedJobZeroPowerSkips(t *testing.T) {
	db, mock := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	refund := mjRefund()
	refund.Power = 0
	refunded, err := ledger.RefundFailedJob(refund)
	if err != nil {
		t.Fatalf("power=0 不应报错：%v", err)
	}
	if refunded {
		t.Fatal("power=0 不应产生退款")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestLedgerRefundFailedJobRejectsIncomplete 幂等键不完整必须拒绝，避免退成"每次都能退"
func TestLedgerRefundFailedJobRejectsIncomplete(t *testing.T) {
	db, _ := newLedgerTestDB(t)
	ledger := NewLedgerService(db)

	refund := mjRefund()
	refund.RefType = ""
	if _, err := ledger.RefundFailedJob(refund); err == nil {
		t.Fatal("缺少业务类型应被拒绝")
	}

	if _, err := NewLedgerService(nil).RefundFailedJob(mjRefund()); err == nil {
		t.Fatal("账本服务未初始化应被拒绝")
	}
}
