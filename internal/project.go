package internal

import (
	"strings"

	"github.com/pkg/errors"
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

func (p *project) Render(tui *Tui) (err error) {
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
		if err := tui.AddRow(); err != nil {
			return err
		}
	}

	return
}

func (p *project) addTitle(tui *Tui) error {
	return tui.AddProjectTitle(p.name, p.titleOptions)
}
