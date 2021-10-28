package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	githubToken = envString("GITHUB_TOKEN", "")
	httpAddress = envString("HTTP_ADDRESS", ":8042")
	slackAPI    = envString("SLACK_API_TOKEN", "")
)

type Service struct {
	Logger       log.Logger
	GithubClient *GithubClient
	SlackClient  *SlackClient
}

func main() {
	logger := log.Logger{}

	service := Service{
		GithubClient: NewGithubClient(githubToken),
		SlackClient:  NewSlackClient(slackAPI),
	}

	errorChannel := make(chan error)
	// HTTP transport.
	go func() {
		httpServer := NewServer(
			logger,
			httpAddress,
			service,
		)
		errorChannel <- httpServer.Open()
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errorChannel <- fmt.Errorf("got signal: %s", <-c)
	}()

	// Wait for any error.
	if err := <-errorChannel; err != nil {
		logger.Print(err)
		os.Exit(1)
	}
}

func envString(key string, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
