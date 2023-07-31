package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if sugarLogger != nil {
		return sugarLogger
	}

	logLevel := zap.NewAtomicLevelAt(getLogLevel(os.Getenv("LOG_LEVEL")))
	encoder := getEncoder()
	writerSyncer := getLogWriter()
	fileCore := zapcore.NewCore(encoder, writerSyncer, logLevel)
	consoleOutput := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(
		encoder,
		consoleOutput,
		logLevel,
	)
	core := zapcore.NewTee(fileCore, consoleCore)
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
	return sugarLogger
}

// core 三个参数之  编码
func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLogLevel(level string) zapcore.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
