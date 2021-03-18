package daily

import (
	"bot/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"strings"
	"time"
)

const InfographicsSlackId = "C01HPPN012P"

func (d *Daily) StartReport() error {
	d.logger.Info("Start report")

	if err := d.Init(); err != nil {
		return err
	}

	for _, project := range d.projects {
		d.projectsToReport <- project
	}

	if err := d.SendReportToInfographics(); err != nil {
		d.logger.Error(err)
	}

	return nil
}

func (d Daily) DropReports() error {
	d.logger.Info("Drop reports")

	reports, err  := d.database.SlackReport().GetAllByDate(time.Now())

	if err != nil {
		return err
	}

	for _, report := range reports {
		if _, _, err := d.slack.Client().DeleteMessage(report.SlackChannelId, report.Ts); err != nil {
			d.logger.Error(err)
		}
	}

	return nil
}

func (d *Daily) SendUpdatingReportByUser(user model.User) {
	if user.IsInfographic {
		if _, err := d.database.SlackReport().FindBySlackChannelAndDate(InfographicsSlackId, time.Now()); err == nil {
			if err := d.SendReportToInfographics(); err != nil {
				d.logger.Error(err)
			}
		}
	}

	for _, project := range d.projects {
		for _, u := range project.Users {
			if u.Id == user.Id {
				_, err := d.database.SlackReport().FindBySlackChannelAndDate(project.SlackId, time.Now())

				if err == nil {
					d.projectsToReport <- project
				}
			}
		}
	}
}

func (d *Daily) startSendingReports()  {
	for project := range d.projectsToReport {
		if project.IsInfographics == false {
			if err := d.SendReportToProject(project); err != nil {
				d.logger.Error(err)
			}
		}
	}
}

func (d *Daily) SendReportToProject(project model.Project) error {
	d.logger.Infof("Sending report to %s (%d, %s)", project.Name, project.Id, project.SlackId)

	var ids []int

	for _, user := range project.Users {
		ids = append(ids, user.Id)
	}

	reports, err := d.database.DailyReport().FindByUsersAndDate(ids, time.Now())

	if err != nil {
		return err
	}

	report, err := d.database.SlackReport().FindBySlackChannelAndDate(project.SlackId, time.Now())

	var replaceTS string

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.logger.Error(err)
	}

	if report != nil {
		replaceTS = report.Ts
	}

	d.sendSlackReportToChannel(
		project.SlackId,
		project,
		reports,
		replaceTS,
	)

	return nil
}

func (d *Daily) SendReportToInfographics() error {
	d.logger.Infof("Sending report to Infographics")

	var ids []int
	var users []model.User

	for _, user := range d.users {
		if user.IsInfographic {
			ids = append(ids, user.Id)
			users = append(users, user)
		}
	}

	reports, err := d.database.DailyReport().FindByUsersAndDate(ids, time.Now())

	if err != nil {
		return err
	}

	var replaceTS string

	slackReport, err := d.database.SlackReport().FindBySlackChannelAndDate(InfographicsSlackId, time.Now())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.logger.Error(err)
	}

	if slackReport != nil {
		replaceTS = slackReport.Ts
	}

	d.sendSlackReportToInfographics(users, reports, replaceTS)

	return nil
}

func (d Daily) sendSlackReportToInfographics(users []model.User, reports []model.DailyReport, replace string)  {
	messageBlocks := []slack.Block{
		getHeaderSection(),
		slack.NewDividerBlock(),
		getWillDoSection(),
	}

	messageBlocks = append(messageBlocks, getReportsSection(users, reports)...)

	msg := slack.MsgOptionCompose(
		slack.MsgOptionText("Отчет для Инфографики", false),
		slack.MsgOptionBlocks(messageBlocks...),
	)

	var reportBlockOptions []slack.MsgOption

	reportBlockOptions = append(reportBlockOptions, msg)

	if replace != "" {
		reportBlockOptions = append(reportBlockOptions, slack.MsgOptionUpdate(replace))
	}

	d.slack.SendMessage(
		InfographicsSlackId,
		func(ts string) {
			report := &model.SlackReport{
				SlackChannelId: InfographicsSlackId,
				Ts: ts,
				Date: time.Now(),
			}
			if err := d.database.SlackReport().UpdateOrCreate(report); err != nil {
				d.logger.Error(err)
			}
		},
		reportBlockOptions...,
	)
}

