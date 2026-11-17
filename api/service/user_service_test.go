package service

import (
	"errors"
	"geekai/store/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestPowerAmounts(t *testing.T) {
	s := NewUserService(nil)
	if s.DecreasePower(7, 0, model.PowerLog{}) != nil || s.IncreasePower(7, 0, model.PowerLog{}) != nil {
		t.Fatal("zero cost should not change the ledger")
	}
	if s.DecreasePower(7, -1, model.PowerLog{}) == nil || s.IncreasePower(7, -1, model.PowerLog{}) == nil {
		t.Fatal("negative amounts must be rejected")
	}
}

func TestDecreasePowerAtomicGuardAndRollback(t *testing.T) {
	for _, failure := range []string{"insufficient", "read", "log", "commit", "success"} {
		t.Run(failure, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer sqlDB.Close()
			db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
			if err != nil {
				t.Fatal(err)
			}
			mock.ExpectBegin()
			rows := int64(1)
			if failure == "insufficient" {
				rows = 0
			}
			mock.ExpectExec("UPDATE `geekai_users` SET `power`=power - \\? WHERE id = \\? AND power >= \\?").WithArgs(2, 7, 2).WillReturnResult(sqlmock.NewResult(0, rows))
			if failure != "insufficient" {
				query := mock.ExpectQuery("SELECT .*geekai_users").WithArgs(7)
				if failure == "read" {
					query.WillReturnError(errors.New("read failed"))
				} else {
					query.WillReturnRows(sqlmock.NewRows([]string{"id", "username", "power"}).AddRow(7, "alice", 8))
					insert := mock.ExpectExec("INSERT INTO `geekai_power_logs`")
					if failure == "log" {
						insert.WillReturnError(errors.New("log failed"))
					} else {
						insert.WillReturnResult(sqlmock.NewResult(1, 1))
					}
				}
			}
			switch failure {
			case "success":
				mock.ExpectCommit()
			case "commit":
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			default:
				mock.ExpectRollback()
			}
			err = NewUserService(db).DecreasePower(7, 2, model.PowerLog{})
			if (err == nil) != (failure == "success") {
				t.Fatalf("transaction error: %v", err)
			}
			if failure == "insufficient" && !errors.Is(err, ErrInsufficientPower) {
				t.Fatalf("wrong quota error: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
