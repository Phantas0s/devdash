// This package is an abstraction for any Terminal UI you want to use.
// It is not necessary since I don't need it and add complexity. It is kind of an
// experiment how to use interface effectively.
package internal

type renderer interface {
	Render()
	Close()
}

type drawer interface {
	TextBox(data string, fg uint16, bd uint16, bdlabel string, h int)
	LineChart(data []float64, dimensions []string)
	BarChart(data []int, dimensions []string, barWidth int)
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
}

type lineChartAttr struct {
	Data       []float64
	Dimensions []string
}

type barChartAttr struct {
	Data       []int
	Dimensions []string
	BarWidth   int
}

func NewTUI(instance manager) *Tui {
	return &Tui{
		instance: instance,
	}
}

type Tui struct {
	instance manager
}

func (t *Tui) AddTextBox(attr textBoxAttr) {
	t.instance.TextBox(attr.Data, attr.Fg, attr.Bd, attr.Bdlabel, attr.H)
}

func (t *Tui) AddLineChart(attr lineChartAttr) {
	t.instance.LineChart(attr.Data, attr.Dimensions)
}

func (t *Tui) AddBarChart(attr barChartAttr) {
	t.instance.BarChart(attr.Data, attr.Dimensions, attr.BarWidth)
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
