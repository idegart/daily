package main

func (a *App) initChats() error {

	channels, err := a.slack.GetActiveProjectChannels(true)

	if err != nil {
		return err
	}

	a.slackProjects = channels

	return nil
}
