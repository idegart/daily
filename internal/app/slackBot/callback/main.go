package callback

import (
	"SlackBot/internal/app/dailyBot"
	"SlackBot/internal/server"
	"SlackBot/internal/slack"
	"github.com/sirupsen/logrus"
)

type Receiver struct {
	logger   *logrus.Logger
	server   *server.Server
	slack    *slack.Slack
	dailyBot *dailyBot.DailyBot
}

func NewCallbackReceiver(
	logger *logrus.Logger,
	server *server.Server,
	slack *slack.Slack,
	dailyBot *dailyBot.DailyBot,
) *Receiver {
	return &Receiver{
		logger:   logger,
		server:   server,
		slack:    slack,
		dailyBot: dailyBot,
	}
}

func (rec *Receiver) Configure() {
	s := rec.server.Router.PathPrefix("/callback").Subrouter()

	s.HandleFunc("/event", rec.handleCallbackEvents())
	s.HandleFunc("/slash-command", rec.handleCallbackSlashCommand())
	s.HandleFunc("/interaction", rec.handleCallbackInteractive())
}
