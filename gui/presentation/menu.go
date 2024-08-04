package presentation

import "github.com/g-e-e-z/cucu/gui/types"

func GetMenuItemDisplayStrings(menuItem *types.MenuItem) []string {
	return menuItem.LabelColumns
}

