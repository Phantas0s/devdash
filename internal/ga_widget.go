package internal

import (
	"strconv"
	"strings"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/pkg/errors"
)

const (
	realtime       = "ga.realtime"
	users          = "ga.users"
	pages          = "ga.pages"
	new_returning  = "ga.new_returning"
	traffic_source = "ga.traffic_source"
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
		err = g.newVsReturning(widget)
	case traffic_source:
		err = g.trafficSource(widget)
	default:
		return errors.Errorf("can't find the widget %s", widget.Name)
	}

	return
}

// GaRTActiveUser get the real time active users from Google Analytics
func (g *gaWidget) realTimeUser(widget Widget) error {
	title := " Real time users : "

	if _, ok := widget.Options["title"]; ok {
		title = widget.Options["title"]
	}

	users, err := g.client.RealTimeUsers(g.viewID)
	if err != nil {
		return err
	}

	foreground := uint16(3)
	if users == "0" {
		foreground = uint16(2)
	}

	g.tui.AddTextBox(textBoxAttr{
		Data:       users,
		Foreground: foreground,
		Background: 5,
		Title:      title,
		H:          3,
	})

	return nil
}

// users get the number of users the 7 last days on your website
func (g *gaWidget) users(widget Widget) error {
	// defaults
	title := "Users"
	startDate := "7daysAgo"
	endDate := "today"

	if _, ok := widget.Options["start_date"]; ok {
		startDate = widget.Options["start_date"]
	}

	if _, ok := widget.Options["end_date"]; ok {
		endDate = widget.Options["end_date"]
	}

	if _, ok := widget.Options["title"]; ok {
		title = widget.Options["title"]
	}

	rep, err := g.client.Users(g.viewID, startDate, endDate)
	if err != nil {
		return err
	}

	// this will extract the different dimensions and data associated
	var dates []string
	var u []int
	dateSeparator := "-"
	for _, v := range rep.Reports {
		for l := 0; l < len(v.Data.Rows); l++ {
			dates = append(dates, v.Data.Rows[l].Dimensions[0]+dateSeparator+v.Data.Rows[l].Dimensions[1])
			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]
				if v, err := strconv.ParseInt(value, 0, 0); err == nil {
					u = append(u, int(v))
				}
				if err != nil {
					return err
				}
			}
		}
	}

	g.tui.AddBarChart(barChartAttr{
		Data:       u,
		Dimensions: dates,
		BarWidth:   6,
		Background: 5,
		Title:      title,
	})

	return nil
}

func (g *gaWidget) pages(widget Widget) (err error) {
	// defaults
	var nbrPages int64 = 5
	var pLen int64 = 20
	title := "Most page viewed"
	global := false
	startDate := "7daysAgo"
	endDate := "today"

	if _, ok := widget.Options["title"]; ok {
		title = widget.Options["title"]
	}

	if _, ok := widget.Options["page_limit"]; ok {
		nbrPages, err = strconv.ParseInt(widget.Options["page_limit"], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options["page_limit"])
		}
	}

	if _, ok := widget.Options["global"]; ok {
		global, err = strconv.ParseBool(widget.Options["global"])
		if err != nil {
			return errors.Wrapf(err, "could not parse string %s to bool", widget.Options["global"])
		}
	}

	if _, ok := widget.Options["start_date"]; ok {
		startDate = widget.Options["start_date"]
	}

	if _, ok := widget.Options["end_date"]; ok {
		endDate = widget.Options["end_date"]
	}

	rep, err := g.client.Pages(g.viewID, startDate, endDate, global)
	if err != nil {
		return err
	}

	var pages []string
	var u []int
	for _, v := range rep.Reports {
		for l := 0; l < len(v.Data.Rows); l++ {
			p := v.Data.Rows[l].Dimensions[0]
			pages = append(pages, p)
			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]
				if v, err := strconv.ParseInt(value, 0, 0); err == nil {
					u = append(u, int(v))
				}
				if err != nil {
					return err
				}
			}
		}
	}

	table := [][]string{{"Page", "Page Views"}}
	for i := 0; i < int(nbrPages); i++ {
		p := strings.Trim(pages[i], " ")
		if _, ok := widget.Options["page_length"]; ok {
			pLen, err = strconv.ParseInt(widget.Options["page_length"], 0, 0)
		}
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options["page_length"])
		}
		if len(p) > int(pLen) {
			p = p[:pLen]
		}
		table = append(table, []string{p, strconv.Itoa(u[i])})
	}

	g.tui.AddTable(tableAttr{
		Data:            table,
		Background:      5,
		BackgroundLabel: title,
	})

	return nil
}

