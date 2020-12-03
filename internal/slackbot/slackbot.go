package slackbot

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

type SlackBot struct {
	config *Config
	logger *logrus.Logger
	Api    *slack.Client
}

func NewSlackBot(config *Config, logger *logrus.Logger) (*SlackBot, error) {
	api := slack.New(config.ApiToken)

	if _, err := api.AuthTest(); err != nil {
		return nil, err
	}

	bot := &SlackBot{
		config: config,
		logger: logger,
		Api:    api,
	}

	return bot, nil
}

func (bot *SlackBot) HandleEvent(body string) (*slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, e := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: bot.config.VerificationToken}),
	)

	return &eventsAPIEvent, e
}

func (bot *SlackBot) HandleVerification(body string, w http.ResponseWriter) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text")
	if _, err := w.Write([]byte(r.Challenge)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (bot *SlackBot) GetActiveUsers() ([]slack.User, error) {
	bot.logger.Info("Get active users")

	var activeUsers []slack.User

	users, err := bot.Api.GetUsers()

	if err != nil {
		return nil, err
	}

	bot.logger.Infof("Total users: %d", len(users))

	for i := 0; i < len(users); i++ {
		if users[i].Deleted == false && users[i].IsBot == false {
			if bot.config.TestUser != "" {
				if users[i].ID == bot.config.TestUser {
					activeUsers = append(activeUsers, users[i])
				}
			} else {
				activeUsers = append(activeUsers, users[i])
			}
		}
	}

	bot.logger.Infof("Total active users: %d", len(activeUsers))

	return activeUsers, nil
}
