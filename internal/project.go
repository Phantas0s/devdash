package internal

import (
	"strings"

	"github.com/pkg/errors"
)

type service interface {
	createWidgets(widget Widget, tui *Tui) (err error)
}

type Widget struct {
	Name string
	Size string
}

type project struct {
	name          string
	widgets       [][][]Widget
	sizes         [][]string
	gaWidget      service
	monitorWidget service
}

func NewProject(
	name string,
	widgets [][][]Widget,
	sizes [][]string,
) *project {
	return &project{
		name:    name,
		widgets: widgets,
		sizes:   sizes,
	}
}

func (p *project) WithGa(ga *gaWidget) {
	p.gaWidget = ga
}

func (p *project) WithMonitor(mon *monitorWidget) {
	p.monitorWidget = mon
}

func (p *project) Render(tui *Tui) (err error) {
	p.addTitle(tui)

	for r, row := range p.widgets {
		for c, col := range row {
			for _, w := range col {
				serviceID := strings.Split(w.Name, ".")[0]
				switch serviceID {
				case "ga":
					if p.gaWidget == nil {
						return errors.Errorf("can't use the widget %s without the service GoogleAnalytics - please fix your configuration file.", w.Name)
					}

					err = p.gaWidget.createWidgets(w, tui)
				case "mon":
					if p.monitorWidget == nil {
						return errors.Errorf("can't use the widget %s without the service Monitor - please fix your configuration file.", w.Name)
					}

					err = p.monitorWidget.createWidgets(w, tui)
				default:
					return errors.Errorf("could not find the service for widget %s - wrong name - please verify your configuration file", w.Name)
				}
			}
			if len(col) > 0 {
				tui.AddCol(p.sizes[r][c])
			}
		}
		tui.AddRow()
	}

	return
}

func (p *project) addTitle(tui *Tui) {
	tui.AddText(textAttr{
		Text: p.name,
		Fg:   5,
		Size: "XL",
	})
}
