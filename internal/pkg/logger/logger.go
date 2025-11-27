package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init() error {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	// Извлекаем request_id из контекста
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, "request_id", requestID)
	}
	globalLogger.Sugar().Infow(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, "request_id", requestID)
	}
	globalLogger.Sugar().Errorw(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, "request_id", requestID)
	}
	globalLogger.Sugar().Fatalw(msg, fields...)
}

func Sync() {
	err := globalLogger.Sync()
	if err != nil {
		return
	}
}
