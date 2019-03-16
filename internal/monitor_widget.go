package internal

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	goping "github.com/sparrc/go-ping"
)

const (
	boxPing         = "mon.box_ping"
	boxAvailability = "mon.box_availability"
)

type monitorWidget struct {
	tui     *Tui
	address string
}

func NewMonitorWidget(address string) (*monitorWidget, error) {
	return &monitorWidget{
		address: address,
	}, nil
}

func (m *monitorWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	m.tui = tui

	switch widget.Name {
	case boxPing:
		err = m.pingWidget(widget)
	case boxAvailability:
		err = m.availabilityWidget(widget)
	default:
		return errors.New("can't find the widget " + widget.Name)
	}

	return
}

func (m *monitorWidget) pingWidget(widget Widget) error {
	URL, err := url.Parse(m.address)
	if err != nil {
		return err
	}

	pinger, err := goping.NewPinger(URL.Host)
	if err != nil {
		return err
	}
	pinger.Count = 1
	// fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	title := " Availability "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	m.tui.AddTextBox(
		fmt.Sprintf("Sent: %d / Received: %d / Time: %d", stats.PacketsSent, stats.PacketsRecv, stats.AvgRtt),
		title,
		widget.Options,
	)

	return nil
}

func (m *monitorWidget) availabilityWidget(widget Widget) error {
	req, err := http.NewRequest(http.MethodGet, m.address, nil)
	if err != nil {
		return err
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
	}

	title := " Availability "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	m.tui.AddTextBox(
		fmt.Sprintf("%s (%d)", status, statusCode),
		title,
		widget.Options,
	)

	return nil
}
