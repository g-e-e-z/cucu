package gui

import (
	"net/http"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/gui/components"
	"github.com/g-e-e-z/cucu/gui/presentation"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getRequestsPanel() *components.ListComponent[*commands.Request] {
	return &components.ListComponent[*commands.Request]{
		Log:            gui.Log,
		View:           gui.Views.Requests,
		ListPanel:      components.ListPanel[*commands.Request]{
			List:        components.NewFilteredList[*commands.Request](),
			View:        gui.Views.Requests,
		},
		Gui:            gui.toInterface(),
		NoItemsMessage: "No Requests",
        GetTableCells: func(request *commands.Request) []string {
			return presentation.GetRequestStrings(request)
		},
	}
}

func (gui *Gui) renderRequests() error {
	requests, err := gui.HttpCommands.GetRequests()
	if err != nil {
		return err
	}
	gui.Components.Requests.SetItems(requests)
	return gui.Components.Requests.Rerender()
}

func (gui *Gui) reRenderRequests() error {
    requests := gui.Components.Requests.GetItems()
	gui.Components.Requests.SetItems(requests)
	return gui.Components.Requests.Rerender()
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
    newRequestList := append(gui.Components.Requests.GetItems(), newRequest)
    gui.Components.Requests.SetItems(newRequestList)

	return gui.reRenderRequests()
}

func (gui *Gui) handleRequestSend(g *gocui.Gui, v *gocui.View) error {
    // TODO: This is a weird way to handle the no items string, fix later
	request, err := gui.Components.Requests.GetSelectedItem(gui.Components.Requests.NoItemsMessage)
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
	return gui.Components.Requests.Rerender()
}
