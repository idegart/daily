package logger

import (
	"bot/internal/config"
)

func NewConfig(level string) *config.Logger {
	return &config.Logger{
		LogLevel: level,
	}
}