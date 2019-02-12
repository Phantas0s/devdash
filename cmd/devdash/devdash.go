package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/Phantas0s/devdash/internal"
	"github.com/Phantas0s/devdash/internal/plateform"
)

func main() {
	file := flag.String("config", ".devdash.yml", "The config file")
	flag.Parse()

	data, _ := ioutil.ReadFile(*file)
	cfg := mapConfig(data)

	termui, err := plateform.NewTermUI()
	if err != nil {
		fmt.Println(err)
		return
	}

	tui := internal.NewTUI(termui)
	tui.AddKQuit(cfg.KQuit())
	defer tui.Close()

	err = run(cfg.Projects, tui)
	if err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(time.Duration(cfg.General.Refresh) * time.Minute)
	go func() {
		for range ticker.C {
			cmd := exec.Command("clear") // for linux...
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}

			err = run(cfg.Projects, tui)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	tui.Render()
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

		monService := p.Services.Monitor
		if !monService.empty() {
			monWidget, err := internal.NewMonitorWidget(monService.Address)
			if err != nil {
				return err
			}
			project.WithMonitor(monWidget)
		}

		err = project.Render(tui)
		if err != nil {
			return err
		}
	}

	return nil
}
