package database

type Database interface {
	Open() error
	Close() error
	User() UserRepository
	DailyReport() DailyReportRepository
	SlackReport() SlackReportRepository
	AbsentUser() AbsentUserRepository
	Project() ProjectRepository
}
