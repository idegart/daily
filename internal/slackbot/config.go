package slackbot

import (
	"SlackBot/internal/env"
)

type Config struct {
	ApiToken string
	VerificationToken string
	TestUser string
}

func NewConfig() *Config {
	return &Config{
		ApiToken: env.Get("SLACK_API_TOKEN", ""),
		VerificationToken: env.Get("SLACK_VERIFICATION_TOKEN", ""),
		TestUser: env.Get("SLACK_TEST_USER", ""),
	}
}