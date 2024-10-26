package jsonrpc

import "time"

type Config struct {
	DefaultTimeout time.Duration
	MaxRequestSize int64
}

func DefaultConfig() *Config {
	return &Config{
		DefaultTimeout: time.Second,
		MaxRequestSize: 1024 * 1024, // 1MB
	}
}
