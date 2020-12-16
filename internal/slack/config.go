package slack

import (
	"SlackBot/internal/env"
)

type Config struct {
	ApiToken string
	VerificationToken string
	SigningSecret string
}

func NewConfig() *Config {
	return &Config{
		ApiToken: env.Get("SLACK_API_TOKEN", ""),
		VerificationToken: env.Get("SLACK_VERIFICATION_TOKEN", ""),
		SigningSecret: env.Get("SLACK_SIGNING_SECRET", ""),
	}
}