package slack

import (
	"bot/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Slack struct {
	config *config.Slack
	logger *logrus.Logger
	client *slack.Client

	users []slack.User
	activeUsers []slack.User

	channels []slack.Channel
	activeChannels []slack.Channel
}

func NewSlack(config *config.Slack, logger *logrus.Logger) *Slack {
	return &Slack{
		config: config,
		logger: logger,
		client: slack.New(config.ApiToken),
	}
}

func (s *Slack) Client() *slack.Client {
	return s.client
}



