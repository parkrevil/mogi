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
	Env                                  Environment
	ServerHost                           string
	ServerPort                           int
	ServerStreamBufferSize               int
	ServerMaxIdleTimeout                 int64
	ServerKeepAlivePeriod                int64
	ServerMaxIncomingStreams             int64
	ServerMaxIncomingUniStreams          int64
	ServerInitialStreamReceiveWindow     uint64
	ServerMaxStreamReceiveWindow         uint64
	ServerInitialConnectionReceiveWindow uint64
	ServerMaxConnectionReceiveWindow     uint64
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

	serverHost, err := validateAndGetEnv("SUCTION_SERVER_HOST", "string")
	if err != nil {
		return nil, err
	}

	serverPort, err := validateAndGetEnv("SUCTION_SERVER_PORT", "int")
	if err != nil {
		return nil, err
	}

	serverStreamBufferSize, err := validateAndGetEnv("SUCTION_SERVER_STREAM_BUFFER_SIZE", "int")
	if err != nil {
		return nil, err
	}

	serverMaxIdleTimeout, err := validateAndGetEnv("SUCTION_SERVER_MAX_IDLE_TIMEOUT", "int64")
	if err != nil {
		return nil, err
	}

	serverKeepAlivePeriod, err := validateAndGetEnv("SUCTION_SERVER_KEEP_ALIVE_PERIOD", "int64")
	if err != nil {
		return nil, err
	}

	serverMaxIncomingStreams, err := validateAndGetEnv("SUCTION_SERVER_MAX_INCOMING_STREAMS", "int64")
	if err != nil {
		return nil, err
	}

	serverMaxIncomingUniStreams, err := validateAndGetEnv("SUCTION_SERVER_MAX_INCOMING_UNI_STREAMS", "int64")
	if err != nil {
		return nil, err
	}

	serverInitialStreamReceiveWindow, err := validateAndGetEnv("SUCTION_SERVER_INITIAL_STREAM_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	serverMaxStreamReceiveWindow, err := validateAndGetEnv("SUCTION_SERVER_MAX_STREAM_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	serverInitialConnectionReceiveWindow, err := validateAndGetEnv("SUCTION_SERVER_INITIAL_CONNECTION_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	serverMaxConnectionReceiveWindow, err := validateAndGetEnv("SUCTION_SERVER_MAX_CONNECTION_RECEIVE_WINDOW", "uint64")
	if err != nil {
		return nil, err
	}

	return &Config{
		Env:                                  env,
		ServerHost:                           serverHost.(string),
		ServerPort:                           serverPort.(int),
		ServerStreamBufferSize:               serverStreamBufferSize.(int),
		ServerMaxIdleTimeout:                 serverMaxIdleTimeout.(int64),
		ServerKeepAlivePeriod:                serverKeepAlivePeriod.(int64),
		ServerMaxIncomingStreams:             serverMaxIncomingStreams.(int64),
		ServerMaxIncomingUniStreams:          serverMaxIncomingUniStreams.(int64),
		ServerInitialStreamReceiveWindow:     serverInitialStreamReceiveWindow.(uint64),
		ServerMaxStreamReceiveWindow:         serverMaxStreamReceiveWindow.(uint64),
		ServerInitialConnectionReceiveWindow: serverInitialConnectionReceiveWindow.(uint64),
		ServerMaxConnectionReceiveWindow:     serverMaxConnectionReceiveWindow.(uint64),
	}, nil
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
