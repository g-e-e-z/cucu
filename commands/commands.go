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

func (hc *HttpCommand) GetRequests() ([]*Request, error) {
    requests, err := hc.OSCommand.GetRequests()
    if err != nil {
        return nil, err
    }
    for _, request := range requests {
        request.HttpCommand = hc
    }
    return requests, nil
}
