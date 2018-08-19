package internal

// TODO should be dcoupled, not termui here
import "github.com/gizak/termui"

type TB interface {
	TextBox(data string, fg uint16, bd uint16, bdlabel string, h int) *termui.Par
}

// Or property with ALL properties
type TextBoxAttr struct {
	data    string
	fg      uint16
	bd      uint16
	bdlabel string
	h       int
}

type TUI struct {
	// TODO general interface for now, to change to something more specific
	widgets []interface{}
	textbox TB
}

func (t TUI) AddTextBox(attr TextBoxAttr) {
	t.widgets = append(t.widgets, t.textbox.TextBox(attr.data, attr.fg, attr.bd, attr.bdlabel, attr.h))
}
