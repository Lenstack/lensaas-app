package infrastructure

import (
	"go.uber.org/zap"
	"strings"
)

type Logger struct {
	environment string
	Log         *zap.Logger
}

func NewLogger(environment string) *Logger {
	if strings.ToLower(environment) == "production" {
		logger, _ := zap.NewProduction()
		return &Logger{environment, logger}
	}

	logger, _ := zap.NewDevelopment()
	return &Logger{Log: logger}
}
