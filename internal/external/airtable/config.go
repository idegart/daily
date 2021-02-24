package airtable

import "bot/internal/config"

type TableView struct {
	table string
	view string
}

type ActiveConfig struct {
	projects TableView
	team TableView
}

type InfographicsConfig struct {
	team TableView
}

func NewConfig(apiKey, appId string) *config.Airtable {
	return &config.Airtable{
		ApiKey:        apiKey,
		AppId:         appId,
	}
}

func NewActiveConfig(projectsTable, projectsView, teamTable, teamView string) *ActiveConfig {
	return &ActiveConfig{
		projects: TableView{
			table: projectsTable,
			view: projectsView,
		},
		team: TableView{
			table: teamTable,
			view: teamView,
		},
	}
}

func NewInfographicsConfig(teamTable, teamView string) *InfographicsConfig {
	return &InfographicsConfig{
		team: TableView{
			table: teamTable,
			view: teamView,
		},
	}
}
