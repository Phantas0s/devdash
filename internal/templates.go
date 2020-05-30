package internal

func DefaultTemplate() string {
	return `---
general:
  refresh: 600
  keys:
    quit: "C-c"
    hot_reload: "C-r"

projects:
  - name: Default dashboard located at $HOME/.config/devdash/default.yml
    services:
      monitor:
        address: "https://thevaluable.dev"
    widgets:
      - row:
          - col:
              size: "M"
              elements:
                - name: mon.box_availability
                  options:
                    title: " thevaluable.dev status "
                    color: yellow`
}
