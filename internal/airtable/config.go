package airtable

type Config struct {
	apiKey     string
	baseId     string
	usersTable string
}

func NewConfig(apiKey, baseId, usersTable string) *Config {
	return &Config{
		apiKey:     apiKey,
		baseId:     baseId,
		usersTable: usersTable,
	}
}
