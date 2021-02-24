package main

import (
	"bot/internal/model"
)

func (a *App) PrepareTeam() error {
	teamProjects, err := a.airtable.GetActiveProjects(true)

	if err != nil {
		return err
	}

	a.teamProjects = teamProjects

	teamUsers, err := a.airtable.GetActiveUsers(true)

	if err != nil {
		return err
	}

	a.teamUsers = teamUsers

	infographicsUsers, err := a.airtable.GetInfographicsUsers(true)

	if err != nil {
		return err
	}

	a.infographicsUsers = infographicsUsers

	if err := a.GenerateDatabaseUsersFromAirtable(); err != nil {
		return err
	}

	return nil
}

func (a *App) GenerateDatabaseUsersFromAirtable() error  {
	var users []model.User

	for _, airtableUser := range a.teamUsers{
		user := &model.User{
			Name: airtableUser.Fields.Name,
			Email: airtableUser.Fields.Email,
			AirtableId: airtableUser.Fields.ID,
			SlackId: airtableUser.Fields.SlackUserID,
		}

		if err := a.database.User().UpdateOrCreate(user); err != nil {
			return err
		}

		users = append(users, *user)
	}

	a.users = users

	return nil
}