package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if logger != nil {
		return logger
	}

	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(zap.InfoLevel)
	log, _ := zap.Config{
		Level:            logLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	logger = log.Sugar()
	return logger
}
