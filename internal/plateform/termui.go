package plateform

import (
	"fmt"

	"github.com/gizak/termui"
	"github.com/pkg/errors"
)

const maxRowSize = 12

type termUI struct {
	body    *termui.Grid
	widgets []termui.GridBufferer
	col     []*termui.Row
	row     []*termui.Row
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
		row:  []*termui.Row{},
	}, nil
}

func (t *termUI) Init() {
	// set the basic properties
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	t.body = body
	t.widgets = []termui.GridBufferer{}
	t.col = []*termui.Row{}
	t.row = []*termui.Row{}
}

func (termUI) Close() {
	termui.Close()
}

func (t *termUI) AddCol(size int) {
	t.col = append(t.col, termui.NewCol(size, 0, t.widgets...))
	t.widgets = []termui.GridBufferer{}
}

func (t *termUI) AddRow() error {
	t.body.AddRows(termui.NewRow(t.col...))
	t.body.Align()
	termui.Render(t.body)

	// clean the internal row
	t.row = []*termui.Row{}
	t.col = []*termui.Row{}

	return nil
}

func (t termUI) validateRowSize() error {
	var ts int
	for _, r := range t.row {
		for _, c := range r.Cols {
			ts += c.Offset
		}
	}

	if ts > maxRowSize {
		return errors.Errorf("could not create row: size %d too big", ts)
	}

	return nil
}

func (t *termUI) Render() {
	termui.Loop()
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

	t.widgets = append(t.widgets, textBox)
}

func (t *termUI) Text(text string, fg uint16, size int) {
	pro := termui.NewPar(text)
	pro.Border = false
	pro.TextFgColor = termui.Attribute(fg)

	t.body.AddRows(termui.NewCol(size, 0, pro))
}

func (t *termUI) BarChart(data []int, dimensions []string, barWidth int, bdLabel string, size int) {
	bc := termui.NewBarChart()
	bc.BorderLabel = bdLabel
	bc.Data = data
	bc.BarWidth = barWidth
	bc.BarGap = 0
	bc.DataLabels = dimensions
	bc.Width = 200
	bc.Height = 10
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorBlue
	bc.NumColor = termui.ColorWhite

	t.widgets = append(t.widgets, bc)
}

// KQuit set a key to quit the application.
func (termUI) KQuit(key string) {
	termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
		termui.StopLoop()
	})
}
