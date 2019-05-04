package plateform

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Github struct {
	client *github.Client
	repo   *github.Repository
}

func NewGithubClient(token string, user string, repo string) (*Github, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// get go-github client
	client := github.NewClient(tc)

	r, _, err := client.Repositories.Get(ctx, user, repo)
	if err != nil {
		return nil, errors.Wrapf(err, "can't find repo %s with owner %s", repo, user)
	}

	return &Github{
		client: client,
		repo:   r,
	}, nil
}

func (g *Github) GetStars() *int {
	return g.repo.StargazersCount
}
