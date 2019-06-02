package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Phantas0s/devdash/internal"
	"github.com/Phantas0s/devdash/internal/plateform"
)

// debug mode
var debug *bool

func main() {
	file := flag.String("config", ".devdash.yml", "The config file")
	debug = flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	cfg, tui, err := loadFile(*file)
	if err != nil {
		fmt.Println(err)
	}
	defer tui.Close()

	err = run(cfg.Projects, tui)
	if err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(time.Duration(cfg.General.Refresh) * time.Second)
	go func() {
		for range ticker.C {
			tui.Clean()

			err = run(cfg.Projects, tui)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	tui.Loop()
}

func loadFile(file string) (config, *internal.Tui, error) {
	data, _ := ioutil.ReadFile(file)
	cfg := mapConfig(data)

	termui, err := plateform.NewTermUI(*debug)
	if err != nil {
		return config{}, nil, err
	}

	tui := internal.NewTUI(termui)
	tui.AddKQuit(cfg.KQuit())

	return cfg, tui, nil
}

func run(projects []Project, tui *internal.Tui) (err error) {
	for _, p := range projects {
		rows, sizes := p.OrderWidgets()
		if err != nil {
			return err
		}
		project := internal.NewProject(p.Name, p.TitleOptions, rows, sizes)

		gaService := p.Services.GoogleAnalytics
		if !gaService.empty() {
			gaWidget, err := internal.NewGaWidget(gaService.Keyfile, gaService.ViewID)
			if err != nil {
				return err
			}
			project.WithGa(gaWidget)
		}

		gscService := p.Services.GoogleSearchConsole
		if !gscService.empty() {
			gscWidget, err := internal.NewGscWidget(gscService.Keyfile, gscService.ViewID, gscService.Address)
			if err != nil {
				return err
			}
			project.WithGoogleSearchConsole(gscWidget)
		}

		monService := p.Services.Monitor
		if !monService.empty() {
			monWidget, err := internal.NewMonitorWidget(monService.Address)
			if err != nil {
				return err
			}
			project.WithMonitor(monWidget)
		}

		githubService := p.Services.Github
		if !githubService.empty() {
			githubWidget, err := internal.NewGithubWidget(
				githubService.Token,
				githubService.Owner,
				githubService.Repository,
			)
			if err != nil {
				return err
			}
			project.WithGithub(githubWidget)
		}

		err = project.Render(tui, *debug)
		if err != nil {
			return err
		}
	}

	return nil
}
