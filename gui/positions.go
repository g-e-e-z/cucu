package gui

import (
	"log"
)

const (
	REQUESTS_VIEW = "requests"
	PARAMS_VIEW   = "params"
	RESPONSE_VIEW = "response"
)

type viewPosition struct {
	X0, Y0, X1, Y1 int
}

var VIEW_POSITIONS = map[string]viewPosition{
	REQUESTS_VIEW: {
		1, 1, 10, 10,
	},
	PARAMS_VIEW: {
		10, 1, 20, 10,
	},
	RESPONSE_VIEW: {
		20, 1, 30, 10,
	},
}

func (gui *Gui) getWindowDimensions() map[string]viewPosition {
	minimumHeight := 9
	minimumWidth := 10
	width, height := gui.g.Size()
	if width < minimumWidth || height < minimumHeight {
		log.Panic("Terminal is too small")
	}

	return VIEW_POSITIONS
}
