package airtable

import (
	"bot/internal/config"
	"github.com/brianloveswords/airtable"
	"github.com/sirupsen/logrus"
)

type Airtable struct {
	config *config.Airtable
	logger *logrus.Logger
	client *airtable.Client

	users       []User
	activeUsers []User

	projects []Project
	activeProjects []Project
}

func NewAirtable(config *config.Airtable, logger *logrus.Logger) *Airtable {
	return &Airtable{
		logger: logger,
		config: config,
		client: &airtable.Client{
			APIKey: config.ApiKey,
			BaseID: config.AppId,
		},
	}
}