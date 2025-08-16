package main

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"go-shared/common"
	"go-shared/pb"

	"github.com/cenkalti/backoff/v5"
	"github.com/golang/snappy"
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

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			qc.logger.Info("handleConnection received shutdown signal, stopping...",
				zap.String("reason", ctx.Err().Error()))
			return ctx.Err()
		case <-ticker.C:
			// Generate simple test data
			sensorReadings := make([]float32, 1000) // 1000 sensor readings
			for i := 0; i < 1000; i++ {
				sensorReadings[i] = rand.Float32() * 100.0
			}
			
			// Create protobuf message with simple data
			clientData := &pb.ClientData{
				Timestamp:      time.Now().Unix(),
				Message:        "Test message from client",
				SensorReadings: sensorReadings,
			}
			
			// Marshal to protobuf
			protobufData, err := proto.Marshal(clientData)
			if err != nil {
				qc.logger.Error("Failed to marshal protobuf message", zap.Error(err))
				return err
			}

			// TODO: Snappy Compression Testing Requirements
			// 1. Data Size Testing:
			//    - Test with 1KB, 5KB, 10KB, 20KB, 50KB, 100KB data sizes
			//    - Compare compression ratios for different data volumes
			//    - Find optimal data size threshold for compression
			//
			// 2. Data Pattern Testing:
			//    - Test with repetitive data patterns (e.g., constant sensor values)
			//    - Test with random data patterns (current implementation)
			//    - Test with real sensor data patterns (temperature, humidity, pressure)
			//    - Test with structured data patterns
			//
			// 3. Performance Testing:
			//    - Measure compression time for different data sizes
			//    - Measure decompression time on server side
			//    - Monitor memory usage during compression/decompression
			//    - Test CPU usage patterns
			//
			// 4. Network Performance Testing:
			//    - Compare transmission time with/without compression
			//    - Measure bandwidth usage reduction
			//    - Test latency impact of compression
			//
			// 5. Error Handling Testing:
			//    - Test with corrupted compressed data
			//    - Test memory exhaustion scenarios
			//    - Test network error handling
			//
			// 6. Optimization Opportunities:
			//    - Implement conditional compression (only compress if beneficial)
			//    - Compare with other compression algorithms (gzip, lz4)
			//    - Test different compression levels
			//    - Implement adaptive compression based on data characteristics
			
			// Compress data with snappy
			compressedData := snappy.Encode(nil, protobufData)
			compressionRatio := float64(len(compressedData)) / float64(len(protobufData)) * 100

			_, err = stream.Write(compressedData)
			if err != nil {
				qc.logger.Error("Failed to write to stream", zap.Error(err))

				return err
			}

			qc.logger.Info("Snappy compressed protobuf transmission", 
				zap.Int("sensor_count", len(sensorReadings)),
				zap.Int("protobuf_size", len(protobufData)),
				zap.Int("compressed_size", len(compressedData)),
				zap.Float64("compression_ratio_percent", compressionRatio),
				zap.Int("size_reduction_bytes", len(protobufData)-len(compressedData)))
		}
	}
}
// test comment
