package airtable

import "bot/internal/config"

func NewConfig(apiKey, appId, teamTable, projectsTable string) *config.Airtable {
	return &config.Airtable{
		ApiKey:        apiKey,
		AppId:         appId,
		TeamTable:     teamTable,
		ProjectsTable: projectsTable,
	}
}
