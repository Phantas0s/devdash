package cmd

import (
	"fmt"
	"time"

	"github.com/Phantas0s/devdash/internal"
	"github.com/Phantas0s/devdash/internal/platform"
)

func begin(file string, debug bool) {
	termui, err := platform.NewTermUI(debug)
	if err != nil {
		fmt.Println(err)
	}

	tui := internal.NewTUI(termui)
	defer tui.Close()

	cfg := mapConfig(file)

	// first display
	display(file, tui)

	tui.AddKQuit(cfg.KQuit())
	hotReload := make(chan time.Time)
	tui.AddKHotReload(cfg.KHotReload(), hotReload)

	ticker := time.NewTicker(time.Duration(cfg.RefreshTime()) * time.Second)

	go func(c <-chan time.Time) {
		for {
			val := <-c
			hotReload <- val
		}
	}(ticker.C)

	go func() {
		for range hotReload {
			if cfg.General.HotReload {
				tui.HotReload()
			} else {
				tui.Clean()
			}

			display(file, tui)
		}
	}()

	tui.Loop()
}

// TODO separate render from parsing projects
func display(file string, tui *internal.Tui) {
	cfg := mapConfig(file)
	for _, p := range cfg.Projects {
		rows, sizes := p.OrderWidgets()
		project := internal.NewProject(p.Name, p.NameOptions, rows, sizes, p.Themes, tui)

		gaService := p.Services.GoogleAnalytics
		if !gaService.empty() {
			gaWidget, err := internal.NewGaWidget(gaService.Keyfile, gaService.ViewID)
			if err != nil {
				internal.DisplayError(tui, err)
			}
			project.WithGa(gaWidget)
		}

		gscService := p.Services.GoogleSearchConsole
		if !gscService.empty() {
			gscWidget, err := internal.NewGscWidget(gscService.Keyfile, gscService.Address)
			if err != nil {
				internal.DisplayError(tui, err)
			}
			project.WithGoogleSearchConsole(gscWidget)
		}

		monService := p.Services.Monitor
		if !monService.empty() {
			monWidget, err := internal.NewMonitorWidget(monService.Address)
			if err != nil {
				internal.DisplayError(tui, err)
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
				internal.DisplayError(tui, err)
			}
			project.WithGithub(githubWidget)
		}

		travisService := p.Services.TravisCI
		if !travisService.empty() {
			travisWidget := internal.NewTravisCIWidget(travisService.Token)
			project.WithTravisCI(travisWidget)
		}

		feedlyService := p.Services.Feedly
		if !feedlyService.empty() {
			feedlyService := internal.NewFeedlyWidget(feedlyService.Address)
			project.WithFeedly(feedlyService)
		}

		gitService := p.Services.Git
		if !gitService.empty() {
			gitService := internal.NewGitWidget(gitService.Path)
			project.WithGit(gitService)
		}

		remoteHostService := p.Services.RemoteHost
		if !remoteHostService.empty() {
			remoteHostService, err := internal.NewRemoteHostWidget(
				remoteHostService.Username,
				remoteHostService.Address,
			)
			if err != nil {
				fmt.Println(err)
				internal.DisplayError(tui, err)
			}
			project.WithRemoteHost(remoteHostService)
		}

		renderFuncs := project.CreateWidgets()
		if !debug {
			project.Render(renderFuncs)
		}
	}
}
