package gui

import (
	"github.com/g-e-e-z/cucu/commands"
	lc "github.com/g-e-e-z/cucu/gui/list_components"
)

func GetRequestString(request *commands.Request) []string {
    return []string{
        request.Method,
		request.Name,
	}
}

func (gui *Gui) getRequestsComponent() *lc.ListComponent[*commands.Request]{
    return &lc.ListComponent[*commands.Request]{
    	SelectedIdx:    0,
    	View:           gui.Views.Requests,
    	Gui:            gui.intoInterface(),
    	GetRenderList:  GetRequestString,
    	NoItemsMessage: "No Requests Configured",
    }
}

func (gui *Gui) refreshRequests() error {
    // Skip if Requests hasn't been created
    if gui.Views.Requests == nil {
        return nil
    }

    requests, err := gui.OSCommands.GetRequests()
    if err != nil {
        return err
    }

    gui.Components.Requests.SetItems(requests)

    if err := gui.Components.Requests.RerenderList(); err != nil {
		return err
	}

    return nil
}

