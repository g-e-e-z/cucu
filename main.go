package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

type App struct {
	viewIndex    int
	currentPopup string
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
	PARAMS_VIEW   = "params"
	RESPONSE_VIEW = "response"

	ERROR_VIEW = "error_view"
)

var VIEW_TITLES = map[string]string{
	REQUESTS_VIEW: "Requests",
	PARAMS_VIEW:   "Request Params",
	RESPONSE_VIEW: "Response",

	ERROR_VIEW: "Error",
}
var VIEWS = []string{
	REQUESTS_VIEW,
	PARAMS_VIEW,
	RESPONSE_VIEW,
	ERROR_VIEW,
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
	title    string
	frame    bool
	editable bool
	wrap     bool
	editor   gocui.Editor
	text     string
}

var VIEW_PROPERTIES = map[string]viewProperties{
	REQUESTS_VIEW: {
		title:    "Requests",
		frame:    true,
		editable: true,
		wrap:     false,
		editor:   defaultEditor.origEditor,
	},
	PARAMS_VIEW: {
		title:    "Request Params",
		frame:    true,
		editable: true,
		wrap:     false,
		editor:   nil,
	},
	RESPONSE_VIEW: {
		title:    "Response",
		frame:    true,
		editable: true,
		wrap:     false,
		editor:   nil,
	},
}

func initApp(a *App, g *gocui.Gui) {
	g.Cursor = true
	g.InputEsc = false
	g.BgColor = gocui.ColorDefault
	g.FgColor = gocui.ColorDefault
	g.SetManagerFunc(a.Layout)
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

func (a *App) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if maxX < MIN_WIDTH || maxY < MIN_HEIGHT {
		if v, err := setView(g, ERROR_VIEW); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			setViewDefaults(v)
			v.Title = VIEW_TITLES[ERROR_VIEW]
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
	// refreshStatusLine(a, g)

	return nil
}

func main() {
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
	app := &App{}

	// overwrite default editor
	defaultEditor = ViewEditor{app, g, false, gocui.DefaultEditor}

	initApp(app, g)

	// load config (must be done *before* app.ParseArgs, as arguments
	// should be able to override config values). An empty string passed
	// to LoadConfig results in LoadConfig loading the default config
	// location. If there is no config, the values in
	// config.DefaultConfig will be used.
	// err = app.LoadConfig(configPath)
	// if err != nil {
	// 	g.Close()
	// 	log.Fatalf("Error loading config file: %v", err)
	// }
	//
	// err = app.ParseArgs(g, args)

	// Some of the values in the config need to have some startup
	// behavior associated with them. This is run after ParseArgs so
	// that command-line arguments can override configuration values.
	// app.InitConfig()

	if err != nil {
		g.Close()
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	// err = app.SetKeys(g)

	if err != nil {
		g.Close()
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	defer g.Close()

    // Temporary Keybinding to exit
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
