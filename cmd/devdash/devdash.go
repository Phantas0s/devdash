package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/Phantas0s/devdash/internal"
	"github.com/Phantas0s/devdash/internal/platform"
	"golang.org/x/crypto/ssh/terminal"
)

// debug mode
var debug *bool

func main() {
	file := flag.String("config", ".devdash.yml", "The config file")
	debug = flag.Bool("debug", false, "Debug mode")
	term := flag.Bool("term", false, "Display terminal dimensions")
	flag.Parse()

	if *term {
		width, height, _ := terminal.GetSize(0)
		fmt.Printf("Width: %d, Height: %d", width, height)
		return
	}

	termui, err := platform.NewTermUI(*debug)
	if err != nil {
		fmt.Println(err)
	}

	tui := internal.NewTUI(termui)
	defer tui.Close()

	cfg := loadFile(*file)
	run(*file, tui)()

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		tui.AddKQuit("C-c")
		internal.DisplayNoFile(tui)
		err := tui.AddCol("5")
		if err != nil {
			fmt.Println(err)
		}

		tui.AddRow()
		tui.Render()
	} else {
		var m sync.Mutex
		tui.AddKQuit(cfg.KQuit())
		tui.AddKHotReload(cfg.KHotReload(), run(*file, tui), &m)
		ticker := time.NewTicker(time.Duration(cfg.RefreshTime()) * time.Second)
		go func() {
			for range ticker.C {
				m.Lock()
				if cfg.General.HotReload {
					tui.HotReload()
				} else {
					tui.Clean()
				}

				run(*file, tui)()
				m.Unlock()
			}
		}()
	}

	tui.Loop()
}

func loadFile(file string) config {
	data, _ := ioutil.ReadFile(file)
	cfg := mapConfig(data)

	return cfg
}

func run(file string, tui *internal.Tui) func() {
	return func() {
		cfg := loadFile(file)
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
				travisWidget := internal.NewTravisCIWidget(
					travisService.Token,
				)
				project.WithTravisCI(travisWidget)
			}

			feedlyService := p.Services.Feedly
			if !feedlyService.empty() {
				feedlyService := internal.NewFeedlyWidget(
					feedlyService.Address,
				)
				project.WithFeedly(feedlyService)
			}

			gitService := p.Services.Git
			if !gitService.empty() {
				gitService := internal.NewGitWidget(
					gitService.Path,
				)
				project.WithGit(gitService)
			}

			project.Render(*debug)
		}
	}
}
