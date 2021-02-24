package database

import (
	"bot/internal/model"
	"time"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	UpdateOrCreate(user *model.User) error
}

type DailyReportRepository interface {
	Create(report *model.DailyReport) error
	Update(report *model.DailyReport) error
	GetByDate(time time.Time) ([]model.DailyReport, error)
	FindByUserAndDate(userId int, time time.Time) (*model.DailyReport, error)
	FindByUsersAndDate(usersId []int, time time.Time) ([]model.DailyReport, error)
	UpdateOrCreate(report *model.DailyReport) error
}

type SlackReportRepository interface {
	Create(report *model.SlackReport) error
	Update(report *model.SlackReport) error
	FindBySlackChannelAndDate(slackChannelId string, time time.Time) (*model.SlackReport, error)
	UpdateOrCreate(report *model.SlackReport) error
}