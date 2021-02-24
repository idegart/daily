package slack

import (
	"github.com/slack-go/slack"
)

func (s *Slack) GetChannels(force bool) ([]slack.Channel, error) {
	if !force && s.channels != nil {
		return s.channels, nil
	}

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

	s.channels = channels

	s.logger.Info("Total slack channels: ", len(s.channels))

	return s.channels, nil
}

func (s *Slack) GetActiveProjectChannels(force bool) ([]slack.Channel, error) {
	if !force && s.activeChannels != nil {
		return s.activeChannels, nil
	}

	if s.channels == nil {
		if _, err := s.GetChannels(force); err != nil {
			return nil, err
		}
	}

	var channels []slack.Channel

	for i := range s.channels {
		//m, err := regexp.MatchString(`^(p\d{4}-|bot-test)`, s.channels[i].Name)
		//m, err := regexp.MatchString(`^id-732-fullhouse`, s.channels[i].Name)

		//if err != nil {
		//	return nil, err
		//}

		if !s.channels[i].IsArchived {
			channels = append(channels, s.channels[i])
		}
	}

	s.activeChannels = channels

	s.logger.Info("Total active slack channels: ", len(s.activeChannels))

	return s.activeChannels, nil
}
