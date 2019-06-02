package plateform

import (
	"context"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	githubRepoStars      = "stars"
	githubRepoWatchers   = "watchers"
	githubRepoForks      = "forks"
	githubRepoOpenIssues = "open_issues"
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

func (g *Github) ListBranches(repository string, limit int) ([][]string, error) {
	headers := []string{"name"}

	bs, err := g.fetchBranches(repository, limit)
	if err != nil {
		return nil, err
	}

	if limit > len(bs) {
		limit = len(bs)
	}

	branches := make([][]string, limit+1)
	branches[0] = headers
	for k, v := range bs {
		n := "unknown"
		if v.Name != nil {
			n = *v.Name
		}

		if k < limit {
			branches[k+1] = append(branches[k+1], n)
		}
	}

	return branches, nil
}

func (g *Github) ListRepo(limit int, order string, metrics []string) ([][]string, error) {
	headers := []string{"name"}

	rs, err := g.fetchAllRepo(order)
	if err != nil {
		return nil, err
	}

	if limit > len(rs) {
		limit = len(rs)
	}

	repos := make([][]string, limit+1)

	stars := false
	watchers := false
	forks := false
	openIssues := false

	if contains(metrics, githubRepoStars) {
		headers = append(headers, githubRepoStars)
		stars = true
	}
	if contains(metrics, githubRepoWatchers) {
		headers = append(headers, githubRepoWatchers)
		watchers = true
	}
	if contains(metrics, githubRepoForks) {
		headers = append(headers, githubRepoForks)
		forks = true
	}
	if contains(metrics, githubRepoOpenIssues) {
		headers = append(headers, githubRepoOpenIssues)
		openIssues = true
	}

	repos[0] = headers

	for k, v := range rs {
		if k < limit {
			repos[k+1] = append(repos[k+1], v.GetName())
			if stars {
				repos[k+1] = append(repos[k+1], strconv.FormatInt(int64(v.GetStargazersCount()), 10))
			}
			if watchers {
				repos[k+1] = append(repos[k+1], strconv.FormatInt(int64(v.GetSubscribersCount()), 10))
			}
			if forks {
				repos[k+1] = append(repos[k+1], strconv.FormatInt(int64(v.GetForksCount()), 10))
			}
			if openIssues {
				repos[k+1] = append(repos[k+1], strconv.FormatInt(int64(v.GetOpenIssuesCount()), 10))
			}
		}
	}

	return repos, nil
}

func (g *Github) ListIssues(repository string, limit int) ([][]string, error) {
	headers := []string{"name", "state"}

	is, err := g.fetchIssues(repository, limit)
	if err != nil {
		return nil, err
	}

	if limit > len(is) {
		limit = len(is)
	}

	issues := make([][]string, limit+1)
	issues[0] = headers
	for k, v := range is {
		n := "unknown"
		if v.Title != nil {
			n = *v.Title
		}

		state := "unknown"
		if v.State != nil {
			state = *v.State
		}

		if k < limit {
			issues[k+1] = append(issues[k+1], n)
			issues[k+1] = append(issues[k+1], state)
		}
	}

	return issues, nil
}

func (g *Github) ListPullRequests(repository string, limit int) ([][]string, error) {
	headers := []string{"title", "state", "created at", "merged", "commits"}

	is, err := g.fetchPullRequests(repository, limit)
	if err != nil {
		return nil, err
	}

	if limit > len(is) {
		limit = len(is)
	}

	prs := make([][]string, limit+1)
	prs[0] = headers
	for k, v := range is {
		n := "unknown"
		if v.Title != nil {
			n = *v.Title
		}

		state := "unknown"
		if v.State != nil {
			state = *v.State
		}

		createdAt := "unknown"
		if v.CreatedAt != nil {
			t := *v.CreatedAt
			createdAt = t.String()
		}

		merged := "unknown"
		if v.Merged != nil {
			m := *v.Merged
			merged = strconv.FormatBool(m)
		}

		commits := "unknown"
		if v.Commits != nil {
			c := *v.Commits
			commits = strconv.FormatInt(int64(c), 10)
		}

		if k < limit {
			prs[k+1] = append(prs[k+1], n)
			prs[k+1] = append(prs[k+1], state)
			prs[k+1] = append(prs[k+1], createdAt)
			prs[k+1] = append(prs[k+1], merged)
			prs[k+1] = append(prs[k+1], commits)
		}
	}

	return prs, nil
}

func (g *Github) TrafficView(repository string, days int) ([]string, []int, error) {
	tv, err := g.fetchTrafficViews(repository)
	if err != nil {
		return nil, nil, err
	}

	counts := []int{}
	for _, v := range tv.Views {
		counts = append(counts, v.GetCount())
	}

	dimension := []string{}
	for _, v := range tv.Views {
		dimension = append(dimension, v.GetTimestamp().Format("01-02"))
	}

	return dimension, counts, nil
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
		return nil, errors.Wrapf(err, "can't find repo %s of owner %s", repo, g.owner)
	}

	return r, nil

}

func (g *Github) fetchTrafficViews(repository string) (*github.TrafficViews, error) {
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	if repo == "" {
		return nil, errors.New("you need to specify a repository in the github service or in the widget")
	}

	t, _, err := g.client.Repositories.ListTrafficViews(context.Background(), g.owner, repo, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't find repo %s of owner %s", repo, g.owner)
	}

	return t, nil
}

// Fetch the branches of a repo
func (g *Github) fetchBranches(repository string, limit int) ([]*github.Branch, error) {
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	if repo == "" {
		return nil, errors.New("you need to specify a repository in the github service or in the widget")
	}

	opt := github.ListOptions{PerPage: limit}
	bs, _, err := g.client.Repositories.ListBranches(context.Background(), g.owner, repo, &opt)
	if err != nil {
		return nil, errors.Wrapf(err, "can't find branches of owner %s for repo %s", g.owner, repo)
	}

	return bs, nil
}

func (g *Github) fetchIssues(repository string, limit int) ([]*github.Issue, error) {
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	if repo == "" {
		return nil, errors.New("you need to specify a repository in the github service or in the widget")
	}

	opt := github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: limit},
	}
	is, _, err := g.client.Issues.ListByRepo(context.Background(), g.owner, repo, &opt)
	if err != nil {
		return nil, errors.Wrapf(err, "can't find branches of owner %s for repo %s", g.owner, repo)
	}

	return is, nil
}

func (g *Github) fetchPullRequests(repository string, limit int) ([]*github.PullRequest, error) {
	ctx := context.Background()
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	prs, _, err := g.client.PullRequests.List(ctx, g.owner, repo, &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: limit},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "can't find all repo of owner %s", g.owner)
	}

	return prs, nil
}

func (g *Github) fetchAllRepo(order string) ([]*github.Repository, error) {
	ctx := context.Background()

	r, _, err := g.client.Repositories.List(ctx, g.owner, &github.RepositoryListOptions{Sort: order})
	if err != nil {
		return nil, errors.Wrapf(err, "can't find all repo of owner %s", g.owner)
	}

	return r, nil
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
