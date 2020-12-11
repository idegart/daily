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
	if err := b.dailyBot.Init(); err != nil {
		return err
	}

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

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			res, err := b.slack.HandleVerification(body);
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			_, err = w.Write(res)
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

		switch command.Command {
		case "/start-daily":
			if err := b.dailyBot.StartForUser(command.UserID); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/send-reports":
			if err := b.dailyBot.SendReports(); err != nil {
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
			if err := b.dailyBot.StartUserReport(interaction); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "daily_report_finish":
			b.logger.Info("daily_report_finish")
			if err := b.dailyBot.FinishUserReport(interaction); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
