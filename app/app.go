package app

import (
	"net/http"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/config"
	"github.com/g-e-e-z/cucu/gui"
	"github.com/g-e-e-z/cucu/log"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config    *config.AppConfig
	Log       *logrus.Entry
	Client    *http.Client
	OSCommand *commands.OSCommand
	Gui       *gui.Gui
}

func NewApp(config *config.AppConfig) (*App, error) {
	app := &App{
		Config: config,
	}
	app.Log = log.NewLogger(config, "23432119147a4367abf7c0de2aa99a2d")
	app.OSCommand = commands.NewOSCommand(app.Log, config)
	app.Client = &http.Client{Timeout: 0}

	app.Gui = gui.NewGuiWrapper(app.Log, config, app.OSCommand, app.Client)

	app.Log.Info("APP INITIALIZED")

	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}

// func (app *App) Close() error {
//     return app.Gui.Close()
// }
