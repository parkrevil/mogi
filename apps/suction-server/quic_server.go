package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"go/common"
	"go/pb"

	"github.com/golang/snappy"
	"github.com/quic-go/quic-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type QuicServer struct {
	listener   *quic.Listener
	logger     *zap.Logger
	BufferSize int
}

func NewQuicServer(logger *zap.Logger, config *common.Config, tls *common.Tls, lifecycle fx.Lifecycle) (*QuicServer, error) {
	quicConfig := &quic.Config{
		MaxIdleTimeout:                 time.Duration(config.QuicMaxIdleTimeout) * time.Second,
		KeepAlivePeriod:                time.Duration(config.QuicKeepAlivePeriod) * time.Second,
		MaxIncomingStreams:             config.QuicServerMaxIncomingStreams,
		MaxIncomingUniStreams:          config.QuicServerMaxIncomingUniStreams,
		InitialStreamReceiveWindow:     config.QuicServerInitialStreamReceiveWindow * 1024 * 1024,     // 스트림 당 초기 버퍼(n MB)
		MaxStreamReceiveWindow:         config.QuicServerMaxStreamReceiveWindow * 1024 * 1024,         // 스트림 당 최대 버퍼(n MB)
		InitialConnectionReceiveWindow: config.QuicServerInitialConnectionReceiveWindow * 1024 * 1024, // 연결 당 초기 버퍼(n MB)
		MaxConnectionReceiveWindow:     config.QuicServerMaxConnectionReceiveWindow * 1024 * 1024,     // 연결 당 최대 버퍼(n MB)
	}

	listener, err := quic.ListenAddr(config.QuicServerListeningAddress, tls.Config, quicConfig)
	if err != nil {
		return nil, errors.New("Failed to start QUIC server: " + err.Error())
	}

	server := &QuicServer{
		listener:   listener,
		logger:     logger,
		BufferSize: config.QuicServerStreamBufferSize,
	}

	var cancel context.CancelFunc

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting QUIC server", zap.String("address", config.QuicServerListeningAddress))

			quicCtx, c := context.WithCancel(context.Background())
			cancel = c

			go server.acceptConnections(quicCtx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping QUIC server")

			if cancel != nil {
				cancel()
			}

			return listener.Close()
		},
	})

	return server, nil
}

func (qs *QuicServer) acceptConnections(ctx context.Context) {
	for {
		conn, err := qs.listener.Accept(ctx)

		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, context.Canceled) {
				qs.logger.Info("QUIC listener closed", zap.Error(err))

				return
			}

			qs.logger.Error("Failed to accept connection", zap.Error(err))

			continue
		}

		qs.logger.Debug("New QUIC connection accepted", zap.String("remote", conn.RemoteAddr().String()))

		go qs.handleConnection(conn)
	}
}

func (qs *QuicServer) handleConnection(conn *quic.Conn) {
	defer func() {
		if err := conn.CloseWithError(0, "connection closed"); err != nil {
			qs.logger.Debug("Failed to close connection", zap.Error(err))
		}
	}()

	qs.logger.Debug("Handling new connection", zap.String("remote", conn.RemoteAddr().String()))

	for {
		stream, err := conn.AcceptStream(conn.Context())
		if err != nil {
			qs.logger.Error("Failed to accept stream", zap.Error(err))
			return
		}

		qs.logger.Debug("New stream accepted", zap.Uint64("stream_id", uint64(stream.StreamID())))

		go qs.handleStream(stream)
	}
}

func (qs *QuicServer) handleStream(stream *quic.Stream) {
	defer stream.Close()

	buffer := make([]byte, qs.BufferSize)

	for {
		n, err := stream.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				qs.logger.Debug("Stream closed by peer", zap.Uint64("stream_id", uint64(stream.StreamID())))

				return
			}

			qs.logger.Error("Failed to read from stream", zap.Error(err))

			return
		}

		compressedData := buffer[:n]

		// TODO: Snappy Decompression Testing Requirements
		// 1. Decompression Performance Testing:
		//    - Measure decompression time for different data sizes
		//    - Test with various compression ratios
		//    - Monitor memory usage during decompression
		//    - Test CPU usage patterns during decompression
		//
		// 2. Error Handling Testing:
		//    - Test with corrupted compressed data
		//    - Test with incomplete compressed data
		//    - Test memory exhaustion during decompression
		//    - Test with malformed protobuf data after decompression
		//
		// 3. Concurrent Processing Testing:
		//    - Test multiple concurrent decompression operations
		//    - Test memory usage under high load
		//    - Test performance degradation with multiple streams
		//
		// 4. Monitoring and Metrics:
		//    - Track decompression success/failure rates
		//    - Monitor average decompression time
		//    - Track compression ratio statistics
		//    - Monitor memory usage patterns
		//
		// 5. Optimization Opportunities:
		//    - Implement decompression buffer pooling
		//    - Test different buffer sizes for optimal performance
		//    - Consider async decompression for large data
		//    - Implement decompression timeout handling

		// Decompress data with snappy
		decompressedData, err := snappy.Decode(nil, compressedData)
		if err != nil {
			qs.logger.Error("Failed to decompress snappy data", zap.Error(err))
			continue
		}

		// Unmarshal protobuf message
		var clientData pb.ClientData
		if err := proto.Unmarshal(decompressedData, &clientData); err != nil {
			qs.logger.Error("Failed to unmarshal protobuf message", zap.Error(err))
			continue
		}

		// Calculate compression statistics
		compressionRatio := float64(len(compressedData)) / float64(len(decompressedData)) * 100
		sizeReduction := len(decompressedData) - len(compressedData)

		qs.logger.Info("Received snappy compressed protobuf message",
			zap.Int64("timestamp", clientData.Timestamp),
			zap.Int("message_length", len(clientData.Message)),
			zap.Int("sensor_readings_count", len(clientData.SensorReadings)),
			zap.Int("original_size", len(decompressedData)),
			zap.Int("compressed_size", len(compressedData)),
			zap.Float64("compression_ratio_percent", compressionRatio),
			zap.Int("size_reduction_bytes", sizeReduction),
			zap.Uint64("stream_id", uint64(stream.StreamID())))

		// Send response back to client
		response := fmt.Sprintf("Server received compressed protobuf (original: %d bytes, compressed: %d bytes, ratio: %.1f%%)",
			len(decompressedData), len(compressedData), compressionRatio)
		_, err = stream.Write([]byte(response))
		if err != nil {
			qs.logger.Error("Failed to write to stream", zap.Error(err))
			return
		}
	}
}
