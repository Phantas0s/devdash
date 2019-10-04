package internal

import (
	"github.com/pkg/errors"
)

type service interface {
	CreateWidgets(widget Widget, tui *Tui) (err error)
}

type project struct {
	name          string
	nameOptions   map[string]string
	widgets       [][][]Widget
	sizes         [][]string
	themes        map[string]map[string]string
	gaWidget      service
	monitorWidget service
	gscWidget     service
	githubWidget  service
	tui           *Tui
}

// NewProject for the dashboard.
func NewProject(
	name string,
	nameOptions map[string]string,
	widgets [][][]Widget,
	sizes [][]string,
	themes map[string]map[string]string,
	tui *Tui,
) *project {
	return &project{
		name:        name,
		nameOptions: nameOptions,
		widgets:     widgets,
		sizes:       sizes,
		themes:      themes,
		tui:         tui,
	}
}

// WithGa add Google Analytics service to the project.
func (p *project) WithGa(ga *gaWidget) {
	p.gaWidget = ga
}

// WithMonitor add the Monitor service to the project.
func (p *project) WithMonitor(mon *monitorWidget) {
	p.monitorWidget = mon
}

// WithGoogleSearchConsole add the Google Search Console service to the project
func (p *project) WithGoogleSearchConsole(gsc *gscWidget) {
	p.gscWidget = gsc
}

// WithGithub add the Github service to the project
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

// Render all the services' widgets.
func (p *project) Render(debug bool) {
	err := p.addTitle(p.tui)
	if err != nil {
		err = errors.Wrapf(err, "can't add project title %s", p.name)
		DisplayError(p.tui, err)
	}

	for r, row := range p.widgets {
		for c, col := range row {
			for _, w := range col {
				w = p.addDefaultTheme(w)

				switch w.serviceID() {
				case "ga":
					displayWidget(p.gaWidget, "Google Analytics", w, p.tui)
				case "mon":
					displayWidget(p.monitorWidget, "Monitor", w, p.tui)
				case "gsc":
					displayWidget(p.gscWidget, "Google Search Console", w, p.tui)
				case "github":
					displayWidget(p.githubWidget, "Github", w, p.tui)
				default:
					DisplayError(p.tui, errors.Errorf("The service %s doesn't exist (yet?)", w.Name))
				}
			}
			if len(col) > 0 {
				if err = p.tui.AddCol(p.sizes[r][c]); err != nil {
					DisplayError(p.tui, err)
				}
			}
		}
		p.tui.AddRow()
		if !debug {
			p.tui.Render()
		}
	}

	return
}

func (p *project) addTitle(tui *Tui) error {
	return tui.AddProjectTitle(p.name, p.nameOptions)
}

func displayWidget(s service, name string, w Widget, tui *Tui) {
	if s == nil {
		DisplayError(tui, errors.Errorf("Configuration error - you can't use the widget %s without the service %s.", w.Name, name))
	} else if err := s.CreateWidgets(w, tui); err != nil {
		DisplayError(tui, err)
	}
}
