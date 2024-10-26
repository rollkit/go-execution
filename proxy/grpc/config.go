package grpc

import "time"

type Config struct {
	JWTSecret      []byte
	DefaultTimeout time.Duration
	MaxRequestSize int
}

func DefaultConfig() *Config {
	return &Config{
		DefaultTimeout: time.Second,
		MaxRequestSize: 1024 * 1024,
	}
}
