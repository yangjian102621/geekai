package vo

type Parameters struct {
	Type       string              `json:"type"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Function struct {
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Label       string     `json:"label"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
	Required    []string   `json:"required"`
	Action      string     `json:"action"`
	Token       string     `json:"token"`
	Enabled     bool       `json:"enabled"`
}
