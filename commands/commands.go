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
        hc.Log.Error("Error reading requests:", err)
		return nil, err
	}
	for _, request := range requests {
		request.Log = hc.Log
		request.HttpCommand = hc
        // TODO: This is kinda gross, and might just be a bad approach in general. Working for now
        request.createHash()
	}
	return requests, nil
}

func (hc *HttpCommand) SaveRequest(r *Request) error {
	requests, err := hc.GetRequests()
	if err != nil {
		return  err
	}
    for i, request := range requests {
        if request.hash == r.hash {
            r.Modified = false
            r.createHash()
            requests[i] = r
            r.Log.Info("request: commands.go", request.hash)
            r.Log.Info("r: commands.go", r.hash)
            break
        }
    }
    hc.OSCommand.SaveRequests(requests)
    return nil
}

