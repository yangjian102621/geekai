package model

// 用户打赏

type Reward struct {
	BaseModel
	UserId   uint    // 用户 ID
	TxId     string  // 交易ID
	Amount   float64 // 打赏金额
	Remark   string  // 打赏备注
	Status   bool    // 核销状态
	Exchange string  // 众筹兑换详情，JSON
}
