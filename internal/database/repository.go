package database

import (
	"bot/internal/external/airtable"
	"bot/internal/model"
	"time"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	UpdateOrCreate(user *model.User) error
	GenerateFromAirtable(airUsers []airtable.User) ([]model.User, error)
	GetInfographicsUsers() ([]model.User, error)
}

type DailyReportRepository interface {
	Create(report *model.DailyReport) error
	Update(report *model.DailyReport) error
	GetByDate(time time.Time) ([]model.DailyReport, error)
	FindByUserAndDate(userId int, time time.Time) (*model.DailyReport, error)
	FindByUsersAndDate(usersId []int, time time.Time) ([]model.DailyReport, error)
	UpdateOrCreate(report *model.DailyReport) error
	GetLastUserReport(userId int) (*model.DailyReport, error)
}

type SlackReportRepository interface {
	Create(report *model.SlackReport) error
	Update(report *model.SlackReport) error
	GetAllByDate(time time.Time) ([]model.SlackReport, error)
	FindBySlackChannelAndDate(slackChannelId string, time time.Time) (*model.SlackReport, error)
	UpdateOrCreate(report *model.SlackReport) error
}

type AbsentUserRepository interface {
	Create(absentUser *model.AbsentUser) error
	GetAllByDate(time time.Time) ([]model.AbsentUser, error)
	GenerateFromAirtableForDate(airUsers []airtable.AbsentUser, users []model.User, date time.Time) ([]model.AbsentUser, error)
}

type ProjectRepository interface {
	Create(project *model.Project) error
	Update(project *model.Project) error
	GetAll() ([]model.Project, error)
	GenerateFromAirtable(airProjects []airtable.Project) ([]model.Project, error)
	GetUsersForProject(project model.Project) ([]model.User, error)
	AttachUsers(project model.Project, users []model.User) error
	AttachUser(project model.Project, user model.User) error
	DettachUsers(project model.Project, users []model.User) error
	SyncUsers(project model.Project, users []model.User) error
}
