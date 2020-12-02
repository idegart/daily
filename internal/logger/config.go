package logger

import "os"

type Config struct {
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}