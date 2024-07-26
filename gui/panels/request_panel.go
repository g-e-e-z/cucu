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
	Update(func() error)
}

type RequestPanel struct {
	Log            *logrus.Entry
	View           *gocui.View
	Requests       []*commands.Request
	Gui            IGui
	NoItemsMessage string

	indices  []int
	ReqIndex int
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

	item, ok := rq.TryGet(rq.ReqIndex)
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

	// rq.Refocus()

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

	return nil
}
