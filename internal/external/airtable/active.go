package airtable

import "github.com/brianloveswords/airtable"

func (a *Airtable) GetActiveProjects(force bool) ([]Project, error) {
	if !force && a.active.projects != nil {
		return a.active.projects, nil
	}

	return a.loadActiveProjects()
}

func (a *Airtable) loadActiveProjects() ([]Project, error) {
	a.logger.Info("Load active airtable projects")

	var projects []Project

	projectsTable := a.client.Table(a.active.config.projects.table)

	if err := projectsTable.List(&projects, &airtable.Options{
		View: a.active.config.projects.view,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.active.projects = projects

	a.logger.Info("Total active airtable projects: ", len(a.active.projects))

	return a.active.projects, nil
}

func (a *Airtable) GetActiveUsers(force bool) ([]User, error) {
	if !force && a.active.users != nil {
		return a.active.users, nil
	}

	return a.loadActiveUsers()
}

func (a *Airtable) loadActiveUsers() ([]User, error) {
	a.logger.Info("Load active airtable users")

	var users []User

	usersTable := a.client.Table(a.active.config.team.table)

	if err := usersTable.List(&users, &airtable.Options{
		View: a.active.config.team.view,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.active.users = users

	a.logger.Info("Total active airtable users: ", len(a.active.users))

	return a.active.users, nil
}