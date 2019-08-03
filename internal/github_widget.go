package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/internal/plateform"
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
	client *plateform.Github
}

func NewGithubWidget(token string, owner string, repo string) (*githubWidget, error) {
	g, err := plateform.NewGithubClient(token, owner, repo)
	if err != nil {
		return nil, err
	}
	return &githubWidget{
		client: g,
	}, nil
}

func (g *githubWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	g.tui = tui

	switch widget.Name {
	case githubBoxStars:
		err = g.boxStars(widget)
	case githubBoxWatchers:
		err = g.boxWatchers(widget)
	case githubBoxOpenIssues:
		err = g.boxOpenIssues(widget)
	case githubTableRepositories:
		err = g.tableRepo(widget)
	case githubTableBranches:
		err = g.tableBranches(widget)
	case githubTableIssues:
		err = g.tableIssues(widget)
	case githubTablePullRequests:
		err = g.tablePullRequests(widget)
	case githubBarViews:
		err = g.barViews(widget)
	case githubBarCommits:
		err = g.barCommits(widget)
	case githubBarStars:
		err = g.barStars(widget)
	default:
		return errors.Errorf("can't find the widget %s for service github", widget.Name)
	}

	return
}

func (g *githubWidget) boxStars(widget Widget) error {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := fmt.Sprintf(" Github Stars for %s", repo)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	stars, err := g.client.Stars(repo)
	if err != nil {
		return err
	}

	s := strconv.FormatInt(int64(stars), 10)

	g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)

	return nil
}

func (g *githubWidget) boxWatchers(widget Widget) error {
	title := " Github Watchers "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	w, err := g.client.Watchers(repo)
	if err != nil {
		return err
	}

	s := strconv.FormatInt(int64(w), 10)

	g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)

	return nil
}

func (g *githubWidget) boxOpenIssues(widget Widget) error {
	title := " Github Open Issues "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	w, err := g.client.OpenIssues(repo)
	if err != nil {
		return err
	}

	s := strconv.FormatInt(int64(w), 10)

	g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)

	return nil
}

func (g *githubWidget) tableRepo(widget Widget) (err error) {
	title := " Github Repositories "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
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
		order = widget.Options[optionRowLimit]
	}

	rs, err := g.client.ListRepo(int(limit), order, metrics)
	if err != nil {
		return err
	}

	g.tui.AddTable(rs, title, widget.Options)

	return nil
}

func (g *githubWidget) tableBranches(widget Widget) (err error) {
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
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	bs, err := g.client.ListBranches(repo, int(limit))
	if err != nil {
		return err
	}

	g.tui.AddTable(bs, title, widget.Options)

	return nil
}

func (g *githubWidget) tableIssues(widget Widget) (err error) {
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
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	is, err := g.client.ListIssues(repo, int(limit))
	if err != nil {
		return err
	}

	g.tui.AddTable(is, title, widget.Options)

	return nil
}

func (g *githubWidget) tablePullRequests(widget Widget) (err error) {
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
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	is, err := g.client.ListPullRequests(repo, int(limit))
	if err != nil {
		return err
	}

	g.tui.AddTable(is, title, widget.Options)

	return nil
}

func (g *githubWidget) barViews(widget Widget) (err error) {
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
		return err
	}

	g.tui.AddBarChart(counts, dim, title, widget.Options)

	return nil
}

// TODO to refactor
func (g *githubWidget) barCommits(widget Widget) (err error) {
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

	if !strings.Contains(sd, "weeks_ago") {
		return errors.New("the widget bar_commits require you to indicate a week range, ie startDate: 5_weeks_ago, endDate: 0_weeks_ago ")
	}

	t := strings.Split(sd, "_")
	sw, err := strconv.ParseInt(t[0], 0, 0)
	if err != nil {
		return errors.Wrapf(err, "%s is not a valid date", sd)
	}

	if !strings.Contains(ed, "weeks_ago") {
		return errors.New("the widget bar_commits require you to indicate a week range, ie startDate: 5_weeks_ago, endDate: 0_weeks_ago ")
	}

	t = strings.Split(ed, "_")
	ew, err := strconv.ParseInt(t[0], 0, 0)
	if err != nil {
		return errors.Wrapf(err, "%s is not a valid date", ed)
	}

	dim, counts, err := g.client.CommitCounts(repo, sw, ew, time.Now())
	if err != nil {
		return err
	}

	g.tui.AddBarChart(counts, dim, title, widget.Options)

	return nil
}

func (g *githubWidget) barStars(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Stars "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	dim, counts, err := g.client.CountStars(repo)
	if err != nil {
		return err
	}

	g.tui.AddBarChart(counts, dim, title, widget.Options)

	return nil
}
