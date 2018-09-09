package internal

import (
	"errors"
	"strings"
)

type Widget struct {
	Size string
}

type project struct {
	name     string
	widgets  [][]map[string]Widget
	gaWidget gaWidget
}

func NewProject(
	name string,
	widgets [][]map[string]Widget,
	gaWidget gaWidget,
) project {
	return project{
		name:     name,
		widgets:  widgets,
		gaWidget: gaWidget,
	}
}

func (p project) Render(tui *Tui) error {
	for i := 0; i < len(p.widgets); i++ {
		for _, ws := range p.widgets[i] {
			// parse widgets for one row
			for wn, w := range ws {
				serviceID := strings.Split(wn, ".")[0]
				switch serviceID {
				case "ga":
					p.gaWidget.createWidgets(wn, w, tui)
				default:
					return errors.New("could not find the service - please verify your configuration file")
				}
			}
		}
		tui.AddRow()
	}

	return nil
}