func (d *Daily) sendSlackReportToChannel(channelId string, project model.Project, reports []model.DailyReport, replace string) {
	messageBlocks := []slack.Block{
		getHeaderSection(),
		slack.NewDividerBlock(),
		getAbsentSection(project, d.absentUsers),
		slack.NewDividerBlock(),
		getIgnoreSection(project, d.absentUsers, reports),
		slack.NewDividerBlock(),
		getWillDoSection(),
	}

	messageBlocks = append(messageBlocks, getReportsSection(project.Users, reports)...)
	messageBlocks = append(messageBlocks, slack.NewDividerBlock())

	msg := slack.MsgOptionCompose(
		slack.MsgOptionText("Отчет для " + project.Name, false),
		slack.MsgOptionBlocks(messageBlocks...),
	)

	var reportBlockOptions []slack.MsgOption

	reportBlockOptions = append(reportBlockOptions, msg)

	if replace != "" {
		reportBlockOptions = append(reportBlockOptions, slack.MsgOptionUpdate(replace))
	}

	d.slack.SendMessage(
		channelId,
		func(ts string) {
			report := &model.SlackReport{
				SlackChannelId: channelId,
				Ts: ts,
				Date: time.Now(),
			}
			if err := d.database.SlackReport().UpdateOrCreate(report); err != nil {
				d.logger.Error(err)
			}
		},
		reportBlockOptions...,
	)
}

func getHeaderSection() *slack.SectionBlock {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			"*📊 Гайз, я тут подготовил ежедневный отчет. Чек зис аут*",
			false,
			false,
		),
		nil,
		nil,
	)
}

func getAbsentSection(project model.Project, absentUsers []model.AbsentUser) *slack.SectionBlock  {
	var projectAbsents []model.User
	absentText := "*👨‍👧Кто сегодня отсутствует:*\n"

	for _, user := range project.Users {
		for _, absentUser := range absentUsers {
			if absentUser.UserId == user.Id {
				projectAbsents = append(projectAbsents, user)
			}
		}
	}

	if len(projectAbsents) > 0 {
		for _, user := range projectAbsents {
			absentText += fmt.Sprintf(
				"<https://proscomteam.slack.com/team/%s|%s %s>\n",
				user.SlackId,
				user.Name,
				user.Emoji,
			)
		}
	} else {
		absentText = "*👨‍👩‍👧‍👦 Сегодня все на месте*"
	}

	return slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			absentText,
			false,
			false,
		),
		nil,
		nil,
	)
}

func getIgnoreSection(project model.Project, absentUsers []model.AbsentUser, reports []model.DailyReport) *slack.SectionBlock {
	var badUsersIds []string

	LOOP:
	for _, user := range project.Users {
		for _, absentUser := range absentUsers {
			if absentUser.UserId == user.Id {
				continue LOOP
			}
		}

		for _, report := range reports {
			if report.UserId == user.Id {
				continue LOOP
			}
		}

		badUsersIds = append(badUsersIds, "<@"+user.SlackId+">")
	}

	var ignoreText string

	if len(badUsersIds) < 1 {
		ignoreText = "*❤️ Все сегодня молодцы. Никто не проигнорировал*"
	} else {
		ignoreText = fmt.Sprintf("*💔 Кто меня сегодня проигнорировал:*\n%s", strings.Join(badUsersIds, "\n"))
	}

	return slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			ignoreText,
			false,
			false,
		),
		nil,
		nil,
	)
}

func getWillDoSection() *slack.SectionBlock {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			"*🧑🏻‍💻 Что сегодня будет делать команда:*",
			false,
			false,
		),
		nil,
		nil,
	)
}

func getReportsSection(users []model.User, reports []model.DailyReport) []slack.Block {
	var reportsBlocks []slack.Block

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
				strings.ReplaceAll(reportMessage, "\n", "\n>"),
				false,
				false,
			),
			nil,
			nil,
		)

		reportsBlocks = append(reportsBlocks, reportSection)
	}

	return reportsBlocks
}