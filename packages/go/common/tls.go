package common

import (
	"crypto/tls"
	"errors"
)

type Tls struct {
	Config *tls.Config
}

func NewTLSConfig(config *Config) (*Tls, error) {
	if config.IsLocal() {
		tlsCert, err := tls.LoadX509KeyPair("../../samples/server.crt", "../../samples/server.key")
		if err != nil {
			return nil, err
		}

		return &Tls{Config: &tls.Config{
			Certificates:       []tls.Certificate{tlsCert},
			NextProtos:         []string{"suction-quic"},
			InsecureSkipVerify: true,
		}}, nil
	}

	// TODO: production environment support

	return nil, errors.New("support only local environment")
}
