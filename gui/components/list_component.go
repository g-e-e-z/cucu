// This is copy paste of jesseduffield lazydocker side_list_panel - Side doesnt have much meaning in the context of this project
package components

import (
	"fmt"

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
    // Im in too deep with the generic, RequestContext will always be using a Request
	RequestContext *RequestContext[T]

	Log  *logrus.Entry
	View *gocui.View
	ListPanel[T]

	Gui            IGui
	NoItemsMessage string

	// returns the cells that we render to the view in a table format. The cells will
	// be rendered with padding.
	GetTableCells func(T) []string
}

func (self *ListComponent[T]) GetView() *gocui.View {
	return self.View
}

func (self *ListComponent[T]) HandleSelect() error {
	item, err := self.GetSelectedItem(self.NoItemsMessage)
	if err != nil {
		if err.Error() != self.NoItemsMessage {
			return err
		}

		if self.NoItemsMessage != "" {
			self.Log.Warn(self.NoItemsMessage)
		}

		return nil
	}

	self.Refocus()

	return self.renderContext(item)
}

func (self *ListComponent[T]) Rerender() error {
	self.Gui.Update(func() error {
		self.View.Clear()
		table := lo.Map(self.List.GetItems(), func(item T, index int) []string {
			return self.GetTableCells(item)
		})

		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(self.View, renderedTable)

		// TODO: Find work around to get this back in/ evalute if its problematic being commented out: Figure out all callers
		// if self.Gui.IsCurrentView(self.View) {
		// 	return self.HandleSelect()
		// }
		// return nil
		return self.HandleSelect()
	})

	return nil
}

func (self *ListComponent[T]) renderContext(item T) error {
	if self.RequestContext == nil {
		return nil
	}

	// In lazydocker this is the tab names in the main view, will use something similar in future iterations
	// key := self.ContextState.GetCurrentContextKey(item)
	// if !self.Gui.ShouldRefresh(key) {
	// 	return nil
	// }

	// mainView := self.Gui.GetMainView()
	// mainView.Tabs = self.ContextState.GetMainTabTitles()
	// mainView.TabIndex = self.ContextState.mainTabIdx

	// task := self.ContextState.GetCurrentMainTab().Render(item)
	// return self.Gui.QueueTask(task)

	// Url
	self.RequestContext.RenderUrl(item)

	// RequestInfo
	requestInfoView := self.Gui.GetRequestInfoView()
	requestInfoView.Tabs = self.RequestContext.GetRequestInfoTabTitles()
	requestInfoView.TabIndex = self.RequestContext.requestTabIdx
	self.RequestContext.GetCurrentRequestInfoTab().Render(item)
	// task := self.ContextState.GetCurrentRequestInfoTab().Render(item)

	// ResponseInfo
	responseInfoView := self.Gui.GetResponseInfoView()
	responseInfoView.Tabs = self.RequestContext.GetResponseInfoTabTitles()
	responseInfoView.TabIndex = self.RequestContext.responseTabIdx
	self.RequestContext.GetCurrentResponseInfoTab().Render(item)
	// task := self.ContextState.GetCurrentResponseInfoTab().Render(item)

	// TODO: Don't write directly, this whole block is questionable, TextArea etc..
	// urlView := self.Gui.GetUrlView()
	// urlView.ClearTextArea()
	// output := string(bom.Clean([]byte(item.Url)))
	// s := utils.NormalizeLinefeeds(output)
	// urlView.TextArea.TypeString(s)
	// urlView.SetCursor(len(s), 0)
	// fmt.Fprint(urlView, s)
	//
	// paramsView := self.Gui.GetParamsView()
	// paramsView.Clear()
	// params, err := item.GetParams()
	// if err != nil {
	// 	return err
	// }
	// table := utils.MapToSlice(utils.ValuesToMap(params))
	// renderedTable, err := utils.RenderComponent(table)
	// fmt.Fprint(paramsView, renderedTable)
	//
	// responseView := self.Gui.GetResponseView()
	// responseView.Clear()
	// fmt.Fprint(responseView, item.ResponseBody)

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
