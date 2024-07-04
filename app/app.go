package app

import (
	"github.com/g-e-e-z/cucu/gui"
)

type App struct {
    Gui *gui.Gui
    ErrorChan chan error
}

func NewApp() (*App, error) {
    app := &App{
        ErrorChan: make(chan error),
    }
    var err error
    app.Gui, err = gui.NewGui(app.ErrorChan)
	if err != nil {
		return app, err
	}
    return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}