func (g *gaWidget) newVsReturning(widget Widget) error {
	// defaults
	title := "Sessions vs New"
	startDate := "7daysAgo"
	endDate := "today"

	if _, ok := widget.Options["start_date"]; ok {
		startDate = widget.Options["start_date"]
	}

	if _, ok := widget.Options["end_date"]; ok {
		endDate = widget.Options["end_date"]
	}

	if _, ok := widget.Options["title"]; ok {
		title = widget.Options["title"]
	}

	rep, err := g.client.NewVsReturning(g.viewID, startDate, endDate)
	if err != nil {
		return err
	}

	dateSeparator := "-"

	var dates []string
	var u []int
	for _, v := range rep.Reports {
		for l := 0; l < len(v.Data.Rows); l++ {
			dates = append(dates, v.Data.Rows[l].Dimensions[1]+dateSeparator+v.Data.Rows[l].Dimensions[2])
			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]
				if v, err := strconv.ParseInt(value, 0, 0); err == nil {
					u = append(u, int(v))
				}
				if err != nil {
					return err
				}
			}
		}
	}

	s := len(u) / 2
	sessions := u[:s]
	new := u[s:]

	var data [8][]int

	for i := 0; i < 8; i++ {
		for j := 0; j < len(sessions); j++ {
			data[i] = append(data[i], 0)
		}
	}

	data[3] = sessions
	data[4] = new

	g.tui.AddStackedBarChart(stackedBarChartAttr{
		Data:       data,
		Dimensions: dates,
		BarWidth:   6,
		Background: 5,
		Title:      title,
	})

	return nil
}

// users get the number of users the 7 last days on your website
func (g *gaWidget) trafficSource(widget Widget) error {
	// defaults
	var nbrSources int64 = 5
	var pLen int64 = 20
	title := "Traffic"
	startDate := "7daysAgo"
	endDate := "today"

	if _, ok := widget.Options["start_date"]; ok {
		startDate = widget.Options["start_date"]
	}

	if _, ok := widget.Options["end_date"]; ok {
		endDate = widget.Options["end_date"]
	}

	if _, ok := widget.Options["title"]; ok {
		title = widget.Options["title"]
	}

	rep, err := g.client.TrafficSource(g.viewID, startDate, endDate)
	if err != nil {
		return err
	}

	var pages []string
	var u []int
	for _, v := range rep.Reports {
		for l := 0; l < len(v.Data.Rows); l++ {
			p := v.Data.Rows[l].Dimensions[0]
			pages = append(pages, p)
			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]
				if v, err := strconv.ParseInt(value, 0, 0); err == nil {
					u = append(u, int(v))
				}
				if err != nil {
					return err
				}
			}
		}
	}

	table := [][]string{{"Page", "Page Views"}}
	for i := 0; i < int(nbrSources); i++ {
		p := strings.Trim(pages[i], " ")
		if _, ok := widget.Options["page_length"]; ok {
			pLen, err = strconv.ParseInt(widget.Options["page_length"], 0, 0)
		}
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options["page_length"])
		}
		if len(p) > int(pLen) {
			p = p[:pLen]
		}
		table = append(table, []string{p, strconv.Itoa(u[i])})
	}

	g.tui.AddTable(tableAttr{
		Data:            table,
		Background:      5,
		BackgroundLabel: title,
	})
	return nil
}
