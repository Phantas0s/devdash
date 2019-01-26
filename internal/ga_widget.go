package internal

import (
	"strconv"
	"strings"

	"github.com/Phantas0s/devdash/internal/plateform"
	"github.com/pkg/errors"
)

const (
	realtime = "ga.realtime"
	users    = "ga.users"
	pages    = "ga.pages"
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
		err = g.gaRTActiveUser(widget)
	case users:
		err = g.gaUsers(widget)
	case pages:
		err = g.getFirstPages(widget, 3)
	default:
		return errors.Errorf("can't find the widget %s", widget.Name)
	}

	return
}

// GaRTActiveUser get the real time active users from Google Analytics
func (g *gaWidget) gaRTActiveUser(widget Widget) error {
	users, err := g.client.RealTimeUsers(g.viewID)
	if err != nil {
		return err
	}

	fg := uint16(3)
	if users == "0" {
		fg = uint16(2)
	}

	g.tui.AddTextBox(textBoxAttr{
		Data:    users,
		Fg:      fg,
		Bd:      5,
		Bdlabel: "Real time users: ",
		H:       3,
		Size:    widget.Size,
	})

	return nil
}

// gaUsers get the number of users the 7 last days on your website
func (g *gaWidget) gaUsers(widget Widget) error {
	startDate := "7daysAgo"
	if _, ok := widget.Options["start_date"]; ok {
		startDate = widget.Options["start_date"]
	}

	endDate := "today"
	if _, ok := widget.Options["end_date"]; ok {
		endDate = widget.Options["end_date"]
	}

	rep, err := g.client.GetUserReport(g.viewID, startDate, endDate)
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
		Bdlabel:    "Weekly users",
		Size:       widget.Size,
	})

	return nil
}

func (g *gaWidget) getFirstPages(widget Widget, count int) error {
	rep, err := g.client.GetTopContent(g.viewID)
	if err != nil {
		return err
	}

	var nbrPages int64 = 5
	if _, ok := widget.Options["page_limit"]; ok {
		nbrPages, err = strconv.ParseInt(widget.Options["page_limit"], 0, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options["page_limit"])
		}
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
	var pLen int64 = 20
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
		Data:    table,
		BdLabel: "Most page viewed",
	})

	return nil
}
