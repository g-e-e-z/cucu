package commands

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"

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
    var exists bool
    var err error
    if exists, err = c.FileExists(c.Config.RequestFilename()); err != nil {
        return nil, err
    }
    if !exists {
        c.InitRequests()
    }
	file, err := os.Open(c.Config.RequestFilename())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var requests []*Request
    err = json.Unmarshal(byteValue, &requests)
    if err != nil {
        return nil, err
    }

	return requests, nil
}

func (c *OSCommand) SaveRequests(requests []*Request) error {
    tmpName := filepath.Join(c.Config.ConfigDir, "old_requests.json")
    os.Rename(c.Config.RequestFilename(), tmpName)
    defer os.Remove(tmpName)

	outFile, err := os.Create(c.Config.RequestFilename())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write the opening bracket for the JSON array
	_, err = outFile.WriteString("[\n")
	if err != nil {
		return err
	}

	// Iterate over the requests and write them as JSON
	for i, request := range requests {
		jsonString:= request.toJSON()

		// Add a comma if this isn't the last request
		if i != 0 {
			_, err = outFile.WriteString(",\n")
			if err != nil {
				return err
			}
		}

		_, err = outFile.WriteString(jsonString)
		if err != nil {
			return err
		}
	}

	// Write the closing bracket for the JSON array
	_, err = outFile.WriteString("\n]\n")
	if err != nil {
		return err
	}

	return nil
}

func (c *OSCommand) InitRequests() error {
    jsonString := []byte(`[{"uuid":"123123","name":"Kanye Rest","url":"https://api.kanye.rest/","method":"GET"}]`)
    err := os.WriteFile(c.Config.RequestFilename(), jsonString, os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}
