package internal

import (
	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/pkg/errors"
)

const FeedlySubscribers = "feedly.box_subscribers"

type feedlyWidget struct {
	tui    *Tui
	client *platform.Feedly
}

// NewFeedlyWidget with all information necessary to connect to the Feedly API.
func NewFeedlyWidget(address string) *feedlyWidget {
	client := platform.NewFeedly(address)
	return &feedlyWidget{
		client: client,
	}
}

func (f feedlyWidget) CreateWidgets(widget Widget, tui *Tui) (fu func() error, err error) {
	f.tui = tui

	switch widget.Name {
	case FeedlySubscribers:
		fu, err = f.boxSubscribers(widget)
	default:
		return nil, errors.Errorf("can't find the widget %s for service Feedly", widget.Name)
	}

	return
}

func (f feedlyWidget) boxSubscribers(widget Widget) (fu func() error, err error) {
	title := " Feedly subscribers "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	subs, err := f.client.Subscribers()
	if err != nil {
		return nil, err
	}

	fu = func() error {
		return f.tui.AddTextBox(subs, title, widget.Options)
	}

	return
}
