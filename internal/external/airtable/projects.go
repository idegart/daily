package airtable

import "github.com/brianloveswords/airtable"

type Project struct {
	airtable.Record
	Fields struct {
		ID           string `json:"ID Auto"`
		Project      string
		Status       string
		Type         string
		SlackID      string `json:"Slack ID"`
		SlackUsersID []string `json:"DailyBot Summary"`
	}
}

const (
	activeProjectsTable = "tblQupUCVtZD6GqdK"
	activeProjectsView  = "viw9yTS3LylyX0KW7"
)

func (a *Airtable) GetActiveProjects() ([]Project, error) {
	a.logger.Info("Load active airtable projects")

	var projects []Project

	a.client.BaseID = a.projects.appID

	projectsTable := a.client.Table(activeProjectsTable)

	if err := projectsTable.List(&projects, &airtable.Options{
		View: activeProjectsView,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.logger.Info("Total active airtable projects: ", len(projects))

	return projects, nil
}