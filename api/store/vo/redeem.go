package vo

type Redeem struct {
	Id         uint   `json:"id"`
	UserId     uint   `json:"user_id"` // 用户 ID
	Name       string `json:"name"`
	Username   string `json:"username"`
	Power      int    `json:"power"` // 算力
	Code       string `json:"code"`  // 兑换码
	Enabled    bool   `json:"enabled"`
	RedeemedAt int64  `json:"redeemed_at"` // 兑换时间
	CreatedAt  int64  `json:"created_at"`
}
