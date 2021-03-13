package slack

import (
	"bot/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type message struct {
	channelID string
	options []slack.MsgOption
	cb func(timestamp string)
}

type Slack struct {
	config *config.Slack
	logger *logrus.Logger
	client *slack.Client

	messagesToSend chan message
}

func NewSlack(config *config.Slack, logger *logrus.Logger) *Slack {
	return &Slack{
		config: config,
		logger: logger,
		client: slack.New(config.ApiToken, slack.OptionDebug(false)),
		messagesToSend: make(chan message),
	}
}

func (s *Slack) Client() *slack.Client {
	return s.client
}

func (s *Slack) SendMessage(ChannelID string, cb func(timestamp string), options ...slack.MsgOption) {
	s.messagesToSend <- message{
		channelID: ChannelID,
		options: options,
		cb: cb,
	}
}

func (s *Slack) StartSendingMessages() {
	for message := range s.messagesToSend {
		_,timestamp,_,err := s.client.SendMessage(message.channelID, message.options...)

		if err != nil {
			s.logger.Error(err, timestamp)
			return
		}

		if message.cb != nil {
			message.cb(timestamp)
		}
	}
}



