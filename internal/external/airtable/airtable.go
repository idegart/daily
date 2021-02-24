package airtable

import (
	"bot/internal/config"
	"github.com/brianloveswords/airtable"
	"github.com/sirupsen/logrus"
	"strings"
)

type Project struct {
	airtable.Record
	Fields struct {
		ID           int
		Project      string
		Status       string
		Designer     []string
		SlackID      string `json:"Slack ID"`
		SlackUsersID string `json:"DailyBot Summary"`
	}
}

type User struct {
	airtable.Record
	Fields struct {
		ID          int
		Name        string
		Email       string
		Phone       string
		Status      string
		SlackUserID string `json:"Slack User ID"`
	}
}

type Active struct {
	config *ActiveConfig

	users    []User
	projects []Project
}

type Infographics struct {
	config *InfographicsConfig

	users []User
}

type Airtable struct {
	config *config.Airtable
	logger *logrus.Logger
	client *airtable.Client

	active       *Active
	infographics *Infographics
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

func (a *Airtable) SetupActive(config *ActiveConfig) {
	a.active = &Active{
		config: config,
	}
}

func (a *Airtable) SetupInfographics(config *InfographicsConfig) {
	a.infographics = &Infographics{
		config: config,
	}
}

func (p Project) GetSlackIds() []string {
	var ids []string
	var uniqueIds map[string]bool

	s := strings.Split(p.Fields.SlackUsersID, ",")

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
