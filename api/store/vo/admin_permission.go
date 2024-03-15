package vo

type AdminPermission struct {
	Id       int               `json:"id"`
	Name     string            `json:"name"`
	Slug     string            `json:"slug"`
	Sort     int               `json:"sort"`
	Pid      int               `json:"pid"`
	Children []AdminPermission `json:"children"`
}
