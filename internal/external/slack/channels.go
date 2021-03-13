package slack

import (
	"github.com/slack-go/slack"
)

func (s *Slack) GetChannels() ([]slack.Channel, error) {
	s.logger.Info("Load slack channels")

	var channels []slack.Channel
	var cursor string

	firstRequest := true

	for firstRequest || cursor != "" {
		firstRequest = false

		params := &slack.GetConversationsParameters{
			Cursor: cursor,
		}

		chs, c, err := s.client.GetConversations(params)

		channels = append(channels, chs...)

		if err != nil {
			s.logger.Error(err)
			return nil, err
		}

		cursor = c
	}

	s.logger.Info("Total slack channels: ", len(channels))

	return channels, nil
}

func (s *Slack) GetActiveProjectChannels() ([]slack.Channel, error) {
	var activeChannels []slack.Channel
	channels, err := s.GetChannels()

	if err != nil {
		 return nil, err
	}

	for i := range channels {
		//m, err := regexp.MatchString(`^id-732-fullhouse`, s.channels[i].Name)
		if !channels[i].IsArchived {
			activeChannels = append(activeChannels, channels[i])
		}
	}

	activeChannels = channels

	s.logger.Info("Total active slack channels: ", activeChannels)

	return activeChannels, nil
}
