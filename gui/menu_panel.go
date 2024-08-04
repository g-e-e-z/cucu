package gui

import (
	"github.com/g-e-e-z/cucu/gui/components"
	"github.com/g-e-e-z/cucu/gui/presentation"
	"github.com/g-e-e-z/cucu/gui/types"
	"github.com/g-e-e-z/cucu/utils"
)

type CreateMenuOptions struct {
	Title      string
	Items      []*types.MenuItem
	Index      int
	HideCancel bool
}

func (gui *Gui) getMenuPanel() *components.ListComponent[*types.MenuItem] {
	return &components.ListComponent[*types.MenuItem]{
		View: gui.Views.Menu,
		ListPanel: components.ListPanel[*types.MenuItem]{
			List:           components.NewFilteredList[*types.MenuItem](),
			View:           gui.Views.Menu,
			NoItemsMessage: "This should never be seen",
		},
		Gui:           gui.toInterface(),
		GetTableCells: presentation.GetMenuItemDisplayStrings,
	}
}

func (gui *Gui) onMenuPress(menuItem *types.MenuItem) error {
	if err := gui.handleMenuClose(); err != nil {
		return err
	}

	if menuItem.OnPress != nil {
		return menuItem.OnPress()
	}

	return nil
}

func (gui *Gui) handleMenuPress() error {
	selectedMenuItem, err := gui.Components.Menu.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.onMenuPress(selectedMenuItem)
}

func (gui *Gui) Menu(opts CreateMenuOptions) error {
	if !opts.HideCancel {
		// this is mutative but I'm okay with that for now
		opts.Items = append(opts.Items, &types.MenuItem{
			LabelColumns: []string{"cancel"},
			OnPress: func() error {
				return nil
			},
		})
	}

	maxColumnSize := 1

	for _, item := range opts.Items {
		if item.LabelColumns == nil {
			item.LabelColumns = []string{item.Label}
		}

		// if item.OpensMenu {
		// 	item.LabelColumns[0] = utils.OpensMenuStyle(item.LabelColumns[0])
		// }

		maxColumnSize = utils.Max(maxColumnSize, len(item.LabelColumns))
	}

	for _, item := range opts.Items {
		if len(item.LabelColumns) < maxColumnSize {
			// we require that each item has the same number of columns so we're padding out with blank strings
			// if this item has too few
			item.LabelColumns = append(item.LabelColumns, make([]string, maxColumnSize-len(item.LabelColumns))...)
		}
	}
	gui.Components.Menu.SetItems(opts.Items)
	gui.Components.Menu.SetSelectedLineIdx(opts.Index)

	if err := gui.Components.Menu.RerenderList(); err != nil {
		return err
	}

	gui.Views.Menu.Title = opts.Title
	gui.Views.Menu.Visible = true

	gui.g.SetCurrentView(gui.Views.Menu.Name())
	// return gui.switchFocus(gui.Views.Menu)
	return nil
}

func (gui *Gui) handleMenuClose() error {
	gui.Views.Menu.Visible = false

	// this code is here for when we do add filter ability to the menu panel,
	// though it's currently disabled
	// if gui.State.Filter.panel == gui.Panels.Menu {
	// 	if err := gui.clearFilter(); err != nil {
	// 		return err
	// 	}
	//
	// 	// we need to remove the view from the view stack because we're about to
	// 	// return focus and don't want to land in the search view when it was searching
	// 	// the menu in the first place
	// 	gui.removeViewFromStack(gui.Views.Filter)
	// }

	// return gui.returnFocus()
	_, err := gui.g.SetCurrentView(gui.Views.Requests.Name())
	if err != nil {
		return err
	}
	gui.Components.Requests.RerenderList()
	return nil
}
