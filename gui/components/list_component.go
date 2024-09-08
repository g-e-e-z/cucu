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

type ISideListPanel interface {
	HandleSelect() error
	GetView() *gocui.View
	Refocus()
	RerenderList() error
	HandleNextLine() error
	HandlePrevLine() error
}


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

	RequestContext *RequestContext[*commands.Request]
	ListPanel[T]

	Gui            IGui

	// returns the cells that we render to the view in a table format. The cells will
	// be rendered with padding.
	GetTableCells func(T) []string
}

func (self *ListComponent[T]) GetView() *gocui.View {
	return self.View
}

func (self *ListComponent[T]) HandleSelect() error {
	_, err := self.GetSelectedItem()
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

	return self.renderContext()
}

func (self *ListComponent[T]) RerenderList() error {
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
		if self.Gui.IsCurrentView(self.View) {
			return self.HandleSelect()
		}
		return nil
		// return self.HandleSelect()
	})

	return nil
}

func (self *ListComponent[T]) renderContext() error {
	if self.RequestContext == nil {
		return nil
	}

	self.RequestContext.RenderUrl()

	requestInfoView := self.Gui.GetRequestInfoView()
	requestInfoView.Tabs = self.RequestContext.GetRequestInfoTabTitles()
	requestInfoView.TabIndex = self.RequestContext.requestTabIdx
	self.RequestContext.GetCurrentRequestInfoTab().Render()
	// task := self.RequestContext.GetCurrentRequestInfoTab().Render(item)

	responseInfoView := self.Gui.GetResponseInfoView()
	responseInfoView.Tabs = self.RequestContext.GetResponseInfoTabTitles()
	responseInfoView.TabIndex = self.RequestContext.responseTabIdx
	self.RequestContext.GetCurrentResponseInfoTab().Render()
	// task := self.RequestContext.GetCurrentResponseInfoTab().Render(item)

	return nil
}

// Keybinding related
func (self *ListComponent[T]) Refocus() {
	self.Gui.FocusY(self.SelectedIdx, self.List.Len(), self.View)
}

func (self *ListComponent[T]) HandleNextLine() error {
	self.SelectNextLine()

	return self.HandleSelect()
}

func (self *ListComponent[T]) HandlePrevLine() error {
	self.SelectPrevLine()

	return self.HandleSelect()
}

func (self *ListComponent[T]) SelectNextLine() {
	self.moveSelectedLine(1)
}

func (self *ListComponent[T]) SelectPrevLine() {
	self.moveSelectedLine(-1)
}

func (self *ListComponent[T]) moveSelectedLine(delta int) {
	self.SetSelectedLineIdx(self.SelectedIdx + delta)
}

func (self *ListComponent[T]) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if self.List.Len() > 0 {
		clampedValue = Clamp(value, 0, self.List.Len()-1)
	}

	self.SelectedIdx = clampedValue
}

func (self *ListComponent[T]) HandleNextReqTab() error {
	if self.RequestContext == nil {
		return nil
	}

	self.RequestContext.HandleNextReqTab()

	return self.HandleSelect()
}

func (self *ListComponent[T]) HandlePrevReqTab() error {
	if self.RequestContext == nil {
		return nil
	}

	self.RequestContext.HandlePrevReqTab()

	return self.HandleSelect()
}

func (self *ListComponent[T]) HandleNextResTab() error {
	if self.RequestContext == nil {
		return nil
	}

	self.RequestContext.HandleNextResTab()

	return self.HandleSelect()
}

func (self *ListComponent[T]) HandlePrevResTab() error {
	if self.RequestContext == nil {
		return nil
	}

	self.RequestContext.HandlePrevResTab()

	return self.HandleSelect()
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
