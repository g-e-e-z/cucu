package gui

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/gui/components"
	"github.com/g-e-e-z/cucu/gui/presentation"
	"github.com/g-e-e-z/cucu/gui/types"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jesseduffield/gocui"
	"github.com/spkg/bom"
)

func (gui *Gui) getRequestsPanel() *components.ListComponent[*commands.Request] {
	return &components.ListComponent[*commands.Request]{
		Log:  gui.Log,
		View: gui.Views.Requests,
		RequestContext: &components.RequestContext[*commands.Request]{
            GetUrlTab:  func() components.Tab[*commands.Request] {
                return components.Tab[*commands.Request]{
                	Key:   "url",
                	Title: "Url",
                	Render: gui.renderUrl,
                }
            },
			GetRequestInfoTabs: func() []components.Tab[*commands.Request] {
				return []components.Tab[*commands.Request]{
					// {
					// 	Key:    "headers",
					// 	Title:  "Params",
					// 	Render: func(*commands.Request) {},
					// },
					{
						Key:    "params",
						Title:  "Params",
						Render: gui.renderRequestParams,
					},
					{
						Key:    "body",
						Title:  "Body",
						Render: gui.renderRequestBody,
					},
				}
			},
			GetResponseInfoTabs: func() []components.Tab[*commands.Request] {
				return []components.Tab[*commands.Request]{
					{
						Key:    "body",
						Title:  "Body",
						Render: gui.renderResponseBody,
					},
					{
						Key:    "headers",
						Title:  "Headers",
						Render: gui.renderResponseHeaders,
					},
				}
			},
		},
		ListPanel: components.ListPanel[*commands.Request]{
			List:           components.NewFilteredList[*commands.Request](),
			View:           gui.Views.Requests,
			NoItemsMessage: "No Requests Loaded",
		},
		Gui:            gui.toInterface(),
		GetTableCells: func(request *commands.Request) []string {
			return presentation.GetRequestStrings(request)
		},
	}
}

// Rendering
func (gui *Gui) renderRequests() error {
	requests, err := gui.HttpCommands.GetRequests()
	if err != nil {
		return err
	}
	gui.Components.Requests.SetItems(requests)
	return gui.Components.Requests.RerenderList()
}

func (gui *Gui) reRenderRequests() error {
	requests := gui.Components.Requests.GetItems()
	gui.Components.Requests.SetItems(requests)
	return gui.Components.Requests.RerenderList()
}

func (gui *Gui) renderUrl() error {
    request, err := gui.Components.Requests.GetSelectedItem()
    if err != nil {
        return err
    }

    urlView := gui.Views.Url
	urlView.ClearTextArea()
	output := string(bom.Clean([]byte(request.Url)))
	s := utils.NormalizeLinefeeds(output)
	urlView.TextArea.TypeString(s)
	urlView.SetCursor(len(s), 0)
	fmt.Fprint(urlView, s)
    // Below sets origin and cursor to 0, not friendly with editing
    // gui.renderString(gui.g, gui.Views.Url.Name(), s)
    return nil
}

func (gui *Gui) renderRequestParams() error {
    request, err := gui.Components.Requests.GetSelectedItem()
    if err != nil {
        return err
    }
    params, err := request.GetParams()
    if err != nil {
        // TODO: This better
        gui.renderString(gui.g, gui.Views.RequestInfo.Name(), "")
        return err
    }
    table := utils.MapToSlice(utils.ValuesToMap(params))
    renderedTable, err := utils.RenderComponent(table)

    gui.renderString(gui.g, gui.Views.RequestInfo.Name(), renderedTable)
    return nil
}

func (gui *Gui) renderRequestBody() error {
    request, err := gui.Components.Requests.GetSelectedItem()
    if err != nil {
        return err
    }
    if request.Data == nil {
        // TODO: This better
        gui.renderString(gui.g, gui.Views.RequestInfo.Name(), "")
        return errors.New("No Request Data")
    }
    params, err := request.GetData()
    if err != nil {
        return err
    }
    table := utils.MapToSlice(params)
    renderedTable, err := utils.RenderComponent(table)

    gui.renderString(gui.g, gui.Views.RequestInfo.Name(), renderedTable)
    return nil
}
func (gui *Gui) renderResponseHeaders() error {
    // if request.ResponseBody == "" {
    //     // TODO: This better
    //     gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), "")
    //     return
    // }
    // gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), request.ResponseBody)
    return nil
}

func (gui *Gui) renderResponseBody() error {
    request, err := gui.Components.Requests.GetSelectedItem()
    if err != nil {
        return err
    }
    if request.ResponseBody == "" {
        // TODO: This better
        gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), "")
        return errors.New("No Response Body")
    }
    gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), request.ResponseBody)
    return nil
}


// Keybind Actions

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func (gui *Gui) handleEditMethod(_ *gocui.Gui, v *gocui.View) error {
    request, err := gui.Components.Requests.GetSelectedItem()
    if err != nil {
        return err
    }

    handleMenuPress := func (method string) error {
        request.Method = method
        return nil

    }

    var menuItems []*types.MenuItem
    var currentIndex int
    for i, method := range httpMethods {
        if method == request.Method {
            currentIndex = i
        }
        menuItems = append(menuItems, &types.MenuItem{
        		Label:        method,
        		OnPress: func() error {return handleMenuPress(method) },
        	},
        )
    }

	return gui.Menu(CreateMenuOptions{
		Title:      "Change Request Method",
		Items:      menuItems,
		Index:      currentIndex,
		HideCancel: true,
	})
}

func (gui *Gui) handleNewRequest(g *gocui.Gui, v *gocui.View) error {
	gui.Log.Info("Creating New Request")
	newRequest := &commands.Request{
		Name:        "NewRequest!",
		Url:         "this is a placeholder string",
		Method:      http.MethodGet,
		Log:         gui.Log,
		HttpCommand: gui.HttpCommands,
	}
	newRequestList := append(gui.Components.Requests.GetItems(), newRequest)
	gui.Components.Requests.SetItems(newRequestList)

	return gui.reRenderRequests()
}

func (gui *Gui) handleRequestSend(g *gocui.Gui, v *gocui.View) error {
	// TODO: This is a weird way to handle the no items string, fix later
	request, err := gui.Components.Requests.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.SendRequest(request)
}

func (gui *Gui) SendRequest(request *commands.Request) error {
	err := request.Send()
	if err != nil {
		return err
	}
	return gui.Components.Requests.RerenderList()
}
