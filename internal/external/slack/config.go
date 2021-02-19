package slack

import "bot/internal/config"

func NewConfig(apiToken, verificationToken, signingSecret string) *config.Slack {
	return &config.Slack{
		ApiToken: apiToken,
		VerificationToken: verificationToken,
		SigningSecret: signingSecret,
	}
}