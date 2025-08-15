package common

import (
	"errors"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Local      Environment = "local"
	Production Environment = "production"
)

type Config struct {
	Env                                      Environment
	QuicMaxIdleTimeout                       int64
	QuicKeepAlivePeriod                      int64
	QuicServerListeningAddress               string
	QuicServerStreamBufferSize               int
	QuicServerMaxIncomingStreams             int64
	QuicServerMaxIncomingUniStreams          int64
	QuicServerInitialStreamReceiveWindow     uint64
	QuicServerMaxStreamReceiveWindow         uint64
	QuicServerInitialConnectionReceiveWindow uint64
	QuicServerMaxConnectionReceiveWindow     uint64
	QuicClientConnectionAddress              string
}

func NewConfig(logger *zap.Logger) (*Config, error) {
	envStr := os.Getenv("ENV")
	if envStr == "" {
		return nil, errors.New("ENV environment variable is required")
	}

	env := Environment(envStr)
	if env != Local && env != Production {
		return nil, errors.New("Invalid ENV value: " + envStr)
	}

	if env == Local {
		if err := godotenv.Load("../../.env.local"); err != nil {
			return nil, errors.New("Could not load .env.local file: " + err.Error())
		}
	}

	quicMaxIdleTimeout, err := validateAndGetEnv("SUCTION_QUIC_MAX_IDLE_TIMEOUT", "int64")
	if err != nil {
		return nil, err
	}

	quicKeepAlivePeriod, err := validateAndGetEnv("SUCTION_QUIC_KEEP_ALIVE_PERIOD", "int64")
	if err != nil {
		return nil, err
	}

	quicServerListeningAddress, err := validateAndGetEnv("SUCTION_QUIC_SERVER_LISTENING_ADDRESS", "string")
	if err != nil {
		return nil, err
	}

	quicServerStreamBufferSize, err := validateAndGetEnv("SUCTION_QUIC_SERVER_STREAM_BUFFER_SIZE", "int")
	if err != nil {
		return nil, err
	}

	quicServerMaxIncomingStreams, err := validateAndGetEnv("SUCTION_QUIC_SERVER_MAX_INCOMING_STREAMS", "int64")
	if err != nil {
		return nil, err
	}

	quicServerMaxIncomingUniStreams, err := validateAndGetEnv("SUCTION_QUIC_SERVER_MAX_INCOMING_UNI_STREAMS", "int64")
	if err != nil {
		return nil, err
	}

	quicServerInitialStreamReceiveWindow, err := validateAndGetEnv("SUCTION_QUIC_SERVER_INITIAL_STREAM_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	quicServerMaxStreamReceiveWindow, err := validateAndGetEnv("SUCTION_QUIC_SERVER_MAX_STREAM_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	quicServerInitialConnectionReceiveWindow, err := validateAndGetEnv("SUCTION_QUIC_SERVER_INITIAL_CONNECTION_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	quicServerMaxConnectionReceiveWindow, err := validateAndGetEnv("SUCTION_QUIC_SERVER_MAX_CONNECTION_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	quicClientConnectionAddress, err := validateAndGetEnv("SUCTION_QUIC_CLIENT_CONNECTION_ADDRESS", "string")
	if err != nil {
		return nil, err
	}

	return &Config{
		Env:                                      env,
		QuicMaxIdleTimeout:                       quicMaxIdleTimeout.(int64),
		QuicKeepAlivePeriod:                      quicKeepAlivePeriod.(int64),
		QuicServerListeningAddress:               quicServerListeningAddress.(string),
		QuicServerStreamBufferSize:               quicServerStreamBufferSize.(int),
		QuicServerMaxIncomingStreams:             quicServerMaxIncomingStreams.(int64),
		QuicServerMaxIncomingUniStreams:          quicServerMaxIncomingUniStreams.(int64),
		QuicServerInitialStreamReceiveWindow:     quicServerInitialStreamReceiveWindow.(uint64),
		QuicServerMaxStreamReceiveWindow:         quicServerMaxStreamReceiveWindow.(uint64),
		QuicServerInitialConnectionReceiveWindow: quicServerInitialConnectionReceiveWindow.(uint64),
		QuicServerMaxConnectionReceiveWindow:     quicServerMaxConnectionReceiveWindow.(uint64),
		QuicClientConnectionAddress:              quicClientConnectionAddress.(string),
	}, nil
}

func (c *Config) IsLocal() bool {
	return c.Env == Local
}

func (c *Config) IsProduction() bool {
	return c.Env == Production
}

func validateAndGetEnv(key string, valType string) (interface{}, error) {
	val := os.Getenv(key)

	if val == "" {
		return nil, errors.New("\"" + key + "\" is required")
	}

	switch valType {
	case "int":
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("Invalid " + key + ": " + err.Error())
		}
		return intVal, nil
	case "int64":
		int64Val, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, errors.New("Invalid " + key + ": " + err.Error())
		}
		return int64Val, nil
	case "uint64":
		uint64Val, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, errors.New("Invalid " + key + ": " + err.Error())
		}
		return uint64Val, nil
	case "string":
		return val, nil
	default:
		return nil, errors.New("Invalid value type: " + valType)
	}
}
