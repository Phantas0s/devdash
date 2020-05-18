package internal

import (
	"strconv"

	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/pkg/errors"
)

const (
	rhUptime = "rh.box_uptime"
)

type remoteHostWidget struct {
	tui     *Tui
	service *platform.RemoteHost
}

func NewRemoteHostWidget(username, addr string) (*remoteHostWidget, error) {
	service, err := platform.NewRemoteHost(username, addr)
	if err != nil {
		return nil, err
	}

	return &remoteHostWidget{
		service: service,
	}, nil
}

func (ms *remoteHostWidget) CreateWidgets(widget Widget, tui *Tui) (f func() error, err error) {
	ms.tui = tui

	switch widget.Name {
	case rhUptime:
		f, err = ms.Uptime(widget)
	default:
		return nil, errors.Errorf("can't find the widget %s", widget.Name)
	}

	return
}

func (ms *remoteHostWidget) Uptime(widget Widget) (f func() error, err error) {
	title := "Uptime"
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	uptime, err := ms.service.Uptime()
	if err != nil {
		return nil, err
	}

	f = func() error {
		return ms.tui.AddTextBox(strconv.FormatInt(uptime, 10), title, widget.Options)
	}

	return
}

// TODO not sure how to do that yet, see remote_host.go
func (ms *remoteHostWidget) GetMemory(widget Widget, tui *Tui) (f func() error, err error) {
	// headers MemTotal, MemFree, Buffers, Cached, SwapTotal, SwapFree
	return func() error {
		return nil
	}, nil
}

// func (ms *monitorServerWidget) table(widget Widget, firstHeader string) (f func() error, err error) {
// }
