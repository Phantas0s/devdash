// This package is an abstraction for any Terminal UI you want to use.
package internal

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// 8 colors
const (
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

	optionBorderColor = "border_color"
	optionTextColor   = "text_color"
	optionNumColor    = "num_color"

	optionBold = "bold"

	optionFirstColor  = "first_color"
	optionSecondColor = "second_color"

	optionHeight = "height"

	optionBarGap   = "bar_gap"
	optionBarWidth = "bar_width"
	optionBarColor = "bar_color"
)

// map config size to ui size
var sizeLookup = map[string]int{
	"xs":  2,
	"s":   4,
	"m":   6,
	"l":   8,
	"xl":  10,
	"xxl": 12,
}

// map config color to ui color
var colorLookUp = map[string]uint16{
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
		foreground uint16,
		background uint16,
		title string,
		h int,
	)
	BarChart(
		data []int,
		dimensions []string,
		title string,
		bd uint16,
		fg uint16,
		nc uint16,
		height int,
		gap int,
		barWidth int,
		barColor uint16,
	)

	StackedBarChart(
		data [8][]int,
		dimensions []string,
		title string,
		bd uint16,
		fg uint16,
		nc uint16,
		height int,
		gap int,
		barWidth int,
	)

	Table(
		data [][]string,
		bdLabel string,
		bd uint16,
		fg uint16,
	)
	AddCol(size int)
	AddRow()
}

type keyManager interface {
	KQuit(key string)
}

type looper interface {
	Loop()
}

type manager interface {
	keyManager
	renderer
	drawer
	looper
}

// Value objects
type textBoxAttr struct {
	Data       string
	Foreground uint16
	Background uint16
	Title      string
	H          int
}

func (t *Tui) AddCol(size string) error {
	s, err := MapSize(size)
	if err != nil {
		return err
	}
	t.instance.AddCol(s)

	return nil
}

func (t *Tui) AddRow() {
	t.instance.AddRow()
}

func (t *Tui) Render() {
	t.instance.Render()
}

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

func MapSize(size string) (int, error) {
	s := strings.ToLower(size)
	if size, ok := sizeLookup[s]; ok {
		return size, nil
	}

	return 0, errors.Errorf("size %s doesn't exists - please verify your configuration file", s)
}

func (t *Tui) AddProjectTitle(title string, options map[string]string) (err error) {
	size := "XXL"
	if _, ok := options[optionSize]; ok {
		size = options[optionSize]
	}

	textColor := blue
	if _, ok := options[optionTextColor]; ok {
		textColor = colorLookUp[options[optionTextColor]]
	}

	borderColor := blue
	if _, ok := options[optionBorderColor]; ok {
		borderColor = colorLookUp[options[optionBorderColor]]
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

	t.instance.Title(
		title,
		textColor,
		borderColor,
		bold,
		int(height),
		s,
	)

	return nil
}

func (t *Tui) AddTextBox(attr textBoxAttr) {
	t.instance.TextBox(
		attr.Data,
		attr.Foreground,
		attr.Background,
		attr.Title,
		attr.H,
	)
}

func (t *Tui) AddBarChart(
	data []int,
	dimensions []string,
	title string,
	options map[string]string,
) {
	// defaults
	borderColor := blue
	if _, ok := options[optionBorderColor]; ok {
		borderColor = colorLookUp[options[optionBorderColor]]
	}

	textColor := blue
	if _, ok := options[optionTextColor]; ok {
		textColor = colorLookUp[options[optionTextColor]]
	}

	numColor := white
	if _, ok := options[optionNumColor]; ok {
		numColor = colorLookUp[options[optionNumColor]]
	}

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

	var barColor = blue
	if _, ok := options[optionBarColor]; ok {
		barColor = colorLookUp[options[optionBarColor]]
	}

	t.instance.BarChart(
		data,
		dimensions,
		title,
		borderColor,
		textColor,
		numColor,
		int(height),
		int(gap),
		int(barWidth),
		barColor,
	)
}

func (t *Tui) AddStackedBarChart(
	data [8][]int,
	dimensions []string,
	title string,
	options map[string]string,
) {
	// defaults
	borderColor := blue
	if _, ok := options[optionBorderColor]; ok {
		borderColor = colorLookUp[options[optionBorderColor]]
	}

	textColor := blue
	if _, ok := options[optionTextColor]; ok {
		textColor = colorLookUp[options[optionTextColor]]
	}

	numColor := white
	if _, ok := options[optionNumColor]; ok {
		numColor = colorLookUp[options[optionNumColor]]
	}

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

	t.instance.StackedBarChart(
		data,
		dimensions,
		title,
		borderColor,
		textColor,
		numColor,
		int(height),
		int(gap),
		int(barWidth),
	)
}

func (t *Tui) AddTable(data [][]string, title string, options map[string]string) {
	// defaults

	borderColor := blue
	if _, ok := options[optionBorderColor]; ok {
		borderColor = colorLookUp[options[optionBorderColor]]
	}

	textColor := blue
	if _, ok := options[optionTextColor]; ok {
		textColor = colorLookUp[options[optionTextColor]]
	}

	t.instance.Table(
		data,
		title,
		borderColor,
		textColor,
	)
}

func (t *Tui) AddKQuit(key string) {
	t.instance.KQuit(key)
}

func (t *Tui) Loop() {
	t.instance.Loop()
}
