package f

import (
	"os"
)

type Message struct {
	Data []byte `json:"data"`
}

var tagChannelID string = os.Getenv("TAG_CHANNEL_ID")
var mergeChannelID string = os.Getenv("MERGE_CHANNEL_ID")
var slackToken string = os.Getenv("SLACK_TOKEN")
var gitlabToken string = os.Getenv("GITLAB_TOKEN")
