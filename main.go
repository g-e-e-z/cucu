package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/g-e-e-z/cucu/config"
	"github.com/jroimartin/gocui"
)

const VERSION = "0.0.1"

type App struct {
	viewIndex    int
	reqIndex     int
	currentPopup string
	config       *config.Config
	requests     []*Request
}

type Request struct {
	Name   string
	Url    string
	Method string
	// GetParams       string
	// Data            string
	// Headers         string
	// ResponseHeaders string
	// RawResponseBody []byte
	// ContentType     string
	// Duration        time.Duration
	// Formatter       formatter.ResponseFormatter
}

func (a *App) setView(g *gocui.Gui) error {
	// a.closePopup(g, a.currentPopup)
	_, err := g.SetCurrentView(VIEWS[a.viewIndex])
	return err
}

type ViewEditor struct {
	app           *App
	g             *gocui.Gui
	backTabEscape bool
	origEditor    gocui.Editor
}

var defaultEditor ViewEditor

const (
	MIN_WIDTH  = 60
	MIN_HEIGHT = 20
)

const (
	ALL_VIEWS = ""

	REQUESTS_VIEW = "requests"
	PARAMS_VIEW   = "request_params"
	RESPONSE_VIEW = "response"

	ERROR_VIEW = "error_view"
)

var VIEWS = []string{
	REQUESTS_VIEW,
	PARAMS_VIEW,
	RESPONSE_VIEW,
}

type position struct {
	pct float32
	abs int
}

type viewPosition struct {
	x0, y0, x1, y1 position
}

var VIEW_POSITIONS = map[string]viewPosition{
	REQUESTS_VIEW: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.2, -1},
		position{1.0, -2},
	},
	PARAMS_VIEW: {
		position{0.2, 0},
		position{0.0, 0},
		position{0.6, -1},
		position{1.0, -2},
	},
	RESPONSE_VIEW: {
		position{0.6, 0},
		position{0.0, 0},
		position{1.0, -2},
		position{1.0, -2},
	},
	ERROR_VIEW: {
		position{0.0, 0},
		position{0.0, 0},
		position{1.0, -1},
		position{1.0, -1}},
}

type viewProperties struct {
	title     string
	frame     bool
	editable  bool
	wrap      bool
	highlight bool
	editor    gocui.Editor
	text      string
}

var VIEW_PROPERTIES = map[string]viewProperties{
	REQUESTS_VIEW: {
		title:     "Requests",
		frame:     true,
		editable:  true,
		wrap:      false,
		highlight: true,
		editor:    defaultEditor.origEditor,
	},
	PARAMS_VIEW: {
		title:     "Params", // When title is Request Params, the space is an m???
		frame:     true,
		editable:  true,
		wrap:      false,
		highlight: true,
		editor:    nil,
	},
	RESPONSE_VIEW: {
		title:     "Response",
		frame:     true,
		editable:  true,
		wrap:      false,
		highlight: true,
		editor:    nil,
	},
}

func initApp(a *App, g *gocui.Gui) {
	g.Cursor = true
	g.InputEsc = false
	g.BgColor = gocui.ColorDefault
	g.FgColor = gocui.ColorDefault
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(a.Layout)
}

func getViewValue(g *gocui.Gui, name string) string {
	v, err := g.View(name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v.Buffer())
}

func (p position) getCoordinate(max int) int {
	return int(p.pct*float32(max)) + p.abs
}

func setView(g *gocui.Gui, viewName string) (*gocui.View, error) {
	maxX, maxY := g.Size()
	position := VIEW_POSITIONS[viewName]
	return g.SetView(viewName,
		position.x0.getCoordinate(maxX+1),
		position.y0.getCoordinate(maxY+1),
		position.x1.getCoordinate(maxX+1),
		position.y1.getCoordinate(maxY+1))
}

func setViewProperties(v *gocui.View, name string) {
	v.Title = VIEW_PROPERTIES[name].title
	v.Frame = VIEW_PROPERTIES[name].frame
	v.Editable = VIEW_PROPERTIES[name].editable
	v.Wrap = VIEW_PROPERTIES[name].wrap
	v.Editor = VIEW_PROPERTIES[name].editor
	setViewTextAndCursor(v, VIEW_PROPERTIES[name].text)
}

func setViewTextAndCursor(v *gocui.View, s string) {
	v.Clear()
	fmt.Fprint(v, s)
	v.SetCursor(len(s), 0)
}

func setViewDefaults(v *gocui.View) {
	v.Frame = true
	v.Wrap = false
}

func (a *App) LoadConfig(configPath string) error {
	if configPath == "" {
		// Load config from default path
		configPath = config.GetDefaultConfigLocation()
	}

	// If the config file doesn't exist, load the default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		a.config = &config.DefaultConfig
		a.config.Keys = config.DefaultKeys
		// a.statusLine, _ = NewStatusLine(a.config.General.StatusLine)
		return nil
	}

	// XXX: Add conf back when creating non default config
	_, err := config.LoadConfig(configPath)
	if err != nil {
		a.config = &config.DefaultConfig
		a.config.Keys = config.DefaultKeys
		return err
	}

	// a.config = conf
	// sl, err := NewStatusLine(conf.General.StatusLine)
	// if err != nil {
	// 	a.config = &config.DefaultConfig
	// 	a.config.Keys = config.DefaultKeys
	// 	return err
	// }
	// a.statusLine = sl
	return nil
}

