package slack

import "os"

type Config struct {
	ApiToken string
}

func NewConfig() *Config {
	return &Config{
		ApiToken: os.Getenv("SLACK_API_TOKEN"),
	}
}