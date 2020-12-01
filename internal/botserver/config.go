package botserver

import (
	"SlackBot/internal/slack"
	"SlackBot/internal/store"
	"os"
)

type Config struct {
	BindAddr string
	LogLevel string
	Store    *store.Config
	Slack    *slack.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":" + os.Getenv("BIND_INTERNAL_ADDR"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		Store:    store.NewConfig(),
		Slack:    slack.NewConfig(),
	}
}
