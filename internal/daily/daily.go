package daily

import (
	"SlackBot/internal/database"
	"SlackBot/internal/models"
	"SlackBot/internal/slackbot"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"time"
)

type Bot struct {
	logger   *logrus.Logger
	database *database.Database
	bot      *slackbot.SlackBot

	IsEnabled bool

	session  *models.DailySession
	users    map[string]*models.User
	sessions map[int]*models.UserDailySession
}

func NewDailyBot(l *logrus.Logger, d *database.Database, b *slackbot.SlackBot) *Bot {
	return &Bot{
		logger:   l,
		database: d,
		bot:      b,
	}
}

func (b *Bot) Start() (s *models.DailySession, u []slack.User, err error) {
	b.IsEnabled = true

	session, err := b.database.DailySession().FindOrCreateByDate(time.Now())

	if err != nil {
		return nil, nil, err
	}

	b.session = session

	slackUsers, err := b.bot.GetActiveUsers()

	if err != nil {
		return nil, nil, err
	}

	b.users = make(map[string]*models.User, len(slackUsers))
	b.sessions = make(map[int]*models.UserDailySession, len(slackUsers))

	for i := 0; i < len(slackUsers); i++ {
		go b.initSlackUser(slackUsers[i])
	}

	return session, slackUsers, nil
}

func (b *Bot) HandleNewMessage(event *slackevents.MessageEvent) {
	if b.IsEnabled == false {
		return
	}

	if event.ChannelType != slack.TYPE_IM || event.BotID != "" {
		return
	}

	user, ok := b.users[event.User]

	if ok == false {
		b.logger.Warnf("Handle new message from user which is not in daily: %s", event.User)
		return
	}

	userSession, ok := b.sessions[user.Id]

	if ok == false {
		b.logger.Warnf("User session not initialized: %s", event.User)
		return
	}

	if userSession.Done == "" {
		b.setUserDone(userSession, event)
		return
	}

	if userSession.WillDo == "" {
		b.setUserWillDo(userSession, event)
		b.removeUserFromDaily(event)
		return
	}
}

func (b *Bot) initSlackUser(slackUser slack.User) {
	user, err := b.initUser(slackUser)

	if err != nil {
		b.logger.Error(err)
		return
	}

	_, err = b.initUserSession(user)

	if err != nil {
		b.logger.Error(err)
		return
	}

	err = b.sendInitialMessageToUser(slackUser)

	if err != nil {
		b.logger.Errorf("Error when send daily initial message to %s: %s", slackUser.ID, err)
	}
}

func (b *Bot) initUser(user slack.User) (*models.User, error) {
	realUser, err := b.database.User().FindOrCreateBySlackUser(&user)

	if err != nil {
		return nil, err
	}

	b.users[user.ID] = realUser

	return realUser, nil
}

func (b *Bot) initUserSession(user *models.User) (*models.UserDailySession, error) {
	session, err := b.database.UserDailySession().FindOrCreateByDailySessionAndUser(b.session, user)

	if err != nil {
		return nil, err
	}

	b.sessions[user.Id] = session

	return session, nil
}

func (b *Bot) sendInitialMessageToUser(user slack.User) error {
	_, _, err := b.bot.Api.PostMessage(
		user.ID,
		slack.MsgOptionText("Привет, настало время чтобы рассказать чем ты занимался вчера", false),
	)

	return err
}

func (b *Bot) setUserDone(userSession *models.UserDailySession, event *slackevents.MessageEvent) {
	userSession.Done = event.Text
	if err := b.database.UserDailySession().Update(userSession); err != nil {
		b.logger.Error("Can not update done for user session: ", err)
		return
	}

	if _, _, err := b.bot.Api.PostMessage(
		event.User,
		slack.MsgOptionText("Отлично, а чем планируешь заняться?", false),
	); err != nil {
		b.logger.Errorf("Error when send daily initial message to %s: %s", event.User, err)
	}
}

func (b *Bot) setUserWillDo(userSession *models.UserDailySession, event *slackevents.MessageEvent) {
	userSession.WillDo = event.Text
	if err := b.database.UserDailySession().Update(userSession); err != nil {
		b.logger.Error("Can not update done for user session: ", err)
		return
	}

	if _, _, err := b.bot.Api.PostMessage(
		event.User,
		slack.MsgOptionText("Спасибо, можешь работать дальше", false),
	); err != nil {
		b.logger.Errorf("Error when send daily initial message to %s: %s", event.User, err)
	}
}

func (b *Bot) removeUserFromDaily(event *slackevents.MessageEvent) {
	delete(b.sessions, b.users[event.User].Id)
	delete(b.users, event.User)

	if len(b.users) < 1 {
		b.IsEnabled = false
	}
}