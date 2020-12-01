package botserver

import (
	"SlackBot/internal/store"
	"os"
)

type Config struct {
	BindAddr string
	ApiToken string
	LogLevel string
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":" + os.Getenv("BIND_INTERNAL_ADDR"),
		ApiToken: os.Getenv("API_TOKEN"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		Store:    store.NewConfig(),
	}
}
