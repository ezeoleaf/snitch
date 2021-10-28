package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// NewServer creates a new Server instance.
func NewServer(
	logger log.Logger, addr string, service Service,
) Server {
	return Server{
		logger:  logger,
		address: addr,
		service: service,
	}
}

// Server is the HTTP server used to serve requests.
type Server struct {
	address string
	logger  log.Logger

	service Service
}

// Open will setup a tcp listener and serve the http requests.
func (s Server) Open() error {
	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints
	router.HandleFunc("/my-prs", s.GetMyPendingPRs).Methods("POST")
	router.HandleFunc("/prs", s.GetPRsByRepository).Methods("POST")
	router.HandleFunc("/user-prs", s.GetPRsByUsername).Methods("POST")
	router.HandleFunc("/owned-prs", s.GetOwnedPRs).Methods("POST")

	// serve the app
	return http.ListenAndServe(":8042", router)
}

// GetMyPendingPRs post a message with all the PRs I still have to review
func (s Server) GetMyPendingPRs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelID := r.FormValue("channel_id")
	username := strings.ReplaceAll(r.FormValue("user_name"), ".", "-")
	channelName := r.FormValue("channel_name")

	if channelName == DM { // This should be done so we can publish both IM and public messages
		channelID = r.FormValue("user_id")
	}

	// Workaround for not being able to use Schibsted Slack space to get the username
	if username == "ezequiel-olea-figuero" || username == "ezeoleaf" {
		username = "ezequiel-olea-figueroa"
	}

	prs, err := s.service.GithubClient.GetPRsByUser(ctx, username)

	if err != nil {
		fmt.Println(err)
	}

	err = s.service.SlackClient.SendMessage(ctx, channelID, prs, nil, nil)

	if err != nil {
		json.NewEncoder(w).Encode("Slack message not posted")
	}
}

// GetOwnedPRs post a message with all the PRs created by me and that still missing reviews
func (s Server) GetOwnedPRs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelID := r.FormValue("channel_id")
	username := strings.ReplaceAll(r.FormValue("user_name"), ".", "-")
	channelName := r.FormValue("channel_name")

	if channelName == DM { // This should be done so we can publish both IM and public messages
		channelID = r.FormValue("user_id")
	}

	// Workaround for not being able to use Schibsted Slack space to get the username
	if username == "ezequiel-olea-figuero" || username == "ezeoleaf" {
		username = "ezequiel-olea-figueroa"
	}

	prs, err := s.service.GithubClient.GetPRsOwned(ctx, username)

	if err != nil {
		fmt.Println(err)
	}

	err = s.service.SlackClient.SendMessage(ctx, channelID, prs, nil, nil)

	if err != nil {
		json.NewEncoder(w).Encode("Slack message not posted")
	}
}

// GetPRsByUsername post a message with all the PRs in which the username is missing reviews
func (s Server) GetPRsByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelID := r.FormValue("channel_id")
	channelName := r.FormValue("channel_name")
	username := r.FormValue("text")

	if channelName == DM { // This should be done so we can publish both IM and public messages
		channelID = r.FormValue("user_id")
	}

	prs, err := s.service.GithubClient.GetPRsByUser(ctx, username)

	if err != nil {
		fmt.Println(err)
	}

	err = s.service.SlackClient.SendMessage(ctx, channelID, prs, nil, &username)

	if err != nil {
		json.NewEncoder(w).Encode("Slack message not posted")
	}
}

// GetPRsByRepository post a message with all the PRs from a repository that are missing reviews
func (s Server) GetPRsByRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelID := r.FormValue("channel_id")
	channelName := r.FormValue("channel_name")
	repository := r.FormValue("text")

	if channelName == DM { // This should be done so we can publish both IM and public messages
		channelID = r.FormValue("user_id")
	}

	prs, err := s.service.GithubClient.GetPRsByRepository(ctx, repository)

	if err != nil {
		fmt.Println(err)
	}

	err = s.service.SlackClient.SendMessage(ctx, channelID, prs, &repository, nil)

	if err != nil {
		json.NewEncoder(w).Encode("Slack message not posted")
	}
}
