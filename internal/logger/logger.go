package logger

import "github.com/sirupsen/logrus"

func NewLogger(config *Config) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(config.LogLevel)

	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	logger.SetLevel(level)

	return logger, nil
}
