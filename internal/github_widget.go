package internal

import (
	"fmt"
	"strconv"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/pkg/errors"
)

const (
	githubBoxStars          = "github.box_stars"
	githubBoxWatchers       = "github.box_watchers"
	githubBoxOpenIssues     = "github.box_open_issues"
	githubTableRepositories = "github.table_repositories"
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

	order := "pushed"
	if _, ok := widget.Options[optionOrder]; ok {
		order = widget.Options[optionRowLimit]
	}

	rs, err := g.client.ListRepo(int(limit), order)
	if err != nil {
		return err
	}

	g.tui.AddTable(rs, title, widget.Options)

	return nil
}
