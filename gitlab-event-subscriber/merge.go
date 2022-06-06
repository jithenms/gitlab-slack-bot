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

func MergeEventHandler(ctx context.Context, m Message) error {
	var event gitlab.MergeEvent

	if err := json.Unmarshal(m.Data, &event); err != nil {
		log.Fatalf("Failed to parse PubSub message: %v", err)
		return err
	}

	g, err := gitlab.NewClient(gitlabToken)

	if err != nil {
		log.Fatalf("Failed to create gitlab client: %v", err)
		return err
	}

	s := slack.New(slackToken)

	switch event.ObjectAttributes.Action {
	case "open":
		{
			if err := handleMergeOpen(event, mergeChannelID, g, s); err != nil {
				return err
			}
		}
	default:
		{
			return nil
		}
	}
	return nil
}

func handleMergeOpen(event gitlab.MergeEvent, channelID string, g *gitlab.Client, s *slack.Client) error {
	notificationOption := GenerateNotificationOption("New MR %v in %v", event.ObjectAttributes.Title, event.Project.Name)

	headerBlock := GenerateHeaderBlock("merge-header", "New MR in %v", event.Project.Name)

	mergeSection := generateMergeSection(event)

	assigneeBlock, err := generateAssigneeBlock(event, g, s)

	if err != nil {
		return err
	}

	reviewerBlock, err := generateReviewerBlock(event, g, s)

	if err != nil {
		return err
	}

	collaboratorSection := slack.NewSectionBlock(nil, []*slack.TextBlockObject{assigneeBlock, reviewerBlock}, nil, slack.SectionBlockOptionBlockID("collaborator-section"))

	channel, timestamp, text, err := s.SendMessage(
		channelID,
		notificationOption,
		slack.MsgOptionBlocks(
			headerBlock,
			slack.NewDividerBlock(),
			mergeSection,
			collaboratorSection,
		),
	)

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Merge event message: %v\n successfully sent to channel %s at %s\n",
		text, channel, timestamp)
	return nil
}

func generateMergeSection(event gitlab.MergeEvent) *slack.SectionBlock {

	nameBlock := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Name:*\n<%v|%v>",
			event.ObjectAttributes.URL,
			event.ObjectAttributes.Title), false, false)

	serviceBlock := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Service:*\n<%v|%v>",
			event.Repository.Homepage,
			event.Repository.Name), false, false)

	return slack.NewSectionBlock(nil, []*slack.TextBlockObject{nameBlock, serviceBlock}, nil, slack.SectionBlockOptionBlockID("merge-section"))
}

func generateAssigneeBlock(event gitlab.MergeEvent, g *gitlab.Client, s *slack.Client) (*slack.TextBlockObject, error) {
	assignee, _, err := g.Users.GetUser(event.ObjectAttributes.AssigneeID, gitlab.GetUsersOptions{})

	if err != nil {
		log.Fatalf("Failed to lookup gitlab reviewer: %v", err)
		return nil, err
	}

	if assignee.Email != "[REDACTED]" && assignee.Email != "" {
		owner, err := s.GetUserByEmail(assignee.Email)

		if err != nil {
			log.Fatalf("Failed to lookup slack assignee user: %v", err)
			return nil, err
		}

		return slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Assignee:*\n<@%v>",
			owner.ID), false, false), nil
	} else if assignee.PublicEmail != "[REDACTED]" && assignee.PublicEmail != "" {
		owner, err := s.GetUserByEmail(assignee.PublicEmail)

		if err != nil {
			log.Fatalf("Failed to lookup slack assignee user: %v", err)
			return nil, err
		}

		return slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Assignee:*\n<@%v>",
			owner.ID), false, false), nil
	} else {
		return slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Assignee:*\n%v",
			assignee.Name), false, false), nil
	}
}

func generateReviewerBlock(event gitlab.MergeEvent, g *gitlab.Client, s *slack.Client) (*slack.TextBlockObject, error) {
	mergeRequest, _, err := g.MergeRequests.GetMergeRequest(event.Project.PathWithNamespace, event.ObjectAttributes.IID,
		&gitlab.GetMergeRequestsOptions{})

	if err != nil {
		log.Fatalf("Failed to lookup merge request: %v", err)
		return nil, err
	}

	reviewers := []string{}

	if len(mergeRequest.Reviewers) == 0 {
		return slack.NewTextBlockObject("mrkdwn", "*Reviewers:*\nNone", false, false), nil
	}

	for _, user := range mergeRequest.Reviewers {
		reviewer, _, err := g.Users.GetUser(user.ID, gitlab.GetUsersOptions{})

		if err != nil {
			log.Fatalf("Failed to lookup reviewer: %v", err)
			return nil, err
		}

		if reviewer.Email != "[REDACTED]" && reviewer.Email != "" {
			user, err := s.GetUserByEmail(reviewer.Email)

			if err != nil {
				return nil, err
			}

			reviewers = append(reviewers, fmt.Sprintf("<@%s>", user.ID))
		} else if reviewer.PublicEmail != "[REDACTED]" && reviewer.PublicEmail != "" {
			user, err := s.GetUserByEmail(reviewer.PublicEmail)

			if err != nil {
				return nil, err
			}

			reviewers = append(reviewers, fmt.Sprintf("<@%s>", user.ID))
		} else {
			reviewers = append(reviewers, fmt.Sprint(reviewer.Name))
		}
	}

	return slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Reviewers:*\n%v",
		strings.Trim(fmt.Sprint(reviewers), "[]")), false, false), nil
}
