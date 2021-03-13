package config

type Logger struct {
	LogLevel string
}

type DB struct {
	Driver string
	Url    string
}

type Server struct {
	BindAddr string
}

type Airtable struct {
	ApiKey string
}

type Slack struct {
	ApiToken          string
	VerificationToken string
	SigningSecret     string
}
