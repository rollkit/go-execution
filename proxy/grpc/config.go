package grpc

import "time"

// Config holds configuration settings for the gRPC proxy.
type Config struct {
	JWTSecret      []byte
	DefaultTimeout time.Duration
	MaxRequestSize int
}

// DefaultConfig returns a Config instance populated with default settings.
func DefaultConfig() *Config {
	return &Config{
		DefaultTimeout: time.Second,
		MaxRequestSize: 1024 * 1024,
	}
}
