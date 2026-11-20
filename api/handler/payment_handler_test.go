package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * 支付回调 paySuccess 的幂等/金额校验回归
// *
// * 这两条是"钱的口子"上最容易出事的地方：
// *  1. 渠道重复回调 / 轮询补偿并发进来时，只能有一个执行者发积分；
// *  2. 渠道回传金额与订单金额不一致时，绝不能发货。
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/service"
	"geekai/service/payment"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// newPaymentTestHandler 直接构造 PaymentHandler（paySuccess 只依赖 DB 与 LedgerService）
func newPaymentTestHandler(t *testing.T) (*PaymentHandler, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	db, err := gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}),
		&gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)},
	)
	if err != nil {
		t.Fatal(err)
	}
	h := &PaymentHandler{
		BaseHandler:   BaseHandler{DB: db},
		ledgerService: service.NewLedgerService(db),
	}
	return h, mock
}

// expectPendingOrder 期望一次"按订单号查订单"，返回待支付订单
func expectPendingOrder(mock sqlmock.Sqlmock, amount float64) {
	mock.ExpectQuery("SELECT .*geekai_orders").
		WithArgs("ORDER-2026").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "product_id", "username", "order_no", "subject", "amount", "status", "remark", "channel", "pay_way",
		}).AddRow(
			1024, 7, 3, "alice", "ORDER-2026", "标准版", amount, 0,
			`{"power":500,"name":"标准版","price":9.9}`, "epay", "alipay",
		))
}

const orderPaidGuardSQL = "UPDATE `geekai_orders` SET .* WHERE id = \\? AND status <> \\?"

// TestPaySuccessGrantsOnceOnFirstCallback 首次回调：订单置成功 + 发积分 + 销量全部在同一事务
func TestPaySuccessGrantsOnceOnFirstCallback(t *testing.T) {
	h, mock := newPaymentTestHandler(t)
	expectPendingOrder(mock, 9.9)

	mock.ExpectBegin()
	mock.ExpectExec(orderPaidGuardSQL).WillReturnResult(sqlmock.NewResult(0, 1))
	// 发积分：行锁读用户 -> 插幂等日志 -> 加余额
	mock.ExpectQuery("SELECT .*geekai_users").WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 100))
	mock.ExpectExec("INSERT INTO `geekai_power_logs`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE `geekai_users` SET `power`=power \\+ \\?").WithArgs(500, 7).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE `geekai_products` SET `sales`=sales \\+ \\?").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	info := payment.OrderInfo{OutTradeNo: "ORDER-2026", TradeId: "TRADE-1", Amount: "9.90", Status: payment.Success, PayTime: "2026-01-02 03:04:05"}
	if err := h.paySuccess(info); err != nil {
		t.Fatalf("首次回调应当成功，实际 %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestPaySuccessSkipsWhenGuardMisses 幂等闸门命中 0 行（并发回调/渠道重投）：
// 事务正常提交，但【不发积分、不加余额、不改销量】——这就是"第二次调用余额不变"的证明。
func TestPaySuccessSkipsWhenGuardMisses(t *testing.T) {
	h, mock := newPaymentTestHandler(t)
	// 模拟另一个回调已把订单改成成功之前，本次读到的仍是待支付快照
	expectPendingOrder(mock, 9.9)

	mock.ExpectBegin()
	// RowsAffected == 0 表示"状态已不是待支付"，发货权归别人
	mock.ExpectExec(orderPaidGuardSQL).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	info := payment.OrderInfo{OutTradeNo: "ORDER-2026", TradeId: "TRADE-1", Amount: "9.90", Status: payment.Success}
	if err := h.paySuccess(info); err != nil {
		t.Fatalf("重复回调不应报错，实际 %v", err)
	}
	// ExpectationsWereMet 同时证明：整个过程没有任何 INSERT geekai_power_logs、没有任何 UPDATE geekai_users
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestPaySuccessRejectsAmountMismatch 金额不一致必须拒绝发货（且不开启事务）
func TestPaySuccessRejectsAmountMismatch(t *testing.T) {
	h, mock := newPaymentTestHandler(t)
	expectPendingOrder(mock, 9.9)

	info := payment.OrderInfo{OutTradeNo: "ORDER-2026", Amount: "0.01", Status: payment.Success}
	err := h.paySuccess(info)
	if err == nil {
		t.Fatal("金额不一致时必须返回错误让渠道重试/人工介入")
	}
	if !strings.Contains(err.Error(), "支付金额不匹配") {
		t.Fatalf("错误信息应说明金额不匹配，实际 %v", err)
	}
	// 没有任何写操作：ExpectBegin 都没设置，若实现里开了事务这里会报错
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
