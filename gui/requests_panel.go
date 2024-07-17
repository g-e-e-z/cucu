package gui


type RequestPanel struct {

}

func (rp *RequestPanel) SetItems() {

}
func (rp *RequestPanel) RenderList() {

}

func (gui *Gui) setRequestsPanel() {
}


func (gui *Gui) refreshRequests() error {
    return nil
}
// func (gui *Gui) refreshRequests() error {
// 	if gui.Views.Requests == nil {
// 		// if the containersView hasn't been instantiated yet we just return
// 		return nil
// 	}
//     requests, err := gui.OSCommands.RefreshRequests()
//     if err != nil {
//         return err
//     }
//     gui.RequestPanel.SetItems(requests)
//
//     return gui.renderRequests()
// }
//
// func (gui *Gui) renderRequests() error {
// 	if err := gui.RequestPanel.RerenderList(); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
