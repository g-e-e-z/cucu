package commands

import (
	"net/http"
	"time"

	"github.com/g-e-e-z/cucu/config"
	"github.com/sirupsen/logrus"
)

// HttpCommand main http interface
type HttpCommand struct {
	Log       *logrus.Entry
	Config    *config.AppConfig
	OSCommand *OSCommand
	Client    *http.Client
}

func NewHttpCommands(log *logrus.Entry, config *config.AppConfig, osCommand *OSCommand) (*HttpCommand, error) {
	command := &HttpCommand{
		Log:       log,
		Config:    config,
		OSCommand: osCommand,
		Client:    &http.Client{Timeout: 5 * time.Second}, // TODO: Configure Client
	}
	return command, nil
}

func (hc *HttpCommand) GetRequests() ([]*Request, error) {
	requests, err := hc.OSCommand.GetRequests()
	if err != nil {
		return nil, err
	}
	for _, request := range requests {
		request.Log = hc.Log
		request.HttpCommand = hc
	}
	return requests, nil
}
