package main

import (
	"bot/internal/model"
	"github.com/slack-go/slack"
)

//
//import (
//	"bot/internal/model"
//	"fmt"
//	"github.com/slack-go/slack"
//	"strings"
//	"time"
//)
//
func (a *App) SendSlackInitialMessageToUser(user model.User) error {
	attachment := slack.Attachment{
		Pretext:    "Привет, настало время чтобы рассказать чем ты занимаешься",
		CallbackID: "daily_report_start",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			{
				Name:  "accept",
				Text:  "Рассказать",
				Style: "primary",
				Type:  "button",
				Value: "accept",
			},
		},
	}

	message := slack.MsgOptionAttachments(attachment)

	_, _, err := a.slack.Client().PostMessage(
		user.SlackId,
		slack.MsgOptionText("", false),
		message,
	)

	return err
}
//
//func (a *App) sendReportModal(callback *slack.InteractionCallback) {
//	if err := a.slack.Client().OpenDialog(
//		callback.TriggerID,
//		slack.Dialog{
//			CallbackID:  "daily_report_finish",
//			Title:       "Daily report",
//			SubmitLabel: "Отправить",
//			Elements: []slack.DialogElement{
//				slack.NewTextAreaInput("Done", "Опиши, что было сделано", ""),
//				slack.NewTextAreaInput("WillDo", "Опиши, что было сделано", ""),
//				slack.NewTextAreaInput("Blocker", "Расскажи, что тебе может помешать в твоей работе", ""),
//			},
//			State: fmt.Sprintf(
//				"%s", time.Now().Format("2006-01-02"),
//			),
//		},
//	); err != nil {
//		a.logger.Error(err)
//	}
//}
//
//func (a *App) sendThanksForReport(callback *slack.InteractionCallback) {
//	if _, _, err := a.slack.Client().PostMessage(callback.User.ID,
//		slack.MsgOptionText("Спасибо, можешь продолжить свою работу! Если захочешь изменить свой ответ, просто нажми на кнопку 'Рассказать' соответствующего дня", false),
//		slack.MsgOptionAttachments(),
//		//slack.MsgOptionReplaceOriginal(replaceUrl),
//	); err != nil {
//		a.logger.Error(err)
//	}
//}
//
//func (a *App) sendReportToChannel(channelId string, users []model.User, badUsers []model.User, reports []model.DailyReport, replace string) (string, string, error) {
//	headerText := slack.NewTextBlockObject("mrkdwn", "*Гайз, я тут подготовил ежедневный отчет. Чек зис аут*", false, false)
//	headerSection := slack.NewSectionBlock(headerText, nil, nil)
//
//	var badUsersIds []string
//
//	for _, user := range badUsers {
//		badUsersIds = append(badUsersIds, "<@"+user.SlackId+">")
//	}
//
//	var ignoreText string
//
//	if len(badUsersIds) < 1 {
//		ignoreText = "*Все сегодня молодцы. Ни одного игнорирования*"
//	} else {
//		ignoreText = fmt.Sprintf("*Кто меня сегодня проигнорировал:*\n%s", strings.Join(badUsersIds, "\n"))
//	}
//
//	ignoreTextBlock := slack.NewTextBlockObject(
//		"mrkdwn",
//		ignoreText,
//		false,
//		false)
//	ignoreSection := slack.NewSectionBlock(ignoreTextBlock, nil, nil)
//
//	reportsSlice := make([]*slack.TextBlockObject, 0)
//
//	for _, report := range reports {
//
//		var user model.User
//
//		for _, u := range users {
//			if u.Id == report.UserId {
//				user = u
//				break
//			}
//		}
//
//		reportField := slack.NewTextBlockObject(
//			"mrkdwn",
//			fmt.Sprintf(
//				"*%s*\n```Делал вчера:\n%s\n\nДелает сегодня:\n%s\n\nБлокеры:\n%s```",
//				user.Name,
//				report.Done,
//				report.WillDo,
//				report.Blocker,
//			),
//			false,
//			false,
//		)
//
//		reportsSlice = append(reportsSlice, reportField)
//	}
//
//	willDoText := slack.NewTextBlockObject("mrkdwn", "*Что сегодня будет делать команда:*", false, false)
//	reportsSection := slack.NewSectionBlock(willDoText, reportsSlice, nil)
//
//	msg := slack.MsgOptionBlocks(
//		headerSection,
//		slack.NewDividerBlock(),
//		ignoreSection,
//		slack.NewDividerBlock(),
//		reportsSection,
//	)
//
//	var options []slack.MsgOption
//
//	options = append(options, msg)
//
//	if replace != "" {
//		s1, s2, _, err := a.slack.Client().UpdateMessage(
//			channelId,
//			replace,
//			msg,
//		)
//
//		return s1, s2, err
//	}
//
//	return a.slack.Client().PostMessage(
//		channelId,
//		options...,
//	)
//}