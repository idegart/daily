package staff

import (
	"bot/internal/database"
	"bot/internal/external/airtable"
	"bot/internal/external/slack"
	"bot/internal/model"
	"github.com/sirupsen/logrus"
	slackGo "github.com/slack-go/slack"
	"sync"
)

type Staff struct {
	logger   *logrus.Logger
	database database.Database

	airtable *airtable.Airtable
	slack    *slack.Slack

	airtableUsers []airtable.User
	slackUsers    []slackGo.User
	users         []model.User
	slackProjects []slackGo.Channel
}

func New(logger *logrus.Logger, database database.Database, airtable *airtable.Airtable, slack *slack.Slack) *Staff {
	return &Staff{
		logger:   logger,
		database: database,
		airtable: airtable,
		slack:    slack,
	}
}

func (s *Staff) Initialize() error {
	if err := s.loadAirtableUsers(); err != nil {
		return err
	}

	if err := s.loadSlackUsers(); err != nil {
		return err
	}

	if err := s.loadSlackProjects(); err != nil {
		return err
	}

	s.initializeUsers()

	return nil
}

func (s *Staff) GetUsers(force bool) ([]model.User, error) {
	if !force && s.users != nil {
		return s.users, nil
	}

	if err := s.Initialize(); err != nil {
		return nil, err
	}

	return s.users, nil
}

func (s *Staff) loadAirtableUsers() error {
	users, err := s.airtable.GetActiveUsers(true)

	if err != nil {
		return err
	}

	s.airtableUsers = users

	return nil
}

func (s *Staff) loadSlackUsers() error {
	users, err := s.slack.GetActiveUsers(true)

	if err != nil {
		return err
	}

	s.slackUsers = users

	return nil
}

func (s *Staff) loadSlackProjects() error {
	channels, err := s.slack.GetActiveProjectChannels(true)

	if err != nil {
		return err
	}

	s.slackProjects = channels

	return nil
}

func (s *Staff) initializeUsers() {
	var users []model.User

	for i := range s.airtableUsers {
		for j := range s.slackUsers {
			if s.airtableUsers[i].Fields.Email == s.slackUsers[j].Profile.Email {
				var user = &model.User{
					Email:      s.airtableUsers[i].Fields.Email,
					Name:       s.airtableUsers[i].Fields.Name,
					AirtableId: s.airtableUsers[i].Fields.ID,
					SlackId:    s.slackUsers[j].ID,
				}

				users = append(users, *user)
			}
		}
	}

	wg := &sync.WaitGroup{}

	for i := range users {
		wg.Add(1)
		go func(wg *sync.WaitGroup, user *model.User) {
			defer wg.Done()
			if err := s.database.User().UpdateOrCreate(user); err != nil {
				s.logger.Error(err)
			}
		}(wg, &users[i])
	}

	wg.Wait()

	s.users = users
}
