package internal

import (
	"strconv"
	"strings"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/Phantas0s/devdash/totime"
	"github.com/pkg/errors"
)

const (
	gsc_pages   = "gsc.pages"
	gsc_queries = "gsc.queries"
)

type gscWidget struct {
	tui     *Tui
	client  *plateform.SearchConsole
	viewID  string
	address string
}

func NewGscWidget(keyfile string, viewID string, address string) (*gscWidget, error) {
	sc, err := plateform.NewSearchConsoleClient(keyfile)
	if err != nil {
		return nil, err
	}

	return &gscWidget{
		client:  sc,
		viewID:  viewID,
		address: address,
	}, nil
}

func (s *gscWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	s.tui = tui
	switch widget.Name {
	case gsc_pages:
		err = s.pages(widget)
	case gsc_queries:
		err = s.table(widget)
	default:
		return errors.Errorf("can't find the widget %s", widget.Name)
	}

	return
}

func (s *gscWidget) pages(widget Widget) error {
	if widget.Options == nil {
		widget.Options = map[string]string{optionMetric: "page"}
	} else {
		widget.Options[optionMetric] = "page"
	}

	return s.table(widget)
}

func (s *gscWidget) table(widget Widget) (err error) {
	startDate, endDate := totime.NPrevMonth(1)
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	var elLimit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		elLimit, err = strconv.ParseInt(widget.Options[optionRowLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	title := "Search Console"
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	charLimit := 1000
	if _, ok := widget.Options[optionCharLimit]; ok {
		c, err := strconv.ParseInt(widget.Options[optionCharLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionCharLimit])
		}
		charLimit = int(c)
	}

	metric := "query"
	if _, ok := widget.Options[optionMetric]; ok {
		metric = widget.Options[optionMetric]
	}

	table, err := s.client.Pages(
		s.viewID,
		startDate,
		endDate,
		elLimit,
		s.address,
		metric,
	)
	if err != nil {
		return err
	}

	// Shorten the URL of the page.
	// Begins the loop to 1 not to shorten the headers.
	for i := 1; i < len(table); i++ {
		URL := strings.TrimPrefix(table[i][0], s.address)

		if charLimit < len(URL) {
			URL = URL[:charLimit]
		}

		table[i][0] = URL
	}

	s.tui.AddTable(table, title, widget.Options)

	return nil

}
