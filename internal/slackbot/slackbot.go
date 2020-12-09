package slackbot

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io"
	"io/ioutil"
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

func (bot *SlackBot) WriteSlashResponse(w http.ResponseWriter, params *slack.Msg) {
	b, err := json.Marshal(params)
	if err != nil {
		bot.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (bot *SlackBot) HandleEvent(body string) (*slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, e := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: bot.config.VerificationToken}),
	)

	return &eventsAPIEvent, e
}

func (bot *SlackBot) HandleInteraction(r *http.Request) (*slack.InteractionCallback, error) {

	var payload slack.InteractionCallback

	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)

	if err != nil {
		return nil, err
	}

	if payload.Token != bot.config.VerificationToken {
		return nil, errors.New("bad verification")
	}

	return &payload, err
}

func (bot *SlackBot) HandleSlashCommand(r *http.Request) (*slack.SlashCommand, error) {
	verifier, err := slack.NewSecretsVerifier(r.Header, bot.config.SigningSecret)
	if err != nil {
		bot.logger.Error(err)
		return nil, err
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		bot.logger.Error(err)
		return nil, err
	}

	if err = verifier.Ensure(); err != nil {
		bot.logger.Error(err)
		return nil, err
	}

	return &s, nil
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
				if users[i].Profile.Email == bot.config.TestUser {
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

func (bot *SlackBot) GetActiveProjectChats() ([]slack.Channel, error) {
	params := &slack.GetConversationsParameters{}

	var cursor string
	var channels []slack.Channel

	isFirst := true

	for cursor != "" || isFirst {
		isFirst = false
		params.Cursor = cursor
		ch, c, err := bot.Api.GetConversations(params)
		if err != nil {
			bot.logger.Error(err)
			return nil, err
		}

		cursor = c

		for i := 0; i < len(ch); i++ {
			if ch[i].IsArchived == false {
				if bot.config.TestChat != "" {
					if ch[i].ID == bot.config.TestChat {
						channels = append(channels, ch[i])
					}
				} else {
					channels = append(channels, ch[i])
				}
			}
		}
	}

	return channels, nil
}