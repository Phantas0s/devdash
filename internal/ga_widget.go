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
	gaBoxRealtime         = "ga.box_real_time"
	gaBoxTotal            = "ga.box_total"
	gaBar                 = "ga.bar"
	gaBarSessions         = "ga.bar_sessions"
	gaBarBounces          = "ga.bar_bounces"
	gaBarUsers            = "ga.bar_users"
	gaBarReturning        = "ga.bar_returning"
	gaBarNewReturning     = "ga.bar_new_returning"
	gaBarPages            = "ga.bar_pages"
	gaTablePages          = "ga.table_pages"
	gaTableTrafficSources = "ga.table_traffic_sources"

	// format for every start date / end date
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
	case gaBoxRealtime:
		err = g.realTimeUser(widget)
	case gaBoxTotal:
		err = g.totalMetric(widget)
	case gaBarSessions:
		err = g.barMetric(widget)
	case gaBarUsers:
		err = g.users(widget)
	case gaBar:
		err = g.barMetric(widget)
	case gaTablePages:
		err = g.table(widget, "Page")
	case gaTableTrafficSources:
		err = g.trafficSource(widget)
	case gaBarNewReturning:
		err = g.newVsReturning(widget)
	case gaBarReturning:
		err = g.barReturning(widget)
	case gaBarPages:
		err = g.barPages(widget)
	case gaBarBounces:
		err = g.barBounces(widget)
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

	return g.barMetric(widget)
}

func (g *gaWidget) barReturning(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options[optionMetric] = "users"
	widget.Options[optionDimensions] = "user_type"
	widget.Options[optionTitle] = " Returning users "

	return g.barMetric(widget)
}

func (g *gaWidget) barPages(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}
	widget.Options[optionDimensions] = "page_path"
	widget.Options[optionMetric] = "page_views"
	widget.Options[optionTitle] = " Page views "

	if _, ok := widget.Options[optionFilters]; !ok {
		return errors.New("The widget ga.bar_pages require a filter (relative url of your page, i.e '/my-super-page/')")
	}

	widget.Options[optionTitle] += " - filter " + widget.Options[optionFilters] + " "

	return g.barMetric(widget)
}

func (g *gaWidget) barBounces(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}
	widget.Options[optionMetric] = "bounces"
	widget.Options[optionTitle] += " Bounces "

	return g.barMetric(widget)
}

func (g *gaWidget) barMetric(widget Widget) error {
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
		metric = widget.Options[optionMetric]
	}

	dimensions := []string{}
	if _, ok := widget.Options[optionDimensions]; ok {
		if len(widget.Options[optionDimensions]) > 0 {
			dimensions = strings.Split(strings.TrimSpace(widget.Options[optionDimensions]), ",")
		}
	}

	filters := []string{}
	if _, ok := widget.Options[optionFilters]; ok {
		if len(widget.Options[optionFilters]) > 0 {
			filters = strings.Split(strings.TrimSpace(widget.Options[optionFilters]), ",")
		}
	}

	timePeriod := "day"
	if _, ok := widget.Options[optionTimePeriod]; ok {
		timePeriod = strings.TrimSpace(widget.Options[optionTimePeriod])
	}

	title := fmt.Sprintf(" %s per %s ", strings.Title(metric), timePeriod)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	dim, val, err := g.analytics.BarMetric(
		g.viewID,
		startDate.Format(gaTimeFormat),
		endDate.Format(gaTimeFormat),
		metric,
		dimensions,
		timePeriod,
		filters,
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

	filters := []string{}
	if _, ok := widget.Options[optionFilters]; ok {
		if len(widget.Options[optionFilters]) > 0 {
			filters = strings.Split(strings.TrimSpace(widget.Options[optionFilters]), ",")
		}
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
		filters,
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

	// Filter out every empty columns
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

	widget.Options[optionDimension] = "traffic_source"

	return g.table(widget, "Source")
}

func (g *gaWidget) newVsReturning(widget Widget) error {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options[optionDimensions] = "user_type"

	return g.stackedBar(widget)
}

func (g *gaWidget) stackedBar(widget Widget) error {
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

	timePeriod := "day"
	if _, ok := widget.Options[optionTimePeriod]; ok {
		timePeriod = strings.TrimSpace(widget.Options[optionTimePeriod])
	}

	dimensions := []string{}
	if _, ok := widget.Options[optionDimensions]; ok {
		if len(widget.Options[optionDimensions]) > 0 {
			dimensions = strings.Split(strings.TrimSpace(widget.Options[optionDimensions]), ",")
		}
	}

	// this should return new and ret instead of a unique slice val...
	dim, new, ret, err := g.analytics.StackedBar(
		g.viewID,
		startDate.Format(gaTimeFormat),
		endDate.Format(gaTimeFormat),
		metric,
		timePeriod,
		dimensions,
	)
	if err != nil {
		return err
	}

	var data [8][]int
	// need to fill data with []int containing 0
	for i := 0; i < 8; i++ {
		for j := 0; j < len(ret); j++ {
			data[i] = append(data[i], 0)
		}
	}
	data[0] = new
	data[1] = ret

	colors := []uint16{blue, green}
	if _, ok := widget.Options[optionFirstColor]; ok {
		colors[0] = colorLookUp[widget.Options[optionFirstColor]]
	}
	if _, ok := widget.Options[optionSecondColor]; ok {
		colors[1] = colorLookUp[widget.Options[optionSecondColor]]
	}

	title := fmt.Sprintf(
		" %s: Returning (%s) vs New (%s) ",
		strings.Trim(strings.Title(metric), "_"),
		colorStr(colors[0]),
		colorStr(colors[1]),
	)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	g.tui.AddStackedBarChart(data, dim, title, colors, widget.Options)

	return nil
}
