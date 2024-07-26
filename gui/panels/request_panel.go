package panels

import (
	"fmt"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jroimartin/gocui"
	"github.com/samber/lo"
)

type IGui interface {
	Update(func() error)
}

type RequestPanel struct {
	View     *gocui.View
	Requests []*commands.Request
	Gui      IGui

	indicies []int
	reqIndex int
}

func (rq *RequestPanel) SetRequests() error {
	rq.Requests = []*commands.Request{
		{
			Name:   "Req1",
			Url:    "localhost:6969",
			Method: "GET",
		},
		{
			Name:   "Link1",
			Url:    "myspace.com",
			Method: "GET",
		},
	}
	return nil
}

func (rq *RequestPanel) GetRequests() []*commands.Request {
    return rq.Requests
}

func (rq *RequestPanel) GetRequestStrings(request *commands.Request) []string {
	return []string{request.Method, request.Name}
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

		// if rq.Gui.IsCurrentView(rq.View) {
		// 	return rq.HandleSelect()
		// }
		return nil
	})

	return nil
}
