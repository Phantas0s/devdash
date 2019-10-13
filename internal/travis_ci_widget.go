package internal

import (
	"strconv"

	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/pkg/errors"
)

const (
	travisCITableBuilds = "travis.table_builds"
)

type travisCIWidget struct {
	tui    *Tui
	client *platform.TravisCI
}

// NewTravisCIWidget with all information necessary to connect to the Github API.
func NewTravisCIWidget(token string) *travisCIWidget {
	c := platform.NewTravisCI(token)
	return &travisCIWidget{
		client: c,
	}
}

func (tc travisCIWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	tc.tui = tui

	switch widget.Name {
	case travisCITableBuilds:
		err = tc.tableBuilds(widget)
	default:
		return errors.Errorf("can't find the widget %s for service travis ci", widget.Name)
	}

	return
}

func (tc travisCIWidget) tableBuilds(widget Widget) (err error) {
	title := " Tracis CI builds "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	repo := ""
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	owner := ""
	if _, ok := widget.Options[optionOwner]; ok {
		owner = widget.Options[optionOwner]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	builds, err := tc.client.Builds(repo, owner, limit)
	if err != nil {
		return err
	}

	tc.tui.AddTable(builds, title, widget.Options)

	return nil
}
