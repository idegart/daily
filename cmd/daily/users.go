package main

import (
	"bot/internal/airtable"
	"bot/internal/model"
	"github.com/slack-go/slack"
	"sync"
)

func (a *App) initUsers() error {
	airtableUsers, slackUsers := a.loadUsers()

	users, err := a.saveUsers(airtableUsers, slackUsers)

	if err != nil {
		return err
	}

	a.users = users

	return nil
}

func (a *App) loadUsers() ([]airtable.User, []slack.User) {
	type Users struct {
		airtableUsers []airtable.User
		slackUsers []slack.User
	}

	var users = &Users{}

	var wg = &sync.WaitGroup{}

	wg.Add(2)

	go func(wg *sync.WaitGroup, users *Users) {
		defer wg.Done()

		airtableUsers, err := app.airtable.GetActiveUsers(true)

		if err != nil {
			app.logger.Fatal(err)
		}

		users.airtableUsers = airtableUsers

	}(wg, users)

	go func(wg *sync.WaitGroup, users *Users) {
		defer wg.Done()

		slackUsers, err := app.slack.GetActiveUsers(true)

		if err != nil {
			app.logger.Fatal(err)
		}

		users.slackUsers = slackUsers

	}(wg, users)

	wg.Wait()

	return users.airtableUsers, users.slackUsers
}

func (a *App) saveUsers(airtableUsers []airtable.User, slackUsers []slack.User) ([]model.User, error) {
	a.logger.Info("Save users")

	var users []model.User

	for i := range airtableUsers {
		for j := range slackUsers {
			if airtableUsers[i].Fields.Email == slackUsers[j].Profile.Email {
				var user = &model.User{
					Email: airtableUsers[i].Fields.Email,
					Name: airtableUsers[i].Fields.Name,
					AirtableId: airtableUsers[i].Fields.ID,
					SlackId: slackUsers[j].ID,
				}

				users = append(users, *user)
			}
		}
	}

	a.logger.Info("Total users to save: ", len(users))

	wg := &sync.WaitGroup{}

	for i := range users {
		wg.Add(1)
		go func(wg *sync.WaitGroup, user *model.User) {
			defer wg.Done()
			if err := a.database.User().UpdateOrCreate(user); err != nil {
				a.logger.Error(err)
			}
		}(wg, &users[i])
	}

	wg.Wait()

	a.logger.Info("Users saved")

	return users, nil
}