package airtable

import "SlackBot/internal/env"

type Config struct {
	APIKey     string
	BaseID     string
	UsersTable string
}

func NewConfig() *Config {
	return &Config{
		APIKey:     env.Get("AIRTABLE_API_KEY", ""),
		BaseID:     env.Get("AIRTABLE_BASE_ID", ""),
		UsersTable: env.Get("AIRTABLE_USERS_TABLE", ""),
	}
}
