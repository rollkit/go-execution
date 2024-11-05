package jsonrpc

import "time"

// Config represents configuration settings for server/client.
type Config struct {
	DefaultTimeout time.Duration
	MaxRequestSize int64
}

// DefaultConfig returns Config struct initialized with default settings.
func DefaultConfig() *Config {
	return &Config{
		DefaultTimeout: time.Second,
		MaxRequestSize: 1024 * 1024, // 1MB
	}
}
