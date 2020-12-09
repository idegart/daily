package main

import (
	"bytes"
	"github.com/slack-go/slack/slackevents"
	"io"
	"net/http"
)

func (a *App) configureRouter() {
	a.server.Router.HandleFunc("/", a.handleHello())

	callbackRoute := a.server.Router.PathPrefix("/callback").Subrouter()
	callbackRoute.HandleFunc("/event", a.handleCallbackEvents())
	callbackRoute.HandleFunc("/slash-command", a.handleCallbackSlashCommand())
	callbackRoute.HandleFunc("/interaction", a.handleCallbackInteraction())
}

func (a *App) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle hello")
		io.WriteString(w, "Hello world")
	}
}

func (a *App) handleCallbackEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle new slack event")

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		eventsAPIEvent, err := a.slackBot.HandleEvent(body)

		if err != nil {
			a.logger.Error("Event API error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			a.slackBot.HandleVerification(body, w)
			return
		}
	}
}

func (a *App) handleCallbackSlashCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle new slack slash command")

		s, err := a.slackBot.HandleSlashCommand(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch s.Command {
		case "/start-daily":
			msg, err := a.dailyBot.Start()
			if err != nil {
				a.logger.Error(err)
			}
			a.slackBot.WriteSlashResponse(w, msg)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (a *App) handleCallbackInteraction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle new slack interaction")

		payload, err := a.slackBot.HandleInteraction(r)

		if err != nil {
			a.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if payload.CallbackID == "daily_init" {
			a.dailyBot.HandleStartSlackUser(payload)
		}

		if payload.CallbackID == "daily_report" {
			a.dailyBot.HandleReport(payload)
		}
	}
}