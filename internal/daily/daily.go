package daily

import (
	"SlackBot/internal/airtable"
	"SlackBot/internal/database"
	"SlackBot/internal/models"
	"SlackBot/internal/slackbot"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"strings"
	"time"
)

type Bot struct {
	logger   *logrus.Logger
	database *database.Database
	slack    *slackbot.SlackBot
	airtable *airtable.Airtable

	session  *models.DailySession
	users    map[string]*models.User
	sessions map[int]*models.UserDailySession

	chats []slack.Channel
}

func NewDailyBot(l *logrus.Logger, d *database.Database, b *slackbot.SlackBot, a *airtable.Airtable) *Bot {
	return &Bot{
		logger:   l,
		database: d,
		slack:    b,
		airtable: a,
	}
}

func (b *Bot) Start() (*slack.Msg, error) {
	err := b.initSession()

	if err != nil {
		return nil, err
	}

	slackUsers, err := b.getActiveSlackUsers()

	if err != nil {
		return nil, err
	}

	slackChats, err := b.slack.GetActiveProjectChats()

	if err != nil {
		return nil, err
	}

	b.chats = slackChats

	msg := &slack.Msg{}
	msg.Text = fmt.Sprintf("Daily started. Will be notified %d users", len(slackUsers))
	msg.Hidden = false

	b.users = make(map[string]*models.User, len(slackUsers))
	b.sessions = make(map[int]*models.UserDailySession, len(slackUsers))

	for i := 0; i < len(slackUsers); i++ {
		go b.sendInitialMessage(&slackUsers[i])
	}

	return msg, nil
}

func (b *Bot) initSession() error {
	session, err := b.database.DailySession().FindOrCreateByDate(time.Now())

	if err != nil {
		return err
	}

	b.session = session

	return nil
}

func (b *Bot) getActiveSlackUsers() ([]slack.User, error) {
	activeUsers, err := b.airtable.ActiveUsers()

	if err != nil {
		return nil, err
	}

	slackUsers, err := b.slack.GetActiveUsers()

	if err != nil {
		return nil, err
	}

	var users []slack.User

START:
	for i := 0; i < len(slackUsers); i++ {
		for j := 0; j < len(activeUsers); j++ {
			if slackUsers[i].Profile.Email == activeUsers[j].Fields.Email {
				users = append(users, slackUsers[i])
				continue START
			}
		}
	}

	b.logger.Infof("Total active users: %d", len(activeUsers))
	b.logger.Infof("Total slack users: %d", len(slackUsers))
	b.logger.Infof("Total users: %d", len(users))

	return users, nil
}

func (b *Bot) sendInitialMessage(user *slack.User) {
	attachment := slack.Attachment{
		Pretext:    "Привет, настало время чтобы рассказать чем ты занимался вчера",
		CallbackID: "daily_init",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "accept",
				Text:  "Приступить",
				Style: "primary",
				Type:  "button",
				Value: "accept",
			},
		},
	}

	message := slack.MsgOptionAttachments(attachment)

	b.slack.Api.PostMessage(user.ID, slack.MsgOptionText("", false), message)
}

func (b *Bot) HandleStartSlackUser(payload *slack.InteractionCallback) {
	b.slack.Api.OpenDialog(
		payload.TriggerID,
		slack.Dialog{
			CallbackID:  "daily_report",
			Title:       "Daily report",
			SubmitLabel: "Отправить",
			Elements: []slack.DialogElement{
				slack.NewTextAreaInput("Done", "Опиши что было сделано", ""),
				slack.NewTextAreaInput("WillDo", "Опиши что было сделано", ""),
				slack.NewTextAreaInput("Blocker", "Расскажи что тебе может помешать в твоей работе", ""),
			},
			State: payload.ResponseURL,
		},
	)
}

func (b *Bot) HandleReport(payload *slack.InteractionCallback) {
	data := payload.DialogSubmissionCallback.Submission

	msgText := fmt.Sprintf("User: %s, done: %s, will do: %s, blocker: %s", payload.User.Name, data["Done"], data["WillDo"], data["Blocker"])

	b.logger.Infof("Handle new daily report: %s", msgText)

	url := strings.ReplaceAll(payload.State, "\\", "")
	url = strings.ReplaceAll(url, "\"", "")

	_, _, err := b.slack.Api.PostMessage(payload.User.ID,
		slack.MsgOptionText(msgText, false),
		slack.MsgOptionAttachments(),
		slack.MsgOptionReplaceOriginal(url),
	)

	if err != nil {
		b.logger.Error(err)
	}

	go b.sendReportToChats(payload.User, msgText)
}

func (b *Bot) sendReportToChats(user slack.User, report string) {
	for i := 0; i < len(b.chats); i++ {
		b.slack.Api.PostMessage(b.chats[i].ID, slack.MsgOptionText(report, false))
	}
}