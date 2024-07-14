package gui

import (
	"log"
)

const (
	REQUESTS_VIEW = "requests"
	PARAMS_VIEW   = "params"
	RESPONSE_VIEW = "response"
)

type position struct {
	pct float32
    abs int
}

func (p position) GetCoordinate(max int) int {
	return int(p.pct*float32(max)) + p.abs
}


type viewPosition struct {
	X0, Y0, X1, Y1 position
}

var VIEW_POSITIONS = map[string]viewPosition{
	REQUESTS_VIEW: {
        position{0.0, 0},
        position{0.0, 0},
        position{0.2, -1},
        position{1.0, -1},
	},
	PARAMS_VIEW: {
        position{0.2, 0},
        position{0.0, 0},
        position{0.6, -1},
        position{1.0, -1},
	},
	RESPONSE_VIEW: {
        position{0.6, 0},
        position{0.0, 0},
        position{1.0, -1},
        position{1.0, -1},
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
