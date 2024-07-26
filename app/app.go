package app

import (
	"github.com/g-e-e-z/cucu/log"
	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/config"
	"github.com/g-e-e-z/cucu/gui"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config       *config.AppConfig
	Log          *logrus.Entry
	OSCommand    *commands.OSCommand
	HttpCommands *commands.HttpCommand
	Gui          *gui.Gui
}

func NewApp(config *config.AppConfig) (*App, error) {
	app := &App{
		Config: config,
	}
	var err error
	app.Log = log.NewLogger(config, "23432119147a4367abf7c0de2aa99a2d")
	app.OSCommand = commands.NewOSCommand(config)
	app.HttpCommands, err = commands.NewHttpCommands(app.Log, app.Config, app.OSCommand)
	if err != nil {
		return app, err
	}

	app.Gui = gui.NewGuiWrapper(app.Log, config, app.OSCommand, app.HttpCommands)

    app.Log.Info("APP INITIALIZED")

	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}

// func (app *App) Close() error {
//     return app.Gui.Close()
// }
