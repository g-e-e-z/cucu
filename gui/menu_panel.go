package gui

import (
	"github.com/g-e-e-z/cucu/gui/components"
	"github.com/g-e-e-z/cucu/gui/types"
)
type CreateMenuOptions struct {
	Title      string
	Items      []*types.MenuItem
	HideCancel bool
}

func (gui *Gui) getMenuPanel() *components.ListComponent[*types.MenuItem] {
	return &components.ListComponent[*types.MenuItem]{
		ListPanel: components.ListPanel[*types.MenuItem]{
			List: components.NewFilteredList[*types.MenuItem](),
			View: gui.Views.Menu,
		},
		NoItemsMessage: "",
		Gui:            gui.intoInterface(),
		OnClick:        gui.onMenuPress,
		Sort:           nil,
		GetTableCells:  presentation.GetMenuItemDisplayStrings,
		OnRerender: func() error {
			return gui.resizePopupPanel(gui.Views.Menu)
		},
		// so that we can avoid some UI trickiness, the menu will not have filtering
		// abillity yet. To support it, we would need to have filter state against
		// each panel (e.g. for when you filter the images panel, then bring up
		// the options menu, then try to filter that too.
		DisableFilter: true,
	}
}

