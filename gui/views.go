package gui

import (
	"github.com/jesseduffield/gocui"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

type Views struct {
	Requests *gocui.View
	Url      *gocui.View
	RequestInfo   *gocui.View
	ResponseInfo *gocui.View

    Menu *gocui.View
}

type viewNameMapping struct {
	viewPtr **gocui.View
	name    string
}

func (gui *Gui) orderedViewNameMappings() []viewNameMapping {
	return []viewNameMapping{
		{viewPtr: &gui.Views.Requests, name: "requests"},
		{viewPtr: &gui.Views.Url, name: "url"},
		{viewPtr: &gui.Views.RequestInfo, name: "params"},
		{viewPtr: &gui.Views.ResponseInfo, name: "response"},

        {viewPtr: &gui.Views.Menu, name: "menu"},
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

    gui.Views.Url.Editor = gocui.EditorFunc(gui.wrapEditor(gocui.SimpleEditor))


	gui.Views.RequestInfo.Highlight = false
	gui.Views.RequestInfo.Title = "Params"
	// gui.Views.RequestInfo.SelBgColor = gocui.ColorRed // Come back later

	gui.Views.ResponseInfo.Highlight = false
	gui.Views.ResponseInfo.Title = "Response"
    gui.Views.ResponseInfo.Wrap = true

	gui.Views.Menu.Visible = false
	gui.Views.Menu.Highlight = true
	gui.Views.Menu.Title = "Choose Http Method"

	return nil
}

func (gui *Gui) wrapEditor(f func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
		matched := f(v, key, ch, mod)
        request, _:= gui.Components.Requests.GetSelectedItem()
        request.Url = v.TextArea.GetContent()
        gui.Components.Requests.RerenderList()
        // TODO: Handle Modified Better, save initial state? hash?
        request.Modified = true

		return matched
	}
}
