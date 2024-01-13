package xxl

import (
	"fmt"
	"log"
)

// LogFunc 应用日志
type LogFunc func(req LogReq, res *LogRes) []byte

// Logger 系统日志
type Logger interface {
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
}

type logger struct {
}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a...))
}
