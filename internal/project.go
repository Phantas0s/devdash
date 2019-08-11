package internal

import (
	"github.com/pkg/errors"
)

type service interface {
	CreateWidgets(widget Widget, tui *Tui) (err error)
}

type project struct {
	name          string
	titleOptions  map[string]string
	widgets       [][][]Widget
	sizes         [][]string
	themes        map[string]map[string]string
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
	themes map[string]map[string]string,
) *project {
	return &project{
		name:         name,
		titleOptions: titleOptions,
		widgets:      widgets,
		sizes:        sizes,
		themes:       themes,
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

func (p *project) addDefaultTheme(w Widget) Widget {
	t := w.typeID()

	theme := map[string]string{}
	if _, ok := p.themes[t]; ok {
		theme = p.themes[t]
	}

	if w.Theme != "" {
		if _, ok := p.themes[w.Theme]; ok {
			theme = p.themes[w.Theme]
		}
	}

	if len(theme) > 0 {
		for k, v := range theme {
			if len(w.Options) == 0 {
				w.Options = map[string]string{}
			}
			if _, ok := w.Options[k]; !ok {
				w.Options[k] = v
			}
		}
	}

	return w
}

func (p *project) Render(tui *Tui, debug bool) {
	err := p.addTitle(tui)
	if err != nil {
		err = errors.Wrapf(err, "can't add project title %s", p.name)
		DisplayError(tui, err)
	}

	for r, row := range p.widgets {
		for c, col := range row {
			for _, w := range col {
				w = p.addDefaultTheme(w)

				switch w.serviceID() {
				case "ga":
					displayWidget(p.gaWidget, "Google Analytics", w, tui)
				case "mon":
					displayWidget(p.monitorWidget, "Monitor", w, tui)
				case "gsc":
					displayWidget(p.gscWidget, "Google Search Console", w, tui)
				case "github":
					displayWidget(p.githubWidget, "Github", w, tui)
				default:
					DisplayError(tui, errors.Errorf("The service %s doesn't exist (yet?)", w.Name))
				}
			}
			if len(col) > 0 {
				if err = tui.AddCol(p.sizes[r][c]); err != nil {
					DisplayError(tui, err)
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

func displayWidget(s service, name string, w Widget, tui *Tui) {
	if s == nil {
		DisplayError(tui, errors.Errorf("Configuration error - you can't use the widget %s without the service %s.", w.Name, name))
	} else if err := s.CreateWidgets(w, tui); err != nil {
		DisplayError(tui, err)
	}
}
