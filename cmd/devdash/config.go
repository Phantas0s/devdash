package main

import (
	"bytes"
	"fmt"

	"github.com/Phantas0s/devdash/internal"
	"github.com/spf13/viper"
)

const (
	// keys
	kquit = "C-c"
)

type config struct {
	General  General   `mapstructure:"general"`
	Projects []Project `mapstructure:"projects"`
}

type General struct {
	Keys    map[string]string `mapstructure:"keys"`
	Refresh int64             `mapstructure:"refresh"`
}

type Project struct {
	Name         string            `mapstructure:"name"`
	Services     Services          `mapstructure:"services"`
	Widgets      []Row             `mapstructure:"widgets"`
	TitleOptions map[string]string `mapstructure:"options"`
}

// Row is constitued of columns
type Row struct {
	Row []Column `mapstructure: "row"`
}

// Col is constitued of widgets
type Column struct {
	Col []Widgets `mapstructure: "col"`
}

type Widgets struct {
	Size     string            `mapstructure:"size"`
	Elements []internal.Widget `mapstructure:"elements"`
}

type Services struct {
	GoogleAnalytics GoogleAnalytics `mapstructure:"google_analytics"`
	Monitor         Monitor         `mapstructure:"monitor"`
}

type GoogleAnalytics struct {
	Keyfile string `mapstructure:"keyfile"`
	ViewID  string `mapstructure:"view_id"`
}

type Monitor struct {
	Address string `mapstructure:"address"`
}

// OrderWidgets add the widgets to a three dimensional slice.
// First dimension: index of the rows (ir or indexRows)
// Second dimension: index of the columns (ic or indexColumn)
// Third dimension: index of the widget
func (p Project) OrderWidgets() ([][][]internal.Widget, [][]string) {
	rowLen := len(p.Widgets)

	rows := make([][][]internal.Widget, rowLen)
	sizes := make([][]string, rowLen)
	for ir, r := range p.Widgets {
		for ic, c := range r.Row {
			rows[ir] = append(rows[ir], []internal.Widget{}) // add columns to rows
			for _, ws := range c.Col {
				// keep sizes of columns and good order of widgets in a separate slice
				sizes[ir] = append(sizes[ir], ws.Size)

				// add widgets to columns
				rows[ir][ic] = append(rows[ir][ic], ws.Elements...)
			}
		}
	}

	return rows, sizes
}

func (g GoogleAnalytics) empty() bool {
	return g == GoogleAnalytics{}
}

func (m Monitor) empty() bool {
	return m == Monitor{}
}

func mapConfig(data []byte) config {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("could not read config data %s: %s", string(data), err))
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}

func (c config) KQuit() string {
	if ok := c.General.Keys["quit"]; ok != "" {
		return c.General.Keys["quit"]
	}

	return kquit
}
