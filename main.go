package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"snitch/provider"
	"snitch/publisher"
	"strconv"
	"syscall"
)

var (
	isEnterpriseGithubString = envString("ENTERPRISE_GITHUB", "false")
	githubAddress            = envString("GITHUB_ADDRES", "https://github.com/")
	githubToken              = envString("GITHUB_TOKEN", "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")
	httpAddress              = envString("HTTP_ADDRESS", ":8042")
	slackAPI                 = envString("SLACK_API_TOKEN", "xoxb-2646296415622-2676762513744-TwY7b8bl9GP0aZTz5hEzEMdv")
)

const DM = "directmessage"

type Service struct {
	Logger       log.Logger
	GithubClient *provider.GithubClient
	SlackClient  *publisher.SlackClient
}

func main() {
	logger := log.Default()

	isEnterpriseGithub, err := strconv.ParseBool(isEnterpriseGithubString)
	if err != nil {
		logger.Fatalln(errors.New("error parsing ENTERPRISE_GITHUB as bool"))
		os.Exit(1)
	}

	service := Service{
		GithubClient: provider.NewGithubClient(githubToken, githubAddress, isEnterpriseGithub),
		SlackClient:  publisher.NewSlackClient(slackAPI),
	}

	errorChannel := make(chan error)
	// HTTP transport.
	go func() {
		httpServer := NewServer(
			*logger,
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
