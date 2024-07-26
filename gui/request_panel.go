package gui

import (
	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/gui/panels"
	"github.com/jroimartin/gocui"
)

func (gui *Gui) createRequestsPanel() *panels.RequestPanel {
	return &panels.RequestPanel{
		Log:            gui.Log,
		View:           gui.Views.Requests,
		Gui:            gui.toInterface(),
		NoItemsMessage: "No Requests",
	}
}

func (gui *Gui) renderRequests() error {
    requests, err := gui.HttpCommands.GetRequests()
    if err != nil {
        return err
    }
	gui.RequestPanel.SetRequests(requests)
	return gui.RequestPanel.Rerender()
}

func (gui *Gui) handleRequestSend(g *gocui.Gui, v *gocui.View) error {
    request, err := gui.RequestPanel.GetSelectedRequest()
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
	// return gui.WithWaitingStatus(gui.Tr.PausingStatus, func() (err error) {
	// 	if container.Details.State.Paused {
	// 		err = container.Unpause()
	// 	} else {
	// 		err = container.Pause()
	// 	}
	//
	// 	if err != nil {
	// 		return gui.createErrorPanel(err.Error())
	// 	}
	//
	// 	return gui.refreshContainersAndServices()
	// })
    return nil
}

