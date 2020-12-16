package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Slack struct {
	config *Config
	logger *logrus.Logger
	client *slack.Client
}

func NewSlack(config *Config, logger *logrus.Logger) *Slack {
	return &Slack{
		config: config,
		logger: logger,
		client: slack.New(config.ApiToken),
	}
}

func (s *Slack) HandleCallbackEvent(body string) (*slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: s.config.VerificationToken},
		),
	)

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &eventsAPIEvent, nil
}

func (s *Slack) HandleVerification(body string) ([]byte, error) {
	s.logger.Info("Handle verification")

	var resp *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &resp)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return []byte(resp.Challenge), nil
}

func (s *Slack) HandleCallbackSlashCommand(r *http.Request) (*slack.SlashCommand, error) {
	verifier, err := slack.NewSecretsVerifier(r.Header, s.config.SigningSecret)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if err = verifier.Ensure(); err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &command, nil
}

func (s *Slack) HandleInteraction(r *http.Request) (*slack.InteractionCallback, error) {
	var payload slack.InteractionCallback

	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if payload.Token != s.config.VerificationToken {
		return nil, errors.New("bad verification")
	}

	return &payload, err
}

func (s *Slack) GetBodyFromRequest(r *http.Request) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		s.logger.Error(err)
		return "", err
	}
	body := buf.String()

	return body, nil
}

func (s *Slack) GetActiveUsers() ([]slack.User, error) {
	s.logger.Info("Get active slack users")

	var activeUsers []slack.User

	users, err := s.client.GetUsers()

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Infof("Total slack users: %d", len(users))

	for i := 0; i < len(users); i++ {
		if users[i].Deleted == false && users[i].IsBot == false {
			activeUsers = append(activeUsers, users[i])
		}
	}

	s.logger.Infof("Total active slack users: %d", len(activeUsers))

	return activeUsers, nil
}

func (s *Slack) GetActiveProjectConversations() ([]slack.Channel, error) {
	s.logger.Info("Get active project slack conversations")

	var channels []slack.Channel
	var cursor string

	firstRequest := true

	for firstRequest || cursor != "" {
		firstRequest = false

		params := &slack.GetConversationsParameters{
			Cursor: cursor,
		}

		chs, c, err := s.client.GetConversations(params)

		if err != nil {
			s.logger.Error(err)
			return nil, err
		}

		for _, channel := range chs {
			//m, err := regexp.MatchString(`^(p\d{4}-|bot-test)`, channel.Name)
			m, err := regexp.MatchString(`^bot-test`, channel.Name)

			if err != nil {
				s.logger.Error(err)
				return nil, err
			}

			if !channel.IsArchived && m {
				channels = append(channels, channel)
			}
		}

		cursor = c
	}

	return channels, nil
}

func (s *Slack) Client() *slack.Client {
	return s.client
}

func (s *Slack) GetActiveProjectConversationsForUser(user *slack.User) ([]slack.Channel, error) {

	var channels []slack.Channel
	var cursor string

	firstRequest := true

	for firstRequest || cursor != "" {
		firstRequest = false

		params := &slack.GetConversationsForUserParameters{
			UserID: user.ID,
			Cursor: cursor,
		}

		chs, c, err := s.client.GetConversationsForUser(params)

		if err != nil {
			s.logger.Error(err)
			return nil, err
		}

		for _, channel := range chs {
			m, err := regexp.MatchString(`^bot-test`, channel.Name)

			if err != nil {
				s.logger.Error(err)
				return nil, err
			}

			if !channel.IsArchived && m {
				channels = append(channels, channel)
			}
		}

		cursor = c
	}

	return channels, nil
}