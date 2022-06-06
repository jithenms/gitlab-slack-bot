package f

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
)

func TagEventHandler(ctx context.Context, m Message) error {
	var event gitlab.TagEvent

	if err := json.Unmarshal(m.Data, &event); err != nil {
		log.Fatalf("Failed to parse PubSub message: %v", err)
		return err
	}

	var s = slack.New(slackToken)

	ref := strings.Split(event.Ref, "/")
	notificationOption := GenerateNotificationOption("New Tag %v in %v", ref[len(ref)-1], event.Project.Name)
	headerBlock := GenerateHeaderBlock("tag-header", "New Tag in %v", event.Project.Name)
	tagSection := generateTagSection(event)

	channel, timestamp, text, err := s.SendMessage(
		tagChannelID,
		notificationOption,
		slack.MsgOptionBlocks(
			headerBlock,
			tagSection,
		),
	)

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Tag Event message: %v\n successfully sent to channel %s at %s\n",
		text, channel, timestamp)

	return nil
}

func generateTagSection(event gitlab.TagEvent) *slack.SectionBlock {
	ref := strings.Split(event.Ref, "/")
	tagBlock := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Tag:*\n%v", ref[len(ref)-1]), false, false)
	serviceBlock := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Service:*\n%v", event.Project.Name), false, false)
	return slack.NewSectionBlock(nil, []*slack.TextBlockObject{tagBlock, serviceBlock}, nil)
}
