package app

import (
	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/config"
	"github.com/g-e-e-z/cucu/gui"
)

type App struct {
	Config   *config.AppConfig
	Commands *commands.Command
	Gui      *gui.Gui
}

func NewApp(config *config.AppConfig) (*App, error) {
	app := &App{
		Config:   config,
		Commands: &commands.Command{},
	}
	// Curl/OS Commands?

	// Gui
	app.Gui = gui.NewGuiWrapper(config)

	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}

// func (app *App) Close() error {
//     return app.Gui.Close()
// }
