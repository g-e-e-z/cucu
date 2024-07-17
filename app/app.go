package app

import (
	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/config"
	"github.com/g-e-e-z/cucu/gui"
)

type App struct {
	Config    *config.AppConfig
	OSCommand *commands.OSCommand
	HttpCommands *commands.HttpCommand
	Gui          *gui.Gui
}

func NewApp(config *config.AppConfig) (*App, error) {
	app := &App{
		Config:       config,
	}
	var err error

	app.OSCommand = commands.NewOSCommand(config)
	app.HttpCommands, err = commands.NewHttpCommands(app.Config, app.OSCommand)
	if err != nil {
		return app, err
	}

	app.Gui = gui.NewGuiWrapper(config, app.OSCommand, app.HttpCommands)

	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}

// func (app *App) Close() error {
//     return app.Gui.Close()
// }
