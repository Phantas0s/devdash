package internal

import (
	"strings"

	"github.com/pkg/errors"
)

var debug bool = false

const (
	// option config names
	optionTitle      = "title"
	optionTitleColor = "title_color"

	// time
	optionStartDate = "start_date"
	optionEndDate   = "end_date"
	optionGlobal    = "global"

	// For tables
	optionRowLimit  = "limit_row"
	optionCharLimit = "character_limit"

	// Metrics
	optionDimension  = "dimension"
	optionMetrics    = "metrics"
	optionMetric     = "metric"
	optionTimePeriod = "time_period"

	optionOrder = "order"

	// filtering
	optionMustContain = "must_contain"
	optionFilters     = "filters"
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
						return errors.Errorf("can't use the widget %s without the service Search - please fix your configuration file.", w.Name)
					}

					if err = p.gscWidget.CreateWidgets(w, tui); err != nil {
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
