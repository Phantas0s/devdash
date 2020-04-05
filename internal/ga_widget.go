package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/internal/platform"
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
	gaBarCountry          = "ga.bar_country"
	gaTablePages          = "ga.table_pages"
	gaTableTrafficSources = "ga.table_traffic_sources"
	gaTable               = "ga.table"

	// format for every start date / end date
	gaTimeFormat = "2006-01-02"
)

type gaWidget struct {
	tui       *Tui
	analytics *platform.Analytics
	viewID    string
}

// NewGaWidget including all information to connect to the Google Analytics API.
func NewGaWidget(keyfile string, viewID string) (*gaWidget, error) {
	an, err := platform.NewAnalyticsClient(keyfile)
	if err != nil {
		return nil, err
	}

	return &gaWidget{
		analytics: an,
		viewID:    viewID,
	}, nil
}

// CreateWidgets for Google Analytics.
func (g *gaWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	g.tui = tui

	switch widget.Name {
	case gaBoxRealtime:
		err = g.realTimeUser(widget)
	case gaBoxTotal:
		err = g.totalMetric(widget)
	case gaBarSessions:
		err = g.barMetric(widget, platform.XHeaderTime)
	case gaBarUsers:
		err = g.users(widget)
	case gaBar:
		err = g.barMetric(widget, platform.XHeaderTime)
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
	case gaBarCountry:
		err = g.barCountry(widget)
	case gaBarBounces:
		err = g.barBounces(widget)
	case gaTable:
		err = g.table(widget, widget.Options[optionDimension])
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

	sd, ed, err := platform.ConvertDates(time.Now(), startDate, endDate)
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

	users, err := g.analytics.SimpleMetric(
		platform.AnalyticValues{
			ViewID:    g.viewID,
			StartDate: sd.Format(gaTimeFormat),
			EndDate:   ed.Format(gaTimeFormat),
			Global:    global,
			Metrics:   []string{metric},
		},
	)
	if err != nil {
		return err
	}

	err = g.tui.AddTextBox(users, title, widget.Options)
	if err != nil {
		return err
	}

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

	err = g.tui.AddTextBox(
		users,
		title,
		widget.Options,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *gaWidget) users(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options[optionMetric] = "users"
	xHeader := platform.XHeaderTime

	return g.barMetric(widget, xHeader)
}

func (g *gaWidget) barReturning(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}

	widget.Options[optionMetric] = "users"
	widget.Options[optionDimensions] = "user_type"
	widget.Options[optionTitle] = " Returning users "

	return g.barMetric(widget, platform.XHeaderTime)
}

func (g *gaWidget) barPages(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}
	widget.Options[optionDimensions] = "page_path"
	widget.Options[optionMetric] = "page_views"

	if _, ok := widget.Options[optionFilters]; !ok {
		return errors.New("The widget ga.bar_pages require a filter (relative url of your page, i.e '/my-super-page/')")
	}

	if _, ok := widget.Options[optionTitle]; !ok {
		widget.Options[optionTitle] = widget.Options[optionFilters]
	}

	return g.barMetric(widget, platform.XHeaderTime)
}

func (g *gaWidget) barCountry(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}
	widget.Options[optionDimensions] = "country"
	widget.Options[optionMetric] = "sessions"

	if _, ok := widget.Options[optionTitle]; !ok {
		widget.Options[optionTitle] = widget.Options[optionFilters]
	}

	return g.barMetric(widget, platform.XHeaderOtherDim)
}
func (g *gaWidget) barBounces(widget Widget) (err error) {
	if widget.Options == nil {
		widget.Options = map[string]string{}
	}
	widget.Options[optionMetric] = "bounces"
	widget.Options[optionTitle] += " Bounces "

	return g.barMetric(widget, platform.XHeaderTime)
}

func (g *gaWidget) barMetric(widget Widget, xHeader uint16) (err error) {
	global := false
	if _, ok := widget.Options[optionGlobal]; ok {
		global, err = strconv.ParseBool(widget.Options[optionGlobal])
		if err != nil {
			return errors.Wrapf(err, "could not parse string %s to bool", widget.Options[optionGlobal])
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

	startDate, endDate, err := platform.ConvertDates(time.Now(), sd, ed)
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
		platform.AnalyticValues{
			ViewID:     g.viewID,
			StartDate:  startDate.Format(gaTimeFormat),
			EndDate:    endDate.Format(gaTimeFormat),
			TimePeriod: timePeriod,
			Global:     global,
			Metrics:    []string{metric},
			Dimensions: dimensions,
			Filters:    filters,
			XHeaders:   xHeader,
		},
	)
	if err != nil {
		return err
	}

	g.tui.AddBarChart(val, dim, title, widget.Options)

	return nil
}

func (g *gaWidget) table(widget Widget, firstHeader string) (err error) {
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

	startDate, endDate, err := platform.ConvertDates(time.Now(), sd, ed)
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

	var rowLimit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		rowLimit, err = strconv.ParseInt(widget.Options[optionRowLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	headers, dim, val, err := g.analytics.Table(
		platform.AnalyticValues{
			ViewID:     g.viewID,
			StartDate:  startDate.Format(gaTimeFormat),
			EndDate:    endDate.Format(gaTimeFormat),
			Global:     global,
			Metrics:    metrics,
			Dimensions: []string{dimension},
			Filters:    filters,
			Orders:     orders,
			RowLimit:   rowLimit,
		},
		firstHeader,
	)
	if err != nil {
		return err
	}

	if int(rowLimit) > len(dim) {
		rowLimit = int64(len(dim))
	}

	var charLimit int64 = 20
	if _, ok := widget.Options[optionCharLimit]; ok {
		charLimit, err = strconv.ParseInt(widget.Options[optionCharLimit], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionCharLimit])
		}
	}

	finalTable := formatTable(rowLimit, dim, val, charLimit, headers)
	g.tui.AddTable(finalTable, title, widget.Options)

	return nil
}

func formatTable(
	rowLimit int64,
	dim []string,
	val [][]string,
	charLimit int64,
	headers []string,
) [][]string {
	table := [][]string{headers}

	for k, v := range val {
		if k == int(rowLimit) {
			break
		}

		p := strings.Trim(dim[k], " ")
		if len(p) > int(charLimit) {
			p = p[:charLimit]
		}

		// Add dimension header
		row := []string{p}
		row = append(row, v...)
		table = append(table, row)
	}

	return table
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

	startDate, endDate, err := platform.ConvertDates(time.Now(), sd, ed)
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
		platform.AnalyticValues{
			ViewID:     g.viewID,
			StartDate:  startDate.Format(gaTimeFormat),
			EndDate:    endDate.Format(gaTimeFormat),
			TimePeriod: timePeriod,
			Metrics:    []string{metric},
			Dimensions: dimensions,
		},
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
