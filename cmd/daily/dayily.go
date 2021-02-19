package main
//
//import (
//	"bot/internal/model"
//	"database/sql"
//	"errors"
//	"github.com/slack-go/slack"
//	"strings"
//	"time"
//)
//
//func (a *App) sendInitialMessages() {
//	if err := app.initUsers(); err != nil {
//		app.logger.Fatal(err)
//	}
//
//	for i := range app.users {
//		//if app.users[i].Email == "a.degtyarev@proscom.ru" {
//			a.sendInitialMessageToUser(app.users[i].SlackId)
//		//}
//	}
//}
//
//func (a *App) sendReports() {
//	if err := app.initUsers(); err != nil {
//		app.logger.Fatal(err)
//	}
//
//	if err := app.initChats(); err != nil {
//		app.logger.Fatal(err)
//	}
//
//	reports, err := a.database.DailyReport().GetByDate(time.Now())
//
//	if err != nil {
//		app.logger.Fatal(err)
//	}
//
//	for i := range a.slackProjects {
//		a.sendReportForSlackChannel(a.slackProjects[i], reports)
//	}
//}
//
//func (a *App) sendReportForSlackChannel(channel slack.Channel, reports []model.DailyReport) {
//	slackUsers, err := a.slack.GetActiveUsersInChannel(channel.ID)
//
//	if err != nil {
//		a.logger.Error(err)
//	}
//
//	channelUsers := a.GetUsersBySlackUsers(slackUsers)
//
//	var channelReports []model.DailyReport
//	var badChannelUsers []model.User
//
//out:
//	for i := range channelUsers {
//		for j := range reports {
//			if channelUsers[i].Id == reports[j].UserId {
//				channelReports = append(channelReports, reports[j])
//				continue out
//			}
//		}
//
//		badChannelUsers = append(badChannelUsers, channelUsers[i])
//	}
//
//	slackReport, err := a.database.SlackReport().FindBySlackChannelAndDate(channel.ID, time.Now())
//
//	if err != nil && !errors.Is(err, sql.ErrNoRows) {
//		a.logger.Error(err)
//	}
//
//	var replace string
//
//	if slackReport != nil {
//		replace = slackReport.Ts
//	} else {
//		slackReport = &model.SlackReport{
//			SlackChannelId: channel.ID,
//			Date:           time.Now(),
//		}
//	}
//
//	_, ts, err := a.sendReportToChannel(channel.ID, channelUsers, badChannelUsers, reports, replace)
//
//	if err != nil {
//		a.logger.Error(err)
//	}
//
//	slackReport.Ts = ts
//
//	if err := a.database.SlackReport().UpdateOrCreate(slackReport); err != nil {
//		a.logger.Error(err)
//	}
//}
//
//func (a *App) startForUserByCallback(callback *slack.InteractionCallback) {
//	a.sendReportModal(callback)
//}
//
//func (a *App) finishUserReportByCallback(callback *slack.InteractionCallback, user *model.User) {
//	data := callback.DialogSubmissionCallback.Submission
//
//	date, err := time.Parse("2006-01-02", strings.ReplaceAll(callback.State, "\"", ""))
//
//	if err != nil {
//		a.logger.Error(err)
//		return
//	}
//
//	var report = &model.DailyReport{
//		UserId:  user.Id,
//		Date:    date,
//		Done:    data["Done"],
//		WillDo:  data["WillDo"],
//		Blocker: data["Blocker"],
//	}
//
//	if err := a.database.DailyReport().UpdateOrCreate(report); err != nil {
//		a.logger.Error(err)
//		return
//	}
//
//	a.sendThanksForReport(callback)
//}
