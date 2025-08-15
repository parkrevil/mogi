package main

import (
	"context"
	"go-shared/common"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		common.Module,
		fx.Provide(NewQuicServer),
		fx.Invoke(func(logger *zap.Logger, quicServer *QuicServer) {
			logger.Info("Starting application")
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
