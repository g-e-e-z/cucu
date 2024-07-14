package gui

import (
	"math"

	// "github.com/jroimartin/gocui"
)

const (
    ScrollHeight = 2
)

func (gui *Gui) scrollUpRequests() error {
	requestsView := gui.Views.Requests
	requestsView.Autoscroll = false
	ox, oy := requestsView.Origin()
	newOy := int(math.Max(0, float64(oy-ScrollHeight)))
	return requestsView.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownRequests() error {
	requestsView := gui.Views.Requests
	requestsView.Autoscroll = false
	ox, oy := requestsView.Origin()

	reservedLines := 0
    _, sizeY := requestsView.Size()
    reservedLines = sizeY

	totalLines := len(requestsView.BufferLines()) // Probably Wrong
	if oy+reservedLines >= totalLines {
		return nil
	}

	return requestsView.SetOrigin(ox, oy+ScrollHeight)
}

