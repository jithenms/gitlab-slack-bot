package f

import (
	"fmt"

	"github.com/slack-go/slack"
)

func GenerateNotificationOption(text string, params ...interface{}) slack.MsgOption {
	return slack.MsgOptionText(fmt.Sprintf(text, params...), false)
}

func GenerateHeaderBlock(id string, text string, params ...interface{}) *slack.HeaderBlock {
	return slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", fmt.Sprintf(text, params...), false, false), slack.HeaderBlockOptionBlockID(id))
}
