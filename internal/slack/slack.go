package slack

import (
	"github.com/slack-go/slack"
)

type Slack struct {
	config *Config
	api    *slack.Client
}

func New(config *Config) *Slack {
	return &Slack{
		config: config,
	}
}

func (s *Slack) SetApi() error {

	s.api = slack.New(s.config.ApiToken)

	if _, err := s.api.AuthTest(); err != nil {
		return err
	}

	return nil
}
