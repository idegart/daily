package daily

import (
	"github.com/slack-go/slack"
	"os"
)

func (d *Daily) GoSmoke() error {
	messageOptions := []slack.MsgOption{
		slack.MsgOptionText("Пора покурить (или нет)", false),
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					"mrkdwn",
					"*🚬🤖 Бот сказал - пора курить. (или нет)*",
					false,
					false,
				),
				nil,
				nil,
			),
			slack.NewContextBlock(
				"",
				slack.NewTextBlockObject(
					"mrkdwn",
					"🔞 Курение вредит Вашему здоровью 🚭",
					false,
					false,
				),
			),
		),
	}

	_, _, _, err := d.slack.Client().SendMessage(os.Getenv("SLACK_CLUB_KURILKA"), messageOptions...)

	return err
}
