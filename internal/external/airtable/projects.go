package airtable

import "github.com/brianloveswords/airtable"

type Project struct {
	airtable.Record
	Fields struct {
		ID       int
		Project  string
		Status   string
		Designer []string
	}
}

func (a *Airtable) GetActiveProjects(force bool) ([]Project, error) {
	if !force && a.activeProjects != nil {
		return a.activeProjects, nil
	}

	projects, err := a.loadProjects(force)
	if err != nil {
		return nil, err
	}

	var activeProjects []Project

	var allowedStatuses = []string{
		projectStatusPresale,
		projectStatusInProgress,
		projectStatusJaupaIsReady,
		projectStatusPostProduction,
	}

	for i := range projects {
		for _, status := range allowedStatuses {
			if projects[i].Fields.Status == status {
				activeProjects = append(activeProjects, projects[i])
			}
		}
	}

	a.activeProjects = activeProjects

	a.logger.Info("Total active airtable projects: ", len(a.activeProjects))

	return a.activeProjects, nil
}

func (a Airtable) loadProjects(force bool) ([]Project, error) {
	if !force && a.projects != nil {
		return a.projects, nil
	}

	a.logger.Info("Load airtable projects")

	var projects []Project

	projectsTable := a.client.Table(a.config.ProjectsTable)

	if err := projectsTable.List(&projects, &airtable.Options{}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.projects = projects

	return a.projects, nil
}
