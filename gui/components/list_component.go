// This is copy paste of jesseduffield lazydocker side_list_panel - Side doesnt have much meaning in the context of this project
package components

import (
	"fmt"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type IGui interface {
	IsCurrentView(*gocui.View) bool
	GetUrlView() *gocui.View
	GetRequestInfoView() *gocui.View
	GetResponseInfoView() *gocui.View
	FocusY(selectedLine int, itemCount int, view *gocui.View)
	Update(func() error)
}

type ListComponent[T comparable] struct {
	Log  *logrus.Entry
	View *gocui.View

    // Im in too deep with the generic, RequestContext will always be using a Request
	RequestContext *RequestContext[*commands.Request]
	ListPanel[T]

	Gui            IGui
	NoItemsMessage string

	// returns the cells that we render to the view in a table format. The cells will
	// be rendered with padding.
	GetTableCells func(T) []string
}

func (rp *ListComponent[T]) GetView() *gocui.View {
	return rp.View
}

func (rp *ListComponent[T]) HandleSelect() error {
	_, err := rp.GetSelectedItem(rp.NoItemsMessage)
	if err != nil {
		if err.Error() != rp.NoItemsMessage {
			return err
		}

		if rp.NoItemsMessage != "" {
			rp.Log.Warn(rp.NoItemsMessage)
		}

		return nil
	}

	rp.Refocus()

	return rp.renderContext()
}

func (rp *ListComponent[T]) Rerender() error {
	rp.Gui.Update(func() error {
		rp.View.Clear()
		table := lo.Map(rp.List.GetItems(), func(item T, index int) []string {
			return rp.GetTableCells(item)
		})

		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(rp.View, renderedTable)

		// TODO: Find work around to get this back in/ evalute if its problematic being commented out: Figure out all callers
		// if rp.Gui.IsCurrentView(rp.View) {
		// 	return rp.HandleSelect()
		// }
		// return nil
		return rp.HandleSelect()
	})

	return nil
}

func (rp *ListComponent[T]) renderContext() error {
	if rp.RequestContext == nil {
		return nil
	}

	rp.RequestContext.RenderUrl()

	requestInfoView := rp.Gui.GetRequestInfoView()
	requestInfoView.Tabs = rp.RequestContext.GetRequestInfoTabTitles()
	requestInfoView.TabIndex = rp.RequestContext.requestTabIdx
	rp.RequestContext.GetCurrentRequestInfoTab().Render()
	// task := rp.RequestContext.GetCurrentRequestInfoTab().Render(item)

	responseInfoView := rp.Gui.GetResponseInfoView()
	responseInfoView.Tabs = rp.RequestContext.GetResponseInfoTabTitles()
	responseInfoView.TabIndex = rp.RequestContext.responseTabIdx
	rp.RequestContext.GetCurrentResponseInfoTab().Render()
	// task := rp.RequestContext.GetCurrentResponseInfoTab().Render(item)

	return nil
}

// Keybinding related
func (rp *ListComponent[T]) Refocus() {
	rp.Gui.FocusY(rp.SelectedIdx, rp.List.Len(), rp.View)
}

func (rp *ListComponent[T]) HandleNextLine() error {
	rp.SelectNextLine()

	return rp.HandleSelect()
}

func (rp *ListComponent[T]) HandlePrevLine() error {
	rp.SelectPrevLine()

	return rp.HandleSelect()
}

func (rp *ListComponent[T]) SelectNextLine() {
	rp.moveSelectedLine(1)
}

func (rp *ListComponent[T]) SelectPrevLine() {
	rp.moveSelectedLine(-1)
}

func (rp *ListComponent[T]) moveSelectedLine(delta int) {
	rp.SetSelectedLineIdx(rp.SelectedIdx + delta)
}

func (rp *ListComponent[T]) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if rp.List.Len() > 0 {
		clampedValue = Clamp(value, 0, rp.List.Len()-1)
	}

	rp.SelectedIdx = clampedValue
}

func (rp *ListComponent[T]) HandleNextTab() error {
	if rp.RequestContext == nil {
		return nil
	}

	rp.RequestContext.HandleNextTab()

	return rp.HandleSelect()
}

func (rp *ListComponent[T]) HandlePrevTab() error {
	if rp.RequestContext == nil {
		return nil
	}

	rp.RequestContext.HandlePrevTab()

	return rp.HandleSelect()
}


// Get this from lazycore
// Clamp returns a value x restricted between min and max
func Clamp(x, min, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}
