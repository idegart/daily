package dailyBot

import (
	"SlackBot/internal/models"
	"errors"
	"sync"
)

func (b *DailyBot) initUsers() error {
	b.logger.Info("Init users in daily bot")
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		b.initAirtableUsers()
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		b.initSlackUsers()
	}(&wg)

	wg.Wait()

	if b.airtableUsers == nil {
		return errors.New("airtable users not set")
	}

	if b.slackUsers == nil {
		return errors.New("slack users not set")
	}

	b.users = []models.User{}

	mu := sync.Mutex{}

	for _, au := range b.airtableUsers {
		for _, su := range b.slackUsers {
			if au.Fields.Email == su.Profile.Email {
				wg.Add(1)

				user := &models.User{
					Email: au.Fields.Email,
					Name: au.Fields.Name,
					AirtableId: au.Fields.ID,
					SlackId: su.ID,
				}

				go func(user *models.User) {
					defer wg.Done()

					if err := b.initUser(user); err != nil {
						b.logger.Error(err)
						return
					}

					mu.Lock()
					b.users = append(b.users, *user)
					mu.Unlock()
				}(user)
			}
		}
	}

	wg.Wait()

	b.logger.Info("Total users: ", len(b.users))

	return nil
}

func (b *DailyBot) initAirtableUsers() {
	users, err := b.airtable.ActiveUsers()

	if err != nil {
		b.logger.Error(err)
		return
	}

	b.airtableUsers = users
}

func (b *DailyBot) initSlackUsers() {
	users, err := b.slack.GetActiveUsers()

	if err != nil {
		b.logger.Error(err)
		return
	}

	b.slackUsers = users
}

func (b *DailyBot) initUser(user *models.User) error {
	return b.database.UserRepository().FindByEmailOrCreate(user)
}