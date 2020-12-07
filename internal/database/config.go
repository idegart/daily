package database

import (
	"SlackBot/internal/env"
	"fmt"
)

type Config struct {
	DatabaseUrl string
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			env.Get("DB_HOST", ""),
			env.Get("DB_PORT", ""),
			env.Get("DB_USERNAME", ""),
			env.Get("DB_PASSWORD", ""),
			env.Get("DB_DATABASE", ""),
		),
	}
}