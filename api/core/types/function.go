package types

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Function struct {
	Name        string
	Description string
	Parameters  []Parameter
}

type Parameter struct {
	Type       string
	Required   []string
	Properties map[string]interface{}
}
