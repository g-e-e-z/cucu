package gui

import "github.com/g-e-e-z/cucu/gui/panels"

func (gui *Gui) createRequestsPanel() *panels.RequestPanel {
	return &panels.RequestPanel{
		Log:            gui.Log,
		View:           gui.Views.Requests,
		Gui:            gui.toInterface(),
		NoItemsMessage: "No Requests",
	}
}

func (gui *Gui) renderRequests() error {
	gui.RequestPanel.SetRequests()
	return gui.RequestPanel.Rerender()
}