func (a *App) LoadRequests(g *gocui.Gui) error {
    v, err := g.View(REQUESTS_VIEW)
    // Handle Error Better
	if err != nil {
		rv, _ := g.View(ERROR_VIEW)
		rv.Clear()
		fmt.Fprintf(rv, "Editor open error: %v", err)
		return nil
	}
	a.requests = []*Request{
		{
			Name: "Vinyasa",
            Url:  "localhost:3000/health",
			Method: "GET",
		},
		{
			Name: "Kanye",
			Url:  "https://api.kanye.rest/",
			Method: "GET",
		},
		{
            Name: "Post Echo",
			Url:  "https://httpbin.org/post",
			Method: "POST",
		},
	}
    v.Clear()
    for _, rq := range a.requests {
        fmt.Fprintln(v, rq.Name)}
	return nil
}

func (a *App) SetKeys(g *gocui.Gui) error {
	// load config keybindings
	for viewName, keys := range a.config.Keys {
		if viewName == "global" {
			viewName = ALL_VIEWS
		}
		for keyStr, commandStr := range keys {
			if err := a.setKey(g, keyStr, commandStr, viewName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) NextView(g *gocui.Gui, v *gocui.View) error {
	a.viewIndex = (a.viewIndex + 1) % len(VIEWS)
	return a.setView(g)
}

func (a *App) PrevView(g *gocui.Gui, v *gocui.View) error {
	a.viewIndex = (a.viewIndex - 1 + len(VIEWS)) % len(VIEWS)
	return a.setView(g)
}

func parseKey(k string) (interface{}, gocui.Modifier, error) {
	mod := gocui.ModNone
	if strings.Index(k, "Alt") == 0 {
		mod = gocui.ModAlt
		k = k[3:]
	}
	switch len(k) {
	case 0:
		return 0, 0, errors.New("Empty key string")
	case 1:
		if mod != gocui.ModNone {
			k = strings.ToLower(k)
		}
		return rune(k[0]), mod, nil
	}

	key, found := KEYS[k]
	if !found {
		return 0, 0, fmt.Errorf("Unknown key: %v", k)
	}
	return key, mod, nil
}

func (a *App) setKey(g *gocui.Gui, keyStr, commandStr, viewName string) error {
	if commandStr == "" {
		return nil
	}
	key, mod, err := parseKey(keyStr)
	if err != nil {
		return err
	}
	commandParts := strings.SplitN(commandStr, " ", 2)
	command := commandParts[0]
	var commandArgs string
	if len(commandParts) == 2 {
		commandArgs = commandParts[1]
	}
	keyFnGen, found := COMMANDS[command]
	if !found {
		return fmt.Errorf("Unknown command: %v", command)
	}
	keyFn := keyFnGen(commandArgs, a)
	if err := g.SetKeybinding(viewName, key, mod, keyFn); err != nil {
		return fmt.Errorf("Failed to set key '%v': %v", keyStr, err)
	}
	return nil
}

func (a *App) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Highlight = true

	if maxX < MIN_WIDTH || maxY < MIN_HEIGHT {
		if v, err := setView(g, ERROR_VIEW); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			setViewDefaults(v)
			v.Title = "Error"
			g.Cursor = false
			fmt.Fprintln(v, "Terminal is too small")
		}
		return nil
	}
	if _, err := g.View(ERROR_VIEW); err == nil {
		g.DeleteView(ERROR_VIEW)
		g.Cursor = true
		a.setView(g)
	}

	// for _, name := range []string{RESPONSE_HEADERS_VIEW, RESPONSE_BODY_VIEW} {
	// 	vp := VIEW_PROPERTIES[name]
	// 	vp.editor = a.getResponseViewEditor(g)
	// 	VIEW_PROPERTIES[name] = vp
	// }

	// if a.config.General.DefaultURLScheme != "" && !strings.HasSuffix(a.config.General.DefaultURLScheme, "://") {
	// 	p := VIEW_PROPERTIES[URL_VIEW]
	// 	p.text = a.config.General.DefaultURLScheme + "://"
	// 	VIEW_PROPERTIES[URL_VIEW] = p
	// }

	for _, name := range []string{
		REQUESTS_VIEW,
		PARAMS_VIEW,
		RESPONSE_VIEW,
	} {
		if v, err := setView(g, name); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			setViewProperties(v, name)
		}
	}
	g.SetCurrentView(VIEWS[a.viewIndex])
	// refreshStatusLine(a, g)

	a.LoadRequests(g) // Working on populating the requests

	return nil
}

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
	configPath := ""
	for i, arg := range os.Args {
		switch arg {
		case "-h", "--help":
			help()
			return
		case "-v", "--version":
			fmt.Printf("cucu %v\n", VERSION)
			return
		case "-c", "--config":
			configPath = os.Args[i+1]
			// args := append(os.Args[:i], os.Args[i+2:]...)
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				log.Fatal("Config file specified but does not exist: \"" + configPath + "\"")
			}
		}
	}

	var g *gocui.Gui
	var err error
	for _, outputMode := range []gocui.OutputMode{gocui.Output256, gocui.OutputNormal} {
		g, err = gocui.NewGui(outputMode)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Panicln(err)
	}
	app := &App{
		viewIndex: 0,
		reqIndex:  0,
		requests:  make([]*Request, 0, 31),
	}

	// overwrite default editor
	defaultEditor = ViewEditor{app, g, false, gocui.DefaultEditor}

	initApp(app, g)

	// load config (must be done *before* app.ParseArgs, as arguments
	// should be able to override config values). An empty string passed
	// to LoadConfig results in LoadConfig loading the default config
	// location. If there is no config, the values in
	// config.DefaultConfig will be used.
	err = app.LoadConfig(configPath)
	if err != nil {
		g.Close()
		log.Fatalf("Error loading config file: %v", err)
	}

	err = app.SetKeys(g)
	if err != nil {
		g.Close()
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	defer g.Close()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
