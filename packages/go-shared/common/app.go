package common

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Name string
	Port int
}

// NewAppConfig creates a new app configuration
func NewAppConfig() *AppConfig {
	return &AppConfig{
		Name: "default-app",
		Port: 8080,
	}
}

// StartApp starts the fx application
func StartApp(logger *zap.Logger, config *AppConfig) error {
	logger.Info("Starting application",
		zap.String("name", config.Name),
		zap.Int("port", config.Port))
	return nil
}

// StopApp stops the fx application
func StopApp(logger *zap.Logger) error {
	logger.Info("Stopping application")
	return nil
}

// NewApp creates a new fx application with common providers
func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			NewLoggerProvider,
			NewAppConfig,
		),
		fx.Invoke(StartApp),
		fx.Invoke(StopApp),
	)
}
