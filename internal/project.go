package internal

// Dispatch the different widgets from the config to their widget managers and get back the render functions.

// TODO it's a mess in there between concurrent / non concurrent ways of getting widget render functions.
// TODO Initially it was made to go around the limitation of concurrent connection for Google Analytics.
// TODO I feel the absence of generics here...  To refactor somehow (using reflection?).

import (
	"github.com/pkg/errors"
)

type service interface {
	CreateWidgets(widget Widget, tui *Tui) (f func() error, err error)
}

type project struct {
	name        string
	nameOptions map[string]string
	widgets     [][][]Widget
	sizes       [][]string
	themes      map[string]map[string]string
	tui         *Tui

	gaWidget         service
	monitorWidget    service
	gscWidget        service
	githubWidget     service
	travisCIWidget   service
	feedlyWidget     service
	gitWidget        service
	remoteHostWidget service
	localhostWidget  service
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

func (p *project) WithGit(git *gitWidget) {
	p.gitWidget = git
}

func (p *project) WithRemoteHost(remoteHost *HostWidget) {
	p.remoteHostWidget = remoteHost
}

func (p *project) WithLocalhost(localhost *HostWidget) {
	p.localhostWidget = localhost
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

func (p *project) mapServiceID(serviceID string) (service, error) {
	services := map[string]service{
		"display": NewDisplayWidget(),
		"ga":      p.gaWidget,
		"mon":     p.monitorWidget,
		"gsc":     p.gscWidget,
		"github":  p.githubWidget,
		"travis":  p.travisCIWidget,
		"feedly":  p.feedlyWidget,
		"git":     p.gitWidget,
		"rh":      p.remoteHostWidget,
		"lh":      p.localhostWidget,
	}

	if _, ok := services[serviceID]; ok {
		return services[serviceID], nil
	}

	return nil, errors.Errorf("Impossible to find the service with ID %s", serviceID)
}

func mapServiceName(serviceID string) (string, error) {
	services := map[string]string{
		"display": "Display",
		"ga":      "Google Analytics",
		"mon":     "Monitor",
		"gsc":     "Google Search Console",
		"github":  "Github",
		"travis":  "Travis",
		"feedly":  "Feedly",
		"git":     "Git",
		"rh":      "Remote Host",
		"lh":      "Localhost",
	}

	if _, ok := services[serviceID]; ok {
		return services[serviceID], nil
	}

	return "", errors.Errorf("Impossible to find the service with ID %s", serviceID)
}

// Create all the widgets and populate them with data.
// Return channels with render functions
func (p *project) CreateWidgets() [][][]func() error {
	// TODO: use display.box instead of this shortcut
	err := p.addTitle(p.tui)
	if err != nil {
		err = errors.Wrapf(err, "can't add project title %s", p.name)
		DisplayError(p.tui, err)()
	}

	chs := make([][][]chan func() error, len(p.widgets))

	for ir, row := range p.widgets {
		for ic, col := range row {
			chs[ir] = append(chs[ir], []chan func() error{})
			for _, w := range col {
				w = p.addDefaultTheme(w)
				ch := make(chan func() error)
				chs[ir][ic] = append(chs[ir][ic], ch)

				service, err := p.mapServiceID(w.serviceID())
				if err != nil {
					go func(c chan<- func() error) {
						c <- DisplayError(p.tui, err)
					}(ch)
					continue
				}

				serviceName, err := mapServiceName(w.serviceID())
				if err != nil {
					go func(c chan<- func() error) {
						c <- DisplayError(p.tui, err)
					}(ch)
					continue
				}

				go getChannelRenderers(service, serviceName, w, p.tui, ch)
			}
		}
	}

	funcs := make([][][]func() error, len(chs))
	for r, _ := range chs {
		for c, _ := range chs[r] {
			funcs[r] = append(funcs[r], []func() error{})
			for _, chann := range chs[r][c] {
				f, ok := <-chann
				if ok {
					funcs[r][c] = append(funcs[r][c], f)
				}
			}
		}
	}

	return funcs
}

// getConcurentRenderers and fetch information via different ways depending on Widget (API / SSH / ...)
// A function to display the widget will be send to a channel.
// One channel per widget to keep the widget order in a slice.
func getChannelRenderers(s service, name string, w Widget, tui *Tui, c chan<- func() error) {
	if s == nil {
		c <- DisplayError(tui, errors.Errorf("can't use widget %s without service %s.", w.Name, name))
	} else {
		f, err := s.CreateWidgets(w, tui)
		if err != nil {
			c <- DisplayError(tui, errors.Errorf("%s / %s: %s", name, w.Name, err.Error()))
		} else {
			c <- f
		}
	}
	close(c)
}

// Create all the widgets and populate them with data.
// Return slice of widgets
func (p *project) CreateNonConcWidgets() [][][]func() error {
	// TODO: use display.box instead of this shortcut
	err := p.addTitle(p.tui)
	if err != nil {
		err = errors.Wrapf(err, "can't add project title %s", p.name)
		DisplayError(p.tui, err)()
	}

	funcs := make([][][]func() error, len(p.widgets))

	for ir, row := range p.widgets {
		for ic, col := range row {
			funcs[ir] = append(funcs[ir], []func() error{})
			for _, w := range col {
				w = p.addDefaultTheme(w)

				service, err := p.mapServiceID(w.serviceID())
				if err != nil {
					funcs[ir][ic] = append(funcs[ir][ic], DisplayError(p.tui, err))
					continue
				}

				serviceName, err := mapServiceName(w.serviceID())
				if err != nil {
					funcs[ir][ic] = append(funcs[ir][ic], DisplayError(p.tui, err))
					continue
				}

				funcs[ir][ic] = append(funcs[ir][ic], getFuncRenderers(service, serviceName, w, p.tui))
			}
		}
	}

	return funcs
}

func getFuncRenderers(s service, name string, w Widget, tui *Tui) (f func() error) {
	if s == nil {
		f = DisplayError(tui, errors.Errorf(
			"Configuration error - you can't use the widget %s without the service %s.",
			w.Name,
			name,
		))
	}

	f, err := s.CreateWidgets(w, tui)
	if err != nil {
		f = DisplayError(tui, err)
	}

	return f
}

func (p *project) Render(funcs [][][]func() error) {
	for r, row := range p.widgets {
		for c, col := range row {
			for _, f := range funcs[r][c] {
				err := f()
				if err != nil {
					DisplayError(p.tui, err)()
				}
			}
			if len(col) > 0 {
				if err := p.tui.AddCol(p.sizes[r][c]); err != nil {
					DisplayError(p.tui, err)()
				}
			}
		}
		p.tui.AddRow()
		p.tui.Render()
	}
}

func (p *project) addTitle(tui *Tui) error {
	return tui.AddProjectTitle(p.name, p.nameOptions)
}
