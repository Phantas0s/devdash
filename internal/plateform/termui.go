package plateform

import (
	"github.com/gizak/termui"
)

type termUI struct {
}

// NewGui returns a new Gui object with a given output mode.
func NewTermUI() (*termUI, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	return &termUI{}, nil
}

func (termUI) Close() {
	termui.Close()
}

func (t termUI) TextBox(data string, fg uint16, bd uint16, bdlabel string, h int) *termui.Par {
	textBox := termui.NewPar(data)

	textBox.TextFgColor = termui.Attribute(fg)
	textBox.BorderFg = termui.Attribute(bd)
	textBox.BorderLabel = bdlabel // + data.name
	textBox.Height = h

	return textBox
}
