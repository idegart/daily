package dailyBot

import (
	"SlackBot/internal/models"
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

func (b *DailyBot) sendInitialMessageToUser(user *slack.User) error {
	b.logger.Info("Send initial daily message to user ", user.Name)
	attachment := slack.Attachment{
		Pretext:    "Привет, настало время чтобы рассказать чем ты занимался вчера",
		CallbackID: "daily_report_start",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			{
				Name:  "accept",
				Text:  "Приступить",
				Style: "primary",
				Type:  "button",
				Value: "accept",
			},
		},
	}

	message := slack.MsgOptionAttachments(attachment)

	if _, _, err := b.slack.Client().PostMessage(
		user.ID,
		slack.MsgOptionText("", false),
		message,
	); err != nil {
		b.logger.Error(err)
		return err
	}

	return nil
}

func (b *DailyBot) sendDailyModal(triggerID string, responseURL string) error {
	return b.slack.Client().OpenDialog(
		triggerID,
		slack.Dialog{
			CallbackID:  "daily_report_finish",
			Title:       "Daily report",
			SubmitLabel: "Отправить",
			Elements: []slack.DialogElement{
				slack.NewTextAreaInput("Done", "Опиши что было сделано", ""),
				slack.NewTextAreaInput("WillDo", "Опиши что было сделано", ""),
				slack.NewTextAreaInput("Blocker", "Расскажи что тебе может помешать в твоей работе", ""),
			},
			State: responseURL,
		},
	)
}

func (b *DailyBot) sendThanksForReport(userID string, replaceURL string) error {
	_, _, err := b.slack.Client().PostMessage(userID,
		slack.MsgOptionText("Спасибо, можешь продолжить свою работу!", false),
		slack.MsgOptionAttachments(),
		slack.MsgOptionReplaceOriginal(replaceURL),
	)

	return err
}

func (b *DailyBot) sendReportToChannel(channelId string, users []models.User, badUsers []string, reports []models.DailyReport, replaceURL string) error {
	headerText := slack.NewTextBlockObject("mrkdwn", "*Гайз, я тут подготовил ежедневный отчет. Чек зис аут*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	ignoreText := slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf("*Кто меня сегодня проигнорировал:*\n%s", strings.Join(badUsers, "\n")),
		false,
		false)
	ignoreSection := slack.NewSectionBlock(ignoreText, nil, nil)

	reportsSlice := make([]*slack.TextBlockObject, 0)

	for _, report := range reports {

		var user models.User

		for _, u := range users {
			if u.Id == report.UserId {
				user = u
				break
			}
		}

		reportField := slack.NewTextBlockObject(
			"mrkdwn",
			fmt.Sprintf(
				"*%s*\n```Делал вчера:\n%s\n\nДелает сегодня:\n%s\n\nБлокеры:\n%s```",
				user.Name,
				report.Done,
				report.WillDo,
				report.Blocker,
			),
			false,
			false,
		)

		reportsSlice = append(reportsSlice, reportField)
	}

	willDoText := slack.NewTextBlockObject("mrkdwn", "*Что сегодня будет делать команда:*", false, false)
	reportsSection := slack.NewSectionBlock(willDoText, reportsSlice, nil)

	attachment := slack.Attachment{
		Pretext:    "_Если внес изменения, то *не забудь жмакнуть*_",
		CallbackID: "daily_report_refresh",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			{
				Name:  "accept",
				Text:  "Обновить",
				Style: "primary",
				Type:  "button",
				Value: "accept",
			},
		},
	}

	msg := slack.MsgOptionBlocks(
		headerSection,
		slack.NewDividerBlock(),
		ignoreSection,
		slack.NewDividerBlock(),
		reportsSection,
		slack.NewDividerBlock(),
	)

	var options []slack.MsgOption

	options = append(options, msg)
	options = append(options, slack.MsgOptionAttachments(attachment))

	if replaceURL != "" {
		options = append(options, slack.MsgOptionReplaceOriginal(replaceURL))
	}

	_, _, err := b.slack.Client().PostMessage(
		channelId,
		options...,
	)

	if err != nil {
		b.logger.Error(err)
	}

	return err
}
