package function

type Function interface {
	Invoke(...interface{}) (string, error)
	Name() string
}

type resVo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
