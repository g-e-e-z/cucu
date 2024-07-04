package app

import (
	"github.com/jroimartin/gocui"
)

type App struct {
    Gui *gocui.Gui
}

func NewApp() (*App, error) {
    app := &App{}
    var err error
    app.Gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return app, err
	}
    return app, nil
}

