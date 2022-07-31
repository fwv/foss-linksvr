package zlog

import "go.uber.org/zap"

func Info(msg string, fields ...zap.Field) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Fatal(msg, fields...)
}
