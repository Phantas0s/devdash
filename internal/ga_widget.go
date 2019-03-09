package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/pkg/errors"
)

const (
	// widget config names
	realtime      = "ga.realtime"
	sessions      = "ga.sessions"
	users         = "ga.users"
	barMetric     = "ga.bar_metric"
	pages         = "ga.pages"
	newReturning  = "ga.new_returning"
	returning     = "ga.returning"
	trafficSource = "ga.traffic_source"
	totalMetric   = "ga.total_metric"

	gaTimeFormat = "2006-01-02"
)

type gaWidget struct {
	tui       *Tui
	analytics *plateform.Analytics
	viewID    string
}

func NewGaWidget(keyfile string, viewID string) (*gaWidget, error) {
	an, err := plateform.NewAnalyticsClient(keyfile)
	if err != nil {
		return nil, err
	}

	return &gaWidget{
		analytics: an,
		viewID:    viewID,
	}, nil
}

func (g *gaWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	g.tui = tui

	switch widget.Name {
	case realtime:
		err = g.realTimeUser(widget)
	case totalMetric:
		err = g.totalMetric(widget)
	case sessions:
		err = g.barMetric(widget)
	case users:
		err = g.users(widget)
	case barMetric:
		err = g.barMetric(widget)
	case pages:
		err = g.table(widget, "Page")
	case trafficSource:
		err = g.trafficSource(widget)
	case newReturning:
		err = g.NewVsReturning(widget)
	case returning:
		err = g.returningUsers(widget)
	default:
		return errors.Errorf("can't find the widget %s", widget.Name)
	}

	return
}

func (g *gaWidget) totalMetric(widget Widget) (err error) {
	metric := "sessions"
	if _, ok := widget.Options[optionMetric]; ok {
		if len(widget.Options[optionMetric]) > 0 {
			metric = widget.Options[optionMetric]
		}
	}

	startDate := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	endDate := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	sd, ed, err := ConvertDates(time.Now(), startDate, endDate)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Total %s from %s to %s", metric, sd.Format(gaTimeFormat), ed.Format(gaTimeFormat))
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	global := false
	if _, ok := widget.Options[optionGlobal]; ok {
		global, err = strconv.ParseBool(widget.Options[optionGlobal])
		if err != nil {
			return errors.Wrapf(err, "could not parse string %s to bool", widget.Options[optionGlobal])
		}
	}

	users, err := g.analytics.SimpleMetric(g.viewID, metric, sd.Format(gaTimeFormat), ed.Format(gaTimeFormat), global)
	if err != nil {
		return err
	}

	g.tui.AddTextBox(
		users,
		title,
		widget.Options,
	)

	return nil
}

// GaRTActiveUser get the real time active users from Google Analytics
func (g *gaWidget) realTimeUser(widget Widget) error {
	title := " Real time users "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	users, err := g.analytics.RealTimeUsers(g.viewID)
	if err != nil {
		return err
	}

	g.tui.AddTextBox(
		users,
		title,
		widget.Options,
	)

	return nil
}

func (g *gaWidget) users(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options[optionMetric] = "users"
	widget.Options[optionTitle] = " Users "

	return g.barMetric(widget)
}

func (g *gaWidget) returningUsers(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options[optionMetric] = "users"
	widget.Options[optionDimensions] = "user_returning"
	widget.Options[optionTitle] = " Returning users "

	return g.barMetric(widget)
}

// users get the number of users the 7 last days on your website
func (g *gaWidget) barMetric(widget Widget) error {
	// defaults

	sd := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		sd = widget.Options[optionStartDate]
	}

	ed := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		ed = widget.Options[optionEndDate]
	}

	startDate, endDate, err := ConvertDates(time.Now(), sd, ed)
	if err != nil {
		return err
	}

	title := " Sessions "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	metric := "sessions"
	if _, ok := widget.Options[optionMetric]; ok {
		metric = widget.Options[optionMetric]
	}

	dimensions := []string{}
	if _, ok := widget.Options[optionDimensions]; ok {
		if len(widget.Options[optionDimensions]) > 0 {
			dimensions = strings.Split(strings.TrimSpace(widget.Options[optionDimensions]), ",")
		}
	}

	timePeriod := "day"
	if _, ok := widget.Options[optionTimePeriod]; ok {
		timePeriod = strings.TrimSpace(widget.Options[optionTimePeriod])
	}

	dim, val, err := g.analytics.BarMetric(
		g.viewID,
		startDate.Format(gaTimeFormat),
		endDate.Format(gaTimeFormat),
		metric,
		dimensions,
		timePeriod,
	)
	if err != nil {
		return err
	}

	g.tui.AddBarChart(val, dim, title, widget.Options)

	return nil
}

