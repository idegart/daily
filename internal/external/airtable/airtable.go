package airtable

import (
	"bot/internal/config"
	"github.com/brianloveswords/airtable"
	"github.com/sirupsen/logrus"
	"strings"
)

type airtableApp struct {
	appID string
}

type Airtable struct {
	config *config.Airtable
	logger *logrus.Logger
	client *airtable.Client

	team     airtableApp
	projects airtableApp
}

func NewAirtable(config *config.Airtable, logger *logrus.Logger) *Airtable {
	return &Airtable{
		logger: logger,
		config: config,
		client: &airtable.Client{
			APIKey: config.ApiKey,
		},
	}
}

func (a *Airtable) SetupTeam(appID string) {
	a.team = airtableApp{
		appID: appID,
	}
}

func (a *Airtable) SetupProjects(appID string) {
	a.projects = airtableApp{
		appID: appID,
	}
}

func (p Project) GetSlackIds() []string {
	var ids []string
	var uniqueIds map[string]bool

	s := strings.Split(p.Fields.SlackUsersID[0], ",")

	uniqueIds = make(map[string]bool, len(s))

	for _, id := range s {
		if id != "" {
			tid := strings.Trim(id, " ")

			if ok, _ := uniqueIds[tid]; !ok {
				uniqueIds[tid] = true
				ids = append(ids, tid)
			}
		}
	}

	return ids
}
