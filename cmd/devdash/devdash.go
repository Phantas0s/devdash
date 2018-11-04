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
			cmd := exec.Command("clear") //Linux example, its tested
			cmd.Stdout = os.Stdout
			cmd.Run()

			tui.Init()
			err := run(cfg.Projects, tui)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	tui.Render()
}

func run(projects []Project, tui *internal.Tui) (err error) {
	for _, p := range projects {
		pn := p.Name

		if err != nil {
			return err
		}

		wc := 0
		for _, r := range p.Widgets {
			for _, c := range r.Row {
				for _, ws := range c.Col {
					for _ = range ws.Elements {
						wc = len(ws.Elements)
					}
				}
			}
		}

		rows := make([][][]internal.Widget, wc)
		sizes := make([][]string, wc)
		for ir, r := range p.Widgets {
			for ic, c := range r.Row {
				for _, ws := range c.Col {
					sizes[ir] = append(sizes[ir], ws.Size)
					for _, w := range ws.Elements {
						rows[ir] = append(rows[ir], []internal.Widget{})
						rows[ir][ic] = append(rows[ir][ic], w)
					}
				}
			}
		}

		gaService := p.Services.GoogleAnalytics
		gaWidget, err := internal.NewGaWidget(gaService.Keyfile, gaService.ViewID)
		if err != nil {
			return err
		}

		project := internal.NewProject(pn, rows, sizes, gaWidget)
		if err != nil {
			return err
		}

		err = project.Render(tui)
		if err != nil {
			return err
		}
	}

	return nil
}
