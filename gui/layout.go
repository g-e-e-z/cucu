package gui

import (
	"github.com/jesseduffield/gocui"
)

// getFocusLayout returns a manager function for when view gain and lose focus
func (gui *Gui) getFocusLayout() func(g *gocui.Gui) error {
	var previousView *gocui.View
	return func(g *gocui.Gui) error {
		newView := gui.g.CurrentView()
		if err := gui.onFocusChange(); err != nil {
			return err
		}
		// for now we don't consider losing focus to a popup panel as actually losing focus
		if newView != previousView { //&& !gui.isPopupPanel(newView.Name()) {
			gui.onFocusLost(previousView, newView)
			gui.onFocus(newView)
			previousView = newView
		}
		return nil
	}
}

func (gui *Gui) onFocusChange() error {
	currentView := gui.g.CurrentView()
	for _, view := range gui.g.Views() {
		view.Highlight = view == currentView && view.Name() != "main"
	}
	return nil
}

func (gui *Gui) onFocusLost(v *gocui.View, newView *gocui.View) {
	if v == nil {
		return
	}

	// if !gui.isPopupPanel(newView.Name()) {
	// 	v.ParentView = nil
	// }

	// refocusing because in responsive mode (when the window is very short) we want to ensure that after the view size changes we can still see the last selected item
	// gui.focusPointInView(v)

	// gui.Log.Info(v.Name() + " focus lost")
}

func (gui *Gui) onFocus(v *gocui.View) {
	if v == nil {
		return
	}

	// gui.focusPointInView(v)

	// gui.Log.Info(v.Name() + " focus gained")
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true
	width, height := gui.g.Size()

	viewDimensions := gui.getWindowDimensions()
	// we assume that the view has already been created.
	setViewFromDimensions := func(viewName string, windowName string) (*gocui.View, error) {
		viewPositionObj, ok := viewDimensions[windowName]

		view, err := g.View(viewName)
		if err != nil {
			return nil, err
		}

		if !ok {
			// view not specified in dimensions object: so create the view and hide it
			// making the view take up the whole space in the background in case it needs
			// to render content as soon as it appears, because lazyloaded content (via a pty task)
			// cares about the size of the view.
			_, err := g.SetView(viewName, 0, 0, width, height, 0)
			return view, err
		}

		frameOffset := 1
		if view.Frame {
			frameOffset = 0
		}
		_, err = g.SetView(
			viewName,
			viewPositionObj.X0.GetCoordinate(width)-frameOffset,
			viewPositionObj.Y0.GetCoordinate(height)-frameOffset,
			viewPositionObj.X1.GetCoordinate(width)+frameOffset,
			viewPositionObj.Y1.GetCoordinate(height)+frameOffset,
            0,
		)

		return view, err
	}

	for _, view := range g.Views() {
		_, err := setViewFromDimensions(view.Name(), view.Name())
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
	}
	return nil
}

// func (gui *Gui) focusPointInView(view *gocui.View) {
// 	if view == nil {
// 		return
// 	}
//
// 	currentPanel, ok := gui.currentListPanel()
// 	if ok {
// 		currentPanel.Refocus()
// 	}
// }

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 1, 1, 0)
}
