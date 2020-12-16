package callback

import "net/http"

func (rec *Receiver) handleCallbackInteractive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec.logger.Info("Handle new slack interactive")

		interaction, err := rec.slack.HandleInteraction(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch interaction.CallbackID {
		case "daily_report_start":
			rec.logger.Info("daily_report_start")
			if err := rec.dailyBot.StartUserReport(interaction); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "daily_report_finish":
			rec.logger.Info("daily_report_finish")
			if err := rec.dailyBot.FinishUserReport(interaction); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
