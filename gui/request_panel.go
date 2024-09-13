package gui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/gui/components"
	"github.com/g-e-e-z/cucu/gui/presentation"
	"github.com/g-e-e-z/cucu/gui/types"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jesseduffield/gocui"
	"github.com/spkg/bom"
)

func (gui *Gui) getRequestsPanel() *components.RequestComponent {
	return &components.RequestComponent{
		Log:  gui.Log,
		View: gui.Views.Requests,
		RequestContext: &components.RequestContext[*commands.Request]{
			GetUrlTab: func() components.Tab[*commands.Request] {
				return components.Tab[*commands.Request]{
					Key:    "url",
					Title:  "Url",
					Render: gui.renderUrl,
				}
			},
			GetRequestInfoTabs: func() []components.Tab[*commands.Request] {
				return []components.Tab[*commands.Request]{
					{
						Key:    "headers",
						Title:  "Headers",
						Render: gui.renderRequestHeaders,
					},
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
		Gui: gui.toInterface(),
		GetTableCells: func(request *commands.Request) []string {
			return presentation.GetRequestStrings(request)
		},
	}
}

// Rendering
func (gui *Gui) loadRequests() error {
	requests, err := gui.OSCommands.GetRequests()
	if err != nil {
		return err
	}
    // TODO: Sloppy
    for _, req := range requests {
        req.Client = gui.Client
    }
	gui.Requests.SetItems(requests)
	return gui.Requests.RerenderList()
}

func (gui *Gui) reRenderRequests() error {
	requests := gui.Requests.GetItems()
	gui.Requests.SetItems(requests)
	return gui.Requests.RerenderList()
}

func (gui *Gui) renderUrl() error {
	request, err := gui.Requests.GetSelectedItem()
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
func (gui *Gui) renderRequestHeaders() error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
	headers, err := request.GetRequestHeaders()
	if err != nil {
		// TODO: This better
		gui.renderString(gui.g, gui.Views.RequestInfo.Name(), err.Error())
		return err
	}
	var formattedHeaders [][]string
	for k, v := range headers {
		strSlice := []string{k, fmt.Sprintf("%v", v)}
		formattedHeaders = append(formattedHeaders, strSlice)
	}

	renderedTable, err := utils.RenderComponent(formattedHeaders)

	gui.renderString(gui.g, gui.Views.RequestInfo.Name(), renderedTable)
	return nil
}

func (gui *Gui) renderRequestParams() error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
	params, err := request.GetParams()
	if err != nil {
		// TODO: This better
		gui.renderString(gui.g, gui.Views.RequestInfo.Name(), err.Error())
		return err
	}
	renderedTable, err := utils.RenderComponent(params)

	gui.renderString(gui.g, gui.Views.RequestInfo.Name(), renderedTable)
	return nil
}

func (gui *Gui) renderRequestBody() error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
	data, err := request.GetData()
	if err != nil {
		gui.renderString(gui.g, gui.Views.RequestInfo.Name(), err.Error())
		return err
	}

	// TODO: Wont always be rendered as JSON... x-www-form-urlencoded, graphql..
	// This should also modify how its presented
	jsonData, err := utils.ToJSON(data)
	if err != nil {
		return err
	}

	gui.renderString(gui.g, gui.Views.RequestInfo.Name(), jsonData)
	return nil
}
func (gui *Gui) renderResponseHeaders() error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
	// result := make(map[string]string)
    result := ""
	for key, value := range request.ResponseHeaders {
		result += key + " | " + strings.Join(value,",") + "\n"
	}
	gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), result)
	return nil
}

func (gui *Gui) renderResponseBody() error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
	if request.Status == "" {
		// TODO: This better
		gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), "")
		return nil
	}
	renderString := formatBody(request, gui.Views.ResponseInfo.Width())
	gui.renderString(gui.g, gui.Views.ResponseInfo.Name(), renderString)
	return nil
}

func formatBody(request *commands.Request, width int) string {
	// TODO: This is god awful -> make a func in utils that isnt shit
	formattedString := request.Status + " | " + fmt.Sprintf("%f", request.Duration.Seconds()) + " seconds\n"
	formattedString += strings.Repeat("=", width)
	formattedString += "\n"
	formattedString += request.ResponseBody
	return formattedString
}

// Keybind Actions

var httpMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodHead,
	http.MethodPatch,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func (gui *Gui) handleEditField(_ *gocui.Gui, v *gocui.View) error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}
	gui.Log.Info(v.Name())
	var editOpts CreateEditOptions

	if v.Name() == "requests" {
		editOpts = CreateEditOptions{
			Title: "Edit Name",
			Value: request.Name,
		}
	}
	if v.Name() == "url" {
		editOpts = CreateEditOptions{
			Title: "Edit Url",
			Value: request.Url,
		}
	}
	if v.Name() == "params" {
		tabTitle := gui.Requests.RequestContext.GetCurrentRequestInfoTab().Title
		switch tabTitle {
		case "Body":
			editOpts = CreateEditOptions{
				Title: "Edit Body",
				Value: request.DataToJSON(),
			}
		case "Params":
			editOpts = CreateEditOptions{
				Title: "Edit Params",
				Value: request.Url,
			}
		case "Headers":
			editOpts = CreateEditOptions{
				Title: "Edit Headers",
				Value: request.HeadersToJSON(),
			}
		}
	}

	return gui.Edit(editOpts)
}

func (gui *Gui) handleEditMethod(_ *gocui.Gui, v *gocui.View) error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return err
	}

	handleMenuPress := func(method string) error {
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
			Label:   method,
			OnPress: func() error { return handleMenuPress(method) },
		},
		)
	}

	return gui.OpenMenu(CreateMenuOptions{
		Title:      "Change Request Method",
		Items:      menuItems,
		Index:      currentIndex,
		HideCancel: true,
	})
}

func (gui *Gui) handleNewRequest(g *gocui.Gui, v *gocui.View) error {
	gui.Log.Info("Creating New Request")
	newRequest := commands.NewRequest(gui.Log, gui.Client)
	newRequestList := append(gui.Requests.GetItems(), newRequest)
	gui.Requests.SetItems(newRequestList)

	return gui.reRenderRequests()
}

func (gui *Gui) handleRequestSend(g *gocui.Gui, v *gocui.View) error {
	request, err := gui.Requests.GetSelectedItem()
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
	return gui.reRenderRequests()
}

func (gui *Gui) handleDeleteRequest(g *gocui.Gui, v *gocui.View) error {
	// Must delete from memory and file
	request, err := gui.Requests.RemoveSelectedItem()
	if err != nil {
		return nil
	}
	return gui.DeleteRequest(request)
}

func (gui *Gui) handleSaveRequest(g *gocui.Gui, v *gocui.View) error {
	request, err := gui.Requests.GetSelectedItem()
	if err != nil {
		return nil
	}
	return gui.SaveRequest(request)
}

func (gui *Gui) SaveRequest(request *commands.Request) error {
	err := request.Save()
	if err != nil {
		gui.Log.Warn("Error saving request: ", err.Error())
		return err
	}

	return gui.reRenderRequests()
}

func (gui *Gui) DeleteRequest(request *commands.Request) error {
	allRequests := gui.Requests.GetItems()
	err := request.Delete(allRequests)
	if err != nil {
		gui.Log.Warn("Error deleting request: ", err.Error())
		return err
	}

	return gui.reRenderRequests()
}
