package gui

import (
	"github.com/g-e-e-z/cucu/commands"
	rq "github.com/g-e-e-z/cucu/gui/request_panel"
)

func GetRequestString(request *commands.Request) []string {
    return []string{
        request.Method,
		request.Name,
	}
}

func (gui *Gui) getRequestPanel() *rq.RequestPanel {
    return &rq.RequestPanel{
    	SelectedIdx:   0,
    	View:          gui.Views.Requests,
    	Gui:           gui.intoInterface(),
    	GetRenderList: GetRequestString,
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

    gui.RequestPanel.SetItems(requests)

    if err := gui.RequestPanel.RerenderList(); err != nil {
		return err
	}

    return nil
}

