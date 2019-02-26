package plateform

import (
	"fmt"

	"github.com/Phantas0s/termui"
)

var debug bool = false

type termUI struct {
	body    *termui.Grid
	widgets []termui.GridBufferer
	col     []*termui.Row
	row     []*termui.Row
}

// NewTermUI returns a new Terminal Interface object with a given output mode.
func NewTermUI(d bool) (*termUI, error) {
	debug = d

	if err := termui.Init(); err != nil {
		return nil, err
	}

	// set the basic properties
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	debugPrint(body)

	return &termUI{
		body: body,
		row:  []*termui.Row{},
	}, nil
}

func (termUI) Close() {
	termui.Close()
}

func (t *termUI) AddCol(size int) {
	t.col = append(t.col, termui.NewCol(size, 0, t.widgets...))
	t.widgets = []termui.GridBufferer{}
}

func (t *termUI) AddRow() {
	t.body.AddRows(termui.NewRow(t.col...))
	t.body.Align()
}

func (t *termUI) TextBox(
	data string,
	textColor uint16,
	borderColor uint16,
	title string,
	titleColor uint16,
	height int,
) {
	textBox := termui.NewPar(data)

	textBox.TextFgColor = termui.Attribute(textColor)
	textBox.BorderFg = termui.Attribute(borderColor)
	textBox.BorderLabel = title
	textBox.BorderLabelFg = termui.Attribute(titleColor)
	textBox.Height = height

	t.widgets = append(t.widgets, textBox)
}

func (t *termUI) Title(
	title string,
	textColor uint16,
	borderColor uint16,
	bold bool,
	height int,
	size int,
) {
	pro := termui.NewPar(title)
	pro.TextFgColor = termui.Attribute(textColor)
	if bold {
		pro.TextFgColor = termui.Attribute(textColor) | termui.AttrBold
	}
	pro.BorderFg = termui.Attribute(borderColor)
	pro.Height = height

	t.body.AddRows(termui.NewCol(size, 0, pro))
}

func (t *termUI) BarChart(
	data []int,
	dimensions []string,
	title string,
	bd uint16,
	fg uint16,
	nc uint16,
	height int,
	gap int,
	barWidth int,
	barColor uint16,
) {
	bc := termui.NewBarChart()
	bc.BorderLabel = title
	bc.Data = data
	bc.BarGap = gap
	bc.DataLabels = dimensions
	bc.Height = height
	bc.TextColor = termui.Attribute(fg)
	bc.BarColor = termui.Attribute(barColor)
	bc.NumColor = termui.Attribute(nc)
	bc.BorderFg = termui.Attribute(bd)
	bc.BarWidth = barWidth
	bc.Buffer()

	t.widgets = append(t.widgets, bc)
}

func (t *termUI) StackedBarChart(
	data [8][]int,
	dimensions []string,
	title string,
	bd uint16,
	fg uint16,
	nc uint16,
	height int,
	gap int,
	barWidth int,
) {
	bc := termui.NewMBarChart()
	bc.BorderLabel = title
	bc.Data = data
	bc.BarWidth = barWidth
	bc.Height = height
	bc.BarGap = gap
	bc.DataLabels = dimensions
	bc.TextColor = termui.Attribute(fg)
	bc.BorderFg = termui.Attribute(bd)
	bc.NumColor = [8]termui.Attribute{termui.Attribute(nc)}
	bc.ShowScale = true
	// bc.SetMax(10)

	t.widgets = append(t.widgets, bc)
}

func (t *termUI) Table(
	data [][]string,
	title string,
	bd uint16,
	fg uint16,
) {
	ta := termui.NewTable()
	ta.Rows = data
	ta.BorderLabel = title
	ta.FgColor = termui.Attribute(fg)
	ta.BorderFg = termui.Attribute(bd)
	ta.SetSize()

	t.widgets = append(t.widgets, ta)
}

// KQuit set a key to quit the application.
func (termUI) KQuit(key string) {
	termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
		termui.StopLoop()
	})
}

func (t *termUI) Loop() {
	termui.Loop()
}

func (t *termUI) Render() {
	termui.Render(t.body)
	// delete every widget for the row rendered.
	t.removeWidgets()
}

func (t *termUI) removeWidgets() {
	t.row = []*termui.Row{}
	t.col = []*termui.Row{}
}

func (t *termUI) Clean() {
	t.body = termui.NewGrid()
	t.body.X = 0
	t.body.Y = 0
	t.body.BgColor = termui.ThemeAttr("bg")
	t.body.Width = termui.TermWidth()
}

func debugPrint(v interface{}) {
	if debug {
		fmt.Println(v)
	}
}
