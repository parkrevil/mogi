package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger with development configuration
func NewLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config.Build()
}

// NewLoggerProvider provides a logger for fx dependency injection
func NewLoggerProvider() (*zap.Logger, error) {
	return NewLogger()
}
