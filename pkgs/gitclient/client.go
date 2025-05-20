package gitclient

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
)

func NewGitClient(token string) *github.Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

type HookDetails struct {
	RepoName	string
	Events   	[]string
}

func CreateHook(client *github.Client, details HookDetails, owner string) (*github.Hook, error) {
	ctx := context.Background()

	hook := &github.Hook{
		Events: details.Events,
		Config: map[string]interface{}{
			"url":          "https://server.com/webhook",
			"content_type": "json",
			"secret":       os.Getenv("SECRET"),
			"insecure_ssl": "0",
		},
	}

	createdHook, _, err := client.Repositories.CreateHook(ctx, owner, details.RepoName, hook)

	if err != nil {
		log.Fatalf("Erro ao criar webhook: %v", err)
		return nil, err
	}

	log.Printf("Webhook criado com ID: %d\n", *createdHook.ID)
	return createdHook, nil
}
