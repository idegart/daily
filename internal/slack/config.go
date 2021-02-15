package slack

type Config struct {
	apiToken string
	verificationToken string
	signingSecret string
}

func NewConfig(apiToken, verificationToken, signingSecret string) *Config {
	return &Config{
		apiToken: apiToken,
		verificationToken: verificationToken,
		signingSecret: signingSecret,
	}
}