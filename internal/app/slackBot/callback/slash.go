package callback

import "net/http"

func (rec *Receiver) handleCallbackSlashCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec.logger.Info("Handle new slack slash command")

		command, err := rec.slack.HandleCallbackSlashCommand(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch command.Command {
		case "/start-daily":
			if err := rec.dailyBot.StartForUser(command.UserID); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/send-reports":
			if err := rec.dailyBot.SendReports(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
