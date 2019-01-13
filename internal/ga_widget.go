package internal

import (
	"errors"
	"strconv"

	"github.com/Phantas0s/devdash/internal/plateform"
)

const (
	realtime = "ga.realtime"
	users    = "ga.users"
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

func (g *gaWidget) createWidgets(widget Widget, tui *Tui) (err error) {
	g.tui = tui

	switch widget.Name {
	case realtime:
		err = g.gaRTActiveUser(widget)
	case users:
		err = g.gaWeekUsers(widget)
	default:
		return errors.New("can't find the widget " + widget.Name)
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

	err = g.tui.AddTextBox(textBoxAttr{
		Data:    users,
		Fg:      fg,
		Bd:      5,
		Bdlabel: "Real time users: ",
		H:       3,
		Size:    widget.Size,
	})
	if err != nil {
		return err
	}

	return nil
}

// GaWeekUsers get the number of users the 7 last days on your website
func (g *gaWidget) gaWeekUsers(widget Widget) error {
	rep, err := g.client.GetReport(g.viewID)
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
