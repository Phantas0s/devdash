package plateform

import (
	"fmt"

	"github.com/gizak/termui"
)

type termUI struct {
	body *termui.Grid
}

// NewGui returns a new Gui object with a given output mode.
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
	}, nil
}

func (termUI) Close() {
	termui.Close()
}

func (t termUI) TextBox(
	data string,
	fg uint16,
	bd uint16,
	bdlabel string,
	h int,
) {
	textBox := termui.NewPar(data)

	textBox.TextFgColor = termui.Attribute(fg)
	textBox.BorderFg = termui.Attribute(bd)
	textBox.BorderLabel = bdlabel
	textBox.Height = h

	// TODO add row should not appear here - single responsibility principle
	t.body.AddRows(
		termui.NewRow(termui.NewCol(2, 0, textBox)),
	)
}

func (t termUI) LineChart(data []float64, dimensions []string) {
	lc := termui.NewLineChart()
	lc.BorderLabel = "Users of the week"
	lc.Data = data
	lc.DataLabels = dimensions
	lc.Mode = "dot"
	lc.Width = 77
	lc.Height = 20
	lc.X = 0
	lc.Y = 12
	lc.AxesColor = termui.ColorWhite
	lc.LineColor = termui.ColorGreen | termui.AttrBold
}

func (t termUI) BarChart(data []int, dimensions []string, barWidth int) {
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

	// TODO add row should not appear here - single responsibility principle
	t.body.AddRows(
		termui.NewRow(termui.NewCol(10, 0, bc)),
	)
}

// KQuit set a key to quit the application.
func (termUI) KQuit(key string) {
	termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
		termui.StopLoop()
	})
}

func (t termUI) Render() {
	// Calculate the layout.
	t.body.Align()
	// Render the termui body.
	termui.Clear()
	termui.Render(t.body)
	// TODO render and loop are two different things - responsibility principle
	termui.Loop()
}
