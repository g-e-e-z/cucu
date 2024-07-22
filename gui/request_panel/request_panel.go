package request_panel

import (
	"fmt"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jroimartin/gocui"
	"github.com/samber/lo"
)

const EMPTY_REQUESTS = "No requests loaded"

type RequestPanel struct {
	FilteredList[*commands.Request]

	SelectedIdx int
	View        *gocui.View
	Gui         IGui

	GetRenderList func(*commands.Request) []string
}

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


func (self *RequestPanel) RerenderList() error {
	self.Gui.Update(func() error {
		self.View.Clear()
		table := lo.Map(self.GetItems(), func(item *commands.Request, index int) []string {
			return self.GetRenderList(item)
		})
		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
        // Idk why a new line is prepended to the render string, slice is a temporary fix
        fmt.Fprint(self.View, renderedTable[1:])
		// fmt.Fprint(self.View, "\n")

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
