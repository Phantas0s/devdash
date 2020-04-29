package main

import (
	"bytes"
	"fmt"

	"github.com/Phantas0s/devdash/internal"
	"github.com/spf13/viper"
)

const (
	// keys
	kQuit      = "C-c"
	kHotReload = "C-r"
)

type config struct {
	General  General   `mapstructure:"general"`
	Projects []Project `mapstructure:"projects"`
}

type General struct {
	Keys      map[string]string `mapstructure:"keys"`
	Refresh   int64             `mapstructure:"refresh"`
	HotReload bool              `mapstructure:"hot_reload"`
}

// RefreshTime return the duration before refreshing the data of all widgets, in seconds.
func (c config) RefreshTime() int64 {
	if c.General.Refresh == 0 {
		return 60
	}

	return c.General.Refresh
}

type Project struct {
	Name        string                       `mapstructure:"name"`
	NameOptions map[string]string            `mapstructure:"name_options"`
	Services    Services                     `mapstructure:"services"`
	Themes      map[string]map[string]string `mapstructure:"themes"`
	Widgets     []Row                        `mapstructure:"widgets"`
}

// Row is constitued of columns
type Row struct {
	Row []Column `mapstructure:"row"`
}

// Col is constitued of widgets
type Column struct {
	Col []Widgets `mapstructure:"col"`
}

type Widgets struct {
	Size     string            `mapstructure:"size"`
	Elements []internal.Widget `mapstructure:"elements"`
}

type Services struct {
	GoogleAnalytics     GoogleAnalytics `mapstructure:"google_analytics"`
	GoogleSearchConsole SearchConsole   `mapstructure:"google_search_console"`
	Monitor             Monitor         `mapstructure:"monitor"`
	Github              Github          `mapstructure:"github"`
	TravisCI            TravisCI        `mapstructure:"travis"`
	Feedly              Feedly          `mapstructure:"feedly"`
	Git                 Git             `mapstructure:"git"`
}

type GoogleAnalytics struct {
	Keyfile string `mapstructure:"keyfile"`
	ViewID  string `mapstructure:"view_id"`
}

type SearchConsole struct {
	Keyfile string `mapstructure:"keyfile"`
	Address string `mapstructure:"address"`
}

type Monitor struct {
	Address string `mapstructure:"address"`
}

type Github struct {
	Token      string `mapstructure:"token"`
	Owner      string `mapstructure:"owner"`
	Repository string `mapstructure:"repository"`
}

type TravisCI struct {
	Token string `mapstructure:"token"`
}

type Feedly struct {
	Address string `mapstructure:"address"`
}

type Git struct {
	Path string `mapstructure:"path"`
}

func (g GoogleAnalytics) empty() bool {
	return g == GoogleAnalytics{}
}

func (m Monitor) empty() bool {
	return m == Monitor{}
}
func (s SearchConsole) empty() bool {
	return s == SearchConsole{}
}

func (g Github) empty() bool {
	return g == Github{}
}

func (t TravisCI) empty() bool {
	return t == TravisCI{}
}

func (f Feedly) empty() bool {
	return f == Feedly{}
}

func (g Git) empty() bool {
	return g == Git{}
}

// OrderWidgets add the widgets to a three dimensional slice.
// First dimension: index of the rows (ir or indexRows).
// Second dimension: index of the columns (ic or indexColumn).
// Third dimension: index of the widget.
func (p Project) OrderWidgets() ([][][]internal.Widget, [][]string) {
	rows := make([][][]internal.Widget, len(p.Widgets))
	sizes := make([][]string, len(p.Widgets))
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

	return kQuit
}

func (c config) KHotReload() string {
	if ok := c.General.Keys["hot_reload"]; ok != "" {
		return c.General.Keys["hot_reload"]
	}

	return kHotReload
}
