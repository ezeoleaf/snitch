package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// NewGithubClient returns a new Github Client
func NewGithubClient(token string) *GithubClient {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	uClient, err := github.NewEnterpriseClient("https://github.schibsted.io/api/v3", "https://github.schibsted.io/api/uploads", tc)

	if err != nil {
		panic(err)
	}

	c := &GithubClient{
		searchClient: uClient.Search,
		prClient:     uClient.PullRequests,
	}

	return c
}

// GetPRsByUser returns all of the PRs in which the username is missing a review
func (c *GithubClient) GetPRsByUser(ctx context.Context, username string) ([]PR, error) {
	qs := fmt.Sprintf("review-requested:%s state:open", username)

	issues, _, err := c.searchClient.Issues(ctx, qs, nil)

	if err != nil {
		return nil, err
	}

	prs := []PR{}

	for _, i := range issues.Issues {
		pr, err := c.enrichPR(ctx, i)

		if err != nil {
			continue
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

// GetPRsOwned returns all of the PRs created by username and that are missing reviews
func (c *GithubClient) GetPRsOwned(ctx context.Context, username string) ([]PR, error) {
	qs := fmt.Sprintf("author:%s state:open", username)

	issues, _, err := c.searchClient.Issues(ctx, qs, nil)

	if err != nil {
		return nil, err
	}

	prs := []PR{}

	for _, i := range issues.Issues {
		pr, err := c.enrichPR(ctx, i)

		if err != nil {
			continue
		}

		prs = append(prs, pr)
	}

	return prs, nil
}

func (c *GithubClient) enrichPR(ctx context.Context, issue *github.Issue) (PR, error) {
	repo := strings.Split(issue.GetRepositoryURL(), "/")

	org, name := repo[len(repo)-2], repo[len(repo)-1]

	repoName := fmt.Sprintf("%s/%s", org, name)

	pr := PR{
		Title:      issue.GetTitle(),
		Repository: repoName,
		URL:        issue.GetHTMLURL(),
		Owner:      *issue.User.Login,
		CreatedAt:  issue.CreatedAt,
		Reviewers:  []string{},
	}

	reviewers, _, err := c.prClient.ListReviewers(ctx, org, name, *issue.Number, nil)
	if err == nil {
		for _, reviewer := range reviewers.Users {
			pr.Reviewers = append(pr.Reviewers, *reviewer.Login)
		}
	}

	return pr, nil
}

// GetPRsByRepository returns all of the PRs in a repository and that are missing reviews
func (c *GithubClient) GetPRsByRepository(ctx context.Context, repository string) ([]PR, error) {
	qs := fmt.Sprintf("repo:%s state:open", repository)

	issues, _, err := c.searchClient.Issues(ctx, qs, nil)

	if err != nil {
		return nil, err
	}

	prs := []PR{}

	for _, i := range issues.Issues {
		pr, err := c.enrichPR(ctx, i)

		if err != nil {
			continue
		}

		prs = append(prs, pr)
	}

	return prs, nil
}
