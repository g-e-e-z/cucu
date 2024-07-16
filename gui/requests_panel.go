package gui

func (gui *Gui) refreshRequests() error {
	if gui.Views.Requests == nil {
		// if the containersView hasn't been instantiated yet we just return
		return nil
	}
    requests, err := gui.OSCommands.RefreshRequests()
    if err != nil {
        return err
    }
    gui.Panels.Requests.SetItems(requests)

    return gui.renderRequests()
}

func (gui *Gui) renderRequests() error {
	if err := gui.Panels.Requests.RerenderList(); err != nil {
		return err
	}

	return nil
}
