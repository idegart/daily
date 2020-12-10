package logger

import "SlackBot/internal/env"

type Config struct {
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		LogLevel: env.Get("LOG_LEVEL", "debug"),
	}
}