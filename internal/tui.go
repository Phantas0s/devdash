// This package is an abstraction for any Terminal UI you want to use.
// It is not necessary since I don't need it and add complexity. It is kind of an
// experiment how to use interface effectively.
package internal

import (
	"strings"

	"github.com/pkg/errors"
)

type renderer interface {
	Render()
	Close()
}

type drawer interface {
	Text(text string, foreground uint16, size int)
	TextBox(data string, foreground uint16, background uint16, title string, h int)
	BarChart(data []int, dimensions []string, barWidth int, background uint16, backgroundLabel string)
	StackedBarChart(data [8][]int, dimensions []string, barWidth int, background uint16, backgroundLabel string)
	Table(data [][]string, background uint16, backgroundLabel string)
	AddCol(size int)
	AddRow() error
}

type kManager interface {
	KQuit(key string)
}

type manager interface {
	kManager
	renderer
	drawer
	Init()
}

// Value objects
type textBoxAttr struct {
	Data       string
	Foreground uint16
	Background uint16
	Title      string
	H          int
}

type textAttr struct {
	Text       string
	Foreground uint16
	Size       string
}

type barChartAttr struct {
	Data       []int
	Dimensions []string
	Background uint16
	BarWidth   int
	Title      string
}

type stackedBarChartAttr struct {
	Data       [8][]int
	Dimensions []string
	Background uint16
	BarWidth   int
	Title      string
}

type tableAttr struct {
	Data            [][]string
	Background      uint16
	BackgroundLabel string
}

func NewTUI(instance manager) *Tui {
	return &Tui{
		instance: instance,
	}
}

type Tui struct {
	instance manager
}

func (t *Tui) AddProjectTitle(attr textAttr) error {
	size, err := mapSize(attr.Size)
	if err != nil {
		return err
	}
	t.instance.Text(attr.Text, attr.Foreground, size)

	return nil
}

func (t *Tui) AddTextBox(attr textBoxAttr) {
	t.instance.TextBox(
		attr.Data,
		attr.Foreground,
		attr.Background,
		attr.Title,
		attr.H,
	)
}

func (t *Tui) AddBarChart(attr barChartAttr) {
	t.instance.BarChart(attr.Data, attr.Dimensions, attr.BarWidth, attr.Background, attr.Title)
}

func (t *Tui) AddStackedBarChart(attr stackedBarChartAttr) {
	t.instance.StackedBarChart(attr.Data, attr.Dimensions, attr.BarWidth, attr.Background, attr.Title)
}

func (t *Tui) AddTable(attr tableAttr) {
	t.instance.Table(attr.Data, attr.Background, attr.BackgroundLabel)
}

func (t *Tui) AddKQuit(key string) {
	t.instance.KQuit(key)
}

func (t *Tui) Render() {
	t.instance.Render()
}

func (t *Tui) Close() {
	t.instance.Close()
}

func mapSize(size string) (int, error) {
	s := strings.ToLower(size)
	switch s {
	case "xs":
		return 2, nil
	case "s":
		return 4, nil
	case "m":
		return 6, nil
	case "l":
		return 8, nil
	case "xl":
		return 10, nil
	case "xxl":
		return 12, nil
	default:
		return 0, errors.Errorf("could not find size %s", s)
	}
}

func (t *Tui) AddCol(size string) error {
	s, err := mapSize(size)
	if err != nil {
		return err
	}
	t.instance.AddCol(s)

	return nil
}

func (t *Tui) AddRow() {
	t.instance.AddRow()
}

func (t *Tui) Init() {
	t.instance.Init()
}
