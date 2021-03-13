package slack

import "github.com/slack-go/slack"

func (s *Slack) GetUsers() ([]slack.User, error) {
	s.logger.Info("Load slack users")

	users, err := s.client.GetUsers()

	if err != nil {
		return nil, err
	}

	s.logger.Info("Total slack users: ", len(users))

	return users, nil
}

func (s *Slack) GetActiveUsers() ([]slack.User, error) {
	users, err := s.GetUsers()

	if err != nil {
		return nil, err
	}

	var activeUsers []slack.User

	for i := range users {
		if users[i].Deleted == false && users[i].IsBot == false {
			activeUsers = append(activeUsers, users[i])
		}
	}

	s.logger.Info("Total active slack users: ", len(activeUsers))

	return activeUsers, nil
}

//func (s *Slack) GetActiveUsersInChannel(channelId string) ([]slack.User, error) {
//	activeUsers, err := s.GetActiveUsers()
//
//	if err != nil {
//		return nil, err
//	}
//
//	params := slack.GetUsersInConversationParameters{
//		ChannelID: channelId,
//		Limit:     100,
//	}
//
//	userIds, _, err := s.Client().GetUsersInConversation(&params)
//
//	if err != nil {
//		return nil, err
//	}
//
//	var users []slack.User
//
//	for _, id := range userIds {
//		for i := range activeUsers {
//			if activeUsers[i].ID == id {
//				users = append(users, activeUsers[i])
//			}
//		}
//	}
//
//	return users, nil
//}
