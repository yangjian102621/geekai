package vo

import "geekai/core/types"

type PowerLog struct {
	Id          uint            `json:"id"`
	UserId      uint            `json:"user_id"`
	Username    string          `json:"username"`
	Type        types.PowerType `json:"type"`
	TypeStr     string          `json:"type_str"`
	Amount      int             `json:"amount"`
	Mark        types.PowerMark `json:"mark"`
	Balance     int             `json:"balance"`
	Model       string          `json:"model"`
	Remark      string          `json:"remark"`
	CreatedAt   int64           `json:"created_at"`
	QuotaAmount int64           `json:"quota_amount"`
	QuotaStatus string          `json:"quota_status"`
	QuotaError  string          `json:"quota_error"`
}
