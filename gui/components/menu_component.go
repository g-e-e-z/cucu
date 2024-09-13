// This is copy paste of jesseduffield lazydocker side_list_panel - Side doesnt have much meaning in the context of this project
package components

import (
	"fmt"

	"github.com/g-e-e-z/cucu/gui/types"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

// type IGui interface {
// 	IsCurrentView(*gocui.View) bool
// 	GetUrlView() *gocui.View
// 	GetMenuInfoView() *gocui.View
// 	GetResponseInfoView() *gocui.View
// 	FocusY(selectedLine int, itemCount int, view *gocui.View)
// 	Update(func() error)
// }

type MenuComponent struct {
	Log  *logrus.Entry
	View *gocui.View

	ListPanel[*types.MenuItem]

	Gui            IGui

	// returns the cells that we render to the view in a table format. The cells will
	// be rendered with padding.
	GetTableCells func(*types.MenuItem) []string
}

func (mc *MenuComponent) GetView() *gocui.View {
	return mc.View
}

func (mc *MenuComponent) HandleSelect() error {
	_, err := mc.GetSelectedItem()
	if err != nil {
		if err.Error() != mc.NoItemsMessage {
			return err
		}

		if mc.NoItemsMessage != "" {
			mc.Log.Warn(mc.NoItemsMessage)
		}

		return nil
	}

	mc.Refocus()

	return nil
}

func (mc *MenuComponent) RerenderList() error {
	mc.Gui.Update(func() error {
		mc.View.Clear()
		table := lo.Map(mc.List.GetItems(), func(item *types.MenuItem, index int) []string {
			return mc.GetTableCells(item)
		})

		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(mc.View, renderedTable)

		// TODO: Find work around to get this back in/ evalute if its problematic being commented out: Figure out all callers
		if mc.Gui.IsCurrentView(mc.View) {
			return mc.HandleSelect()
		}
		return nil
	})

	return nil
}


// Keybinding related
func (mc *MenuComponent) Refocus() {
	mc.Gui.FocusY(mc.SelectedIdx, mc.List.Len(), mc.View)
}

func (mc *MenuComponent) HandleNextLine() error {
	mc.SelectNextLine()

	return mc.HandleSelect()
}

func (mc *MenuComponent) HandlePrevLine() error {
	mc.SelectPrevLine()

	return mc.HandleSelect()
}

func (mc *MenuComponent) SelectNextLine() {
	mc.moveSelectedLine(1)
}

func (mc *MenuComponent) SelectPrevLine() {
	mc.moveSelectedLine(-1)
}

func (mc *MenuComponent) moveSelectedLine(delta int) {
	mc.SetSelectedLineIdx(mc.SelectedIdx + delta)
}

func (mc *MenuComponent) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if mc.List.Len() > 0 {
		clampedValue = Clamp(value, 0, mc.List.Len()-1)
	}

	mc.SelectedIdx = clampedValue
}
