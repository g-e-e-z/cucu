package config

import (
	"os"
	"path/filepath"
)

// AppConfig contains the base configuration fields required for cucu.
type AppConfig struct {
	Name       string `long:"name" env:"NAME" default:"cucu"`
	Debug      bool   `long:"debug" env:"DEBUG" default:"false"`
	ConfigDir  string
	ProjectDir string
}

func getDefaultConfigDir() string {
	var configFolderLocation string
	// switch runtime.GOOS {
	// case "linux":
	// 	// Use the XDG_CONFIG_HOME variable if it is set, otherwise
	// 	// $HOME/.config/wuzz/config.toml
	// 	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	// 	if xdgConfigHome != "" {
	// 		configFolderLocation = xdgConfigHome
	// 	} else {
	// 		configFolderLocation, _ = homedir.Expand("~/.config/wuzz/")
	// 	}
	//
	// default:
	// 	// On other platforms we just use $HOME/.wuzz
	// 	configFolderLocation, _ = homedir.Expand("~/.wuzz/")
	// }
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		configFolderLocation = xdgConfigHome
	} else {
		configFolderLocation, _ = os.UserHomeDir()
	}

	return filepath.Join(configFolderLocation, "cucu")
}
func findOrCreateConfigDir(folder string) (string, error) {
	if folder == "" {
		folder = getDefaultConfigDir()
	}
	err := os.MkdirAll(folder, 0o755)
	if err != nil {
		return "", err
	}

	return folder, nil
}

func NewAppConfig(configDir string, projectDir string) (*AppConfig, error) {
	configDir, err := findOrCreateConfigDir(configDir)
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		Name:       "cucu",
		Debug:      true,
		ConfigDir:  configDir,
		ProjectDir: projectDir,
	}
	return appConfig, nil
}

// ConfigFilename returns the filename of the current config file
func (c *AppConfig) ConfigFilename() string {
	return filepath.Join(c.ConfigDir, "config.yml")
}

// RequestFilename returns the filename of the requests file
func (c *AppConfig) RequestFilename() string {
	return filepath.Join(c.ConfigDir, "requests.json")
}
