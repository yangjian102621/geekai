package vo

type RewardRule struct {
	Id     int    `json:"id"`     // 规则ID
	Title  string `json:"title"`  // 规则标题
	Desc   string `json:"desc"`   // 规则描述
	Icon   string `json:"icon"`   // 图标类名
	Color  string `json:"color"`  // 图标颜色
	Reward int    `json:"reward"` // 奖励算力
}
