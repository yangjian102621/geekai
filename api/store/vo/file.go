package vo

type File struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Name      string `json:"name"`
	ObjKey    string `json:"obj_key"`
	URL       string `json:"url"`
	Ext       string `json:"ext"`
	Size      int64  `json:"size"`
	CreatedAt int64  `json:"created_at"`
}
