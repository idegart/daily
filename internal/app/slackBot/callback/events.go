package callback

import (
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

func (rec *Receiver) handleCallbackEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec.logger.Info("Handle new slack event")

		body, err := rec.slack.GetBodyFromRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		eventsAPIEvent, err := rec.slack.HandleCallbackEvent(body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			res, err := rec.slack.HandleVerification(body);
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			_, err = w.Write(res)
		}
	}
}