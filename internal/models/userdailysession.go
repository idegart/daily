package models

import "time"

type UserDailySession struct {
	Id             int       `json:"id" db:"id"`
	UserId         int       `json:"user_id" db:"user_id"`
	DailySessionId int       `json:"daily_session_id" db:"daily_session_id"`
	Done           string    `json:"done" db:"done"`
	WillDo         string    `json:"will_do" db:"will_do"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
