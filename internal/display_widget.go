package internal

func displayError(tui *Tui, err error) {
	tui.AddTextBox(err.Error(), "error", map[string]string{
		optionBorderColor: "red",
		optionTextColor:   "red",
		optionTitleColor:  "red",
	})
}
