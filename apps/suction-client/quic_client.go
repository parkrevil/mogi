package main

import (
	"context"
	"errors"
	"time"

	"go-shared/common"
	"go-shared/pb"

	"github.com/cenkalti/backoff/v5"
	"github.com/quic-go/quic-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type QuicClient struct {
	logger *zap.Logger
	config *common.Config
	tls    *common.Tls
	conn   *quic.Conn
}

func NewQuicClient(logger *zap.Logger, config *common.Config, tls *common.Tls, lifecycle fx.Lifecycle) (*QuicClient, error) {
	client := &QuicClient{
		logger: logger,
		config: config,
		tls:    tls,
	}

	var cancel context.CancelFunc

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting QUIC client, attempting to connect...")

			quicCtx, c := context.WithCancel(context.Background())
			cancel = c

			go client.start(quicCtx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping QUIC client")

			if cancel != nil {
				cancel()
			}

			if client.conn != nil {
				return client.conn.CloseWithError(0, "client closing")
			}

			return nil
		},
	})

	return client, nil
}

func (qc *QuicClient) start(ctx context.Context) {
	operation := func() (string, error) {
		qc.logger.Info("Attempting to connect to QUIC server...")

		conn, err := qc.connect(ctx)
		if err != nil {
			return "", err
		}
		qc.conn = conn
		qc.logger.Info("Connected to QUIC server", zap.String("remote", qc.conn.RemoteAddr().String()))

		return "", qc.handleConnection(ctx)
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 5 * time.Second
	b.MaxInterval = 60 * time.Second

	options := []backoff.RetryOption{
		backoff.WithBackOff(b),
		backoff.WithMaxElapsedTime(0),
		backoff.WithNotify(func(err error, d time.Duration) {
			qc.logger.Error("Connection failed. Retrying...", zap.Error(err), zap.String("retry_in", d.String()))
		}),
	}

	if _, err := backoff.Retry(ctx, operation, options...); err != nil {
		if errors.Is(err, context.Canceled) {
			qc.logger.Info("Client shutdown initiated.")
		} else {
			qc.logger.Fatal("Critical: Connection failed permanently.", zap.Error(err))
		}
	}
}

func (qc *QuicClient) connect(ctx context.Context) (*quic.Conn, error) {
	quicConfig := &quic.Config{
		MaxIdleTimeout:                 time.Duration(qc.config.QuicMaxIdleTimeout) * time.Second,
		KeepAlivePeriod:                time.Duration(qc.config.QuicKeepAlivePeriod) * time.Second,
		MaxIncomingStreams:             qc.config.QuicServerMaxIncomingStreams,
		MaxIncomingUniStreams:          qc.config.QuicServerMaxIncomingUniStreams,
		InitialStreamReceiveWindow:     qc.config.QuicServerInitialStreamReceiveWindow,
		MaxStreamReceiveWindow:         qc.config.QuicServerMaxStreamReceiveWindow,
		InitialConnectionReceiveWindow: qc.config.QuicServerInitialConnectionReceiveWindow,
		MaxConnectionReceiveWindow:     qc.config.QuicServerMaxConnectionReceiveWindow,
	}

	return quic.DialAddr(ctx, qc.config.QuicClientConnectionAddress, qc.tls.Config, quicConfig)
}

func (qc *QuicClient) handleConnection(ctx context.Context) error {
	qc.logger.Info("handleConnection started...")

	stream, err := qc.conn.OpenStreamSync(ctx)
	if err != nil {
		qc.logger.Error("Failed to open stream", zap.Error(err))
		return err
	}
	defer stream.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			qc.logger.Info("handleConnection received shutdown signal, stopping...",
				zap.String("reason", ctx.Err().Error()))
			return ctx.Err()
		case <-ticker.C:
			// Create protobuf message with sensor data
			clientData := &pb.ClientData{
				Timestamp:      time.Now().Unix(),
				Message:        "Hello, QUIC server! Message",
				SensorReadings: []float32{23.5, 45.2, 67.8, 89.1}, // Example sensor readings
			}
			
			data, err := proto.Marshal(clientData)
			if err != nil {
				qc.logger.Error("Failed to marshal protobuf message", zap.Error(err))
				return err
			}

			_, err = stream.Write(data)
			if err != nil {
				qc.logger.Error("Failed to write to stream", zap.Error(err))

				return err
			}

			qc.logger.Info("Sent protobuf message to server", 
				zap.Int64("timestamp", clientData.Timestamp),
				zap.String("message", clientData.Message),
				zap.Int("sensor_readings_count", len(clientData.SensorReadings)))
		}
	}
}
// test comment
