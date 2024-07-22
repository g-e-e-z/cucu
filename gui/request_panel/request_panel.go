package request_panel

import (
	"errors"
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

    prevIdx     int
}

type IGui interface {
	// HandleClick(v *gocui.View, itemCount int, selectedLine *int, handleSelect func() error) error
	// NewSimpleRenderStringTask(getContent func() string) tasks.TaskFunc
	FocusY(selectedLine int, itemCount int, view *gocui.View)
	// ShouldRefresh(contextKey string) bool
	GetUrlView() *gocui.View
	GetParamsView() *gocui.View
	IsCurrentView(*gocui.View) bool
	// FilterString(view *gocui.View) string
	// IgnoreStrings() []string
	Update(func() error)

	// QueueTask(f func(ctx context.Context)) error
}

func (rp *RequestPanel) RerenderList() error {
	rp.Gui.Update(func() error {
		rp.View.Clear()
		table := lo.Map(rp.GetItems(), func(item *commands.Request, index int) []string {
			return rp.GetRenderList(item)
		})
		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		// Idk why a new line is prepended to the render string, slice is a temporary fix
		fmt.Fprint(rp.View, renderedTable[1:])

		// if rp.OnRerender != nil {
		// 	if err := rp.OnRerender(); err != nil {
		// 		return err
		// 	}
		// }

		if rp.Gui.IsCurrentView(rp.View) {
			return rp.HandleSelect()
		}
		return nil
	})

	return nil
}

func (rp *RequestPanel) GetView() *gocui.View {
	return rp.View
}

func (rp *RequestPanel) SetItems(items []*commands.Request) {
	rp.allItems = items
	rp.indices = make([]int, len(items))
	for i := range rp.indices {
		rp.indices[i] = i
	}
}

func (rp *RequestPanel) GetItems() []*commands.Request {
	result := make([]*commands.Request, len(rp.indices))
	for i, index := range rp.indices {
		result[i] = rp.allItems[index]
	}
	return result
}

func (rp *RequestPanel) Len() int {
	return len(rp.indices)
}

func (rp *RequestPanel) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if rp.Len() > 0 {
		clampedValue = Clamp(value, 0, rp.Len()-1)
	}

	rp.SelectedIdx = clampedValue
}

func (rp *RequestPanel) clampSelectedLineIdx() {
	clamped := Clamp(rp.SelectedIdx, 0, rp.Len()-1)

	if clamped != rp.SelectedIdx {
		rp.SelectedIdx = clamped
	}
}

func (rp *RequestPanel) HandleNextLine() error {
	rp.SelectNextLine()

	return rp.HandleSelect()
}

func (rp *RequestPanel) HandlePrevLine() error {
	rp.SelectPrevLine()

	return rp.HandleSelect()
}

func (rp *RequestPanel) GetSelectedItem() (*commands.Request, error) {
	var zero *commands.Request

	item, ok := rp.TryGet(rp.SelectedIdx)
	if !ok {
		// could probably have a better error here
		return zero, errors.New(EMPTY_REQUESTS)
	}

	return item, nil
}

func (rp *RequestPanel) TryGet(index int) (*commands.Request, bool) {
	if index < 0 || index >= len(rp.indices) {
		var zero *commands.Request
		return zero, false
	}

	return rp.allItems[rp.indices[index]], true
}

func (rp *RequestPanel) HandleSelect() error {
	item, err := rp.GetSelectedItem()
	if err != nil {
		if err.Error() != EMPTY_REQUESTS {
			return err
		}

		// if rp.NoItemsMessage != "" {
		// 	rp.Gui.NewSimpleRenderStringTask(func() string { return rp.NoItemsMessage })
		// }

		return nil
	}

	rp.Refocus()

	return rp.renderContext(item)
}

func (rp *RequestPanel) renderContext(request *commands.Request) error {
	if rp.prevIdx == rp.SelectedIdx {
		return nil
	}
    rp.prevIdx = rp.SelectedIdx


	urlView := rp.Gui.GetUrlView()
	urlView.Clear()
	fmt.Fprint(urlView, request.Url)

	// task := rp.ContextState.GetCurrentMainTab().Render(request)
	paramsView := rp.Gui.GetParamsView()
	paramsView.Clear()
    params, err := request.GetParams()
    if err != nil {
        return err
    }
    table := utils.MapToSlice(utils.ValuesToMap(params))
    renderedTable, err := utils.RenderComponent(table)
    fmt.Fprint(paramsView, renderedTable[1:])

	return nil //ro.Gui.QueueTask(task)
}

func (rp *RequestPanel) Refocus() {
	rp.Gui.FocusY(rp.SelectedIdx, rp.Len(), rp.View)
}

// moves the cursor up or down by the given amount (up for negative values)
func (rp *RequestPanel) moveSelectedLine(delta int) {
	rp.SetSelectedLineIdx(rp.SelectedIdx + delta)
}

func (rp *RequestPanel) SelectNextLine() {
	rp.moveSelectedLine(1)
}

func (rp *RequestPanel) SelectPrevLine() {
	rp.moveSelectedLine(-1)
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
