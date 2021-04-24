package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"strings"

	"github.com/Phantas0s/devdash/internal"
	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

const (
	// keys
	kQuit      = "C-c"
	kHotReload = "C-r"
	kEdit      = "C-e"
)

type config struct {
	General  General   `mapstructure:"general"`
	Projects []Project `mapstructure:"projects"`
}

type General struct {
	Keys    map[string]string `mapstructure:"keys"`
	Refresh int64             `mapstructure:"refresh"`
	Editor  string            `mapstructure:"editor"`
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
	RemoteHost          RemoteHost      `mapstructure:"remote_host"`
	Localhost           RemoteHost      `mapstructure:"local_host"`
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

type RemoteHost struct {
	Username string `mapstructure:"username"`
	Address  string `mapstructure:"address"`
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

func (g RemoteHost) empty() bool {
	return g == RemoteHost{}
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

func dashPath() string {
	return filepath.Join(xdg.ConfigHome, "devdash")
}

// Map config and return it with the config path
func mapConfig(cfgFile string) (config, string) {
	if cfgFile == "" {
		cfgFile = "default.yml"
		createConfig(dashPath(), cfgFile, defaultConfig())
	}

	// viper.AddConfigPath(home)
	viper.AddConfigPath(dashPath())
	viper.AddConfigPath(".")

	viper.SetConfigName(removeExt(cfgFile))
	err := viper.ReadInConfig()
	if err != nil {
		tryReadFile(cfgFile)
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	prefix := "DEVDASH"
	for k, _ := range cfg.Projects {
		if cfg.Projects[k].Services.GoogleAnalytics.Keyfile == "" {
			cfg.Projects[k].Services.GoogleAnalytics.Keyfile = os.Getenv(prefix + "_GA_KEYFILE")
		}

		if cfg.Projects[k].Services.GoogleSearchConsole.Keyfile == "" {
			cfg.Projects[k].Services.GoogleSearchConsole.Keyfile = os.Getenv(prefix + "_GSC_KEYFILE")
		}

		if cfg.Projects[k].Services.Github.Token == "" {
			cfg.Projects[k].Services.Github.Token = os.Getenv(prefix + "_GITHUB_TOKEN")
		}
	}

	return cfg, viper.ConfigFileUsed()
}

func removeExt(filepath string) string {
	ext := []string{".json", ".yml", ".yaml", ".toml"}
	for _, v := range ext {
		filepath = strings.Replace(filepath, v, "", -1)
	}

	return filepath
}

func createConfig(path string, filename string, template string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}

	f := filepath.Join(path, filename)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		file, _ := os.Create(f)
		defer file.Close()

		if file != nil {
			_, err := file.Write([]byte(template))
			if err != nil {
				panic(err)
			}
		}
	}

	return f
}

func tryReadFile(cfgFile string) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		panic(fmt.Errorf("config %s doesnt exists", cfgFile))
	}

	f, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(fmt.Errorf("could not read file %s", cfgFile))
	}

	viper.SetConfigType(strings.Trim(filepath.Ext(cfgFile), "."))
	err = viper.ReadConfig(bytes.NewBuffer(f))
	if err != nil {
		panic(fmt.Errorf("could not read config %s data", string(f)))
	}
}

// Keyboard events
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

func (c config) KEdit() string {
	if ok := c.General.Keys["edit"]; ok != "" {
		return c.General.Keys["edit"]
	}

	return kEdit
}

func defaultConfig() string {
	return `---
general:
  refresh: 600
  keys:
    quit: "C-c"
    hot_reload: "C-r"


projects:
  - name: Default dashboard located at $HOME/.config/devdash/default.yml
    services:
      monitor:
        address: "https://thevaluable.dev"
    widgets:
      - row:
          - col:
              size: "M"
              elements:
                - name: mon.box_availability
                  options:
                    title: " thevaluable.dev status "
                    color: yellow`
}
