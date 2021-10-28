package main

import (
	"context"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/slack-go/slack"
)

const DM = "directmessage"

// GithubClient is ..
type GithubClient struct {
	searchClient *github.SearchService
	prClient     *github.PullRequestsService
}

// PR is a representation of the Pull request data
type PR struct {
	Title      string     `json:"title"`
	Repository string     `json:"repository"`
	URL        string     `json:"url"`
	Owner      string     `json:"owner"`
	CreatedAt  *time.Time `json:"created_at"`
	Reviewers  []string   `json:"reviewers"`
}

// SlackClient is ..
type SlackClient struct {
	client *slack.Client
}

type SlackService interface {
	SendMessage(ctx context.Context, channelID string, prs []PR) error
}
