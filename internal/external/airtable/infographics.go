package airtable

//func (a *Airtable) GetInfographicsUsers() ([]User, error) {
//	a.logger.Info("Load infographics airtable users")
//
//	var users []User
//
//	usersTable := a.client.Table(a.infographics.config.team.table)
//
//	if err := usersTable.List(&users, &airtable.Options{
//		View: a.infographics.config.team.view,
//	}); err != nil {
//		a.logger.Error(err)
//		return nil, err
//	}
//
//	a.logger.Info("Total infographics airtable users: ", len(users))
//
//	return users, nil
//}