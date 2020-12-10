package slackBot

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/app/dailyBot"
	"SlackBot/internal/database"
	"SlackBot/internal/server"
	"SlackBot/internal/slack"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
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

func NewSlackBot(logger *logrus.Logger, database *database.Database, server *server.Server, slack *slack.Slack, airtable *airtable.Airtable) *SlackBot {
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

	callbackRoute := b.server.Router.PathPrefix("/callback").Subrouter()
	callbackRoute.HandleFunc("/event", b.handleCallbackEvents())
	callbackRoute.HandleFunc("/slash-command", b.handleCallbackSlashCommand())
	callbackRoute.HandleFunc("/interaction", b.handleCallbackInteractive())
}

func (b *SlackBot) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.logger.Info("Handle hello")
		io.WriteString(w, "Hello world")
	}
}

func (b *SlackBot) handleCallbackEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.logger.Info("Handle new slack event")

		body, err := b.slack.GetBodyFromRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		eventsAPIEvent, err := b.slack.HandleCallbackEvent(body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			if err := b.slack.HandleVerification(w, body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}

func (b *SlackBot) handleCallbackSlashCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.logger.Info("Handle new slack slash command")

		command, err := b.slack.HandleCallbackSlashCommand(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b.logger.Info(command)

		switch command.Command {
		case "/start-daily":
			if err := b.dailyBot.Start(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (b *SlackBot) handleCallbackInteractive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.logger.Info("Handle new slack interactive")

		interaction, err := b.slack.HandleInteraction(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch interaction.CallbackID {
		case "daily_report_start":
			b.logger.Info("daily_report_start")
			b.dailyBot.StartUserReport(interaction)
		case "daily_report_finish":
			b.logger.Info("daily_report_finish")
			b.dailyBot.FinishUserReport(interaction)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
