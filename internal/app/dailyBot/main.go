package dailyBot

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/database"
	"SlackBot/internal/models"
	"SlackBot/internal/slack"
	"database/sql"
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

	airtableUsers map[string]*airtable.User
	slackUsers    map[string]*slackgo.User
}

func NewDailyBot(logger *logrus.Logger, database *database.Database, slack *slack.Slack, airtable *airtable.Airtable) *DailyBot {
	return &DailyBot{
		logger:   logger,
		database: database,
		slack:    slack,
		airtable: airtable,
	}
}

func (b *DailyBot) Start() error {
	if err := b.initAirtableUsers(); err != nil {
		return err
	}

	if err := b.initSlackUsers(); err != nil {
		return err
	}

	b.sendInitialMessages()

	return nil
}

func (b *DailyBot) StartUserReport(callback *slackgo.InteractionCallback) {
	err := b.slack.Client().OpenDialog(
		callback.TriggerID,
		slackgo.Dialog{
			CallbackID:  "daily_report_finish",
			Title:       "Daily report",
			SubmitLabel: "Отправить",
			Elements: []slackgo.DialogElement{
				slackgo.NewTextAreaInput("Done", "Опиши что было сделано", ""),
				slackgo.NewTextAreaInput("WillDo", "Опиши что было сделано", ""),
				slackgo.NewTextAreaInput("Blocker", "Расскажи что тебе может помешать в твоей работе", ""),
			},
			State: callback.ResponseURL,
		},
	)

	if err != nil {
		b.logger.Error(err)
	}
}

func (b *DailyBot) FinishUserReport(callback *slackgo.InteractionCallback) {
	data := callback.DialogSubmissionCallback.Submission

	msgText := fmt.Sprintf("User: %s, done: %s, will do: %s, blocker: %s", callback.User.Name, data["Done"], data["WillDo"], data["Blocker"])

	b.logger.Infof("Handle new daily report: %s", msgText)

	url := strings.ReplaceAll(callback.State, "\\", "")
	url = strings.ReplaceAll(url, "\"", "")

	_, _, err := b.slack.Client().PostMessage(callback.User.ID,
		slackgo.MsgOptionText("Спасибо, можешь продолжить свою работу!", false),
		slackgo.MsgOptionAttachments(),
		slackgo.MsgOptionReplaceOriginal(url),
	)

	if err != nil {
		b.logger.Error(err)
	}

	var slackUser *slackgo.User

	for _, su := range b.slackUsers {
		if su.ID == callback.User.ID {
			slackUser = su
			break
		}
	}

	if slackUser == nil {
		b.logger.Errorf("can not find user % in initial array", callback.User.Name)
		return
	}

	user, err := b.database.UserRepository().FindByEmail(slackUser.Profile.Email)

	if err == sql.ErrNoRows {
		err = nil

		airtableUser, ok := b.airtableUsers[slackUser.Profile.Email]

		if !ok {
			b.logger.Errorf("For slack user %s not found airtable user", slackUser.Profile.Email)
			return
		}

		user := &models.User{
			Email:      slackUser.Profile.Email,
			SlackId:    slackUser.ID,
			AirtableId: airtableUser.Fields.ID,
		}

		err = b.database.UserRepository().Create(user)
	}

	if err != nil {
		b.logger.Error(err)
		return
	}

	report, err := b.database.DailyReportRepository().FindByDateAndUser(time.Now(), user.Id)

	if err == sql.ErrNoRows {
		report = &models.DailyReport{
			UserId:  user.Id,
			Date:    time.Now(),
			Done:    data["Done"],
			WillDo:  data["WillDo"],
			Blocker: data["Blocker"],
		}

		b.logger.Warn(report)

		err = b.database.DailyReportRepository().Create(report)

		if err != nil {
			b.logger.Error(err)
			return
		}

		return
	}

	if err != nil {
		b.logger.Error(err)
		return
	}

	report.Done = data["Done"]
	report.WillDo = data["WillDo"]
	report.Blocker = data["Blocker"]

	err = b.database.DailyReportRepository().Update(report)

	if err != nil {
		b.logger.Error(err)
		return
	}
}

func (b *DailyBot) initAirtableUsers() error {
	airtableUsers, err := b.airtable.ActiveUsers()

	if err != nil {
		return err
	}

	b.airtableUsers = make(map[string]*airtable.User, len(airtableUsers))

	for i := 0; i < len(airtableUsers); i++ {
		b.airtableUsers[airtableUsers[i].Fields.Email] = &airtableUsers[i]
	}

	return nil
}

func (b *DailyBot) initSlackUsers() error {
	activeUsers, err := b.slack.GetActiveUsers()

	if err != nil {
		return err
	}

	b.slackUsers = make(map[string]*slackgo.User, len(activeUsers))

	for i := 0; i < len(activeUsers); i++ {
		b.slackUsers[activeUsers[i].Profile.Email] = &activeUsers[i]
	}

	return nil
}

func (b *DailyBot) sendInitialMessages() {
	ws := sync.WaitGroup{}

	for email, _ := range b.airtableUsers {
		if slackUser, ok := b.slackUsers[email]; ok {
			ws.Add(1)
			go b.sendInitialMessageToUser(&ws, slackUser)
		}
	}

	ws.Wait()
}

func (b *DailyBot) sendInitialMessageToUser(ws *sync.WaitGroup, user *slackgo.User) {
	attachment := slackgo.Attachment{
		Pretext:    "Привет, настало время чтобы рассказать чем ты занимался вчера",
		CallbackID: "daily_report_start",
		Color:      "#3AA3E3",
		Actions: []slackgo.AttachmentAction{
			{
				Name:  "accept",
				Text:  "Приступить",
				Style: "primary",
				Type:  "button",
				Value: "accept",
			},
		},
	}

	message := slackgo.MsgOptionAttachments(attachment)

	if _, _, err := b.slack.Client().PostMessage(user.ID, slackgo.MsgOptionText("", false), message); err != nil {
		b.logger.Error(err)
	}

	ws.Done()
}
