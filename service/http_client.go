package service

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"golang.org/x/net/proxy"
)

const (
	tlsHandshakeTimeout = 5 * time.Second
	timeout             = 10 * time.Second
)

// NewHTTPClient creates new http client
func NewHTTPClient(vars Vars, logger *zap.Logger) *http.Client {
	dialer, err := proxy.SOCKS5(
		"tcp",
		vars.ProxyAddress,
		&proxy.Auth{
			User:     vars.ProxyLogin,
			Password: vars.ProxyPass,
		},
		proxy.Direct)
	if err != nil {
		logger.Fatal("Failed to create HTTP client", zap.Error(err))
	}

	netTransport := &http.Transport{
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}
}
