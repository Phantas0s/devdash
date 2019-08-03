package internal

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	// Widget config options

	optionTitle      = "title"
	optionTitleColor = "title_color"

	// Time
	optionStartDate  = "start_date"
	optionEndDate    = "end_date"
	optionTimePeriod = "time_period"
	optionGlobal     = "global"

	// Tables
	optionRowLimit  = "row_limit"
	optionCharLimit = "character_limit"

	// Metrics
	optionDimension  = "dimension"
	optionDimensions = "dimensions"

	optionMetrics = "metrics"
	optionMetric  = "metric"

	// Ordering
	optionOrder = "order"

	// Filtering
	optionFilters = "filters"

	// Repository
	optionRepository = "repository"
)

type service interface {
	CreateWidgets(widget Widget, tui *Tui) (err error)
}

type Widget struct {
	Name    string            `mapstructures:"name"`
	Size    string            `mapstructures:"size"`
	Options map[string]string `mapstructure:"options"`
}

type project struct {
	name          string
	titleOptions  map[string]string
	widgets       [][][]Widget
	sizes         [][]string
	gaWidget      service
	monitorWidget service
	gscWidget     service
	githubWidget  service
}

func NewProject(
	name string,
	titleOptions map[string]string,
	widgets [][][]Widget,
	sizes [][]string,
) *project {
	return &project{
		name:         name,
		titleOptions: titleOptions,
		widgets:      widgets,
		sizes:        sizes,
	}
}

func (p *project) WithGa(ga *gaWidget) {
	p.gaWidget = ga
}

func (p *project) WithMonitor(mon *monitorWidget) {
	p.monitorWidget = mon
}

func (p *project) WithGoogleSearchConsole(gsc *gscWidget) {
	p.gscWidget = gsc
}

func (p *project) WithGithub(github *githubWidget) {
	p.githubWidget = github
}

// TODO to test
func (p *project) Render(tui *Tui, debug bool) (err error) {
	err = p.addTitle(tui)
	if err != nil {
		return errors.Wrapf(err, "can't add project title %s", p.name)
	}

	for r, row := range p.widgets {
		for c, col := range row {
			for _, w := range col {
				serviceID := strings.Split(w.Name, ".")[0]
				switch serviceID {
				case "ga":
					createWidget(p.gaWidget, "Google Analytics", w, tui)
				case "mon":
					createWidget(p.monitorWidget, "Monitor", w, tui)
				case "gsc":
					createWidget(p.gscWidget, "Google Search Console", w, tui)
				case "github":
					createWidget(p.githubWidget, "Githug", w, tui)
				default:
					displayError(tui, errors.Errorf("The service %s doesn't exist (yet?)", w.Name))
				}
			}
			if len(col) > 0 {
				if err = tui.AddCol(p.sizes[r][c]); err != nil {
					return err
				}
			}
		}
		tui.AddRow()
		if !debug {
			tui.Render()
		}
	}

	return
}

func (p *project) addTitle(tui *Tui) error {
	return tui.AddProjectTitle(p.name, p.titleOptions)
}

func createWidget(s service, name string, w Widget, tui *Tui) {
	if s == nil {
		displayError(tui, errors.Errorf("Configuration error - you can't use the widget %s without the service %s.", w.Name, name))
	} else if err := s.CreateWidgets(w, tui); err != nil {
		displayError(tui, err)
	}
}
