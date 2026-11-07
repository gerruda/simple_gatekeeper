package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.Logger
}

func NewLogger(ctx context.Context) *Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return &Logger{
		Logger: logger,
	}
}
