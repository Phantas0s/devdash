package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/pkg/errors"
)

const (
	githubBoxStars          = "github.box_stars"
	githubBoxWatchers       = "github.box_watchers"
	githubBoxOpenIssues     = "github.box_open_issues"
	githubTableRepositories = "github.table_repositories"
	githubTableBranches     = "github.table_branches"
	githubTableIssues       = "github.table_issues"
	githubTablePullRequests = "github.table_pull_requests"
	githubBarViews          = "github.bar_views"
	githubBarCommits        = "github.bar_commits"
	githubBarStars          = "github.bar_stars"
)

type githubWidget struct {
	tui    *Tui
	client *platform.Github
}

// NewGithubWidget with all information necessary to connect to the Github API.
func NewGithubWidget(token string, owner string, repo string) (*githubWidget, error) {
	g, err := platform.NewGithubClient(token, owner, repo)
	if err != nil {
		return nil, err
	}
	return &githubWidget{
		client: g,
	}, nil
}

// CreateWidgets for the Github service.
func (g *githubWidget) CreateWidgets(widget Widget, tui *Tui) (f func() error, err error) {
	g.tui = tui

	switch widget.Name {
	case githubBoxStars:
		f, err = g.boxStars(widget)
	case githubBoxWatchers:
		f, err = g.boxWatchers(widget)
	case githubBoxOpenIssues:
		f, err = g.boxOpenIssues(widget)
	case githubTableRepositories:
		f, err = g.tableRepo(widget)
	case githubTableBranches:
		f, err = g.tableBranches(widget)
	case githubTableIssues:
		f, err = g.tableIssues(widget)
	case githubTablePullRequests:
		f, err = g.tablePullRequests(widget)
	case githubBarViews:
		f, err = g.barViews(widget)
	case githubBarCommits:
		f, err = g.barCommits(widget)
	case githubBarStars:
		f, err = g.barStars(widget)
	default:
		return nil, errors.Errorf("can't find the widget %s for service github", widget.Name)
	}

	return
}

func (g *githubWidget) boxStars(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := fmt.Sprintf(" Github Stars for %s", repo)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	stars, err := g.client.TotalStars(repo)
	if err != nil {
		return nil, err
	}

	s := strconv.FormatInt(int64(stars), 10)

	f = func() error {
		return g.tui.AddTextBox(
			s,
			title,
			widget.Options,
		)
	}

	return
}

func (g *githubWidget) boxWatchers(widget Widget) (f func() error, err error) {
	title := " Github Watchers "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	w, err := g.client.TotalWatchers(repo)
	if err != nil {
		return nil, err
	}

	s := strconv.FormatInt(int64(w), 10)

	f = func() error {
		return g.tui.AddTextBox(
			s,
			title,
			widget.Options,
		)
	}

	return
}

func (g *githubWidget) boxOpenIssues(widget Widget) (f func() error, err error) {
	title := " Github Open Issues "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	w, err := g.client.TotalOpenIssues(repo)
	if err != nil {
		return nil, err
	}

	s := strconv.FormatInt(int64(w), 10)

	f = func() error {
		return g.tui.AddTextBox(
			s,
			title,
			widget.Options,
		)
	}

	return
}

func (g *githubWidget) tableRepo(widget Widget) (f func() error, err error) {
	title := " Github Repositories "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return nil, errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	metrics := []string{"name", "stars", "watchers", "forks", "open_issues"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(strings.TrimSpace(widget.Options[optionMetrics]), ",")
		}
	}

	order := "pushed"
	if _, ok := widget.Options[optionOrder]; ok {
		order = widget.Options[optionOrder]
	}

	rs, err := g.client.ListRepo(int(limit), order, metrics)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddTable(rs, title, widget.Options)
	}

	return
}

func (g *githubWidget) tableBranches(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Branches "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return nil, errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	bs, err := g.client.ListBranches(repo, int(limit))
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddTable(bs, title, widget.Options)
	}

	return
}

// TODO can filter by open or close issue?
func (g *githubWidget) tableIssues(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Issues "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return nil, errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	is, err := g.client.ListIssues(repo, int(limit))
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddTable(is, title, widget.Options)
	}

	return
}

func (g *githubWidget) tablePullRequests(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Pull Requests "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return nil, errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	is, err := g.client.ListPullRequests(repo, int(limit))
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddTable(is, title, widget.Options)
	}

	return
}

func (g *githubWidget) barViews(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Views "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	dim, counts, err := g.client.Views(repo, 0)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddBarChart(counts, dim, title, widget.Options)
	}

	return
}

// TODO to refactor - transforming any date statement (weeks_ago, month_ago) into days weeks_ago in platform.date, and plugt it in.
func (g *githubWidget) barCommits(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Commit Per Week "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	sd := "7_weeks_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		sd = widget.Options[optionStartDate]
	}

	ed := "0_weeks_ago"
	if _, ok := widget.Options[optionEndDate]; ok {
		ed = widget.Options[optionEndDate]
	}

	scope := ownerScope
	if _, ok := widget.Options[optionScope]; ok {
		scope = widget.Options[optionScope]
	}

	if !strings.Contains(sd, "weeks_ago") || !strings.Contains(ed, "weeks_ago") {
		return nil, errors.New("The widget github.bar_commits require you to indicate a week range, ie startDate: 5_weeks_ago, endDate: 1_weeks_ago ")
	}

	sw, err := platform.ExtractCountPeriod(sd)
	if err != nil {
		return nil, err
	}

	ew, err := platform.ExtractCountPeriod(ed)
	if err != nil {
		return nil, err
	}

	dim, counts, err := g.client.CountCommits(repo, scope, sw, ew, time.Now())
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddBarChart(counts, dim, title, widget.Options)
	}

	return
}

func (g *githubWidget) barStars(widget Widget) (f func() error, err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Stars "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	startDate := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	endDate := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	sd, ed, err := platform.ConvertDates(time.Now(), startDate, endDate)
	if err != nil {
		return nil, err
	}

	dim, counts, err := g.client.CountStars(repo, sd, ed)
	if err != nil {
		return nil, err
	}

	f = func() error {
		return g.tui.AddBarChart(counts, dim, title, widget.Options)
	}

	return
}
