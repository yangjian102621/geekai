package model

import (
	"geekai/core/types"
	"time"
)

// PowerLog 积分消费日志
type PowerLog struct {
	Id       uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId   uint            `gorm:"column:user_id;type:int(11);not null;comment:用户ID" json:"user_id"`
	Username string          `gorm:"column:username;type:varchar(30);not null;comment:用户名" json:"username"`
	Type     types.PowerType `gorm:"column:type;type:tinyint(1);not null;comment:类型（1：充值，2：消费，3：退费）" json:"type"`
	Amount   int             `gorm:"column:amount;type:smallint;not null;comment:积分数值" json:"amount"`
	Balance  int             `gorm:"column:balance;type:int;not null;comment:余额" json:"balance"`
	Model    string          `gorm:"column:model;type:varchar(255);not null;comment:模型" json:"model"`
	Remark   string          `gorm:"column:remark;type:varchar(512);not null;comment:备注" json:"remark"`
	Mark     types.PowerMark `gorm:"column:mark;type:tinyint(1);not null;comment:资金类型（0：支出，1：收入）" json:"mark"`
	// RefType / RefId 是账本幂等键：两者组合成唯一索引 uk_power_ref，
	// 让"同一个业务事件"（同一订单、同一任务）在整个数据库里最多入账一次。
	//
	// 这两个字段【必须】是可空指针（*string / *uint），不能是 string / uint：
	// AutoMigrate 在每次服务启动时对存量表执行 ALTER TABLE ... ADD UNIQUE INDEX，
	// 存量历史行没有业务关联键。若列定义为 NOT NULL DEFAULT ''/0，
	// 所有存量行会取到相同的 ('',0)，建索引时互相撞唯一键，迁移直接失败，
	// 结果是整个服务起不来。用指针时 MySQL 允许唯一索引上存在多个 NULL
	// （NULL 之间互不相等），存量行全部为 NULL 即可安全建索引。
	RefType   *string   `gorm:"column:ref_type;type:varchar(32);uniqueIndex:uk_power_ref;comment:关联业务类型" json:"ref_type,omitempty"`
	RefId     *uint     `gorm:"column:ref_id;type:int(11);uniqueIndex:uk_power_ref;comment:关联业务ID" json:"ref_id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`

	RewardKey        *string    `gorm:"column:reward_key;type:varchar(128);uniqueIndex;default:null" json:"-"`
	QuotaAmount      int64      `gorm:"column:quota_amount;type:bigint;not null;default:0" json:"quota_amount"`
	QuotaStatus      string     `gorm:"column:quota_status;type:varchar(24);not null;default:'';index" json:"quota_status"`
	QuotaSubject     string     `gorm:"column:quota_subject;type:varchar(256);not null;default:''" json:"-"`
	QuotaIdentity    string     `gorm:"column:quota_identity;type:char(64);not null;default:''" json:"-"`
	QuotaError       string     `gorm:"column:quota_error;type:varchar(64);not null;default:''" json:"quota_error"`
	QuotaAttemptedAt *time.Time `gorm:"column:quota_attempted_at;type:datetime" json:"-"`
}

func (m *PowerLog) TableName() string {
	return "geekai_power_logs"
}
