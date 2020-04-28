// This package is an abstraction of the Terminal UI itself.
package internal

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	// Colors
	defaultC uint16 = iota
	black
	red
	green
	yellow
	blue
	magenta
	cyan
	white

	optionSize = "size"

	optionColor = "color"

	optionBorderColor   = "border_color"
	optionTextColor     = "text_color"
	optionNumColor      = "num_color"
	optionEmptyNumColor = "empty_num_color"

	optionBold = "bold"

	optionMultiline = "false"

	optionFirstColor  = "first_color"
	optionSecondColor = "second_color"
	optionThirdColor  = "third_color"
	optionFourthColor = "fourth_color"
	optionFifthColor  = "fifth_color"
	optionSixthColor  = "sixth_color"

	optionHeight = "height"

	optionBarGap   = "bar_gap"
	optionBarWidth = "bar_width"
	optionBarColor = "bar_color"
)

// map config size to ui size
var sizeLookup = map[string]int{
	"xxs": 1,
	"xs":  2,
	"s":   4,
	"m":   6,
	"l":   8,
	"xl":  10,
	"xxl": 12,
}

// map config color to ui color
var colorLookUp = map[string]uint16{
	"default": defaultC,
	"black":   black,
	"red":     red,
	"green":   green,
	"yellow":  yellow,
	"blue":    blue,
	"magenta": magenta,
	"cyan":    cyan,
	"white":   white,
}

// colorStr is used to map a color name to an ui color
func colorStr(value uint16) (key string) {
	for k, v := range colorLookUp {
		if v == value {
			key = k
			return
		}
	}
	return
}

type renderer interface {
	Render()
	Close()
	Clean()
}

type drawer interface {
	Title(
		title string,
		textColor uint16,
		borderColor uint16,
		bold bool,
		height int,
		size int,
	)
	TextBox(
		data string,
		textColor uint16,
		borderColor uint16,
		title string,
		titleColor uint16,
		height int,
		multiline bool,
		bold bool,
	)
	BarChart(
		data []int,
		dimensions []string,
		title string,
		tc uint16,
		bd uint16,
		fg uint16,
		nc uint16,
		enc uint16,
		height int,
		gap int,
		barWidth int,
		barColor uint16,
	)

	StackedBarChart(
		data [8][]int,
		dimensions []string,
		title string,
		tc uint16,
		colors []uint16,
		bd uint16,
		fg uint16,
		nc uint16,
		height int,
		gap int,
		barWidth int,
	)

	Table(
		data [][]string,
		title string,
		tc uint16,
		bd uint16,
		fg uint16,
	)
	AddCol(size int)
	AddRow()
}

type keyManager interface {
	KQuit(key string)
	KHotReload(key string, run func(), m *sync.Mutex)
}

type looper interface {
	Loop()
}

type reloader interface {
	HotReload()
}

type manager interface {
	keyManager
	renderer
	drawer
	looper
	reloader
}

type coloredElements struct {
	textColor     uint16
	borderColor   uint16
	titleColor    uint16
	numColor      uint16
	emptyNumColor uint16
	barColor      uint16
}

func createColoredElements(options map[string]string) coloredElements {
	ce := coloredElements{
		textColor:     defaultC,
		borderColor:   defaultC,
		titleColor:    defaultC,
		numColor:      defaultC,
		emptyNumColor: defaultC,
		barColor:      defaultC,
	}

	if _, ok := options[optionColor]; ok {
		color := colorLookUp[options[optionColor]]
		ce = coloredElements{
			textColor:     color,
			borderColor:   color,
			titleColor:    color,
			numColor:      black,
			emptyNumColor: color,
			barColor:      color,
		}
	}

	if _, ok := options[optionBorderColor]; ok {
		ce.borderColor = colorLookUp[options[optionBorderColor]]
	}

	if _, ok := options[optionTextColor]; ok {
		ce.textColor = colorLookUp[options[optionTextColor]]
	}

	if _, ok := options[optionTitleColor]; ok {
		ce.titleColor = colorLookUp[options[optionTitleColor]]
	}

	if _, ok := options[optionNumColor]; ok {
		ce.numColor = colorLookUp[options[optionNumColor]]
	}

	if _, ok := options[optionEmptyNumColor]; ok {
		ce.emptyNumColor = colorLookUp[options[optionEmptyNumColor]]
	}

	if _, ok := options[optionBarColor]; ok {
		ce.barColor = colorLookUp[options[optionBarColor]]
	}

	return ce
}

// AddCol to the TUI grid.
func (t *Tui) AddCol(size string) error {
	s, err := MapSize(size)
	if err != nil {
		return err
	}
	t.instance.AddCol(s)

	return nil
}

// AddRow to the TUI grid.
func (t *Tui) AddRow() {
	t.instance.AddRow()
}

// Render the TUI.
func (t *Tui) Render() {
	t.instance.Render()
}

// Close the TUI.
func (t *Tui) Close() {
	t.instance.Close()
}

func NewTUI(instance manager) *Tui {
	return &Tui{
		instance: instance,
	}
}

type Tui struct {
	instance manager
}

// Map the size of each column if t-shirt size is provided (XXS to XL).
// Otherwise use the numerical value provided in the config directly.
func MapSize(size string) (int, error) {
	s := strings.ToLower(size)
	if size, ok := sizeLookup[s]; ok {
		return size, nil
	}
	si, err := strconv.ParseInt(size, 0, 0)
	if err != nil {
		return 0, err
	}

	return int(si), err
}

