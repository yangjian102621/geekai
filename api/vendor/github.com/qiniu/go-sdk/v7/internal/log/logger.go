// Package log只是SDK本身自己使用，用来调试代码使用，比如输出HTTP请求和响应信息
package log

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	level LogLevel
}

// New 返回一个Logger 指针
func New(out io.Writer, prefix string, flag int, level LogLevel) *Logger {
	return &Logger{
		Logger: log.New(out, prefix, flag),
		level:  level,
	}
}

var (
	std  = New(os.Stdout, DebugPrefix, log.LstdFlags, LogDebug)
	info = New(os.Stdout, InfoPrefix, log.LstdFlags, LogInfo)
	warn = New(os.Stdout, WarnPrefix, log.LstdFlags, LogWarn)
)

type LogLevel int

const (
	// LogDebug 调试模式
	LogDebug LogLevel = iota

	// Info
	LogInfo

	// Warn
	LogWarn
)

const (
	InfoPrefix  = "[I] "
	DebugPrefix = "[D] "
	WarnPrefix  = "[W] "
)

func (l *Logger) Info(v ...interface{}) {
	l.output(LogInfo, v...)
}

func (l *Logger) output(level LogLevel, v ...interface{}) {
	if l.level <= level {
		l.Logger.Println(v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.output(LogDebug, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.output(LogWarn, v...)
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Info(v ...interface{}) {
	info.Info(v...)
}

func Warn(v ...interface{}) {
	warn.Warn(v...)
}
