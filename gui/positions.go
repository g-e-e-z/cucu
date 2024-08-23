package gui

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
)

const (
	REQUESTS_VIEW = "requests"
	URL_VIEW = "url"
	PARAMS_VIEW   = "params"
	RESPONSE_VIEW = "response"

	MENU_VIEW = "menu"
	EDIT_VIEW = "edit"
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
        position{0.0, 0},   // X0
        position{0.0, 0},   // Y0
        position{0.25, -1},  // X1
        position{1.0, -1},  // Y1
	},
	URL_VIEW: {
        position{0.25, 0},   // X0
        position{0.0, 0},   // Y0
        position{1.0, -1},  // X1
        position{0.1, -1},  // Y1
	},
	PARAMS_VIEW: {
        position{0.25, 0},   // X0
        position{0.1, 0},   // Y0
        position{0.6, -1},  // X1
        position{1.0, -1},  // Y1
	},
	RESPONSE_VIEW: {
        position{0.6, 0},   // X0
        position{0.1, 0},   // Y0
        position{1.0, -1},  // X1
        position{1.0, -1},  // Y1
	},
	MENU_VIEW: {
        position{0.25, 0},   // X0
        position{0.25, 0},   // Y0
        position{0.75, -1},  // X1
        position{0.75, -1},  // Y1
	},
	EDIT_VIEW: {
        position{0.25, 0},   // X0
        position{0.25, 0},   // Y0
        position{0.75, -1},  // X1
        position{0.75, -1},  // Y1
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

func (gui *Gui) viewNames() []string {
	visibleViews := lo.Filter(gui.allViews(), func(view *gocui.View, _ int) bool {
		return view.Visible
	})

	return lo.Map(visibleViews, func(view *gocui.View, _ int) string {
		return view.Name()
	})
}

