package internal

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Phantas0s/termetrics/internal/plateform"
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

func NewGaWidget(keyfile string, viewID string, tui *Tui) (gaWidget, error) {
	client, err := plateform.NewGaClient(keyfile)
	if err != nil {
		return gaWidget{}, err
	}

	return gaWidget{
		tui:    tui,
		client: client,
		viewID: viewID,
	}, nil
}

func (g gaWidget) createWidgets(widgetName string, widget interface{}) error {
	switch widgetName {
	case realtime:
		g.gaRTActiveUser()
	case users:
		g.gaWeekUsers()
	default:
		return errors.New("can't find the widget " + widgetName)
	}

	return nil
}

// GaRTActiveUser get the real time active users from Google Analytics
func (g gaWidget) gaRTActiveUser() error {
	users, err := g.client.RealTimeUsers(g.viewID)
	if err != nil {
		return err
	}

	g.tui.AddTextBox(textBoxAttr{
		Data:    users,
		Fg:      2,
		Bd:      2,
		Bdlabel: "Real time users: ",
		H:       3,
	})

	return nil
}

// GaWeekUsers get the number of users the 7 last days on your website
func (g gaWidget) gaWeekUsers() error {
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
	})

	return nil
}
