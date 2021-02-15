package airtable

import (
	"github.com/brianloveswords/airtable"
	"github.com/sirupsen/logrus"
)

type Airtable struct {
	config *Config
	logger *logrus.Logger
	client *airtable.Client

	users       []User
	activeUsers []User
}

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

func NewAirtable(config *Config, logger *logrus.Logger) *Airtable {
	return &Airtable{
		config: config,
		logger: logger,
		client: &airtable.Client{
			APIKey: config.apiKey,
			BaseID: config.baseId,
		},
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

	for i := range users {
		if users[i].Fields.Status == "Inhouse" || users[i].Fields.Status == "Outsource" {
			activeUsers = append(activeUsers, users[i])
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

	usersTable := a.client.Table(a.config.usersTable)

	if err := usersTable.List(&users, &airtable.Options{}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.users = users

	a.logger.Info("Total airtable users: ", len(a.users))

	return a.users, nil
}