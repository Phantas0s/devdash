package internal

import "github.com/pkg/errors"

type displayWidget struct {
	tui *Tui
}

const displayBox = "display.box"

func NewDisplayWidget() *displayWidget {
	return &displayWidget{}
}

func (d displayWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	d.tui = tui

	switch widget.Name {
	case displayBox:
		err = d.box(widget)
	default:
		return errors.Errorf("can't find the widget %s for service display", widget.Name)
	}

	return
}

func (d displayWidget) box(widget Widget) error {
	title := ""
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	content := "No content"
	if _, ok := widget.Options[optionContent]; ok {
		content = widget.Options[optionContent]
	}

	err := d.tui.AddTextBox(
		content,
		title,
		widget.Options,
	)
	if err != nil {
		return err
	}

	return nil
}

func DisplayError(tui *Tui, err error) {
	_ = tui.AddTextBox(err.Error(), " ERROR ", map[string]string{
		optionBorderColor: "red",
		optionTextColor:   "red",
		optionTitleColor:  "red",
		optionMultiline:   "true",
	})
}

func DisplayNoFile(tui *Tui) {
	_ = tui.AddTextBox(
		`
		In order to use DevDash, you need to provide [a configuration file ](fg-bold).

		You can name the configuration file [my-config.yml](fg-blue,fg-bold), and then run [devdash -config my-config.yml](fg-green,fg-bold)

		There are multiple example of configurations there:
		[https://thedevdash.com/getting-started/](fg-blue,fg-bold)

		More complex configuration examples are available here:
		[https://thedevdash.com/getting-started/use-cases/](fg-blue,fg-bold)

		`,
		" Welcome to DevDash! ",
		map[string]string{
			optionBorderColor: "yellow",
			optionTextColor:   "default",
			optionTitleColor:  "yellow",
			optionHeight:      "14",
			optionMultiline:   "true",
		},
	)
}
