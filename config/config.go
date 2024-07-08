package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	// "runtime"
)
type Config struct {
	General GeneralOptions
	Keys    map[string]map[string]string
}

type GeneralOptions struct {
	ContextSpecificSearch  bool
	DefaultURLScheme       string
	Editor                 string
	FollowRedirects        bool
	FormatJSON             bool
	Insecure               bool
	PreserveScrollPosition bool
	// StatusLine             string
	// TLSVersionMax          uint16
	// TLSVersionMin          uint16
	// Timeout                Duration
}
var DefaultConfig = Config{
	General: GeneralOptions{
		DefaultURLScheme:       "https",
		Editor:                 "vim",
		FollowRedirects:        true,
		FormatJSON:             true,
		Insecure:               false,
		PreserveScrollPosition: true,
		// StatusLine:             "[wuzz {{.Version}}]{{if .Duration}} [Response time: {{.Duration}}]{{end}} [Request no.: {{.RequestNumber}}/{{.HistorySize}}] [Search type: {{.SearchType}}]{{if .DisableRedirect}} [Redirects Restricted Mode {{.DisableRedirect}}]{{end}}",
		// Timeout: Duration{
		// 	defaultTimeoutDuration,
		// },
	},
}
var DefaultKeys = map[string]map[string]string{
	"global": {
		// "CtrlR": "submit",
		"CtrlC": "quit",
		// "CtrlS": "saveResponse",
		// "CtrlF": "loadRequest",
		// "CtrlE": "saveRequest",
		// "CtrlD": "deleteLine",
		// "CtrlW": "deleteWord",
		// "CtrlO": "openEditor",
		// "CtrlT": "toggleContextSpecificSearch",
		// "CtrlX": "clearHistory",
		"Tab":   "nextView",
		"l": "nextView",
		"h": "prevView",
		// "AltH":  "history",
		// "F2":    "focus url",
		// "F3":    "focus get",
		// "F4":    "focus method",
		// "F5":    "focus data",
		// "F6":    "focus headers",
		// "F7":    "focus search",
		// "F8":    "focus response-headers",
		// "F9":    "focus response-body",
		// "F11":   "redirectRestriction",
	},
}

func LoadConfig(configFile string) (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("Config file does not exist.")
	} else if err != nil {
		return nil, err
	}

	conf := DefaultConfig
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		return nil, err
	}

	if conf.Keys == nil {
		conf.Keys = DefaultKeys
	} else {
		// copy default keys
		for keyCategory, keys := range DefaultKeys {
			confKeys, found := conf.Keys[keyCategory]
			if found {
				for key, action := range keys {
					if _, found := confKeys[key]; !found {
						conf.Keys[keyCategory][key] = action
					}
				}
			} else {
				conf.Keys[keyCategory] = keys
			}
		}
	}

	return &conf, nil
}


func GetDefaultConfigLocation() string {
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

	return filepath.Join(configFolderLocation, "config.toml")
}

