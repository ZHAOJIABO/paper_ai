package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init 初始化日志
func Init() error {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

// Info 记录info级别日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error 记录error级别日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Warn 记录warn级别日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Debug 记录debug级别日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Fatal 记录fatal级别日志并退出
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync 同步日志缓冲
func Sync() error {
	return log.Sync()
}
