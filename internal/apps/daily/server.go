package daily

import (
	"bot/internal/model"
	"bot/internal/server"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func (d *Daily) configureServer() {
	d.server = server.NewServer(server.NewConfig(os.Getenv("DAILY_PORT")), d.logger)

	d.server.Router().HandleFunc("/health", handleHealth(d))
	d.server.Router().HandleFunc("/callback/interactive", handleSlackInteractiveCallback(d))

	secure := d.server.Router().PathPrefix("/secure").Subrouter()
	secure.Use(authenticationMiddleware)

	secure.HandleFunc("/start-daily", handleStartDaily(d))
	secure.HandleFunc("/send-reports", handleSendReports(d))
	secure.HandleFunc("/drop-reports", handleDropReports(d))
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if token == os.Getenv("DAILY_AUTHENTICATION") {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func handleHealth(d *Daily) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response = struct {
			Ok               bool `json:"ok"`
			TotalUsers       int  `json:"total_users"`
			TotalProjects    int  `json:"total_projects"`
			TotalAbsentUsers int  `json:"total_absent_users"`
		}{
			Ok:               true,
			TotalUsers:       len(d.users),
			TotalProjects:    len(d.projects),
			TotalAbsentUsers: len(d.absentUsers),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			d.logger.Error(err)
		}
	}
}

func handleStartDaily(d *Daily) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := d.StartInitiation()

		var response = struct {
			Ok               bool  `json:"ok"`
			Error            error `json:"error"`
			TotalUsers       int   `json:"total_users"`
			TotalProjects    int   `json:"total_projects"`
			TotalAbsentUsers int   `json:"total_absent_users"`
		}{
			Ok:               err == nil,
			Error:            err,
			TotalUsers:       len(d.users),
			TotalProjects:    len(d.projects),
			TotalAbsentUsers: len(d.absentUsers),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			d.logger.Error(err)
		}
	}
}

func handleSendReports(d *Daily) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := d.StartReport()

		var response = struct {
			Ok               bool  `json:"ok"`
			Error            error `json:"error"`
			TotalUsers       int   `json:"total_users"`
			TotalProjects    int   `json:"total_projects"`
			TotalAbsentUsers int   `json:"total_absent_users"`
		}{
			Ok:               err == nil,
			Error:            err,
			TotalUsers:       len(d.users),
			TotalProjects:    len(d.projects),
			TotalAbsentUsers: len(d.absentUsers),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			d.logger.Error(err)
		}
	}
}

func handleDropReports(d *Daily) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := d.DropReports()

		var response = struct {
			Ok               bool  `json:"ok"`
			Error            error `json:"error"`
		}{
			Ok:               err == nil,
			Error:            err,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			d.logger.Error(err)
		}
	}
}

func handleSlackInteractiveCallback(d *Daily) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		interaction, err := d.slack.HandleInteraction(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		d.logger.
			WithField("user", interaction.User.Name).
			WithField("callback", interaction.CallbackID).
			WithField("message", interaction.Message.Text).
			Info("Handle new slack interactive")

		switch interaction.CallbackID {
		case SIDailyReportCallbackStart:
			if user := d.GetUserBySlackId(interaction.User.ID); user != nil {
				report, _ := d.database.DailyReport().FindByUserAndDate(user.Id, time.Now())
				latestReport, _ := d.database.DailyReport().GetLastUserReport(user.Id)
				d.SendSlackReportModal(interaction, report, latestReport)
				return
			}
			d.logger.Warn("Slack user not founded")
			w.WriteHeader(http.StatusOK)
			return
		case SIDailyReportCallbackFinish:
			if user := d.GetUserBySlackId(interaction.User.ID); user != nil {
				data := interaction.DialogSubmissionCallback.Submission
				report := &model.DailyReport{
					UserId:  user.Id,
					Date:    time.Now(),
					Done:    data[SIDailyReportDone],
					WillDo:  data[SIDailyReportWillDo],
					Blocker: data[SIDailyReportBlocker],
				}
				if err := d.database.DailyReport().UpdateOrCreate(report); err != nil {
					d.logger.Error(err)
				}
				d.SendSlackThanksForReport(interaction, *user, report)
				d.SendUpdatingReportByUser(*user)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}
