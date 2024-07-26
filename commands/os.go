package commands

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/g-e-e-z/cucu/config"
)

// Platform stores the os state
type Platform struct {
	os              string
	shell           string
	shellArg        string
	openCommand     string
	openLinkCommand string
}

// OSCommand holds all the os commands
type OSCommand struct {
	Platform *Platform
	Config   *config.AppConfig
	command  func(string, ...string) *exec.Cmd
	getenv   func(string) string
}

// NewOSCommand os command runner
func NewOSCommand(config *config.AppConfig) *OSCommand {
	return &OSCommand{
		Config:  config,
		command: exec.Command,
		getenv:  os.Getenv,
	}
}

// FileExists checks whether a file exists at the specified path
func (c *OSCommand) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *OSCommand) GetRequests() ([]*Request, error) {
	file, err := os.Open(c.Config.RequestFilename())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var requests []*Request
	json.Unmarshal(byteValue, &requests)

	return requests, nil
}
