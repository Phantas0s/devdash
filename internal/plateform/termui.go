package plateform

import (
	"fmt"

	"github.com/gizak/termui"
	"github.com/pkg/errors"
)

const maxRowSize = 12

type widget struct {
	element termui.GridBufferer
	size    int
}

type termUI struct {
	body *termui.Grid
	row  []widget
}

// NewTermUI returns a new Terminal Interface object with a given output mode.
func NewTermUI() (*termUI, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	// set the basic properties
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	return &termUI{
		body: body,
		row:  []widget{},
	}, nil
}

func (termUI) Close() {
	termui.Close()
}

func (t *termUI) TextBox(
	data string,
	fg uint16,
	bd uint16,
	bdlabel string,
	h int,
	size int,
) {
	textBox := termui.NewPar(data)

	textBox.TextFgColor = termui.Attribute(fg)
	textBox.BorderFg = termui.Attribute(bd)
	textBox.BorderLabel = bdlabel
	textBox.Height = h

	t.row = append(t.row, widget{element: textBox, size: size})
}

func (t *termUI) BarChart(data []int, dimensions []string, barWidth int, size int) {
	bc := termui.NewBarChart()
	bc.BorderLabel = "Bar Chart"
	bc.Data = data
	bc.BarWidth = barWidth
	bc.BarGap = 0
	bc.DataLabels = dimensions
	bc.Width = 200
	bc.Height = 10
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorYellow

	t.row = append(t.row, widget{element: bc, size: size})
}

// KQuit set a key to quit the application.
func (termUI) KQuit(key string) {
	termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
		termui.StopLoop()
	})
}

func (t *termUI) AddRow() error {
	err := t.validateRowSize()
	if err != nil {
		return err
	}

	var col []*termui.Row
	for _, w := range t.row {
		col = append(col, termui.NewCol(w.size, 0, w.element))
	}

	t.body.AddRows(termui.NewRow(col...))
	// clean the internal row
	t.row = []widget{}

	return nil
}

func (t termUI) validateRowSize() error {
	var ts int
	for _, w := range t.row {
		ts += w.size
	}

	if ts > maxRowSize {
		return errors.Errorf("could not create row: size %d too big", ts)
	}

	return nil
}

func (t *termUI) Render() {
	// Calculate the layout.
	t.body.Align()
	// Render the termui body.
	termui.Clear()
	termui.Render(t.body)
	// TODO render and loop are two different things - responsibility principle
	termui.Loop()
}
