package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/pkg/errors"
)

const (
	// widget config names
	realtime       = "ga.realtime"
	users          = "ga.users"
	pages          = "ga.pages"
	new_returning  = "ga.new_returning"
	traffic_source = "ga.traffic_source"

	// option config names
	optionTitle     = "title"
	optionStartDate = "start_date"
	optionEndDate   = "end_date"
	optionLimit     = "limit"
	optionGlobal    = "global"
	optionMetrics   = "metrics"
	optionCharLimit = "character_limit"
	optionOrder     = "order"
)

type gaWidget struct {
	tui    *Tui
	client *plateform.Client
	viewID string
}

func NewGaWidget(keyfile string, viewID string) (*gaWidget, error) {
	client, err := plateform.NewGaClient(keyfile)
	if err != nil {
		return nil, err
	}

	return &gaWidget{
		client: client,
		viewID: viewID,
	}, nil
}

func (g *gaWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	g.tui = tui

	switch widget.Name {
	case realtime:
		err = g.realTimeUser(widget)
	case users:
		err = g.users(widget)
	case pages:
		err = g.pages(widget)
	case new_returning:
		err = g.ReturningVsNew(widget)
	case traffic_source:
		err = g.trafficSource(widget)
	default:
		return errors.Errorf("can't find the widget %s", widget.Name)
	}

	return
}

// GaRTActiveUser get the real time active users from Google Analytics
func (g *gaWidget) realTimeUser(widget Widget) error {
	title := " Real time users "

	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	users, err := g.client.RealTimeUsers(g.viewID)
	if err != nil {
		return err
	}

	foreground := green
	if users == "0" {
		foreground = red
	}

	g.tui.AddTextBox(textBoxAttr{
		Data:       users,
		Foreground: foreground,
		Background: blue,
		Title:      title,
		H:          3,
	})

	return nil
}

// users get the number of users the 7 last days on your website
func (g *gaWidget) users(widget Widget) error {
	// defaults
	startDate := "7daysAgo"
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	endDate := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	title := fmt.Sprintf(" Users from %s to %s ", startDate, endDate)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	dim, val, err := g.client.Users(g.viewID, startDate, endDate)
	if err != nil {
		return err
	}

	g.tui.AddBarChart(val, dim, title, widget.Options)

	return nil
}

func (g *gaWidget) pages(widget Widget) (err error) {
	// defaults
	var elLimit int64 = 5
	var pLen int64 = 20

	title := "Most page viewed"
	startDate := "7daysAgo"
	endDate := "today"

	global := false

	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	if _, ok := widget.Options[optionLimit]; ok {
		elLimit, err = strconv.ParseInt(widget.Options[optionLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionLimit])
		}
	}

	if _, ok := widget.Options[optionGlobal]; ok {
		global, err = strconv.ParseBool(widget.Options[optionGlobal])
		if err != nil {
			return errors.Wrapf(err, "could not parse string %s to bool", widget.Options[optionGlobal])
		}
	}

	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	metrics := []string{"sessions", "page_views", "entrances", "unique_page_views"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(widget.Options[optionMetrics], ",")
		}
	}

	orders := []string{metrics[0] + " desc"}
	if _, ok := widget.Options[optionOrder]; ok {
		if len(widget.Options[optionOrder]) > 0 {
			orders = strings.Split(widget.Options[optionOrder], ",")
		}
	}

	headers, dim, val, err := g.client.Pages(g.viewID, startDate, endDate, global, metrics, orders)
	if err != nil {
		return err
	}

	if int(elLimit) > len(dim) {
		elLimit = int64(len(dim))
	}

	// total of pages + one row for headers
	table := make([][]string, elLimit+1)
	table[0] = headers

	for i := 0; i < int(elLimit); i++ {
		p := strings.Trim(dim[i], " ")
		if _, ok := widget.Options[optionCharLimit]; ok {
			pLen, err = strconv.ParseInt(widget.Options[optionCharLimit], 0, 0)
		}
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionCharLimit])
		}
		if len(p) > int(pLen) {
			p = p[:pLen]
		}

		// first row after headers
		table[i+1] = []string{p}
		for _, v := range val[i] {
			table[i+1] = append(table[i+1], strconv.Itoa(v))
		}
	}

	g.tui.AddTable(table, title, widget.Options)

	return nil
}

func (g *gaWidget) trafficSource(widget Widget) (err error) {
	// defaults
	var elLimit int64 = 5
	var pLen int64 = 20

	startDate := "7daysAgo"
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	endDate := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	title := fmt.Sprintf(" Traffic from %s to %s ", startDate, endDate)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	if _, ok := widget.Options[optionLimit]; ok {
		elLimit, err = strconv.ParseInt(widget.Options[optionLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionLimit])
		}
	}

	metrics := []string{"sessions", "page_views", "entrances", "unique_page_views"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(widget.Options[optionMetrics], ",")
		}
	}

	orders := []string{metrics[0] + " desc"}
	if _, ok := widget.Options[optionOrder]; ok {
		if len(widget.Options[optionOrder]) > 0 {
			orders = strings.Split(widget.Options[optionOrder], ",")
		}
	}

	global := false
	if _, ok := widget.Options[optionGlobal]; ok {
		global, err = strconv.ParseBool(widget.Options[optionGlobal])
		if err != nil {
			return errors.Wrapf(err, "could not parse string %s to bool", widget.Options[optionGlobal])
		}
	}

	headers, dim, val, err := g.client.TrafficSource(g.viewID, startDate, endDate, global, metrics, orders)
	if err != nil {
		return err
	}

	// total of elements + one row for headers
	table := make([][]string, elLimit+1)
	table[0] = headers

	for i := 0; i < int(elLimit); i++ {
		p := strings.Trim(dim[i], " ")
		if _, ok := widget.Options[optionCharLimit]; ok {
			pLen, err = strconv.ParseInt(widget.Options[optionCharLimit], 0, 0)
		}
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionCharLimit])
		}
		if len(p) > int(pLen) {
			p = p[:pLen]
		}

		// first row after headers
		table[i+1] = []string{p}
		for _, v := range val[i] {
			table[i+1] = append(table[i+1], strconv.Itoa(v))
		}
	}

	g.tui.AddTable(table, title, widget.Options)

	return nil
}

func (g *gaWidget) ReturningVsNew(widget Widget) error {
	// defaults
	startDate := "7daysAgo"
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}
	endDate := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	// this should return new and ret instead of a unique slice val...
	dim, val, err := g.client.ReturningVsNew(g.viewID, startDate, endDate)
	if err != nil {
		return err
	}

	s := len(val) / 2
	ret := val[:s]
	new := val[s:]

	var data [8][]int
	// need to fill data with []int containing 0
	for i := 0; i < 8; i++ {
		for j := 0; j < len(ret); j++ {
			data[i] = append(data[i], 0)
		}
	}

	firstColor := green
	if _, ok := widget.Options[optionFirstColor]; ok {
		firstColor = colorLookUp[widget.Options[optionFirstColor]]
	}
	data[firstColor] = new

	secondColor := yellow
	if _, ok := widget.Options[optionSecondColor]; ok {
		secondColor = colorLookUp[widget.Options[optionSecondColor]]
	}
	data[secondColor] = ret

	title := fmt.Sprintf(
		" Sessions (%s) vs Returning (%s) from %s to %s ",
		colorStr(firstColor),
		colorStr(secondColor),
		startDate,
		endDate,
	)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	g.tui.AddStackedBarChart(data, dim, title, widget.Options)

	return nil
}
