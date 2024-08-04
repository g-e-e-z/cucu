package gui

import (
	"github.com/jesseduffield/gocui"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

type Views struct {
	Requests *gocui.View
	Url      *gocui.View
	Params   *gocui.View
	Response *gocui.View

    // popups
    EditMethod *gocui.View
}

type viewNameMapping struct {
	viewPtr **gocui.View
	name    string
}

func (gui *Gui) orderedViewNameMappings() []viewNameMapping {
	return []viewNameMapping{
		{viewPtr: &gui.Views.Requests, name: "requests"},
		{viewPtr: &gui.Views.Url, name: "url"},
		{viewPtr: &gui.Views.Params, name: "params"},
		{viewPtr: &gui.Views.Response, name: "response"},

        {viewPtr: &gui.Views.EditMethod, name: "editMethod"},
	}
}

func (gui *Gui) createAllViews() error {
    frameRunes := []rune{'─', '│', '╭', '╮', '╰', '╯'}
	var err error
	for _, mapping := range gui.orderedViewNameMappings() {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
        (*mapping.viewPtr).FrameRunes = frameRunes
		(*mapping.viewPtr).FgColor = gocui.ColorDefault
	}
	gui.Views.Requests.Highlight = true
	gui.Views.Requests.Title = "Requests"

	gui.Views.Url.Highlight = false
	gui.Views.Url.Title = "Request Url"
    gui.Views.Url.Wrap = false
    gui.Views.Url.Editable = false
    // gui.Views.Url.Editor = gocui.EditorFunc(gocui.SimpleEditor)
    gui.Views.Url.Editor = gocui.EditorFunc(gui.wrapEditor(gocui.SimpleEditor))


	gui.Views.Params.Highlight = false
	gui.Views.Params.Title = "Params"

	gui.Views.Response.Highlight = false
	gui.Views.Response.Title = "Response"
    gui.Views.Response.Wrap = true

	gui.Views.EditMethod.Visible = false
	gui.Views.EditMethod.Highlight = true
	gui.Views.EditMethod.Title = "Choose Http Method"

	return nil
}

func (gui *Gui) wrapEditor(f func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
		matched := f(v, key, ch, mod)
        request, _:= gui.Panels.Requests.GetSelectedItem(gui.Panels.Requests.NoItemsMessage)
        request.Url = v.TextArea.GetContent()
        // gui.Log.Info("Text: ", v.TextArea.GetContent())
        gui.Panels.Requests.Rerender()
		// if matched {
		// 	// if err := gui.onNewFilterNeedle(v.TextArea.GetContent()); err != nil {
		// 	// 	gui.Log.Error(err)
		// 	// }
		// }
		return matched
	}
}

// func (gui *Gui) handleOpenFilter() error {
// 	panel, ok := gui.currentListPanel()
// 	if !ok {
// 		return nil
// 	}
//
// 	if panel.IsFilterDisabled() {
// 		return nil
// 	}
//
// 	gui.State.Filter.active = true
// 	gui.State.Filter.panel = panel
//
// 	return gui.switchFocus(gui.Views.Filter)
// }

