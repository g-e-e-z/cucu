package gui

import (
	"github.com/g-e-e-z/cucu/config"
	"github.com/jroimartin/gocui"
	// "github.com/g-e-e-z/cucu/commands"
)

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g      *gocui.Gui
	Config *config.AppConfig
	// Log      *logrus.Entry
	// Commands *commands.Command
	Views Views
}

func NewGuiWrapper(config *config.AppConfig) *Gui {
	return &Gui{
		Config: config,
	}
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g

	if err := gui.SetColorScheme(); err != nil {
		return err
	}

	g.SetManager(gocui.ManagerFunc(gui.layout)) //, gocui.ManagerFunc(gui.getFocusLayout()))

	if err := gui.createAllViews(); err != nil {
		return err
	}

	if err = gui.keybindings(g); err != nil {
		return err
	}

	if gui.g.CurrentView() == nil {
		viewName := gui.initiallyFocusedViewName()
		// view, err := gui.g.View(viewName)
		// if err != nil {
		// 	return err
		// }

		if _, err := gui.g.SetCurrentView(viewName); err != nil {
			return err
		}
	}

    // TODO: This
	// ctx, finish := context.WithCancel(context.Background())
	// defer finish()

	err = gui.g.MainLoop()
	// if err == gocui.ErrQuit {
	// 	return nil
	// }
	return err
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

