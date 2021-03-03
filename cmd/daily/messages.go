package main

import (
	"bot/internal/model"
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

const (
	SIDailyReportCallbackStart  = "daily_report_start"
	SIDailyReportCallbackFinish = "daily_report_finish"

	SIDailyReportDone    = "Done"
	SIDailyReportWillDo  = "WillDo"
	SIDailyReportBlocker = "Blocker"
)

func (a *App) SendSlackInitialMessageToUser(user model.User, previousReport *model.DailyReport) error {
	var text string

	if previousReport != nil {
		text = fmt.Sprintf("*Твои прошлые планы:*\n%s", previousReport.WillDo)
	}

	_, _, err := a.slack.Client().PostMessage(
		user.SlackId,
		slack.MsgOptionAttachments(slack.Attachment{
			Pretext:    "Привет, настало время, чтобы рассказать чем ты занимаешься",
			Text: text,
			CallbackID: SIDailyReportCallbackStart,
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
		}),
	)

	return err
}

func (a *App) SendSlackReportModal(callback *slack.InteractionCallback, report *model.DailyReport) {
	var doneMessage string
	var willDoMessage string
	var blockerMessage string

	if report != nil {
		doneMessage = report.Done
		willDoMessage = report.WillDo
		blockerMessage = report.Blocker
	}

	blockerInput := slack.NewTextAreaInput(
		SIDailyReportBlocker,
		"Расскажи, что тебе может помешать в твоей работе",
		blockerMessage,
	)
	blockerInput.Optional = true

	if err := a.slack.Client().OpenDialog(
		callback.TriggerID,
		slack.Dialog{
			TriggerID:   callback.TriggerID,
			CallbackID:  SIDailyReportCallbackFinish,
			Title:       "Daily report",
			SubmitLabel: "Отправить",
			Elements: []slack.DialogElement{
				slack.NewTextAreaInput(
					SIDailyReportDone,
					"Опиши, что было сделано",
					doneMessage,
				),
				slack.NewTextAreaInput(
					SIDailyReportWillDo,
					"Опиши, что будешь делать сегодня",
					willDoMessage,
				),
				blockerInput,
			},
			State: callback.ResponseURL,
		},
	); err != nil {
		a.logger.Error(err)
	}
}

func (a *App) SendSlackThanksForReport(callback *slack.InteractionCallback) {
	replaceUrl := strings.ReplaceAll(callback.State, "\\", "")
	replaceUrl = strings.ReplaceAll(replaceUrl, "\"", "")
	if _, _, err := a.slack.Client().PostMessage(callback.User.ID,
		slack.MsgOptionAttachments(slack.Attachment{
			Pretext:    "Спасибо, можешь продолжить свою работу!",
			CallbackID: SIDailyReportCallbackStart,
			Color:      "#3AA3E3",
			Actions: []slack.AttachmentAction{
				{
					Name:  "accept",
					Text:  "Дорассказать",
					Style: "default",
					Type:  "button",
					Value: "accept",
				},
			},
		}),
		slack.MsgOptionReplaceOriginal(replaceUrl),
	); err != nil {
		a.logger.Error(err)
	}
}

func (a *App) SendSlackReportToChannel(channelId string, users []model.User, badUsers []model.User, reports []model.DailyReport, replace string) (string, string, error) {
	var messageBlocks []slack.Block

	headerSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			"*Гайз, я тут подготовил ежедневный отчет. Чек зис аут*",
			false,
			false,
		),
		nil,
		nil,
	)

	messageBlocks = append(messageBlocks, headerSection)
	messageBlocks = append(messageBlocks, slack.NewDividerBlock())

	var badUsersIds []string

	for _, user := range badUsers {
		badUsersIds = append(badUsersIds, "<@"+user.SlackId+">")
	}

	var ignoreText string

	if len(badUsersIds) < 1 {
		ignoreText = "*Все сегодня молодцы. Никто не проигнорировал*"
	} else {
		ignoreText = fmt.Sprintf("*Кто меня сегодня проигнорировал:*\n%s", strings.Join(badUsersIds, "\n"))
	}

	ignoreSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			ignoreText,
			false,
			false,
		),
		nil,
		nil,
	)

	messageBlocks = append(messageBlocks, ignoreSection)
	messageBlocks = append(messageBlocks, slack.NewDividerBlock())

	willDoSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			"*Что сегодня будет делать команда:*",
			false,
			false,
		),
		nil,
		nil,
	)

	messageBlocks = append(messageBlocks, willDoSection)

	for _, report := range reports {
		var user model.User

		for _, u := range users {
			if u.Id == report.UserId {
				user = u
				break
			}
		}

		reportMessage := fmt.Sprintf(
			"<https://proscomteam.slack.com/team/%s|%s> %s\n*Вчера:*\n%s\n\n*Сегодня:*\n%s",
			user.SlackId,
			user.Name,
			user.Emoji,
			strings.Trim(report.Done, "\n"),
			strings.Trim(report.WillDo, "\n"),
		)

		if report.Blocker != "" {
			reportMessage += fmt.Sprintf(
				"\n\n*Блокеры:*\n%s",
				strings.Trim(report.Blocker, "\n"),
			)
		}

		reportSection := slack.NewSectionBlock(
			slack.NewTextBlockObject(
				"mrkdwn",
				strings.ReplaceAll( reportMessage, "\n", "\n>"),
				false,
				false,
			),
			nil,
			nil,
		)

		messageBlocks = append(messageBlocks, reportSection, slack.NewDividerBlock())
	}

	msg := slack.MsgOptionCompose(
		slack.MsgOptionBlocks(messageBlocks...),
	)

	var reportBlockOptions []slack.MsgOption

	reportBlockOptions = append(reportBlockOptions, msg)

	if replace != "" {
		reportBlockOptions = append(reportBlockOptions, slack.MsgOptionUpdate(replace))
	}

	return a.slack.Client().PostMessage(
		channelId,
		reportBlockOptions...,
	)
}
