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
		Client:    &http.Client{Timeout: 5 * time.Second}, // TODO: Revisit once AppConfig is customizatble
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
        // TODO: This is kinda gross, and might just be a bad approach in general. Aint broke dont fix it
        request.Hash = request.CreateHash()
        request.saved = true
	}
	return requests, nil
}

func (hc *HttpCommand) SaveRequest(r *Request) error {
	requests, err := hc.GetRequests()
	if err != nil {
		return  err
	}
    // Update if its existing
    for i, request := range requests {
        if request.Uuid == r.Uuid {
            r.Modified = false
            r.Hash = r.CreateHash()
            requests[i] = r
            r.saved = true
            return hc.OSCommand.SaveRequests(requests)
        }
    }
    // Append if new
    requests = append(requests, r)
    r.saved = true
    return hc.OSCommand.SaveRequests(requests)
}

func (hc *HttpCommand) DeleteRequest(reqToDelete *Request, requests []*Request) error {
    var newRequests []*Request
    for _, request := range requests {
        if request.Uuid != reqToDelete.Uuid {
            newRequests = append(newRequests, request)
        }
    }
    return hc.OSCommand.SaveRequests(newRequests)

}
