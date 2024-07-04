package gui

import "github.com/jroimartin/gocui"

type Gui struct {
	g             *gocui.Gui
	ErrorChan     chan error
}

func NewGui(errorChan chan error) (*Gui, error) {
	gui := &Gui{
		ErrorChan:     errorChan,
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	defer gui.g.Close()
	if err := gui.createAllViews(); err != nil {
		return err
	}
	// if err := gui.setInitialViewContent(); err != nil {
	// 	return err
	// }

	// // TODO: see if we can avoid the circular dependency
	// gui.setPanels()

	// if err = gui.keybindings(g); err != nil {
	// 	return err
	// }

	// if gui.g.CurrentView() == nil {
	// 	viewName := gui.initiallyFocusedViewName()
	// 	view, err := gui.g.View(viewName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if err := gui.switchFocus(view); err != nil {
	// 		return err
	// 	}
	// }

	// ctx, finish := context.WithCancel(context.Background())
	// defer finish()

	// go gui.listenForEvents(ctx, throttledRefresh.Trigger)
	// go gui.monitorContainerStats(ctx)

	// go func() {
	// 	throttledRefresh.Trigger()

	// 	gui.goEvery(time.Millisecond*30, gui.reRenderMain)
	// 	gui.goEvery(time.Millisecond*1000, gui.updateContainerDetails)
	// 	gui.goEvery(time.Millisecond*1000, gui.checkForContextChange)
	// 	// we need to regularly re-render these because their stats will be changed in the background
	// 	gui.goEvery(time.Millisecond*1000, gui.renderContainersAndServices)
	// }()

	// err = g.MainLoop()
	// if err == gocui.ErrQuit {
	// 	return nil
	// }
	// return err

}
