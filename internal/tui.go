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
	TextBox(data string, fg uint16, bd uint16, bdlabel string, h int, size int)
	BarChart(data []int, dimensions []string, barWidth int, bdLabel string, size int)
	AddRow() error
}

type kManager interface {
	KQuit(key string)
}

type manager interface {
	kManager
	renderer
	drawer
}

type textBoxAttr struct {
	Data    string
	Fg      uint16
	Bd      uint16
	Bdlabel string
	H       int
	Size    string
}

type barChartAttr struct {
	Data       []int
	Dimensions []string
	BarWidth   int
	Bdlabel    string
	Size       string
}

func NewTUI(instance manager) *Tui {
	return &Tui{
		instance: instance,
	}
}

type Tui struct {
	instance manager
}

func (t *Tui) AddTextBox(attr textBoxAttr) error {
	size, err := mapSize(attr.Size)
	if err != nil {
		return err
	}

	t.instance.TextBox(
		attr.Data,
		attr.Fg,
		attr.Bd,
		attr.Bdlabel,
		attr.H,
		size,
	)

	return nil
}

func (t *Tui) AddBarChart(attr barChartAttr) error {
	size, err := mapSize(attr.Size)
	if err != nil {
		return err
	}

	t.instance.BarChart(attr.Data, attr.Dimensions, attr.BarWidth, attr.Bdlabel, size)

	return nil
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

func (t *Tui) AddRow() {
	t.instance.AddRow()
}
