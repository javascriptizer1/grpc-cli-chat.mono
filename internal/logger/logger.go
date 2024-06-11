package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(env string) {
	atomic := zap.NewAtomicLevel()

	var enconfig zapcore.EncoderConfig
	var enc zapcore.Encoder

	switch env {
	case "production":
		atomic.SetLevel(zapcore.InfoLevel)
		enconfig = zap.NewProductionEncoderConfig()
		enc = zapcore.NewJSONEncoder(enconfig)

	default:
		atomic.SetLevel(zapcore.DebugLevel)
		enconfig = zap.NewDevelopmentEncoderConfig()
		enc = zapcore.NewConsoleEncoder(enconfig)
	}

	globalLogger = zap.New(zapcore.NewCore(
		enc,
		zapcore.Lock(os.Stdout),
		atomic,
	))
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
