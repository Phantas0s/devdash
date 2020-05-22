package internal

import "github.com/pkg/errors"

type displayWidget struct {
	tui *Tui
}

const displayBox = "display.box"

func NewDisplayWidget() *displayWidget {
	return &displayWidget{}
}

func (d displayWidget) CreateWidgets(widget Widget, tui *Tui) (f func() error, err error) {
	d.tui = tui

	switch widget.Name {
	case displayBox:
		f, err = d.box(widget)
	default:
		return nil, errors.Errorf("can't find the widget %s for service display", widget.Name)
	}

	return
}

func (d displayWidget) box(widget Widget) (f func() error, err error) {
	title := ""
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	content := "No content"
	if _, ok := widget.Options[optionContent]; ok {
		content = widget.Options[optionContent]
	}

	f = func() error {
		return d.tui.AddTextBox(
			content,
			title,
			widget.Options,
		)
	}

	return
}

func DisplayError(tui *Tui, err error) func() error {
	return func() error {
		return tui.AddTextBox(err.Error(), " ERROR ", map[string]string{
			optionBorderColor: "red",
			optionTextColor:   "red",
			optionTitleColor:  "red",
			optionMultiline:   "true",
		})
	}
}
