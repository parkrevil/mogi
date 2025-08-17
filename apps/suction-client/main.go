package main

import (
	"context"
	"go/common"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		common.Module,
		fx.Provide(NewQuicClient),
		fx.Invoke(func(logger *zap.Logger, quicClient *QuicClient) {
			logger.Info("Starting application")
			
			// Start test data generation in background
			go quicClient.generateTestData()
		}),
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					logger.Info("Stopping application")
					return nil
				},
			})
		}),
	)

	app.Run()
}
