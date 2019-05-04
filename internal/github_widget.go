package internal

import (
	"strconv"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/pkg/errors"
)

const (
	githubBoxStars = "github.box_stars"
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
		err = g.getStars(widget)
	default:
		return errors.Errorf("can't find the widget %s for service github", widget.Name)
	}

	return
}

func (g *githubWidget) getStars(widget Widget) error {
	title := " Github Stars "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	stars := g.client.GetStars()

	var s string
	if stars == nil {
		s = "0"
	} else {
		s = strconv.FormatInt(int64(*stars), 10)
	}

	g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)

	return nil
}
