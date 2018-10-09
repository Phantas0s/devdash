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
	widgets  [][][]Widget
	sizes    [][]string
	gaWidget gaWidget
}

func NewProject(
	name string,
	widgets [][][]Widget,
	sizes [][]string,
	gaWidget gaWidget,
) project {
	return project{
		name:     name,
		widgets:  widgets,
		sizes:    sizes,
		gaWidget: gaWidget,
	}
}

func (p project) Render(tui *Tui) (err error) {
	p.addTitle(tui)
	for r, row := range p.widgets {
		for c, col := range row {
			for _, w := range col {
				serviceID := strings.Split(w.Name, ".")[0]
				switch serviceID {
				case "ga":
					err = p.gaWidget.createWidgets(w, tui)
				default:
					return errors.New("could not find the service - please verify your configuration file")
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

func (p project) addTitle(tui *Tui) {
	tui.AddText(textAttr{
		Text: p.name,
		Fg:   2,
		Size: "XL",
	})
}
