package slack

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Slack struct {
	config *Config
	logger *logrus.Logger
	client *slack.Client

	users []slack.User
	activeUsers []slack.User

	channels []slack.Channel
	activeChannels []slack.Channel
}

func NewSlack(config *Config, logger *logrus.Logger) *Slack {
	return &Slack{
		config: config,
		logger: logger,
		client: slack.New(config.apiToken),
	}
}

func (s *Slack) Client() *slack.Client {
	return s.client
}

//func (s *Slack) HandleCallbackEvent(body string) (*slackevents.EventsAPIEvent, error) {
//	eventsAPIEvent, err := slackevents.ParseEvent(
//		json.RawMessage(body),
//		slackevents.OptionVerifyToken(
//			&slackevents.TokenComparator{VerificationToken: s.config.verificationToken},
//		),
//	)
//
//	if err != nil {
//		s.logger.Error(err)
//		return nil, err
//	}
//
//	return &eventsAPIEvent, nil
//}
//
//func (s *Slack) HandleVerification(body string) ([]byte, error) {
//	s.logger.Info("Handle verification")
//
//	var resp *slackevents.ChallengeResponse
//	err := json.Unmarshal([]byte(body), &resp)
//	if err != nil {
//		s.logger.Error(err)
//		return nil, err
//	}
//
//	return []byte(resp.Challenge), nil
//}
//
//func (s *Slack) HandleCallbackSlashCommand(r *http.Request) (*slack.SlashCommand, error) {
//	verifier, err := slack.NewSecretsVerifier(r.Header, s.config.signingSecret)
//	if err != nil {
//		s.logger.Error(err)
//		return nil, err
//	}
//
//	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
//	command, err := slack.SlashCommandParse(r)
//	if err != nil {
//		s.logger.Error(err)
//		return nil, err
//	}
//
//	if err = verifier.Ensure(); err != nil {
//		s.logger.Error(err)
//		return nil, err
//	}
//
//	return &command, nil
//}
//
//func (s *Slack) HandleInteraction(r *http.Request) (*slack.InteractionCallback, error) {
//	var payload slack.InteractionCallback
//
//	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
//
//	if err != nil {
//		s.logger.Error(err)
//		return nil, err
//	}
//
//	if payload.Token != s.config.VerificationToken {
//		return nil, errors.New("bad verification")
//	}
//
//	return &payload, err
//}
//
//func (s *Slack) GetBodyFromRequest(r *http.Request) (string, error) {
//	buf := new(bytes.Buffer)
//	if _, err := buf.ReadFrom(r.Body); err != nil {
//		s.logger.Error(err)
//		return "", err
//	}
//	body := buf.String()
//
//	return body, nil
//}
//



