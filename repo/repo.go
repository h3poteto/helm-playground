package repo

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Repo struct {
	client *github.Client
}

func New(token string) *Repo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return &Repo{client: client}
}

func (r *Repo) GetRevision(owner, repository, branch string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gitBranch, _, err := r.client.Repositories.GetBranch(ctx, owner, repository, branch)
	if err != nil {
		return "", err
	}
	return *gitBranch.Commit.SHA, nil
}
