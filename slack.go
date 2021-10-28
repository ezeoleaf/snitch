package main

import (
	"context"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

// NewSlackClient returns a new Slack Client
func NewSlackClient(token string) *SlackClient {
	api := slack.New(token)

	c := &SlackClient{
		client: api,
	}

	return c
}

func (s *SlackClient) getMessageTitle(prs []PR, repoName, username *string) string {
	title := "PRs with pending reviews"

	if len(prs) == 0 {
		title = "No PRs pending"
	}

	if repoName != nil {
		title += " in " + *repoName
	}

	if username != nil {
		title += " from " + *username
	}
	return title
}

// SendMessage will post a message with all of the PRs in the specified channelID
func (s *SlackClient) SendMessage(ctx context.Context, channelID string, prs []PR, repositoryName, username *string) error {
	title := s.getMessageTitle(prs, repositoryName, username)

	divideBlock := slack.DividerBlock{
		Type: slack.MBTDivider,
	}

	blocks := []slack.Block{
		slack.HeaderBlock{
			Type: "header",
			Text: &slack.TextBlockObject{
				Type: slack.PlainTextType,
				Text: title,
			},
		},
		divideBlock,
	}

	for _, pr := range prs {

		age := int(time.Since(*pr.CreatedAt).Hours()/24) + 1

		reviewers := ""
		if len(pr.Reviewers) > 0 {
			reviewers += "\n*Reviewers:*"
			for _, r := range pr.Reviewers {
				reviewers += "\n - " + r
			}
		}

		block := slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf("%s in *%s*\nCreated *%v days ago* by *%s*%s\n<%s|Review>", pr.Title, pr.Repository, age, pr.Owner, reviewers, pr.URL),
			},
		}

		blocks = append(blocks, block, divideBlock)
	}

	_, _, err := s.client.PostMessage(
		channelID,
		slack.MsgOptionBlocks(blocks...),
	)

	if err != nil {
		return err
	}

	return nil
}
