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
		fx.Provide(NewQuicServer, NewRedisClient),
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger, quicServer *QuicServer) {
			logger.Info("Starting application")

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
