package vo

type File struct {
	Id        uint
	UserId    uint   `json:"user_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Ext       string `json:"ext"`
	Size      int64  `json:"size"`
	CreatedAt int64  `json:"created_at"`
}
