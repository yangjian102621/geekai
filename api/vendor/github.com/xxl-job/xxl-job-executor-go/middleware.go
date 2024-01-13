package xxl

// Middleware 中间件构造函数
type Middleware func(TaskFunc) TaskFunc

func (e *executor) chain(next TaskFunc) TaskFunc {
	for i := range e.middlewares {
		next = e.middlewares[len(e.middlewares)-1-i](next)
	}
	return next
}
