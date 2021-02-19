package airtable

import "github.com/brianloveswords/airtable"

type User struct {
	airtable.Record
	Fields struct {
		ID     int
		Name   string
		Email  string
		Phone  string
		Status string
	}
}

func (a *Airtable) GetActiveUsers(force bool) ([]User, error) {
	if !force && a.activeUsers != nil {
		return a.activeUsers, nil
	}

	users, err := a.loadUsers(force)
	if err != nil {
		return nil, err
	}

	var activeUsers []User

	var allowedStatuses = []string{
		userStatusOutsource,
		userStatusInHouse,
	}

	for i := range users {
		for _, status := range allowedStatuses {
			if users[i].Fields.Status == status {
				activeUsers = append(activeUsers, users[i])
			}
		}
	}

	a.activeUsers = activeUsers

	a.logger.Info("Total active airtable users: ", len(a.activeUsers))

	return a.activeUsers, nil
}

func (a *Airtable) loadUsers(force bool) ([]User, error) {
	if !force && a.users != nil {
		return a.users, nil
	}

	a.logger.Info("Load airtable users")

	var users []User

	usersTable := a.client.Table(a.config.TeamTable)

	if err := usersTable.List(&users, &airtable.Options{}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.users = users

	a.logger.Info("Total airtable users: ", len(a.users))

	return a.users, nil
}
