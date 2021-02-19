package slack

import "github.com/slack-go/slack"

func (s *Slack) GetUsers(force bool) ([]slack.User, error)  {
	if !force && s.users != nil {
		return s.users, nil
	}

	s.logger.Info("Load slack users")

	users, err := s.client.GetUsers()

	s.logger.Info("Total slack users: ", len(users))

	if err != nil {
		return nil, err
	}

	s.users = users

	return s.users, nil
}

func (s *Slack) GetActiveUsers(force bool) ([]slack.User, error) {
	if !force && s.activeUsers != nil {
		return s.activeUsers, nil
	}

	users, err := s.GetUsers(force)

	if err != nil {
		return nil, err
	}

	var activeUsers []slack.User

	for i := range users {
		if users[i].Deleted == false && users[i].IsBot == false {
			activeUsers = append(activeUsers, users[i])
		}
	}

	s.activeUsers = activeUsers

	s.logger.Info("Total active slack users: ", len(s.activeUsers))

	return s.activeUsers, nil
}

func (s *Slack) GetActiveUsersInChannel(channelId string) ([]slack.User, error)  {
	if s.activeUsers == nil {
		if _, err := s.GetActiveUsers(true); err != nil {
			return nil, err
		}
	}

	params := slack.GetUsersInConversationParameters{
		ChannelID: channelId,
		Limit:     100,
	}

	userIds, _, err := s.Client().GetUsersInConversation(&params)

	if err != nil {
		return nil, err
	}

	var users []slack.User

	for _, id := range userIds {
		for i := range s.activeUsers {
			if s.activeUsers[i].ID == id {
				users = append(users, s.activeUsers[i])
			}
		}
	}

	return users, nil
}
