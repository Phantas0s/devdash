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
	gscPages      = "gsc.pages"
	gscQueries    = "gsc.queries"
	gscTimeFormat = "2006-01-02"
)

type gscWidget struct {
	tui     *Tui
	client  *plateform.SearchConsole
	viewID  string
	address string
}

var mappingGscHeader = map[string]string{
	"page":  "Page",
	"query": "Query",
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
	case gscPages:
		err = s.pages(widget)
	case gscQueries:
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

// table of the result of a Google Search Console query.
// If no metric provided, default "query" with no filters.
func (s *gscWidget) table(widget Widget) (err error) {
	sd := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		sd = widget.Options[optionStartDate]
	}

	ed := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		ed = widget.Options[optionEndDate]
	}

	startDate, endDate, err := ConvertDates(time.Now(), sd, ed)
	if err != nil {
		return err
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

	filters := ""
	if _, ok := widget.Options[optionFilters]; ok {
		filters = widget.Options[optionFilters]
	}

	metrics := []string{"clicks", "impressions", "ctr", "position"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(widget.Options[optionMetrics], ",")
		}
	}

	results, err := s.client.Table(
		s.viewID,
		startDate.Format(gscTimeFormat),
		endDate.Format(gscTimeFormat),
		elLimit,
		s.address,
		metric,
		filters,
	)
	if err != nil {
		return err
	}

	table := make([][]string, len(results)+1)
	table[0] = []string{mappingGscHeader[metric]}
	table[0] = append(table[0], metrics...)

	for k, v := range results {
		table[k+1] = append(table[k+1], v.Dimension)
		for _, m := range metrics {
			if m == "clicks" {
				table[k+1] = append(table[k+1], fmt.Sprintf("%g", v.Clicks))
			}
			if m == "impressions" {
				table[k+1] = append(table[k+1], fmt.Sprintf("%g", v.Impressions))
			}
			if m == "ctr" {
				table[k+1] = append(table[k+1], fmt.Sprintf("%.5f", v.Ctr))
			}
			if m == "position" {
				table[k+1] = append(table[k+1], fmt.Sprintf("%.2f", v.Position))
			}
		}
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
