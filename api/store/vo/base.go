package vo

type BaseVo struct {
	Id        uint  `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
