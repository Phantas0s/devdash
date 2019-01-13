package internal

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	goping "github.com/sparrc/go-ping"
)

const (
	ping         = "mon.ping"
	availability = "mon.availability"
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

func (m *monitorWidget) createWidgets(widget Widget, tui *Tui) (err error) {
	m.tui = tui

	switch widget.Name {
	case ping:
		err = m.pingWidget(widget)
	case availability:
		err = m.availabilityWidget(widget)
	default:
		return errors.New("can't find the widget " + widget.Name)
	}

	return
}

func (m *monitorWidget) pingWidget(widget Widget) error {
	url, err := url.Parse(m.address)
	if err != nil {
		return err
	}

	pinger, err := goping.NewPinger(url.Host)
	if err != nil {
		return err
	}
	pinger.Count = 1
	// fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	fg := uint16(3)
	if stats.PacketsRecv == 0 {
		fg = uint16(2)
	}
	err = m.tui.AddTextBox(textBoxAttr{
		Data:    fmt.Sprintf("Sent: %d / Received: %d / Time: %d", stats.PacketsSent, stats.PacketsRecv, stats.AvgRtt),
		Fg:      fg,
		Bd:      5,
		Bdlabel: "Ping:",
		H:       3,
		Size:    widget.Size,
	})
	if err != nil {
		return err
	}

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

	fg := uint16(3)
	if status == "offline" {
		fg = uint16(2)
	}

	err = m.tui.AddTextBox(textBoxAttr{
		Data:    fmt.Sprintf("Status: %s / Status code: %d", status, statusCode),
		Fg:      fg,
		Bd:      5,
		Bdlabel: "Availability:",
		H:       3,
		Size:    widget.Size,
	})
	if err != nil {
		return err
	}

	return nil
}
