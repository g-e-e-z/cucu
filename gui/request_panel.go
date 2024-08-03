package gui

import (
	"net/http"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/gui/panels"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) createRequestsPanel() *panels.RequestPanel {
	return &panels.RequestPanel{
		Log:            gui.Log,
		View:           gui.Views.Requests,
		ListPanel:      panels.ListPanel[*commands.Request]{
			List:        panels.NewFilteredList[*commands.Request](),
			View:        gui.Views.Requests,
		},
		Gui:            gui.toInterface(),
		NoItemsMessage: "No Requests",
	}
}

func (gui *Gui) renderRequests() error {
	requests, err := gui.HttpCommands.GetRequests()
	if err != nil {
		return err
	}
	gui.RequestPanel.SetItems(requests)
	return gui.RequestPanel.Rerender()
}

func (gui *Gui) reRenderRequests() error {
    requests := gui.RequestPanel.GetItems()
	gui.RequestPanel.SetItems(requests)
	return gui.RequestPanel.Rerender()
}

func (gui *Gui) handleNewRequest(g *gocui.Gui, v *gocui.View) error {
	gui.Log.Info("Creating New Request")
	newRequest := &commands.Request{
		Name:        "NewRequest!",
		Url:         "this is a placeholder string",
		Method:      http.MethodGet,
		Log:         gui.Log,
		HttpCommand: gui.HttpCommands,
	}
    newRequestList := append(gui.RequestPanel.GetItems(), newRequest)
    gui.RequestPanel.SetItems(newRequestList)

	return gui.reRenderRequests()
}

func (gui *Gui) handleRequestSend(g *gocui.Gui, v *gocui.View) error {
    // TODO: This is a weird way to handle the no items string, fix later
	request, err := gui.RequestPanel.GetSelectedItem(gui.RequestPanel.NoItemsMessage)
	if err != nil {
		return nil
	}

	return gui.SendRequest(request)
}

func (gui *Gui) SendRequest(request *commands.Request) error {
	err := request.Send()
	if err != nil {
		return err
	}
	return gui.RequestPanel.Rerender()
}
