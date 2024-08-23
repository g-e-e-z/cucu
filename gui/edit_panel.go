package gui

import (
	"fmt"
	"strings"
)

type CreateEditOptions struct {
	Title string
	Value string
}

func (gui *Gui) Edit(opts CreateEditOptions) error {
	gui.Views.Edit.Title = opts.Title
    gui.Views.Edit.Clear()
    gui.Views.Edit.TextArea.Clear()
    fmt.Fprint(gui.Views.Edit, opts.Value)
	gui.Views.Edit.TextArea.TypeString(opts.Value)

	gui.Views.Edit.Visible = true
	gui.Views.Edit.Editable = true

    gui.Views.Edit.SetCursor(len(opts.Value), 0)

	gui.g.SetCurrentView(gui.Views.Edit.Name())
	// return gui.switchFocus(gui.Views.Menu)
	return nil
}

func (gui *Gui) handleEditConfirm() error {
	_, err := gui.g.SetCurrentView(gui.Views.Requests.Name())
	if err != nil {
		return err
	}
	gui.Views.Edit.Visible = false
	gui.Views.Edit.Editable = false

    // Lord forgive me for my sins
    words := strings.Split(gui.Views.Edit.Title, " ")
	request, err := gui.Components.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
    field := words[1]
    if field == "Name" {
        request.Name = gui.Views.Edit.TextArea.GetContent()
    }
    if field == "Url" {
        request.Url = gui.Views.Edit.TextArea.GetContent()
    }

	return gui.Components.Requests.RerenderList()
}

func (gui *Gui) handleEditCancel() error {
	_, err := gui.g.SetCurrentView(gui.Views.Requests.Name())
	if err != nil {
		return err
	}
	gui.Views.Edit.Visible = false
	gui.Views.Edit.Editable = false
	return gui.Components.Requests.RerenderList()
}
