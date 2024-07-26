package gui

import (
	"github.com/jroimartin/gocui"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

type Views struct {
	Requests *gocui.View
	Url      *gocui.View
	Params   *gocui.View
	Response *gocui.View
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
	}
}

func (gui *Gui) createAllViews() error {
	var err error
	for _, mapping := range gui.orderedViewNameMappings() {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		(*mapping.viewPtr).FgColor = gocui.ColorDefault
	}
	gui.Views.Requests.Highlight = true
	gui.Views.Requests.Title = "Requests"

	gui.Views.Url.Highlight = true
	gui.Views.Url.Title = "Request Url"

	gui.Views.Params.Highlight = true
	gui.Views.Params.Title = "Params"

	gui.Views.Response.Highlight = true
	gui.Views.Response.Title = "Response"
    gui.Views.Response.Wrap = true

	return nil
}
