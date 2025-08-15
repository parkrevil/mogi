package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"go-shared/common"

	"github.com/quic-go/quic-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

		data := buffer[:n]
		qs.logger.Debug("Received data from stream",
			zap.String("data", string(data)),
			zap.Uint64("stream_id", uint64(stream.StreamID())))

		response := fmt.Sprintf("Echo: %s", string(data))
		_, err = stream.Write([]byte(response))
		if err != nil {
			qs.logger.Error("Failed to write to stream", zap.Error(err))
			return
		}
	}
}
