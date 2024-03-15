package vo

type AdminRole struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Permissions interface{} `json:"permissions"`
	CreatedAt   string      `json:"created_at"`
}
