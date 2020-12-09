package airtable

import (
	"SlackBot/internal/env"
	"github.com/brianloveswords/airtable"
)

type Airtable struct {
	config *Config
	Client *airtable.Client
}

type User struct {
	airtable.Record // provides ID, CreatedTime
	Fields          struct {
		Name  string
		Email string
		Phone string
		Status string
	}
}

func NewAirtable(config *Config) *Airtable {
	return &Airtable{
		config: config,
		Client: &airtable.Client{
			APIKey: config.APIKey,
			BaseID: config.BaseID,
		},
	}
}

func (a *Airtable) ActiveUsers() ([]User, error) {
	var users []User

	usersTable := a.Client.Table(a.config.UsersTable)

	if err := usersTable.List(&users, &airtable.Options{
		Filter: `OR({Status}='Inhouse',{Status}='Outsource')`,
	}); err != nil {
		return nil, err
	}

	if testUserEmail := env.Get("TEST_USER", ""); testUserEmail != "" {
		for i := 0; i < len(users); i++ {
			if users[i].Fields.Email == testUserEmail {
				var us []User
				return append(us, users[i]), nil
			}
		}
	}

	return users, nil
}
