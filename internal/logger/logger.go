package logger

import (
	"bot/internal/config"
	"github.com/sirupsen/logrus"
)

func NewLogger(config *config.Logger) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(config.LogLevel)

	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	logger.SetLevel(level)

	return logger, nil
}