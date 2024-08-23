// This is copy paste of jesseduffield lazydocker list_panel
package components

import (
	"errors"

	"github.com/jesseduffield/gocui"
	lcUtils "github.com/jesseduffield/lazycore/pkg/utils"
)

// List Panel Extends FilteredList to include a View and Current Index. Has methods for moving the SelectedIndex. The View is not directly accessed here. Consider moving up a level
type ListPanel[T comparable] struct {
	SelectedIdx    int
	List           *FilteredList[T]
	View           *gocui.View
	NoItemsMessage string
}

func (self *ListPanel[T]) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if self.List.Len() > 0 {
		clampedValue = lcUtils.Clamp(value, 0, self.List.Len()-1)
	}

	self.SelectedIdx = clampedValue
}

func (self *ListPanel[T]) clampSelectedLineIdx() {
	clamped := lcUtils.Clamp(self.SelectedIdx, 0, self.List.Len()-1)

	if clamped != self.SelectedIdx {
		self.SelectedIdx = clamped
	}
}

// moves the cursor up or down by the given amount (up for negative values)
func (self *ListPanel[T]) moveSelectedLine(delta int) {
	self.SetSelectedLineIdx(self.SelectedIdx + delta)
}

func (self *ListPanel[T]) SelectNextLine() {
	self.moveSelectedLine(1)
}

func (self *ListPanel[T]) SelectPrevLine() {
	self.moveSelectedLine(-1)
}

func (self *ListPanel[T]) SetItems(items []T) {
	self.List.SetItems(items)
	// self.FilterAndSort()
}

func (self *ListPanel[T]) GetItems() []T {
	return self.List.GetItems()
}

func (self *ListPanel[T]) GetSelectedItem() (T, error) {
	var zero T

	item, ok := self.List.TryGet(self.SelectedIdx)
	if !ok {
		// could probably have a better error here
		return zero, errors.New(self.NoItemsMessage)
	}

	return item, nil
}

func (self *ListPanel[T]) RemoveSelectedItem() (T, error) {
	var zero T

	item, ok := self.List.TryGet(self.SelectedIdx)
	if !ok {
		// could probably have a better error here
		return zero, errors.New(self.NoItemsMessage)
	}

    items := self.GetItems()
    newItems := append(items[:self.SelectedIdx], items[self.SelectedIdx+1:]...)
    self.SetItems(newItems)
    self.SelectedIdx = min(self.SelectedIdx, self.List.Len()-1)

	return item, nil
}
