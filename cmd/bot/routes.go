package main

import (
	"bytes"
	"fmt"
	"github.com/slack-go/slack/slackevents"
	"io"
	"net/http"
)

func (a *App) configureRouter() {
	a.server.Router.HandleFunc("/", a.handleHello())
	a.server.Router.HandleFunc("/start-daily-bot", a.handleStartDailyBot())

	callbackRoute := a.server.Router.PathPrefix("/callback").Subrouter()
	callbackRoute.HandleFunc("/events", a.handleCallbackEvents())
}

func (a *App) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle hello")
		io.WriteString(w, "Hello world")
	}
}

func (a *App) handleStartDailyBot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle start daily bot")

		if a.dailyBot.IsEnabled {
			io.WriteString(w, "Daily bot already started")
			return
		}

		_, users, err := a.dailyBot.Start()

		if err != nil {
			a.logger.Error(err)
			io.WriteString(w, fmt.Sprintf("Daily not started. Error: %s", err))
			return
		}

		io.WriteString(w, fmt.Sprintf("Daily bot started. Will be notified %d users", len(users)))
	}
}

func (a *App) handleCallbackEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle new slack event")

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		a.logger.Info("Body request: ", body)

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

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			switch event := eventsAPIEvent.InnerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				go a.dailyBot.HandleNewMessage(event)
			}
		}
	}
}