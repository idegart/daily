package main

import (
	"bytes"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

func (app *App) handleEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Handle slack event")

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		app.logger.Info("Body request: ", body)

		eventsAPIEvent, err := app.bot.HandleEvent(body)

		if err != nil {
			app.logger.Error("Event API error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			app.bot.HandleVerification(body, w)
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			switch event := eventsAPIEvent.InnerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				app.handleNewUserMessage(event)
			}
		}
	}
}

func (app *App) handleNewUserMessage(event *slackevents.MessageEvent) {
	if event.ChannelType != slack.TYPE_IM || event.BotID != "" {
		return
	}

	realUser, ok := app.users[event.User]

	if ok == false {
		app.logger.Warnf("Handle new message from unknown user: %s", event.User)
		return
	}

	userSession, ok := app.sessions[realUser.Id]

	if ok == false {
		app.logger.Warnf("Not found session for current user: %s", event.User)
		return
	}

	if userSession.Done == "" {
		userSession.Done = event.Text
		if err := app.database.UserDailySession().Update(userSession); err != nil {
			app.logger.Error("Can not update done for user session: ", err)
			return
		}

		if _, _, err := app.bot.Api.PostMessage(
			event.User,
			slack.MsgOptionText("Отлично, а чем планируешь заняться?", false),
		); err != nil {
			app.logger.Errorf("Error when send daily initial message to %s: %s", event.User, err)
		}
		return
	}

	if userSession.WillDo == "" {
		userSession.WillDo = event.Text
		if err := app.database.UserDailySession().Update(userSession); err != nil {
			app.logger.Error("Can not update done for user session: ", err)
			return
		}

		if _, _, err := app.bot.Api.PostMessage(
			event.User,
			slack.MsgOptionText("Спасибо, можешь работать дальше", false),
		); err != nil {
			app.logger.Errorf("Error when send daily initial message to %s: %s", event.User, err)
		}

		delete(app.users, event.User)
		delete(app.sessions, realUser.Id)

		return
	}

	app.logger.Info(userSession)
}
