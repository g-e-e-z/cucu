package main

import (
	"fmt"
	"log"
	"os"

	"github.com/g-e-e-z/cucu/app"
	"github.com/g-e-e-z/cucu/config"
)

const VERSION = "0.0.1"


func help() {
	fmt.Println(`cucu - Interactive cli tool for http requests

Usage: cucu

Other command line options:
  -c, --config PATH        Specify custom configuration file

Key bindings:
  I'm working on it! :(`,
	)
}

func main() {
	// configDir := ""
	// for i, arg := range os.Args {
	// 	switch arg {
	// 	case "-h", "--help":
	// 		help()
	// 		return
	// 	case "-v", "--version":
	// 		fmt.Printf("cucu %v\n", VERSION)
	// 		return
	// 	case "-c", "--config":
	// 		configDir = os.Args[i+1]
	// 		// args := append(os.Args[:i], os.Args[i+2:]...)
	// 		if _, err := os.Stat(configDir); os.IsNotExist(err) {
	// 			log.Fatal("Config file specified but does not exist: \"" + configDir + "\"")
	// 		}
	// 	}
	// }
    // Come back to this later and handle properly, NewAppConfig operates on configDir not Path
    configDir := ""

    projectDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}


	appConfig, err := config.NewAppConfig(configDir, projectDir)
	if err != nil {
		log.Fatal(err.Error())
	}

    app, err := app.NewApp(appConfig)
    if err == nil {
		err = app.Run()
	}
    // app.Close()

}
