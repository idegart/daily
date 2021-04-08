package daily

import (
	"bot/internal/model"
	"github.com/slack-go/slack"
	"time"
)

func (d *Daily) StartInitiation() error {
	d.logger.Info("Start initiation")

	if err := d.Init(); err != nil {
		return err
	}

LOOP:
	for _, user := range d.users {
		for _, absentUser := range d.absentUsers {
			if user.Id == absentUser.User.Id {
				continue LOOP
			}
		}

		d.usersToInitiate <- user
	}

	return nil
}

func (d *Daily) startSendingInitiations() {
	for user := range d.usersToInitiate {
		if err := d.SendInitiationToUser(user); err != nil {
			d.logger.Error(err)
		}
	}
}

func (d *Daily) SendInitiationToUser(user model.User) error {
	d.logger.Infof("Sending initiation to %s (%d, %s)", user.Name, user.Id, user.Email)

	previousReport, _ := d.database.DailyReport().GetLastUserReport(user.Id)
	currentReport, _ := d.database.DailyReport().FindByUserAndDate(user.Id, time.Now())

	d.sendSlackInitialMessageToUser(user, previousReport, currentReport, nil)

	return nil
}

func (d *Daily) sendSlackInitialMessageToUser(user model.User, previousReport *model.DailyReport, currentReport *model.DailyReport, replace *string) {
	var previousReportFields []slack.AttachmentField

	if previousReport != nil && previousReport.WillDo != "" {
		previousReportFields = append(previousReportFields, slack.AttachmentField{
			Title: "То, чем ты планировал заниматься в прошлый раз",
			Value: previousReport.WillDo,
		})
	}

	if currentReport != nil && currentReport.WillDo != "" {
		previousReportFields = append(previousReportFields, slack.AttachmentField{
			Title: "То, чем ты планировал заниматься в этот раз",
			Value: currentReport.WillDo,
		})
	}

	if len(previousReportFields) < 1 {
		previousReportFields = append(previousReportFields, slack.AttachmentField{
			Title: "У тебя еще нет отчетов, пора бы это исправить",
		})
	}

	headerText := "*🤖 Привет, настало время, чтобы рассказать о том, чем занимаешься*"

	actionButton := slack.AttachmentAction{
		Name:  "accept",
		Text:  "Рассказать",
		Style: "primary",
		Type:  "button",
		Value: "accept",
	}

	if currentReport != nil {
		actionButton.Text = "Дорассказать"
		actionButton.Style = "default"

		headerText = "*🙏🏻 Спасибо, что все рассказал :)*"
	}

	var messageOptions = []slack.MsgOption{}

	if replace == nil {
		messageOptions = append(messageOptions, slack.MsgOptionText("Привет, расскажи чем ты занимался", false))
	}

	messageOptions = append(
		messageOptions,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					"mrkdwn",
					headerText,
					false,
					false,
				),
				nil,
				nil,
			),
			slack.NewDividerBlock(),
		),
		slack.MsgOptionAttachments(slack.Attachment{
			CallbackID: SIDailyReportCallbackStart,
			Color:      "#3AA3E3",
			Footer:     "В течение дня ты можешь изменить свой ответ",
			Fields:     previousReportFields,
			Actions: []slack.AttachmentAction{
				actionButton,
			},
		}),
	)

	if replace != nil {
		messageOptions = append(messageOptions, slack.MsgOptionReplaceOriginal(*replace))
	}

	d.slack.SendMessage(
		user.SlackId,
		nil,
		messageOptions...,
	)
}
