package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func New() {
	logger, _ = zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func Fatal(msg string, err error) {
	logger.Fatal(msg, zap.Error(err))
}

func Sync() {
	_ = logger.Sync()
}
