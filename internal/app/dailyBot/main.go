package dailyBot

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/database"
	"SlackBot/internal/models"
	"SlackBot/internal/slack"
	"database/sql"
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

func (b *DailyBot) StartForUser(slackUserId string) error {
	b.logger.Info("Start daily for user ", slackUserId)

	if err := b.initUsers(); err != nil {
		return err
	}

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

	if err := b.initUsers(); err != nil {
		return err
	}

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

	report := &models.DailyReport{
		UserId:  user.Id,
		Date:    time.Now(),
		Done:    data["Done"],
		WillDo:  data["WillDo"],
		Blocker: data["Blocker"],
	}

	if err := b.database.DailyReportRepository().CreateOrUpdateByDateAndUser(report); err != nil {
		b.logger.Error(err)
		return err
	}

	if err := b.sendThanksForReport(callback.User.ID, url); err != nil {
		return err
	}

	_, err := b.database.SlackReportRepository().DateReportsExists(time.Now())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		b.logger.Error(err)
		return err
	}

	return b.refreshReports(&callback.User)
}

func (b *DailyBot) SendReports() error {
	b.logger.Info("Send reports")

	if err := b.initUsers(); err != nil {
		return err
	}

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
			b.reportInChannel(channel.ID)
		}(channel)
	}

	ws.Wait()

	return nil
}

func (b *DailyBot) refreshReports(user *slackgo.User) error {
	b.logger.Info("Refresh reports")

	if err := b.initUsers(); err != nil {
		return err
	}

	channels, err := b.slack.GetActiveProjectConversationsForUser(user)

	if err != nil {
		return err
	}

	reports, err := b.database.SlackReportRepository().FindAllByDateAndChannels(time.Now(), channels)

	if err != nil {
		return err
	}

	for _, report := range reports {
		go b.reportInChannel(report.SlackChannelId)
	}

	return nil
}

func (b *DailyBot) reportInChannel(channelId string) error {
	params := slackgo.GetUsersInConversationParameters{
		ChannelID: channelId,
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

	slackReport, err := b.database.SlackReportRepository().FindByDateAndSlackChannel(time.Now(), channelId)

	if err == sql.ErrNoRows {
		_, ts, err := b.sendReportToChannel(channelId, users, badUsers, reports, "")

		if err != nil {
			b.logger.Error(err)
			return err
		}

		slackReport := &models.SlackReport{
			SlackChannelId: channelId,
			Date: time.Now(),
			Ts: ts,
		}

		return b.database.SlackReportRepository().Create(slackReport)
	}

	if err != nil {
		b.logger.Error(err)
		return err
	}

	_, ts, err := b.sendReportToChannel(channelId, users, badUsers, reports, slackReport.Ts)

	if err != nil {
		b.logger.Error(err)
		return err
	}

	slackReport.Ts = ts

	return b.database.SlackReportRepository().Update(slackReport)
}