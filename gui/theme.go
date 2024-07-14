package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) SetColorScheme() error {
	gui.g.Cursor = true
	gui.g.InputEsc = false
	gui.g.BgColor = gocui.ColorDefault
	gui.g.FgColor = gocui.ColorDefault
	gui.g.SelFgColor = gocui.ColorGreen
	return nil
}

