package commands

import (
	"net/http"
	"time"

	"github.com/g-e-e-z/cucu/config"
)

// HttpCommand main http interface
type HttpCommand struct {
	Config    *config.AppConfig
	OSCommand *OSCommand
	Client    *http.Client
}

func NewHttpCommands(config *config.AppConfig, osCommand *OSCommand) (*HttpCommand, error) {
	command := &HttpCommand{
		Config:    config,
		OSCommand: osCommand,
        Client:    &http.Client{Timeout: 5*time.Second}, // TODO: Configure Client
	}
	return command, nil
}
