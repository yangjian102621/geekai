package vo

type InviteStats struct {
	InviteCount int    `json:"invite_count"` // 累计邀请数
	RewardTotal int    `json:"reward_total"` // 获得奖励总数
	TodayInvite int    `json:"today_invite"` // 今日邀请数
	InviteCode  string `json:"invite_code"`  // 邀请码
	InviteLink  string `json:"invite_link"`  // 邀请链接
}