// AddProjectTitle to the TUI.
func (t *Tui) AddProjectTitle(title string, options map[string]string) (err error) {
	size := "XXL"
	if _, ok := options[optionSize]; ok {
		size = options[optionSize]
	}

	bold := true
	if _, ok := options[optionBold]; ok {
		bold, err = strconv.ParseBool(options[optionBold])
		if err != nil {
			return errors.Wrapf(err, "can't convert %s to bool - please verify your configuration (correct values: true or false)", options[optionBold])
		}
	}

	var height int64 = 3
	if _, ok := options[optionHeight]; ok {
		height, _ = strconv.ParseInt(options[optionHeight], 0, 0)
	}

	s, err := MapSize(size)
	if err != nil {
		return err
	}

	ce := createColoredElements(options)
	t.instance.Title(
		title,
		ce.textColor,
		ce.borderColor,
		bold,
		int(height),
		s,
	)

	return nil
}

// AddTextBox to the TUI.
func (t *Tui) AddTextBox(
	data string,
	title string,
	options map[string]string,
) (err error) {

	var height int64 = 3
	if _, ok := options[optionHeight]; ok {
		height, _ = strconv.ParseInt(options[optionHeight], 0, 0)
	}

	multiline := false
	if _, ok := options[optionMultiline]; ok {
		multiline, err = strconv.ParseBool(options[optionMultiline])
		if err != nil {
			return errors.Wrapf(
				err,
				"can't convert %s to bool - please verify your configuration (correct values: 'true' or 'false')",
				options[optionMultiline],
			)
		}
	}

	bold := false
	if _, ok := options[optionBold]; ok {
		bold, err = strconv.ParseBool(options[optionBold])
		if err != nil {
			return errors.Wrapf(err, "can't convert %s to bool - please verify your configuration (correct values: true or false)", options[optionBold])
		}
	}

	ce := createColoredElements(options)
	t.instance.TextBox(
		data,
		ce.textColor,
		ce.borderColor,
		title,
		ce.titleColor,
		int(height),
		multiline,
		bold,
	)

	return nil
}

// AddBarChart to the TUI, a representation of the evolution of a dataset overtime.
func (t *Tui) AddBarChart(
	data []int,
	dimensions []string,
	title string,
	options map[string]string,
) {
	fmt.Println("BEGIN ADD BAR CHART")
	var height int64 = 10
	if _, ok := options[optionHeight]; ok {
		height, _ = strconv.ParseInt(options[optionHeight], 0, 0)
	}

	var gap int64 = 0
	if _, ok := options[optionBarGap]; ok {
		gap, _ = strconv.ParseInt(options[optionBarGap], 0, 0)
	}

	var barWidth int64 = 6
	if _, ok := options[optionBarWidth]; ok {
		barWidth, _ = strconv.ParseInt(options[optionBarWidth], 0, 0)
	}

	ce := createColoredElements(options)
	t.instance.BarChart(
		data,
		dimensions,
		title,
		ce.titleColor,
		ce.borderColor,
		ce.textColor,
		ce.numColor,
		ce.emptyNumColor,
		int(height),
		int(gap),
		int(barWidth),
		ce.barColor,
	)
	fmt.Println("END ADD BAR CHART")
}

// AddStackedBarChart to the TUI, which represent two or more dataset overtime.
func (t *Tui) AddStackedBarChart(
	data [8][]int,
	dimensions []string,
	title string,
	colors []uint16,
	options map[string]string,
) {
	var height int64 = 10
	if _, ok := options[optionHeight]; ok {
		height, _ = strconv.ParseInt(options[optionHeight], 0, 0)
	}

	var gap int64 = 0
	if _, ok := options[optionBarGap]; ok {
		gap, _ = strconv.ParseInt(options[optionBarGap], 0, 0)
	}

	var barWidth int64 = 6
	if _, ok := options[optionBarWidth]; ok {
		barWidth, _ = strconv.ParseInt(options[optionBarWidth], 0, 0)
	}

	ce := createColoredElements(options)
	t.instance.StackedBarChart(
		data,
		dimensions,
		title,
		ce.titleColor,
		colors,
		ce.borderColor,
		ce.textColor,
		ce.numColor,
		int(height),
		int(gap),
		int(barWidth),
	)
}

// AddTable to the TUI, with a header and the dataset.
func (t *Tui) AddTable(data [][]string, title string, options map[string]string) {
	ce := createColoredElements(options)
	t.instance.Table(
		data,
		title,
		ce.titleColor,
		ce.borderColor,
		ce.textColor,
	)
}

// Add keyboard shortcut from the config to quit DevDash. Default Control C.
func (t *Tui) AddKQuit(key string) {
	t.instance.KQuit(key)
}

func (t *Tui) AddKHotReload(key string, run func(), m *sync.Mutex) {
	t.instance.KHotReload(key, run, m)
}

// Loop the TUI to receive events.
func (t *Tui) Loop() {
	t.instance.Loop()
}

// Clean and reinitialize the TUI.
func (t *Tui) Clean() {
	t.instance.Clean()
}

// Hot reload the whole TUI
func (t *Tui) HotReload() {
	t.instance.HotReload()
}
