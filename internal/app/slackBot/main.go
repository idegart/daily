package slackBot

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/app/dailyBot"
	"SlackBot/internal/app/slackBot/callback"
	"SlackBot/internal/database"
	"SlackBot/internal/server"
	"SlackBot/internal/slack"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type SlackBot struct {
	logger   *logrus.Logger
	database *database.Database
	server   *server.Server
	slack    *slack.Slack
	airtable *airtable.Airtable

	dailyBot *dailyBot.DailyBot
}

func NewSlackBot(
	logger *logrus.Logger,
	database *database.Database,
	server *server.Server,
	slack *slack.Slack,
	airtable *airtable.Airtable,
) *SlackBot {
	return &SlackBot{
		logger:   logger,
		database: database,
		server:   server,
		slack:    slack,
		airtable: airtable,
		dailyBot: dailyBot.NewDailyBot(logger, database, slack, airtable),
	}
}

func (b *SlackBot) Serve() error {
	b.configureRoutes()

	return b.server.Start()
}

func (b *SlackBot) configureRoutes() {
	b.server.Router.HandleFunc("/", b.handleHello())

	callbackReceiver := callback.NewCallbackReceiver(b.logger, b.server, b.slack, b.dailyBot)
	callbackReceiver.Configure()
}

func (b *SlackBot) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.logger.Info("Handle hello")
		_, err := io.WriteString(w, "Hello world")

		if err != nil {
			b.logger.Error(err)
		}
	}
}