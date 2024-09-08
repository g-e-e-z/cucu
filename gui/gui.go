package gui

import (
	"net/http"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/config"
	"github.com/g-e-e-z/cucu/gui/components"
	"github.com/g-e-e-z/cucu/gui/types"
	"github.com/jesseduffield/gocui"
	"github.com/sirupsen/logrus"
)

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g          *gocui.Gui
	Config     *config.AppConfig
	Log        *logrus.Entry
	OSCommands *commands.OSCommand
	Client     *http.Client
	Views      Views

	RequestContext *RequestContext
	Requests       *components.ListComponent[*commands.Request]
	Menu *components.ListComponent[*types.MenuItem]
	// ReqInfo *components.ReqInfo
	// Response *components.Response
}

func NewGuiWrapper(log *logrus.Entry, config *config.AppConfig, osCommands *commands.OSCommand, client *http.Client) *Gui {
	return &Gui{
		Config:     config,
		Log:        log,
		OSCommands: osCommands,
		Client:     client,
	}
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:       gocui.OutputTrue,
		RuneReplacements: map[rune]string{},
	})
	if err != nil {
		return err
	}

	defer g.Close()

	gui.g = g

	if err := gui.SetColorScheme(); err != nil {
		return err
	}

	g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.getFocusLayout()))

	if err := gui.createAllViews(); err != nil {
		return err
	}

	gui.createPanels()

	if err = gui.keybindings(g); err != nil {
		return err
	}

	if gui.g.CurrentView() == nil {
		viewName := gui.initiallyFocusedViewName()

		if _, err := gui.g.SetCurrentView(viewName); err != nil {
			return err
		}
	}

	// Populate Requests Panel
	gui.loadRequests()

	err = gui.g.MainLoop()
	if err == gocui.ErrQuit {
		return nil
	}
	return err
}

func (gui *Gui) createPanels() {
	gui.Requests = gui.getRequestsPanel()
	gui.Menu =     gui.getMenuPanel()
}

func (gui *Gui) Update(f func() error) {
	gui.g.Update(func(*gocui.Gui) error { return f() })
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (gui *Gui) scrollView(v *gocui.View, dy int) error {
	v.Autoscroll = false
	ox, oy := v.Origin()
	if oy+dy < 0 {
		dy = -oy
	}
	if _, err := v.Line(dy); dy > 0 && err != nil {
		dy = 0
	}
	v.SetOrigin(ox, oy+dy)
	return nil
}

func (gui *Gui) scrollViewUp(_ *gocui.Gui, v *gocui.View) error {
	return gui.scrollView(v, -1)
}

func (gui *Gui) scrollViewDown(_ *gocui.Gui, v *gocui.View) error {
	return gui.scrollView(v, 1)
}

func (gui *Gui) initiallyFocusedViewName() string {
	return "requests"
}

// This works but I dont like it :(
func (gui *Gui) handleToggleEdit(_ *gocui.Gui, v *gocui.View) error {
	v.Editable = !v.Editable
	editMessage := " | EDIT"
	if v.Editable {
		v.Title += editMessage
	} else {
		v.Title = v.Title[:len(editMessage)]
	}
	return nil
}
