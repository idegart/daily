package airtable

import "github.com/brianloveswords/airtable"

func (a *Airtable) GetInfographicsUsers(force bool) ([]User, error) {
	if !force && a.infographics.users != nil {
		return a.infographics.users, nil
	}

	return a.loadInfographicsUsers()
}

func (a *Airtable) loadInfographicsUsers() ([]User, error) {
	a.logger.Info("Load infographics airtable users")

	var users []User

	usersTable := a.client.Table(a.infographics.config.team.table)

	if err := usersTable.List(&users, &airtable.Options{
		View: a.infographics.config.team.view,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.infographics.users = users

	a.logger.Info("Total infographics airtable users: ", len(a.infographics.users))

	return a.infographics.users, nil
}