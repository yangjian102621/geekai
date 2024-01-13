package api

// QINIU SDK error type

// 可以根据Code判断是何种类型错误
type QError struct {
	Code    string
	Message string
}

// Error 继承error接口
func (e *QError) Error() string {
	return e.Code + ": " + e.Message
}

// NewError 返回QError指针
func NewError(code, message string) *QError {
	return &QError{
		Code:    code,
		Message: message,
	}
}
