package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Ok bool
	ErrorMessage string
	Data interface{}
}

type StartResponse struct {
	Total int
}

func handleHealth(a *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response = Response{
			Ok: true,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			a.logger.Fatal(err)
		}
	}
}

//func handleSlackInteractiveCallback(a *App) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		a.logger.Info("Handle new slack interactive")
//
//		interaction, err := a.slack.HandleInteraction(r)
//
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		switch interaction.CallbackID {
//		case "daily_report_start":
//			if user := a.FindUserBySlackId(interaction.User.ID); user != nil {
//				a.startForUserByCallback(interaction)
//				return
//			}
//		case "daily_report_finish":
//			if user := a.FindUserBySlackId(interaction.User.ID); user != nil {
//				a.finishUserReportByCallback(interaction, user)
//				return
//			}
//		}
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//}
//
func handleStartDaily(a *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response = Response{}

		users, err := a.SendInitialMessages()

		if err != nil {
			response.Ok = false
			response.ErrorMessage = err.Error()
		} else {
			response.Ok = true
			response.Data = StartResponse{
				Total: len(users),
			}
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			a.logger.Fatal(err)
		}
	}
}
//
//func handleSecretFinishDaily(a *App) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		a.sendReports()
//		if err := json.NewEncoder(w).Encode(map[string]bool{"started": true}); err != nil {
//			a.logger.Fatal(err)
//		}
//	}
//}
