package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go/common"
	"go/pb"

	"github.com/cenkalti/backoff/v5"
	"github.com/golang/snappy"
	"github.com/quic-go/quic-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// DataPool represents a thread-safe pool for collecting sensor data
type DataPool struct {
	mu    sync.RWMutex
	items []*pb.ClientData
}

// NewDataPool creates a new data pool
func NewDataPool() *DataPool {
	return &DataPool{
		items: make([]*pb.ClientData, 0),
	}
}

// AddData adds data to the pool in a thread-safe manner
func (dp *DataPool) AddData(data *pb.ClientData) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.items = append(dp.items, data)
}

// GetAndClearData retrieves all data from the pool and clears it
func (dp *DataPool) GetAndClearData() []*pb.ClientData {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if len(dp.items) == 0 {
		return nil
	}

	// Create a copy of the data
	data := make([]*pb.ClientData, len(dp.items))
	copy(data, dp.items)

	// Clear the pool
	dp.items = dp.items[:0]

	return data
}

// GetDataCount returns the current number of items in the pool
func (dp *DataPool) GetDataCount() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return len(dp.items)
}

type QuicClient struct {
	logger   *zap.Logger
	config   *common.Config
	tls      *common.Tls
	conn     *quic.Conn
	dataPool *DataPool
}

// AddExternalData adds external data to the client's data pool
func (qc *QuicClient) AddExternalData(data *pb.ClientData) {
	qc.dataPool.AddData(data)
	qc.logger.Debug("External data added to pool",
		zap.Int("pool_size", qc.dataPool.GetDataCount()),
		zap.String("message", data.Message))
}

func NewQuicClient(logger *zap.Logger, config *common.Config, tls *common.Tls, lifecycle fx.Lifecycle) (*QuicClient, error) {
	client := &QuicClient{
		logger:   logger,
		config:   config,
		tls:      tls,
		dataPool: NewDataPool(),
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
			// Get all data from the pool
			poolData := qc.dataPool.GetAndClearData()

			if len(poolData) == 0 {
				qc.logger.Debug("No data in pool to send")
				continue
			}

			// Create batch message with all collected data
			batchMessage := &pb.ClientData{
				Timestamp:      time.Now().Unix(),
				Message:        fmt.Sprintf("Batch transmission - %d items", len(poolData)),
				SensorReadings: make([]float32, 0),
			}

			// Combine all sensor readings from collected data
			totalSensors := 0
			for _, data := range poolData {
				totalSensors += len(data.SensorReadings)
				batchMessage.SensorReadings = append(batchMessage.SensorReadings, data.SensorReadings...)
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
			protobufData, err := proto.Marshal(batchMessage)
			if err != nil {
				qc.logger.Error("Failed to marshal protobuf message", zap.Error(err))
				return err
			}

			compressedData := snappy.Encode(nil, protobufData)
			compressionRatio := float64(len(compressedData)) / float64(len(protobufData)) * 100

			_, err = stream.Write(compressedData)
			if err != nil {
				qc.logger.Error("Failed to write to stream", zap.Error(err))
				return err
			}

			qc.logger.Info("Batch snappy compressed protobuf transmission",
				zap.Int("batch_items", len(poolData)),
				zap.Int("total_sensors", totalSensors),
				zap.Int("protobuf_size", len(protobufData)),
				zap.Int("compressed_size", len(compressedData)),
				zap.Float64("compression_ratio_percent", compressionRatio),
				zap.Int("size_reduction_bytes", len(protobufData)-len(compressedData)))
		}
	}
}

// generateTestData simulates external data being added to the pool
func (qc *QuicClient) generateTestData() {
	// TODO 10ms 로 테스트시 서버에서 압축해제 에러가 발생함
	// 100ms 로 테스트시 서버에서 압축해제 에러가 발생하지 않음
	// 실제 데이터 적용할때 로직 개선해야됨
	// TODO 데이터 스키마가 정해진 상태가 아니라 임시로직으로 남겨둠.
	// 스키마가 정해지면 로직 자체가 변경되어야 함
	ticker := time.NewTicker(100 * time.Millisecond) // Add data every 100ms
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate test data
			sensorReadings := make([]float32, 100+rand.Intn(200)) // 100-300 sensors
			for i := 0; i < len(sensorReadings); i++ {
				sensorReadings[i] = rand.Float32() * 100.0
			}

			testData := &pb.ClientData{
				Timestamp:      time.Now().Unix(),
				Message:        fmt.Sprintf("Test data %d", rand.Intn(1000)),
				SensorReadings: sensorReadings,
			}

			qc.AddExternalData(testData)
		}
	}
}
