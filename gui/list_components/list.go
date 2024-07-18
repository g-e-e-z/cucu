package list_component

import (
	"errors"
	"fmt"

	"github.com/g-e-e-z/cucu/utils"
	"github.com/jroimartin/gocui"
)

type IGui interface {
	// HandleClick(v *gocui.View, itemCount int, selectedLine *int, handleSelect func() error) error
	// NewSimpleRenderStringTask(getContent func() string) tasks.TaskFunc
	FocusY(selectedLine int, itemCount int, view *gocui.View)
	// ShouldRefresh(contextKey string) bool
	// GetMainView() *gocui.View
	IsCurrentView(*gocui.View) bool
	// FilterString(view *gocui.View) string
	// IgnoreStrings() []string
	Update(func() error)

	// QueueTask(f func(ctx context.Context)) error
}

type ListComponent[T comparable] struct {
	items   []T
	indices []int

	SelectedIdx int
	View        *gocui.View
	// a representation of the gui.... Why though
	Gui IGui

	GetRenderList func(T) []string

	NoItemsMessage string
}

// TODO: Understand this syntax more
func NewList[T comparable]() *ListComponent[T] {
	return &ListComponent[T]{}
}

func (self *ListComponent[T]) GetView() *gocui.View {
    return self.View
}

func (self *ListComponent[T]) SetItems(items []T) {
	self.items = items
	self.indices = make([]int, len(items))
	for i := range self.indices {
		self.indices[i] = i
	}
}

func (self *ListComponent[T]) GetItems() []T {
	result := make([]T, len(self.indices))
	for i, index := range self.indices {
		result[i] = self.items[index]
	}
	return result
}

// go get lo
func Map[T any, R any](collection []T, iteratee func(T, int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}

func (self *ListComponent[T]) RerenderList() error {
	self.Gui.Update(func() error {
		self.View.Clear()
		// move to lo.Map after internet
		table := Map(self.GetItems(), func(item T, index int) []string {
			return self.GetRenderList(item)
		})
		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(self.View, renderedTable)

		// if self.OnRerender != nil {
		// 	if err := self.OnRerender(); err != nil {
		// 		return err
		// 	}
		// }

		// if self.Gui.IsCurrentView(self.View) {
		// 	return self.HandleSelect()
		// }
		return nil
	})

	return nil
}

func (self *ListComponent[T]) Len() int {
	return len(self.indices)
}

func (self *ListComponent[T]) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if self.Len() > 0 {
		clampedValue = Clamp(value, 0, self.Len()-1)
	}

	self.SelectedIdx = clampedValue
}

func (self *ListComponent[T]) clampSelectedLineIdx() {
	clamped := Clamp(self.SelectedIdx, 0, self.Len()-1)

	if clamped != self.SelectedIdx {
		self.SelectedIdx = clamped
	}
}

func (self *ListComponent[T]) HandleNextLine() error {
	self.SelectNextLine()

	return self.HandleSelect()
}

func (self *ListComponent[T]) HandlePrevLine() error {
	self.SelectPrevLine()

	return self.HandleSelect()
}

func (self *ListComponent[T]) GetSelectedItem() (T, error) {
	var zero T

	item, ok := self.TryGet(self.SelectedIdx)
	if !ok {
		// could probably have a better error here
		return zero, errors.New(self.NoItemsMessage)
	}

	return item, nil
}

func (self *ListComponent[T]) TryGet(index int) (T, bool) {
    if index < 0 || index >= len(self.indices) {
		var zero T
		return zero, false
	}

	return self.items[self.indices[index]], true
}


func (self *ListComponent[T]) HandleSelect() error {
    // When we render the params/ body to that view uncomment item and use
	// item, err := self.GetSelectedItem()
	_, err := self.GetSelectedItem()
	if err != nil {
		if err.Error() != self.NoItemsMessage {
			return err
		}

		// if self.NoItemsMessage != "" {
		// 	self.Gui.NewSimpleRenderStringTask(func() string { return self.NoItemsMessage })
		// }

		return nil
	}

	self.Refocus()

	return nil //self.renderContext(item)
}

func (self *ListComponent[T]) Refocus() {
	self.Gui.FocusY(self.SelectedIdx, self.Len(), self.View)
}


// moves the cursor up or down by the given amount (up for negative values)
func (self *ListComponent[T]) moveSelectedLine(delta int) {
	self.SetSelectedLineIdx(self.SelectedIdx + delta)
}

func (self *ListComponent[T]) SelectNextLine() {
	self.moveSelectedLine(1)
}

func (self *ListComponent[T]) SelectPrevLine() {
	self.moveSelectedLine(-1)
}

// Get this from lazycore
// Clamp returns a value x restricted between min and max
func Clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}
