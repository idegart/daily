package dailyBot

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/database"
	"SlackBot/internal/models"
	"SlackBot/internal/slack"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	slackgo "github.com/slack-go/slack"
	"strings"
	"sync"
	"time"
)

type DailyBot struct {
	logger   *logrus.Logger
	database *database.Database
	slack    *slack.Slack
	airtable *airtable.Airtable

	airtableUsers []airtable.User
	slackUsers    []slackgo.User
	users         []models.User
}

func NewDailyBot(logger *logrus.Logger, database *database.Database, slack *slack.Slack, airtable *airtable.Airtable) *DailyBot {
	return &DailyBot{
		logger:   logger,
		database: database,
		slack:    slack,
		airtable: airtable,
	}
}

func (b *DailyBot) Init() error {
	if err := b.initUsers(); err != nil {
		return err
	}

	return nil
}

func (b *DailyBot) initUsers() error {
	b.logger.Info("Init users in daily bot")
	wg := sync.WaitGroup{}
	wg.Add(2)

	go b.initAirtableUsers(&wg)
	go b.initSlackUsers(&wg)

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
				go func(email string, name string, airtableId int, slackId string) {
					defer wg.Done()

					user, err := b.database.UserRepository().FindByEmailOrCreate(
						email,
						name,
						airtableId,
						slackId,
					)

					if err != nil {
						b.logger.Error(err)
						return
					}

					mu.Lock()
					b.users = append(b.users, *user)
					mu.Unlock()

				}(au.Fields.Email, au.Fields.Name, au.Fields.ID, su.ID)
			}
		}
	}

	wg.Wait()

	b.logger.Info("Total users: ", len(b.users))

	return nil
}

func (b *DailyBot) initAirtableUsers(wg *sync.WaitGroup) {
	defer wg.Done()

	users, err := b.airtable.ActiveUsers()

	if err != nil {
		b.logger.Error(err)
		return
	}

	b.airtableUsers = users
}

func (b *DailyBot) initSlackUsers(wg *sync.WaitGroup) {
	defer wg.Done()

	users, err := b.slack.GetActiveUsers()

	if err != nil {
		b.logger.Error(err)
		return
	}

	b.slackUsers = users
}

func (b *DailyBot) StartForUser(slackUserId string) error {
	b.logger.Info("Start daily for user ", slackUserId)

	var slackUser *slackgo.User

	for _, u := range b.slackUsers {
		if u.ID == slackUserId {
			slackUser = &u
			break
		}
	}

	if slackUser == nil {
		b.logger.Error("slack user not initialized in daily bot")
		return errors.New("slack user not initialized in daily bot")
	}

	return b.sendInitialMessageToUser(slackUser)
}

func (b *DailyBot) StartUserReport(callback *slackgo.InteractionCallback) error {
	b.logger.Info("Start user report")

	return b.sendDailyModal(callback.TriggerID, callback.ResponseURL)
}

func (b *DailyBot) FinishUserReport(callback *slackgo.InteractionCallback) error {
	b.logger.Info("Finish user report")

	data := callback.DialogSubmissionCallback.Submission

	msgText := fmt.Sprintf("User: %s, done: %s, will do: %s, blocker: %s", callback.User.Name, data["Done"], data["WillDo"], data["Blocker"])

	url := strings.ReplaceAll(callback.State, "\\", "")
	url = strings.ReplaceAll(url, "\"", "")

	b.logger.Infof("Handle new daily report: %s", msgText)

	var user *models.User

	for _, u := range b.users {
		if u.SlackId == callback.User.ID {
			user = &u
			break
		}
	}

	if user == nil {
		b.logger.Error("user not founded in initial")
		return errors.New("user not founded in initial")
	}

	_, err := b.database.DailyReportRepository().CreateOrUpdateByDateAndUser(time.Now(), user.Id, data["Done"], data["WillDo"], data["Blocker"])

	if err != nil {
		b.logger.Error(err)
		return err
	}

	return b.sendThanksForReport(callback.User.ID, url)
}

func (b *DailyBot) SendReports() error {
	b.logger.Info("Send reports")

	channels, err := b.slack.GetActiveProjectConversations()

	if err != nil {
		b.logger.Error(err)
		return err
	}

	ws := sync.WaitGroup{}

	for _, channel := range channels {
		ws.Add(1)
		go func(channel slackgo.Channel) {
			defer ws.Done()
			b.reportInChannel(channel, "")
		}(channel)
	}

	ws.Wait()

	return nil
}

func (b *DailyBot) RefreshReport(callback *slackgo.InteractionCallback) error {
	b.reportInChannel(callback.Channel, callback.ResponseURL)
	return nil
}

func (b *DailyBot) reportInChannel(channel slackgo.Channel, replaceURL string) error {
	params := slackgo.GetUsersInConversationParameters{
		ChannelID: channel.ID,
		Limit:     100,
	}

	slackUsersIds, _, err := b.slack.Client().GetUsersInConversation(&params)

	if err != nil {
		b.logger.Error(err)
		return err
	}

	var users []models.User

	for _, slackUserId := range slackUsersIds {
		for _, user := range b.users {
			if user.SlackId == slackUserId {
				users = append(users, user)
				break
			}
		}
	}

	reports, err := b.database.DailyReportRepository().FindAllByDateAndUsers(time.Now(), users)

	if err != nil {
		b.logger.Error(err)
		return err
	}

	var badUsers []string

	out:
	for _, users := range users {
		for _, report := range reports {
			if report.UserId == users.Id {
				continue out
			}
		}

		badUsers = append(badUsers, "<@"+users.SlackId+">")
	}

	return b.sendReportToChannel(channel.ID, users, badUsers, reports, replaceURL)
}