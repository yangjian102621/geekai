package vo

type Reward struct {
	BaseVo
	UserId   uint           `json:"user_id"` // 用户 ID
	Username string         `json:"username"`
	TxId     string         `json:"tx_id"`  // 交易ID
	Amount   float64        `json:"amount"` // 打赏金额
	Remark   string         `json:"remark"` // 打赏备注
	Status   bool           `json:"status"` // 核销状态
	Exchange RewardExchange `json:"exchange"`
}

type RewardExchange struct {
	Power int `json:"power"`
}
