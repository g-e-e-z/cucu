package gui

import (
	"fmt"
	"net/http"

	"github.com/jesseduffield/gocui"
)

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func (gui *Gui) handleEditMethod(_ *gocui.Gui, v *gocui.View) error {
    editView := gui.Views.EditMethod
	editView.Visible = true
    gui.g.SetCurrentView("editMethod")
	req, _ := gui.RequestPanel.GetSelectedRequest()
	currentMethod := req.Method
    editView.Clear()
    fmt.Fprint(editView, "Current Method: ", currentMethod, "\n")
    for _, httpMethod := range httpMethods {
        fmt.Fprint(editView, httpMethod, "\n")
    }

	return nil
}

func (gui *Gui) handleCloseEditMethod(_ *gocui.Gui, v *gocui.View) error {
    // Need to modify keybinds here as well: for now, just making non-competing keybinds
    editView := gui.Views.EditMethod
	editView.Visible = false
    gui.g.SetCurrentView("requests")
	return nil
}
