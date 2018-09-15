package main

import (
	"flag"
	"fmt"
	"io/ioutil"

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

	for pn, p := range cfg.Projects {
		fmt.Println(pn)

		if err != nil {
			fmt.Println(err)
		}

		// create a slice of slice with shape row[rowNbr][WidgetNbr]map[widgetName]Widget
		rows := make([][]map[string]internal.Widget, len(p.Widgets))
		fmt.Println(p.Widgets)
		for i := 0; i < len(p.Widgets); i++ {
			for wn, w := range p.Widgets[i] {
				rows[i] = append(rows[i], map[string]internal.Widget{wn: w})
			}
		}

		gaService := p.Services.GoogleAnalytics
		gaWidget, err := internal.NewGaWidget(gaService.Keyfile, gaService.ViewID)
		if err != nil {
			fmt.Println(err)
		}

		project := internal.NewProject(pn, rows, gaWidget)
		if err != nil {
			fmt.Println(err)
		}

		err = project.Render(tui)
		if err != nil {
			fmt.Println(err)
		}
	}

	tui.Render()
	tui.Close()
}
