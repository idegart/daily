package models

import "time"

type DailyReport struct {
	Id        int       `json:"id" db:"id"`
	Date      time.Time `json:"date" db:"date"`
	UserId    int       `json:"user_id" db:"user_id"`
	Done      string    `json:"done" db:"done"`
	WillDo    string    `json:"will_do" db:"will_do"`
	Blocker   string    `json:"blocker" db:"blocker"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
