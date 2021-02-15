package slack

import (
	"encoding/json"
	"errors"
	"github.com/slack-go/slack"
	"net/http"
)

func (s *Slack) HandleInteraction(r *http.Request) (*slack.InteractionCallback, error) {
	var payload slack.InteractionCallback

	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if payload.Token != s.config.verificationToken {
		return nil, errors.New("bad verification")
	}

	return &payload, err
}