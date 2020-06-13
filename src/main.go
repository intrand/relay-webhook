package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type dockerHubRepository struct {
	CommentCount    int64   `json:"comment_count"`
	DateCreated     float64 `json:"date_created"`
	Description     string  `json:"description"`
	Dockerfile      string  `json:"dockerfile"`
	FullDescription string  `json:"full_description"`
	IsOfficial      bool    `json:"is_official"`
	IsPrivate       bool    `json:"is_private"`
	IsTrusted       bool    `json:"is_trusted"`
	Name            string  `json:"name"`
	Namespace       string  `json:"namespace"`
	Owner           string  `json:"owner"`
	RepoName        string  `json:"repo_name"`
	RepoURL         string  `json:"repo_url"`
	StarCount       int64   `json:"star_count"`
	Status          string  `json:"status"`
}

type dockerHubPushData struct {
	Images   []string `json:"images"`
	PushedAt float64  `json:"pushed_at"`
	Pusher   string   `json:"pusher"`
	Tag      string   `json:"tag"`
}

type dockerHubIncomingWebhook struct {
	CallbackURL string              `json:"callback_url"`
	PushData    dockerHubPushData   `json:"push_data"`
	Repository  dockerHubRepository `json:"repository"`
}

type discordOut struct {
	Content string `json:"content"`
	Name    string `json:"username"`
}

func main() {
	webhookURL := os.Getenv("discord_webhook")
	whURL := flag.String("webhook.url", webhookURL, "")
	flag.Parse()

	if webhookURL == "" && *whURL == "" {
		fmt.Fprintf(os.Stderr, "error: environment variable discord_webhook not found\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "info: Listening on 0.0.0.0:8000\n")
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		dha := dockerHubIncomingWebhook{}
		err = json.Unmarshal(b, &dha)
		if err != nil {
			panic(err)
		}

		DO := discordOut{
			Name: "hub.docker.com",
		}

		DO.Content = fmt.Sprintf("`%s:%s` updated.", dha.Repository.RepoName, dha.PushData.Tag)

		DOD, _ := json.Marshal(DO)
		http.Post(*whURL, "application/json", bytes.NewReader(DOD))
	}))
}
