package cmd

// TODO see gocket to make the command right (with possibility to use env variables)

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Phantas0s/devdash/internal"
	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/spf13/cobra"
)

var (
	// Used for flags
	file    string
	logpath string
	debug   bool

	rootCmd = &cobra.Command{
		Use:   "devdash",
		Short: "DevDash is a highly configurable terminal dashboard for developers and creators",
		Long:  `DevDash is a highly flexible terminal dashboard for developers and creators, which allows you to gather and refresh the data you really need from Google Analytics, Google Search Console, Github, TravisCI, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&file, "config", "c", "", "A valid dashboard configuration")
	// TODO logger
	// rootCmd.Flags().StringVarP(&logpath, "logpath", "l", "", "Path for logging")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug Mode - doesn't display graph")
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(editCmd())
	rootCmd.AddCommand(generateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}

// run the dashboard.
func run() {
	termui, err := platform.NewTermUI(debug)
	if err != nil {
		fmt.Println(err)
	}

	// Create the TUI.
	tui := internal.NewTUI(termui)
	defer tui.Close()

	// Map dashboard config to a struct Config.
	cfg, file := mapConfig(file)

	// Passing a time.Time to this channel reload the entire dashboard.
	hotReload := make(chan time.Time)

	// Add keystrokes to reload and quit
	tui.AddKHotReload(cfg.KHotReload(), hotReload)
	tui.AddKQuit(cfg.KQuit())

	// Passing a bool to this channel stop the automatic reload of the dashboard.
	stopAutoReload := make(chan bool)
	autoReload(cfg.RefreshTime(), stopAutoReload, hotReload)

	editor := os.Getenv("EDITOR")
	if cfg.General.Editor != "" {
		editor = cfg.General.Editor
	}

	// Add keystroke (managed by TUI) to edit the configuration in a CLI editor.
	// Wrap edit config in lambda to defer the execution.
	tui.AddKEdit(
		cfg.KEdit(),
		func() {
			stopReload(stopAutoReload)
			editDashboard(editor, file)
			hotReload <- time.Now()
			autoReload(cfg.RefreshTime(), stopAutoReload, hotReload)
		},
	)

	// First display.
	build(file, tui)

	// Automatic reload
	go func() {
		for hr := range hotReload {
			tui.HotReload()
			build(file, tui)
			if debug {
				fmt.Println("Last reload: " + hr.Format("2006-01-02 15:04:05"))
			}
		}
	}()

	tui.Loop()
}

func autoReload(refresh int64, stopAutoReload <-chan bool, hotReload chan<- time.Time) {
	go func() {
		ticker := time.NewTicker(time.Duration(refresh) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopAutoReload:
				return
			case tick := <-ticker.C:
				hotReload <- tick
			}
		}
	}()
}

func stopReload(stopAutoReload chan<- bool) {
	stopAutoReload <- true
}

// build every services present in the configuration
func build(file string, tui *internal.Tui) {
	cfg, _ := mapConfig(file)
	for _, p := range cfg.Projects {
		rows, sizes := p.OrderWidgets()
		project := internal.NewProject(p.Name, p.NameOptions, rows, sizes, p.Themes, tui)

		gaService := p.Services.GoogleAnalytics
		if !gaService.empty() {
			gaWidget, err := internal.NewGaWidget(gaService.Keyfile, gaService.ViewID)
			if err != nil {
				internal.DisplayError(tui, err)()
			} else {
				project.WithGa(gaWidget)
			}
		}

		gscService := p.Services.GoogleSearchConsole
		if !gscService.empty() {
			gscWidget, err := internal.NewGscWidget(gscService.Keyfile, gscService.Address)
			if err != nil {
				internal.DisplayError(tui, err)()
			} else {
				project.WithGoogleSearchConsole(gscWidget)
			}
		}

		monService := p.Services.Monitor
		if !monService.empty() {
			monWidget, err := internal.NewMonitorWidget(monService.Address)
			if err != nil {
				internal.DisplayError(tui, err)()
			} else {
				project.WithMonitor(monWidget)
			}
		}

		githubService := p.Services.Github
		if !githubService.empty() {
			githubWidget, err := internal.NewGithubWidget(
				githubService.Token,
				githubService.Owner,
				githubService.Repository,
			)
			if err != nil {
				internal.DisplayError(tui, err)()
			} else {
				project.WithGithub(githubWidget)
			}
		}

		travisService := p.Services.TravisCI
		if !travisService.empty() {
			travisWidget := internal.NewTravisCIWidget(travisService.Token)
			project.WithTravisCI(travisWidget)
		}

		feedlyService := p.Services.Feedly
		if !feedlyService.empty() {
			feedlyWidget := internal.NewFeedlyWidget(feedlyService.Address)
			project.WithFeedly(feedlyWidget)
		}

		gitService := p.Services.Git
		if !gitService.empty() {
			gitWidget := internal.NewGitWidget(gitService.Path)
			project.WithGit(gitWidget)
		}

		remoteHostService := p.Services.RemoteHost
		if !remoteHostService.empty() {
			remoteHostWidget, err := internal.NewHostWidget(
				remoteHostService.Username,
				remoteHostService.Address,
			)
			if err != nil {
				fmt.Println(err)
				internal.DisplayError(tui, err)()
			} else {
				project.WithRemoteHost(remoteHostWidget)
			}
		}

		localhost, err := internal.NewHostWidget("localhost", "localhost")
		if err != nil {
			fmt.Println(err)
			internal.DisplayError(tui, err)()
		}
		project.WithLocalhost(localhost)

		// TODO choice between concurency and non concurency
		// renderFuncs := project.CreateNonConcWidgets()
		renderFuncs := project.CreateWidgets()
		if !debug {
			project.Render(renderFuncs)
		}
	}
}

// TODO - Wrap logger. If logger nil, drop the message
func InitLoggerFile(logpath string) *log.Logger {
	if logpath == "" {
		return nil
	}

	file, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return log.New(file, "", 0)
}
