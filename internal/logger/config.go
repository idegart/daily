package logger

type Config struct {
	LogLevel string
}

func NewConfig(level string) *Config {
	return &Config{
		LogLevel: level,
	}
}