package plateform

import (
	termui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

const maxRowSize = 12

type termUI struct {
	body    *termui.Grid
	widgets []interface{}
	col     []interface{}
	row     []interface{}
}

// NewTermUI returns a new Terminal Interface object with a given output mode.
func NewTermUI() (*termUI, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	// set the basic properties
	body := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	body.SetRect(0, 0, termWidth, termHeight)

	// body.Y = 0
	// body.BgColor = termui.ThemeAttr("bg")
	// body.Width = termui.TermWidth()

	return &termUI{
		body: body,
	}, nil
}

func (t *termUI) Init() {
	// set the basic properties
	// body := termui.NewGrid()
	// body.X = 0
	// body.Y = 0
	// body.BgColor = termui.ThemeAttr("bg")
	// body.Width = termui.TermWidth()

	// t.body = body
	// t.widgets = []termui.Buffer{}
	// t.col = []*termui.GridItem{}
	// t.row = []*termui.GridItem{}
}

func (termUI) Close() {
	// termui.Close()
}

func (t *termUI) AddCol(size int) {
	t.col = append(t.col, termui.NewCol(0.3, t.widgets...))
	t.widgets = nil
}

func (t *termUI) AddRow() error {
	t.row = append(t.row, termui.NewRow(0.5, t.col...))
	t.col = nil
	// termui.Render(t.body)

	// clean the internal row
	// t.row = []*termui.Row{}
	// t.col = []*termui.Row{}

	return nil
}

// func (t termUI) validateRowSize() error {
// 	var ts int
// 	for _, r := range t.row {
// 		for _, c := range r.Cols {
// 			ts += c.Offset
// 		}
// 	}

// 	if ts > maxRowSize {
// 		return errors.Errorf("could not create row: size %d too big", ts)
// 	}

// 	return nil
// }

func (t *termUI) TextBox(
	data string,
	fg uint16,
	bd uint16,
	bdlabel string,
	h int,
) {
	p0 := widgets.NewParagraph()
	p0.Text = data
	p0.SetRect(0, 0, 20, 5)
	p0.Border = false

	// textBox.TextFgColor = termui.Attribute(fg)
	// textBox.BorderFg = termui.Attribute(bd)
	// textBox.BorderLabel = bdlabel
	// textBox.Height = h

	t.widgets = append(t.widgets, p0)
}

func (t *termUI) Text(text string, fg uint16, size int) {
	textBox := widgets.NewParagraph()
	textBox.Text = text

	t.widgets = append(t.widgets, textBox)
}

func (t *termUI) BarChart(
	data []int,
	dimensions []string,
	barWidth int,
	bd uint16,
	bdLabel string,
) {
	bc := widgets.NewBarChart()
	// bc.Data = data
	// bc.BorderLabel = bdLabel
	// bc.Data = data
	// bc.BarWidth = barWidth
	// bc.BarGap = 0
	// bc.DataLabels = dimensions
	// bc.Width = 200
	// bc.Height = 10
	// bc.TextColor = termui.ColorGreen
	// bc.BarColor = termui.ColorBlue
	// bc.NumColor = termui.ColorWhite
	// bc.BorderFg = termui.Attribute(bd)

	t.widgets = append(t.widgets, bc)
}

func (t *termUI) StackedBarChart(
	data [8][]int,
	dimensions []string,
	barWidth int,
	bd uint16,
	bdLabel string,
) {
	bc := widgets.NewStackedBarChart()
	// bc.BorderLabel = bdLabel
	// bc.Data = data
	// bc.BarWidth = barWidth
	// bc.DataLabels = dimensions
	// bc.Width = 200
	// bc.Height = 20
	// bc.TextColor = termui.ColorGreen
	// bc.BorderFg = termui.Attribute(bd)
	// bc.SetMax(10)

	t.widgets = append(t.widgets, bc)
}

func (t *termUI) Table(
	data [][]string,
	bd uint16,
	bdLabel string,
) {
	ta := widgets.NewTable()
	ta.Rows = data
	// ta.BorderLabel = bdLabel
	// ta.FgColor = termui.ColorGreen
	// ta.BorderFg = termui.Attribute(bd)
	// ta.SetSize()

	t.widgets = append(t.widgets, ta)
}

// KQuit set a key to quit the application.
func (termUI) KQuit(key string) {
	// termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
	// 	termui.StopLoop()
	// })
}

func (t *termUI) Render() {
	t.body.Set(t.row...)
	termui.Render(t.body)
}
