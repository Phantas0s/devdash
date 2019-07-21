package internal

func displayError(tui *Tui, err error) {
	tui.AddTextBox(err.Error(), " ERROR ", map[string]string{
		optionBorderColor: "red",
		optionTextColor:   "red",
		optionTitleColor:  "red",
	})
}

func displayNoFile(tui *Tui) {
	tui.AddTextBox("There is no file!", " No configuration file found ", map[string]string{
		optionBorderColor: "red",
		optionTextColor:   "red",
		optionTitleColor:  "red",
	})
}
