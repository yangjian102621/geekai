package model

// 用户打赏

type Reward struct {
	BaseModel
	TxId   string  // 交易ID
	Amount float64 // 打赏金额
	Remark string  // 打赏备注
	Status bool    // 核销状态
}
