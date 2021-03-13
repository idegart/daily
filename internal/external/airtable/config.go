package airtable

import "bot/internal/config"

type TableView struct {
	table string
	view string
}

func NewConfig(apiKey string) *config.Airtable {
	return &config.Airtable{
		ApiKey:        apiKey,
	}
}