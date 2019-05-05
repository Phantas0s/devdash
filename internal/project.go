package internal

import (
	"strings"

	"github.com/pkg/errors"
)

var debug bool = false

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
	optionMustContain = "must_contain"
	optionFilters     = "filters"

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

func (p *project) Render(tui *Tui, d bool) (err error) {
	debug = d

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
					if p.gaWidget == nil {
						return errors.Errorf("can't use the widget %s without the service GoogleAnalytics - please fix your configuration file.", w.Name)
					}

					if err = p.gaWidget.CreateWidgets(w, tui); err != nil {
						return err
					}
				case "mon":
					if p.monitorWidget == nil {
						return errors.Errorf("can't use the widget %s without the service Monitor - please fix your configuration file.", w.Name)
					}

					if err = p.monitorWidget.CreateWidgets(w, tui); err != nil {
						return err
					}
				case "gsc":
					if p.gscWidget == nil {
						return errors.Errorf("can't use the widget %s without the service Google Search Console - please fix your configuration file.", w.Name)
					}

					if err = p.gscWidget.CreateWidgets(w, tui); err != nil {
						return err
					}
				case "github":
					if p.githubWidget == nil {
						return errors.Errorf("can't use the widget %s without the service Github - please fix your configuration file.", w.Name)
					}

					if err = p.githubWidget.CreateWidgets(w, tui); err != nil {
						return err
					}
				default:
					return errors.Errorf("could not find the service for widget %s - wrong name - please verify your configuration file", w.Name)
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
