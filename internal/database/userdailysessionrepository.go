package database

import (
	"SlackBot/internal/models"
	"database/sql"
	"time"
)

type UserDailySessionRepository struct {
	database *Database
}

func (u *UserDailySessionRepository) Create(session *models.UserDailySession) error {
	if err := u.database.db.QueryRow(
		"INSERT INTO user_daily_sessions (user_id, daily_session_id) VALUES ($1, $2) RETURNING id",
		session.UserId,
		session.DailySessionId,
	).Scan(&session.Id); err != nil {
		return err
	}

	return nil
}

func (u *UserDailySessionRepository) Update(session *models.UserDailySession) error {
	_, err := u.database.db.NamedExec(
		"UPDATE user_daily_sessions SET done=:done, will_do=:will, updated_at=:updated WHERE id=:id",
		map[string]interface{}{
			"id": session.Id,
			"done": session.Done,
			"will": session.WillDo,
			"updated": time.Now(),
		},
	)

	return err
}

func (u *UserDailySessionRepository) FindByDailySessionUser(dailySession *models.DailySession, user *models.User) (*models.UserDailySession, error) {
	var session models.UserDailySession

	if err := u.database.db.Get(
		&session,
		"SELECT * FROM user_daily_sessions WHERE daily_session_id=$1 and user_id=$2",
		dailySession.Id,
		user.Id,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

func (u *UserDailySessionRepository) FindOrCreateByDailySessionAndUser(dailySession *models.DailySession, user *models.User) (*models.UserDailySession, error) {
	session, err := u.FindByDailySessionUser(dailySession, user)

	if err != nil {
		if err == sql.ErrNoRows {
			session = &models.UserDailySession{
				DailySessionId: dailySession.Id,
				UserId:         user.Id,
			}

			if err := u.Create(session); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return session, nil
}
