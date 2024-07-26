package panels

import (
	"errors"
	"fmt"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jroimartin/gocui"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type IGui interface {
	IsCurrentView(*gocui.View) bool
	GetUrlView() *gocui.View
	GetParamsView() *gocui.View
	GetResponseView() *gocui.View
	FocusY(selectedLine int, itemCount int, view *gocui.View)
	Update(func() error)
}

type RequestPanel struct {
	Log            *logrus.Entry
	View           *gocui.View
	Requests       []*commands.Request
	Gui            IGui
	NoItemsMessage string

	indices     []int
	SelectedIdx int
}

func (rq *RequestPanel) GetView() *gocui.View {
	return rq.View
}

func (rq *RequestPanel) SetRequests(requests []*commands.Request) error {
	rq.Requests = requests
	rq.indices = make([]int, len(rq.Requests))
	for i := range rq.indices {
		rq.indices[i] = i
	}

	return nil
}

func (rq *RequestPanel) GetRequests() []*commands.Request {
	return rq.Requests
}

func (rq *RequestPanel) TryGet(index int) (*commands.Request, bool) {
	if index < 0 || index >= len(rq.indices) {
		var zero *commands.Request
		return zero, false
	}

	return rq.Requests[rq.indices[index]], true
}

func (rq *RequestPanel) GetSelectedRequest() (*commands.Request, error) {
	var zero *commands.Request

	item, ok := rq.TryGet(rq.SelectedIdx)
	if !ok {
		// could probably have a better error here
		return zero, errors.New(rq.NoItemsMessage)
	}

	return item, nil
}

func (rq *RequestPanel) GetRequestStrings(request *commands.Request) []string {
	return []string{request.Method, request.Name}
}

func (rq *RequestPanel) HandleSelect() error {
	item, err := rq.GetSelectedRequest()
	if err != nil {
		if err.Error() != rq.NoItemsMessage {
			return err
		}

		if rq.NoItemsMessage != "" {
			rq.Log.Warn(rq.NoItemsMessage)
		}

		return nil
	}

	rq.Refocus()

	return rq.renderContext(item)
}

func (rq *RequestPanel) Rerender() error {
	rq.Gui.Update(func() error {
		rq.View.Clear()
		table := lo.Map(rq.GetRequests(), func(req *commands.Request, index int) []string {
			return rq.GetRequestStrings(req)
		})
		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(rq.View, renderedTable)

		if rq.Gui.IsCurrentView(rq.View) {
			return rq.HandleSelect()
		}
		return nil
	})

	return nil
}

func (rq *RequestPanel) renderContext(request *commands.Request) error {
	// if rq.ContextState == nil {
	// 	return nil
	// }

	// In lazydocker this is the tab names in the main view, will use something similar in future iterations
	// key := rq.ContextState.GetCurrentContextKey(item)
	// if !rq.Gui.ShouldRefresh(key) {
	// 	return nil
	// }

	// mainView := rq.Gui.GetMainView()
	// mainView.Tabs = rq.ContextState.GetMainTabTitles()
	// mainView.TabIndex = rq.ContextState.mainTabIdx

	// task := rq.ContextState.GetCurrentMainTab().Render(item)
	// return rq.Gui.QueueTask(task)

	// TODO: Don't write directly
	urlView := rq.Gui.GetUrlView()
	urlView.Clear()
	fmt.Fprint(urlView, request.Url)

	paramsView := rq.Gui.GetParamsView()
	paramsView.Clear()
	params, err := request.GetParams()
	if err != nil {
		return err
	}
	table := utils.MapToSlice(utils.ValuesToMap(params))
	renderedTable, err := utils.RenderComponent(table)
	fmt.Fprint(paramsView, renderedTable)

	responseView := rq.Gui.GetResponseView()
	responseView.Clear()
	fmt.Fprint(responseView, request.ResponseBody)

	return nil
}

// Keybinding related
func (rp *RequestPanel) Refocus() {
	// rp.Gui.FocusY(rp.SelectedIdx, rp.Len(), rp.View)
	rp.Gui.FocusY(rp.SelectedIdx, len(rp.Requests), rp.View)
}

func (rp *RequestPanel) HandleNextLine() error {
	rp.SelectNextLine()

	return rp.HandleSelect()
}

func (rp *RequestPanel) HandlePrevLine() error {
	rp.SelectPrevLine()

	return rp.HandleSelect()
}

func (rp *RequestPanel) SelectNextLine() {
	rp.moveSelectedLine(1)
}

func (rp *RequestPanel) SelectPrevLine() {
	rp.moveSelectedLine(-1)
}

func (rp *RequestPanel) moveSelectedLine(delta int) {
	rp.SetSelectedLineIdx(rp.SelectedIdx + delta)
}

func (rp *RequestPanel) SetSelectedLineIdx(value int) {
	clampedValue := 0
	// if rp.Len() > 0 {
	if len(rp.Requests) > 0 {
		clampedValue = Clamp(value, 0, len(rp.Requests)-1)
	}

	rp.SelectedIdx = clampedValue
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
