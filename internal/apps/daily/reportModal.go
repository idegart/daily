package daily

import (
	"bot/internal/model"
	"github.com/slack-go/slack"
	"strings"
)

func (d *Daily) SendSlackReportModal(callback *slack.InteractionCallback, report *model.DailyReport, latestReport *model.DailyReport) {
	var doneMessage string
	var willDoMessage string
	var blockerMessage string

	var isLatestReportUsed bool

	if report != nil {
		doneMessage = report.Done
		willDoMessage = report.WillDo
		blockerMessage = report.Blocker
	}

	if report == nil && latestReport != nil {
		doneMessage = latestReport.WillDo
		isLatestReportUsed = true
	}

	blockerInput := slack.NewTextAreaInput(
		SIDailyReportBlocker,
		"Расскажи, что тебе может помешать в твоей работе",
		blockerMessage,
	)
	blockerInput.Optional = true

	doneInput := slack.NewTextAreaInput(
		SIDailyReportDone,
		"Опиши, что было сделано",
		doneMessage,
	)

	if isLatestReportUsed {
		doneInput.Hint = "Составлено на основе твоего прошлого ответа"
	}

	if err := d.slack.Client().OpenDialog(
		callback.TriggerID,
		slack.Dialog{
			TriggerID:   callback.TriggerID,
			CallbackID:  SIDailyReportCallbackFinish,
			Title:       "Daily report",
			SubmitLabel: "Отправить",
			Elements: []slack.DialogElement{
				doneInput,
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
		d.logger.Error(err)
	}
}

func (d *Daily) SendSlackThanksForReport(callback *slack.InteractionCallback, user model.User, report *model.DailyReport) {
	replaceUrl := strings.ReplaceAll(callback.State, "\\", "")
	replaceUrl = strings.ReplaceAll(replaceUrl, "\"", "")

	previousReport, _ := d.database.DailyReport().GetLastUserReport(user.Id)

	d.sendSlackInitialMessageToUser(user, previousReport, report, &replaceUrl)
}