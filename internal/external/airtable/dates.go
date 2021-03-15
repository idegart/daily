package airtable

import (
	"github.com/brianloveswords/airtable"
	"time"
)

const holidaysTable = "tblFc9I08E7iERCRe"

type Holiday struct {
	airtable.Record
	Fields struct {
		Date        string
		Description string
	}
}

func (h Holiday) IsCurrentDate() bool  {
	return time.Now().Format("2006-01-02") == h.Fields.Date
}

func (a *Airtable) GetHolidayDays() ([]Holiday, error) {
	a.logger.Info("Load holiday days")

	a.client.BaseID = a.team.appID

	usersTable := a.client.Table(holidaysTable)

	var holidays []Holiday

	if err := usersTable.List(&holidays, &airtable.Options{}); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return holidays, nil
}

func (a *Airtable) TodayIsHoliday() (string, error) {
	days, err := a.GetHolidayDays()

	if err != nil {
		return "", err
	}

	for _, day := range days {
		if day.IsCurrentDate() {
			return day.Fields.Description, err
		}
	}

	return "", nil
}
