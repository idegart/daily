package main

import (
	"bot/internal/server"
	"encoding/json"
	"net/http"
	"os"
)

func (a *App) configureServer() {
	a.server = server.NewServer(server.NewConfig(os.Getenv("DAILY_PORT")), a.logger)

	a.server.Router().HandleFunc("/health", a.handleHealth())
	a.server.Router().HandleFunc("/callback/interactive", a.handleSlackInteractiveCallback())

	a.server.Router().HandleFunc("/secret/start-daily", a.handleSecretStartDaily())
	a.server.Router().HandleFunc("/secret/finish-daily", a.handleSecretFinishDaily())
}

func (a App) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
			a.logger.Fatal(err)
		}
	}
}

func (a *App) handleSlackInteractiveCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Handle new slack interactive")

		interaction, err := a.slack.HandleInteraction(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch interaction.CallbackID {
		case "daily_report_start":
			if user := a.FindUserBySlackId(interaction.User.ID); user != nil {
				a.startForUserByCallback(interaction)
				return
			}
		case "daily_report_finish":
			if user := a.FindUserBySlackId(interaction.User.ID); user != nil {
				a.finishUserReportByCallback(interaction, user)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *App) handleSecretStartDaily() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sendInitialMessages()
		if err := json.NewEncoder(w).Encode(map[string]bool{"started": true}); err != nil {
			a.logger.Fatal(err)
		}
	}
}

func (a *App) handleSecretFinishDaily() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sendReports()
		if err := json.NewEncoder(w).Encode(map[string]bool{"started": true}); err != nil {
			a.logger.Fatal(err)
		}
	}
}