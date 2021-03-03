package main

import (
	"bot/internal/model"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Response struct {
	Ok           bool
	ErrorMessage string
	Data         interface{}
}

type StartResponse struct {
	Total int
}

func handleHealth(a *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.ResponseWithOk(w)
	}
}

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

func handleSendReports(a *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response = Response{}
		err := a.SendReports()
		if err != nil {
			response.Ok = false
			response.ErrorMessage = err.Error()
		} else {
			response.Ok = true
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			a.logger.Fatal(err)
		}
	}
}

func handleSlackInteractiveCallback(a *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		interaction, err := a.slack.HandleInteraction(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.logger.WithField("interaction", interaction).Info("Handle new slack interactive")

		switch interaction.CallbackID {
		case SIDailyReportCallbackStart:
			if user := a.GetUserBySlackId(interaction.User.ID); user != nil {
				report, err := a.database.DailyReport().FindByUserAndDate(user.Id, time.Now())
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					a.logger.Error(err)
				}
				a.SendSlackReportModal(interaction, report)
				return
			}
		case SIDailyReportCallbackFinish:
			if user := a.GetUserBySlackId(interaction.User.ID); user != nil {
				data := interaction.DialogSubmissionCallback.Submission
				report := &model.DailyReport{
					UserId:  user.Id,
					Date:    time.Now(),
					Done:    data[SIDailyReportDone],
					WillDo:  data[SIDailyReportWillDo],
					Blocker: data[SIDailyReportBlocker],
				}
				if err := a.database.DailyReport().UpdateOrCreate(report); err != nil {
					a.logger.Error(err)
				}
				a.SendSlackThanksForReport(interaction)
				a.ResendReportsByUser(user)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *App) ResponseWithOk(w http.ResponseWriter) {
	var response = Response{
		Ok: true,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Fatal(err)
	}
}
