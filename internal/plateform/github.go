// To get fixtures in order to write tests
// j, _ := json.Marshal(is)
// fmt.Println(string(j))

package plateform

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/google/go-github/v27/github"
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
	is, err := g.fetchPullRequests(repository, limit)
	if err != nil {
		return nil, err
	}

	lpr := formatListPullRequests(is, limit)

	return lpr, nil
}

func formatListPullRequests(is []*github.PullRequest, limit int) [][]string {
	if limit > len(is) {
		limit = len(is)
	}

	headers := []string{"title", "state", "created at", "merged", "commits"}

	defaultHeader := "unknown"
	prs := make([][]string, limit+1)
	prs[0] = headers
	for k, v := range is {
		n := defaultHeader
		if v.Title != nil {
			n = *v.Title
		}

		state := defaultHeader
		if v.State != nil {
			state = *v.State
		}

		createdAt := defaultHeader
		if v.CreatedAt != nil {
			t := *v.CreatedAt
			createdAt = t.String()
		}

		merged := defaultHeader
		if v.Merged != nil {
			m := *v.Merged
			merged = strconv.FormatBool(m)
		}

		commits := defaultHeader
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

	return prs
}

func (g *Github) Views(repository string, days int) ([]string, []int, error) {
	tv, err := g.fetchViews(repository)
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

// TODO rename countCommits
func (g *Github) CommitCounts(repository string, startWeek int64, endWeek int64, startDate time.Time) ([]string, []int, error) {
	c, err := g.fetchCommitCount(repository)
	if err != nil {
		return nil, nil, err
	}

	d, co := formatCommitCounts(c, startWeek, endWeek, startDate)

	return d, co, nil
}

// TODO rename formatCountCommits
func formatCommitCounts(
	c *github.RepositoryParticipation,
	startWeek int64,
	endWeek int64,
	startDate time.Time,
) ([]string, []int) {
	// Reverse the count of commits (from ASC to DESC).
	for i := len(c.Owner)/2 - 1; i >= 0; i-- {
		opp := len(c.Owner) - 1 - i
		c.Owner[i], c.Owner[opp] = c.Owner[opp], c.Owner[i]
	}

	counts := []int{}
	for _, v := range c.Owner {
		counts = append(counts, v)
	}

	dimension := []string{}
	for k, _ := range c.Owner {
		// Since the startDate is the end of the week,
		// we need to come back to the first day of it
		// and then go back to the number of week (7 days)
		// specified in start date.
		weekDay := int(startDate.Weekday())
		s := startDate.AddDate(0, 0, (-(weekDay) - (7 * k)))
		if weekDay == 0 && k == 0 {
			s = startDate
		}
		dimension = append(dimension, s.Format("01-02"))
	}

	d := dimension[endWeek:startWeek]
	co := counts[endWeek:startWeek]

	// reverse the count of commits (from DESC to ASC)
	for i := len(co)/2 - 1; i >= 0; i-- {
		opp := len(co) - 1 - i
		co[i], co[opp] = co[opp], co[i]
	}

	// reverse the dimensions (from DESC to ASC)
	for i := len(d)/2 - 1; i >= 0; i-- {
		opp := len(d) - 1 - i
		d[i], d[opp] = d[opp], d[i]
	}
	return d, co
}

func (g *Github) CountStars(repository string) (dim []string, val []int, err error) {
	se, err := g.fetchStars(repository)
	if err != nil {
		return nil, nil, err
	}

	dim, val = formatCountStars(se, "01-02", true)
	return dim, val, nil
}

func formatCountStars(stargazers []*github.Stargazer, timeLayout string, addMissingDays bool) (dim []string, val []int) {
	sort.SliceStable(stargazers, func(i, j int) bool {
		return stargazers[i].StarredAt.Time.Before(stargazers[j].StarredAt.Time)
	})
	d, val := aggregateStarResults(stargazers)
	if addMissingDays {
		d, val = fillMissingDays(d, val)
	}

	for _, v := range d {
		dim = append(dim, v.Format(timeLayout))
	}

	return dim, val
}

// aggregateStarResults from Github. We look ahead of one element while the array is parsed, and add up
// if the date is the same as the current element of the array.
// If not, we add the date and value to the returning dimension and value slices.
func aggregateStarResults(stargazers []*github.Stargazer) (dim []time.Time, val []int) {
	var count int
	for k, v := range stargazers {

		count++

		// If last element of the slice processed
		if len(stargazers) == k+1 {
			dim = append(dim, v.StarredAt.Time)
			val = append(val, count)
			return
		}

		processingDate := v.StarredAt.Time.Format("2006-01-02")
		nextDate := stargazers[k+1].StarredAt.Time.Format("2006-01-02")

		if processingDate != nextDate {
			dim = append(dim, v.StarredAt.Time)
			val = append(val, count)
			count = 0
		}
	}

	return
}

// Fetch all events for a repo
func (g *Github) fetchStars(repository string) (s []*github.Stargazer, err error) {
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	if repo == "" {
		return nil, errors.New("you need to specify a repository in the github service or in the widget")
	}

	// TODO implement goroutines here
	count := 1
	for {
		e, _, err := g.client.Activity.ListStargazers(context.Background(), g.owner, repo, &github.ListOptions{
			Page:    count,
			PerPage: 100,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "can't find repo %s of owner %s", repo, g.owner)
		}

		s = append(s, e...)

		if len(e) < 100 {
			return s, nil
		}

		count++
	}
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

// TODO possibility to filter by ALL or OWNER
func (g *Github) fetchCommitCount(repository string) (*github.RepositoryParticipation, error) {
	repo := g.repo
	if repository != "" {
		repo = repository
	}

	if repo == "" {
		return nil, errors.New("you need to specify a repository in the github service or in the widget")
	}

	p, _, err := g.client.Repositories.ListParticipation(context.Background(), g.owner, repo)
	if err != nil {
		return nil, errors.Wrapf(err, "can't find repo %s of owner %s", repo, g.owner)
	}

	return p, nil
}

func (g *Github) fetchViews(repository string) (*github.TrafficViews, error) {
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

// Possibility to add options to filter quite a lot
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

// TODO add sorting
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

// TODO possibility to add filters / ordering
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

// fillMissingDays add the missing dates between two dates and add 0 values
func fillMissingDays(dates []time.Time, values []int) ([]time.Time, []int) {
	d := []time.Time{}
	val := []int{}
	for k, v := range dates {
		d = append(d, v)
		val = append(val, values[k])

		if len(dates) <= k+1 {
			return d, val
		}

		nextDate := dates[k+1]
		mg := missingDays(v, nextDate)
		for _, _ = range mg {
			val = append(val, 0)
		}
		d = append(d, mg...)
	}

	return d, val
}
