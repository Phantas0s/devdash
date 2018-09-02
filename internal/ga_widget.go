package internal

import (
	"errors"
	"fmt"
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

func NewGaWidget(keyfile string, viewID string) (gaWidget, error) {
	client, err := plateform.NewGaClient(keyfile)
	if err != nil {
		return gaWidget{}, err
	}

	return gaWidget{
		client: client,
		viewID: viewID,
	}, nil
}

func (g gaWidget) createWidgets(widgetName string, widget Widget, tui *Tui) error {
	g.tui = tui

	switch widgetName {
	case realtime:
		g.gaRTActiveUser(widget)
	case users:
		g.gaWeekUsers(widget)
	default:
		return errors.New("can't find the widget " + widgetName)
	}

	return nil
}

// GaRTActiveUser get the real time active users from Google Analytics
func (g gaWidget) gaRTActiveUser(widget Widget) error {
	users, err := g.client.RealTimeUsers(g.viewID)
	if err != nil {
		return err
	}

	err = g.tui.AddTextBox(textBoxAttr{
		Data:    users,
		Fg:      2,
		Bd:      2,
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
func (g gaWidget) gaWeekUsers(widget Widget) error {
	rep, err := g.client.GetReport(g.viewID)
	if err != nil {
		return err
	}

	// this will extract the different dimensions and data associated
	var dates []string
	var u []int
	for _, v := range rep.Reports {
		for l := 0; l < len(v.Data.Rows); l++ {
			dates = append(dates, v.Data.Rows[l].Dimensions[0]+v.Data.Rows[l].Dimensions[1])
			for m := 0; m < len(v.Data.Rows[l].Metrics); m++ {
				value := v.Data.Rows[l].Metrics[m].Values[0]
				if v, err := strconv.ParseInt(value, 0, 0); err == nil {
					u = append(u, int(v))
				}
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	g.tui.AddBarChart(barChartAttr{
		Data:       u,
		Dimensions: dates,
		BarWidth:   6,
		Size:       widget.Size,
	})

	return nil
}
