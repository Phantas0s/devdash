package plateform

import (
	"fmt"

	"github.com/Phantas0s/termui"
)

type termUI struct {
	body    *termui.Grid
	widgets []termui.GridBufferer
	col     []*termui.Row
	row     []*termui.Row
}

// NewTermUI returns a new Terminal Interface object with a given output mode.
func NewTermUI(d bool) (*termUI, error) {
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

// Close termui.
func (termUI) Close() {
	termui.Close()
}

// AddCol to the termui grid system.
func (t *termUI) AddCol(size int) {
	t.col = append(t.col, termui.NewCol(size, 0, t.widgets...))
	t.widgets = []termui.GridBufferer{}
}

// AddRow to the termui grid system.
func (t *termUI) AddRow() {
	t.body.AddRows(termui.NewRow(t.col...))
	t.body.Align()
}

// TextBox widget type.
func (t *termUI) TextBox(
	data string,
	textColor uint16,
	borderColor uint16,
	title string,
	tc uint16,
	height int,
	multiline bool,
) {
	textBox := termui.NewPar(data)

	textBox.TextFgColor = termui.Attribute(textColor)
	textBox.BorderFg = termui.Attribute(borderColor)
	textBox.BorderLabel = title
	textBox.BorderLabelFg = termui.Attribute(tc)
	textBox.Height = height
	textBox.Multiline = multiline

	t.widgets = append(t.widgets, textBox)
}

// Title is a special TextBox widget type.
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

// BarChar widget type.
func (t *termUI) BarChart(
	data []int,
	dimensions []string,
	title string,
	tc uint16,
	bd uint16,
	fg uint16,
	nc uint16,
	enc uint16,
	height int,
	gap int,
	barWidth int,
	barColor uint16,
) {
	bc := termui.NewBarChart()
	bc.BorderLabel = title
	bc.Data = data
	bc.BorderLabelFg = termui.Attribute(tc)
	bc.DataLabels = dimensions
	bc.Height = height
	bc.TextColor = termui.Attribute(fg)
	bc.BorderFg = termui.Attribute(bd)
	bc.BarWidth = barWidth
	bc.BarColor = termui.Attribute(barColor)
	bc.BarGap = gap
	bc.NumColor = termui.Attribute(nc)
	bc.EmptyNumColor = termui.Attribute(enc)
	bc.Buffer()

	t.widgets = append(t.widgets, bc)
}

// StackedBarChar widget type.
func (t *termUI) StackedBarChart(
	data [8][]int,
	dimensions []string,
	title string,
	tc uint16,
	colors []uint16,
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
	bc.BorderLabelFg = termui.Attribute(tc)
	bc.BarWidth = barWidth
	bc.Height = height
	bc.BarGap = gap
	bc.DataLabels = dimensions
	bc.TextColor = termui.Attribute(fg)
	bc.BorderFg = termui.Attribute(bd)
	bc.BarColor = [8]termui.Attribute{termui.Attribute(colors[0]), termui.Attribute(colors[1])}
	bc.NumColor = [8]termui.Attribute{termui.Attribute(nc), termui.Attribute(nc)}

	t.widgets = append(t.widgets, bc)
}

// Table widget type.
func (t *termUI) Table(
	data [][]string,
	title string,
	tc uint16,
	bd uint16,
	fg uint16,
) {
	ta := termui.NewTable()
	ta.Rows = data
	ta.BorderLabel = title
	ta.FgColor = termui.Attribute(fg)
	ta.BorderLabelFg = termui.Attribute(tc)
	ta.BorderFg = termui.Attribute(bd)
	ta.SetSize()

	t.widgets = append(t.widgets, ta)
}

// KQuit set a key to quit the application.
func (*termUI) KQuit(key string) {
	termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
		termui.StopLoop()
	})
}

// Loop termui to receive events.
func (t *termUI) Loop() {
	termui.Loop()
}

// Render termui and delete the instance of the widgets rendered.
func (t *termUI) Render() {
	termui.Render(t.body)
	// delete every widget for the row rendered.
	t.removeWidgets()
}

func (t *termUI) removeWidgets() {
	t.row = []*termui.Row{}
	t.col = []*termui.Row{}
}

// Clean and create a new empty grid.
func (t *termUI) Clean() {
	t.body = termui.NewGrid()
	t.body.X = 0
	t.body.Y = 0
	t.body.BgColor = termui.ThemeAttr("bg")
	t.body.Width = termui.TermWidth()
}

func (t *termUI) HotReload() {
	termui.Close()
	_ = termui.Init()
	t.Clean()
}
