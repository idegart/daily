package daily

//func (b *Bot) sendInitialMessageToUser(user slack.User) error {
//	_, _, err := b.bot.Api.PostMessage(
//		user.ID,
//		slack.MsgOptionText("Привет, настало время чтобы рассказать чем ты занимался вчера", false),
//	)
//
//	return err
//}
//
//func (b *Bot) setUserDone(userSession *models.UserDailySession, event *slackevents.MessageEvent) {
//	userSession.Done = event.Text
//	if err := b.database.UserDailySession().Update(userSession); err != nil {
//		b.logger.Error("Can not update done for user session: ", err)
//		return
//	}
//
//	if _, _, err := b.bot.Api.PostMessage(
//		event.User,
//		slack.MsgOptionText("Отлично, а чем планируешь заняться?", false),
//	); err != nil {
//		b.logger.Errorf("Error when send done message to %s: %s", event.User, err)
//	}
//}
//
//func (b *Bot) setUserWillDo(userSession *models.UserDailySession, event *slackevents.MessageEvent) {
//	userSession.WillDo = event.Text
//	if err := b.database.UserDailySession().Update(userSession); err != nil {
//		b.logger.Error("Can not update will do for user session: ", err)
//		return
//	}
//
//	if _, _, err := b.bot.Api.PostMessage(
//		event.User,
//		slack.MsgOptionText("Есть что-то, что мешает тебе работать над этим?", false),
//	); err != nil {
//		b.logger.Errorf("Error when send will do message to %s: %s", event.User, err)
//	}
//}
//
//func (b *Bot) setBlocker(userSession *models.UserDailySession, event *slackevents.MessageEvent) {
//	userSession.Blocker = event.Text
//	if err := b.database.UserDailySession().Update(userSession); err != nil {
//		b.logger.Error("Can not update blocker for user session: ", err)
//		return
//	}
//
//	if _, _, err := b.bot.Api.PostMessage(
//		event.User,
//		slack.MsgOptionText("Спасибо, можешь работать дальше", false),
//	); err != nil {
//		b.logger.Errorf("Error when send blocker message to %s: %s", event.User, err)
//	}
//}
//
//func (b *Bot) sendDailySessionToChats(userSession *models.UserDailySession, event *slackevents.MessageEvent) {
//	params := &slack.GetConversationsParameters{}
//	b.getConversations(params)
//	//params := &slack.GetConversationsParameters{}
//	//channels, _, err := b.bot.Api.GetConversations(params)
//	//
//	//if err != nil {
//	//	b.logger.Error(err)
//	//	return
//	//}
//	//
//	////b.logger.Info(cursor)
//	//
//	//for _, channel := range channels {
//	//
//	//	if channel.IsArchived == false {
//	//		b.logger.Info(channel.Name)
//	//	}
//	//}
//}
//
//func (b *Bot) getConversations(params *slack.GetConversationsParameters) {
//
//	conversations, cursor, err := b.bot.Api.GetConversations(params)
//
//	if err != nil {
//		fmt.Printf("%s\n", err)
//		return
//	}
//
//	for _, conversation := range conversations {
//		fmt.Printf("ID: %s, Name: %s, Archived: %t\n", conversation.ID, conversation.Name, conversation.IsArchived)
//	}
//
//	if cursor != "" {
//		params.Cursor = cursor
//		b.getConversations(params)
//	}
//}