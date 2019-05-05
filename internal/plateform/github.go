package plateform

import (
	"context"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Github struct {
	client *github.Client
	repo   string
	owner  string
}

func NewGithubClient(token string, owner string, repo string) (*Github, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// get go-github client
	client := github.NewClient(tc)
	return &Github{
		client: client,
		repo:   repo,
		owner:  owner,
	}, nil
}

func (g *Github) Stars(repository string) (int, error) {
	r, err := g.fetchRepo(repository)
	if err != nil {
		return 0, err
	}

	return r.GetStargazersCount(), nil
}

func (g *Github) Watchers(repository string) (int, error) {
	r, err := g.fetchRepo(repository)
	if err != nil {
		return 0, err
	}

	return r.GetSubscribersCount(), nil
}

func (g *Github) OpenIssues(repository string) (int, error) {
	r, err := g.fetchRepo(repository)
	if err != nil {
		return 0, err
	}

	return r.GetOpenIssuesCount(), nil
}

func (g *Github) ListRepo(limit int, order string) ([][]string, error) {
	headers := []string{"name", "stars", "watchers", "forks"}

	rs, err := g.fetchAllRepo(order)
	if err != nil {
		return nil, err
	}

	if limit > len(rs) {
		limit = len(rs)
	}

	repos := make([][]string, limit+1)
	repos[0] = headers
	for k, v := range rs {
		if k <= limit && k != 0 {
			repos[k] = append(repos[k], v.GetName())
			repos[k] = append(repos[k], strconv.FormatInt(int64(v.GetStargazersCount()), 10))
			repos[k] = append(repos[k], strconv.FormatInt(int64(v.GetSubscribersCount()), 10))
			repos[k] = append(repos[k], strconv.FormatInt(int64(v.GetForksCount()), 10))
		}
	}

	return repos, nil
}

// Fetch the whole repo per widget since we need to fetch the data during a regular time interval.
func (g *Github) fetchRepo(repository string) (*github.Repository, error) {
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	if repo == "" {
		return nil, errors.New("you need to specify a repository in the github service or in the widget")
	}

	r, _, err := g.client.Repositories.Get(context.Background(), g.owner, repo)
	if err != nil {
		return nil, errors.Wrapf(err, "can't find repo %s of owner %s", repository, g.owner)
	}

	return r, nil
}

func (g *Github) fetchAllRepo(order string) ([]*github.Repository, error) {
	ctx := context.Background()

	r, _, err := g.client.Repositories.List(ctx, g.owner, &github.RepositoryListOptions{Sort: order})
	if err != nil {
		return nil, errors.Wrapf(err, "can't find all repo of owner %s", g.owner)
	}

	return r, nil
}
