package internal

import (
	"fmt"
	"net/http"
	"net/url"

	goping "github.com/go-ping/ping"
	"github.com/pkg/errors"
)

const (
	boxPing         = "mon.box_ping"
	boxAvailability = "mon.box_availability"
)

type monitorWidget struct {
	tui     *Tui
	address string
}

// NewMonitorWidget with the address of the website to monitor.
func NewMonitorWidget(address string) (*monitorWidget, error) {
	return &monitorWidget{
		address: address,
	}, nil
}

// CreateWidgets for the monitor service.
func (m *monitorWidget) CreateWidgets(widget Widget, tui *Tui) (f func() error, err error) {
	m.tui = tui

	switch widget.Name {
	case boxPing:
		f, err = m.pingWidget(widget)
	case boxAvailability:
		f, err = m.availabilityWidget(widget)
	default:
		return nil, errors.New("can't find the widget " + widget.Name)
	}

	return
}

func (m *monitorWidget) pingWidget(widget Widget) (f func() error, err error) {
	u := m.address

	if _, ok := widget.Options[optionAddress]; ok {
		u = widget.Options[optionAddress]
	}

	URL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	pinger, err := goping.NewPinger(URL.Host)
	if err != nil {
		return nil, err
	}
	pinger.Count = 1
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	title := " Availability "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	f = func() error {
		return m.tui.AddTextBox(
			fmt.Sprintf("Sent: %d / Received: %d / Time: %d", stats.PacketsSent, stats.PacketsRecv, stats.AvgRtt),
			title,
			widget.Options,
		)
	}

	return
}

func (m *monitorWidget) availabilityWidget(widget Widget) (f func() error, err error) {
	u := m.address
	if _, ok := widget.Options[optionAddress]; ok {
		u = widget.Options[optionAddress]
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	status := "online"
	statusCode := http.StatusOK
	if err != nil || statusCode != res.StatusCode {
		if res == nil {
			statusCode = 0
		} else {
			statusCode = res.StatusCode
		}
		status = "offline"
	} else {
		defer res.Body.Close()
	}

	title := " Availability "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	f = func() error {
		return m.tui.AddTextBox(
			fmt.Sprintf("%s (%d)", status, statusCode),
			title,
			widget.Options,
		)
	}

	return
}
