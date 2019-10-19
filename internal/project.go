package internal

import (
	"github.com/pkg/errors"
)

type service interface {
	CreateWidgets(widget Widget, tui *Tui) (err error)
}

type project struct {
	name           string
	nameOptions    map[string]string
	widgets        [][][]Widget
	sizes          [][]string
	themes         map[string]map[string]string
	gaWidget       service
	monitorWidget  service
	gscWidget      service
	githubWidget   service
	travisCIWidget service
	feedlyWidget   service
	tui            *Tui
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

// Service builder
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
func (p *project) WithTravisCI(travisCI *travisCIWidget) {
	p.travisCIWidget = travisCI
}

func (p *project) WithFeedly(feedly *feedlyWidget) {
	p.feedlyWidget = feedly
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
	// TODO: use display.box instead of this shortcut
	err := p.addTitle(p.tui)
	if err != nil {
		err = errors.Wrapf(err, "can't add project title %s", p.name)
		DisplayError(p.tui, err)
	}

	for r, row := range p.widgets {
		for c, col := range row {
			for _, w := range col {
				w = p.addDefaultTheme(w)

				// Map widget prefix with serice
				switch w.serviceID() {
				case "ga":
					createWidgets(p.gaWidget, "Google Analytics", w, p.tui)
				case "mon":
					createWidgets(p.monitorWidget, "Monitor", w, p.tui)
				case "gsc":
					createWidgets(p.gscWidget, "Google Search Console", w, p.tui)
				case "github":
					createWidgets(p.githubWidget, "Github", w, p.tui)
				case "travis":
					createWidgets(p.travisCIWidget, "Travis", w, p.tui)
				case "feedly":
					createWidgets(p.feedlyWidget, "Feedly", w, p.tui)
				case "display":
					createWidgets(NewDisplayWidget(), "Display", w, p.tui)
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

func createWidgets(s service, name string, w Widget, tui *Tui) {
	if s == nil {
		DisplayError(tui, errors.Errorf("Configuration error - you can't use the widget %s without the service %s.", w.Name, name))
	} else if err := s.CreateWidgets(w, tui); err != nil {
		DisplayError(tui, err)
	}
}
