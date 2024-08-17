package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{}
	Modifier    gocui.Modifier
	Description string
}

// GetKey is a function.
func (b *Binding) GetKey() string {
	key := 0

	switch b.Key.(type) {
	case rune:
		key = int(b.Key.(rune))
	case gocui.Key:
		key = int(b.Key.(gocui.Key))
	}

	// special keys
	switch key {
	case 27:
		return "esc"
	case 13:
		return "enter"
	case 32:
		return "space"
	case 65514:
		return "►"
	case 65515:
		return "◄"
	case 65517:
		return "▲"
	case 65516:
		return "▼"
	case 65508:
		return "PgUp"
	case 65507:
		return "PgDn"
	}

	return fmt.Sprintf("%c", key)
}

func (gui *Gui) GetInitialKeybindings() []*Binding {
	bindings := []*Binding{
		{
			ViewName: "",
			Key:      gocui.KeyCtrlC,
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlN,
			Modifier: gocui.ModNone,
			Handler:  gui.handleNewRequest,
		},
		{
			ViewName: "requests",
			Key:      gocui.KeyCtrlR,
			Modifier: gocui.ModNone,
			Handler:  gui.handleRequestSend,
		},
		{
			ViewName: "requests",
			Key:      'r',
			Modifier: gocui.ModNone,
			Handler:  gui.handleEditName,
		},
		{
			ViewName: "requests",
			Key:      'm',
			Modifier: gocui.ModNone,
			Handler:  gui.handleEditMethod,
		},
		{
			ViewName: "requests",
			Key:      gocui.KeyCtrlS,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSaveRequest,
		},
		{
			ViewName: "requests",
			Key:      '[',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.Components.Requests.HandleNextTab),
		},
		{
			ViewName: "requests",
			Key:      ']',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.Components.Requests.HandlePrevTab),
		},
		{
			ViewName: "params",
			Key:      '[',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.Components.Requests.HandleNextTab),
		},
		{
			ViewName: "params",
			Key:      ']',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.Components.Requests.HandlePrevTab),
		},
		{
			ViewName: "menu",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuClose),
		},
		{
			ViewName: "menu",
			Key:      gocui.KeySpace,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuPress),
		},
		{
			ViewName: "menu",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuPress),
		},
		{
			ViewName: "menu",
			Key:      'y',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuPress),
		},
		{
			ViewName: "url",
			Key:      gocui.KeyCtrlR,
			Modifier: gocui.ModNone,
			Handler:  gui.handleRequestSend,
		},
		{
			ViewName: "url",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleToggleEdit, // Find a good place for this: Will be applicable to several views
		},
		{
			ViewName: "url",
			Key:      gocui.KeyCtrlE,
			Modifier: gocui.ModNone,
			Handler:  gui.handleToggleEdit, // Find a good place for this: Will be applicable to several views
		},
	}

	for _, view := range gui.allViews() {
		bindings = append(bindings, []*Binding{
			// TODO: Revist once editor is figured out
			{ViewName: view.Name(), Key: gocui.KeyCtrlH, Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: view.Name(), Key: gocui.KeyCtrlL, Modifier: gocui.ModNone, Handler: gui.nextView},
		}...)
	}

	setUpDownClickBindings := func(viewName string, onUp func() error, onDown func() error) {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: wrappedHandler(onUp)},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: wrappedHandler(onDown)},
		}...)
	}

	for _, component := range gui.allListComponents() {
		setUpDownClickBindings(component.GetView().Name(), component.HandlePrevLine, component.HandleNextLine)
	}

	return bindings
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}

	return nil
}

func wrappedHandler(f func() error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f()
	}
}
