package internal

import (
	"errors"
	"strings"
)

type Widget struct {
	Name string
	Size string
}

type project struct {
	name     string
	widgets  [][]Widget
	gaWidget gaWidget
}

func NewProject(
	name string,
	widgets [][]Widget,
	gaWidget gaWidget,
) project {
	return project{
		name:     name,
		widgets:  widgets,
		gaWidget: gaWidget,
	}
}

func (p project) Render(tui *Tui) (err error) {
	for i := 0; i < len(p.widgets); i++ {
		for _, w := range p.widgets[i] {
			// parse widgets for one row
			serviceID := strings.Split(w.Name, ".")[0]
			switch serviceID {
			case "ga":
				err = p.gaWidget.createWidgets(w, tui)
			default:
				return errors.New("could not find the service - please verify your configuration file")
			}
		}
		tui.AddRow()
	}

	return
}
