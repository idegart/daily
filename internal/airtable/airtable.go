package airtable

import (
	"github.com/brianloveswords/airtable"
	"github.com/sirupsen/logrus"
)

type Airtable struct {
	config *Config
	logger *logrus.Logger
	Client *airtable.Client
}

type User struct {
	airtable.Record // provides ID, CreatedTime
	Fields          struct {
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
		Client: &airtable.Client{
			APIKey: config.APIKey,
			BaseID: config.BaseID,
		},
	}
}

func (a *Airtable) ActiveUsers() ([]User, error) {
	a.logger.Info("Get active airtable users")

	var users []User

	usersTable := a.Client.Table(a.config.UsersTable)

	if err := usersTable.List(&users, &airtable.Options{
		Filter: `OR({Status}='Inhouse',{Status}='Outsource')`,
	}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.logger.Infof("Total airtable users: %d", len(users))

	return users, nil
}