func (g *gaWidget) table(widget Widget, firstHeader string) (err error) {
	// defaults
	var pLen int64 = 20

	global := false
	if _, ok := widget.Options[optionGlobal]; ok {
		global, err = strconv.ParseBool(widget.Options[optionGlobal])
		if err != nil {
			return errors.Wrapf(err, "could not parse string %s to bool", widget.Options[optionGlobal])
		}
	}

	dimension := "page_path"
	if _, ok := widget.Options[optionDimension]; ok {
		if len(widget.Options[optionDimension]) > 0 {
			dimension = widget.Options[optionDimension]
		}
	}

	metrics := []string{"sessions", "page_views", "entrances", "unique_page_views"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(strings.TrimSpace(widget.Options[optionMetrics]), ",")
		}
	}

	orders := []string{metrics[0] + " desc"}
	if _, ok := widget.Options[optionOrder]; ok {
		if len(widget.Options[optionOrder]) > 0 {
			orders = strings.Split(strings.TrimSpace(widget.Options[optionOrder]), ",")
		}
	}

	sd := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		sd = widget.Options[optionStartDate]
	}

	ed := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		ed = widget.Options[optionEndDate]
	}

	startDate, endDate, err := ConvertDates(time.Now(), sd, ed)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("%s from %s to %s", firstHeader, startDate.Format(gaTimeFormat), endDate.Format(gaTimeFormat))
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	headers, dim, val, err := g.analytics.Table(
		g.viewID,
		startDate.Format(gaTimeFormat),
		endDate.Format(gaTimeFormat),
		global,
		metrics,
		dimension,
		orders,
		firstHeader,
	)
	if err != nil {
		return err
	}

	mustContain := ""
	if _, ok := widget.Options[optionMustContain]; ok {
		if len(widget.Options[optionMustContain]) > 0 {
			mustContain = widget.Options[optionMustContain]
		}
	}

	var rowLimit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		rowLimit, err = strconv.ParseInt(widget.Options[optionRowLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}
	if int(rowLimit) > len(dim) {
		rowLimit = int64(len(dim))
	}

	// total of pages + one row for headers
	table := make([][]string, rowLimit+1)
	table[0] = headers

	for i := 0; i < int(rowLimit); i++ {
		if mustContain != "" && !strings.Contains(dim[i], mustContain) {
			continue
		}

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
		table[i+1] = append(table[i+1], val[i]...)
	}

	// Filter out every empty columns because of mustContain conditional
	finalTable := [][]string{}
	for _, v := range table {
		if v != nil {
			finalTable = append(finalTable, v)
		}
	}

	g.tui.AddTable(finalTable, title, widget.Options)

	return nil
}

func (g *gaWidget) trafficSource(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options["dimension"] = "traffic_source"

	return g.table(widget, "Source")
}

func (g *gaWidget) NewVsReturning(widget Widget) error {
	// defaults
	sd := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		sd = widget.Options[optionStartDate]
	}

	ed := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		ed = widget.Options[optionEndDate]
	}

	startDate, endDate, err := ConvertDates(time.Now(), sd, ed)
	if err != nil {
		return err
	}

	metric := "sessions"
	if _, ok := widget.Options[optionMetric]; ok {
		if len(widget.Options[optionMetric]) > 0 {
			metric = widget.Options[optionMetric]
		}
	}

	// this should return new and ret instead of a unique slice val...
	dim, val, err := g.analytics.NewVsReturning(
		g.viewID,
		startDate.Format(gaTimeFormat),
		endDate.Format(gaTimeFormat),
		metric,
	)
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

	firstColor := blue
	if _, ok := widget.Options[optionFirstColor]; ok {
		firstColor = colorLookUp[widget.Options[optionFirstColor]]
	}
	data[firstColor] = new

	secondColor := green
	if _, ok := widget.Options[optionSecondColor]; ok {
		secondColor = colorLookUp[widget.Options[optionSecondColor]]
	}
	data[secondColor] = ret

	title := fmt.Sprintf(" Sessions: Returning (%s) vs New (%s) ", colorStr(firstColor), colorStr(secondColor))
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	g.tui.AddStackedBarChart(data, dim, title, widget.Options)

	return nil
}
