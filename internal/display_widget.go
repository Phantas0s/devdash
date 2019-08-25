package internal

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

		You can name the configuration file [my-config.yml](fg-blue,fg-bold), and then run [devdash -config my-config.yml](fg-green,fg-bold).

		There are multiple example of configurations there:
		[https://thedevdash.com/getting-started/](fg-blue,fg-bold).

		More complex configuration examples are available here:
		[https://github.com/Phantas0s/devdash#configuration-examples](fg-blue,fg-bold).

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
