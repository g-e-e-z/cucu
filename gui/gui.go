package gui

import "github.com/jroimartin/gocui"

type Gui struct {
	g           *gocui.Gui
	ErrorChan   chan error

	Views       Views
}

func NewGui(errorChan chan error) (*Gui, error) {
	gui := &Gui{
		ErrorChan: errorChan,
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	defer gui.g.Close()
	if err := gui.createAllViews(); err != nil {
		return err
	}
	// if err := gui.setInitialViewContent(); err != nil {
	// 	return err
	// }

	// // TODO: see if we can avoid the circular dependency
	// gui.setPanels()

	// if err = gui.keybindings(g); err != nil {
	// 	return err
	// }

	// if gui.g.CurrentView() == nil {
	// 	viewName := gui.initiallyFocusedViewName()
	// 	view, err := gui.g.View(viewName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if err := gui.switchFocus(view); err != nil {
	// 		return err
	// 	}
	// }

	// ctx, finish := context.WithCancel(context.Background())
	// defer finish()

	// go gui.listenForEvents(ctx, throttledRefresh.Trigger)
	// go gui.monitorContainerStats(ctx)

	// go func() {
	// 	throttledRefresh.Trigger()

	// 	gui.goEvery(time.Millisecond*30, gui.reRenderMain)
	// 	gui.goEvery(time.Millisecond*1000, gui.updateContainerDetails)
	// 	gui.goEvery(time.Millisecond*1000, gui.checkForContextChange)
	// 	// we need to regularly re-render these because their stats will be changed in the background
	// 	gui.goEvery(time.Millisecond*1000, gui.renderContainersAndServices)
	// }()

	// err = g.MainLoop()
	// if err == gocui.ErrQuit {
	// 	return nil
	// }
	// return err

}

func (gui *Gui) createAllViews() error {
	var err error
	for _, mapping := range gui.orderedViewNameMappings() {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}
		(*mapping.viewPtr).FrameRunes = frameRunes
		(*mapping.viewPtr).FgColor = gocui.ColorDefault
	}

	selectedLineBgColor := GetGocuiStyle(gui.Config.UserConfig.Gui.Theme.SelectedLineBgColor)

	gui.Views.Main.Wrap = gui.Config.UserConfig.Gui.WrapMainPanel
	// when you run a docker container with the -it flags (interactive mode) it adds carriage returns for some reason. This is not docker's fault, it's an os-level default.
	gui.Views.Main.IgnoreCarriageReturns = true

	gui.Views.Project.Title = gui.Tr.ProjectTitle

	gui.Views.Services.Highlight = true
	gui.Views.Services.Title = gui.Tr.ServicesTitle
	gui.Views.Services.SelBgColor = selectedLineBgColor

	gui.Views.Containers.Highlight = true
	gui.Views.Containers.SelBgColor = selectedLineBgColor
	if gui.Config.UserConfig.Gui.ShowAllContainers || !gui.DockerCommand.InDockerComposeProject {
		gui.Views.Containers.Title = gui.Tr.ContainersTitle
	} else {
		gui.Views.Containers.Title = gui.Tr.StandaloneContainersTitle
	}

	gui.Views.Images.Highlight = true
	gui.Views.Images.Title = gui.Tr.ImagesTitle
	gui.Views.Images.SelBgColor = selectedLineBgColor

	gui.Views.Volumes.Highlight = true
	gui.Views.Volumes.Title = gui.Tr.VolumesTitle
	gui.Views.Volumes.SelBgColor = selectedLineBgColor

	gui.Views.Networks.Highlight = true
	gui.Views.Networks.Title = gui.Tr.NetworksTitle
	gui.Views.Networks.SelBgColor = selectedLineBgColor

	gui.Views.Options.Frame = false
	gui.Views.Options.FgColor = gui.GetOptionsPanelTextColor()

	gui.Views.AppStatus.FgColor = gocui.ColorCyan
	gui.Views.AppStatus.Frame = false

	gui.Views.Information.Frame = false
	gui.Views.Information.FgColor = gocui.ColorGreen

	gui.Views.Confirmation.Visible = false
	gui.Views.Confirmation.Wrap = true
	gui.Views.Menu.Visible = false
	gui.Views.Menu.SelBgColor = selectedLineBgColor

	gui.Views.Limit.Visible = false
	gui.Views.Limit.Title = gui.Tr.NotEnoughSpace
	gui.Views.Limit.Wrap = true

	gui.Views.FilterPrefix.BgColor = gocui.ColorDefault
	gui.Views.FilterPrefix.FgColor = gocui.ColorGreen
	gui.Views.FilterPrefix.Frame = false

	gui.Views.Filter.BgColor = gocui.ColorDefault
	gui.Views.Filter.FgColor = gocui.ColorGreen
	gui.Views.Filter.Editable = true
	gui.Views.Filter.Frame = false
	gui.Views.Filter.Editor = gocui.EditorFunc(gui.wrapEditor(gocui.SimpleEditor))

	return nil
}
