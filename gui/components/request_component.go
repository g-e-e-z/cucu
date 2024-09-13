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

type RequestComponent struct {
	Log  *logrus.Entry
	View *gocui.View

	RequestContext *RequestContext[*commands.Request]
	ListPanel[*commands.Request]

	Gui            IGui

	// returns the cells that we render to the view in a table format. The cells will
	// be rendered with padding.
	GetTableCells func(*commands.Request) []string
}

func (rc *RequestComponent) GetView() *gocui.View {
	return rc.View
}

func (rc *RequestComponent) HandleSelect() error {
	_, err := rc.GetSelectedItem()
	if err != nil {
		if err.Error() != rc.NoItemsMessage {
			return err
		}

		if rc.NoItemsMessage != "" {
			rc.Log.Warn(rc.NoItemsMessage)
		}

		return nil
	}

	rc.Refocus()

	return rc.renderContext()
}

func (rc *RequestComponent) RerenderList() error {
	rc.Gui.Update(func() error {
		rc.View.Clear()
		table := lo.Map(rc.List.GetItems(), func(item *commands.Request, index int) []string {
			return rc.GetTableCells(item)
		})

		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(rc.View, renderedTable)

		// TODO: Find work around to get this back in/ evalute if its problematic being commented out: Figure out all callers
		if rc.Gui.IsCurrentView(rc.View) {
			return rc.HandleSelect()
		}
		return nil
		// return rc.HandleSelect()
	})

	return nil
}

func (rc *RequestComponent) renderContext() error {
	if rc.RequestContext == nil {
		return nil
	}

	rc.RequestContext.RenderUrl()

	requestInfoView := rc.Gui.GetRequestInfoView()
	requestInfoView.Tabs = rc.RequestContext.GetRequestInfoTabTitles()
	requestInfoView.TabIndex = rc.RequestContext.requestTabIdx
	rc.RequestContext.GetCurrentRequestInfoTab().Render()
	// task := rc.RequestContext.GetCurrentRequestInfoTab().Render(item)

	responseInfoView := rc.Gui.GetResponseInfoView()
	responseInfoView.Tabs = rc.RequestContext.GetResponseInfoTabTitles()
	responseInfoView.TabIndex = rc.RequestContext.responseTabIdx
	rc.RequestContext.GetCurrentResponseInfoTab().Render()
	// task := rc.RequestContext.GetCurrentResponseInfoTab().Render(item)

	return nil
}

// Keybinding related
func (rc *RequestComponent) Refocus() {
	rc.Gui.FocusY(rc.SelectedIdx, rc.List.Len(), rc.View)
}

func (rc *RequestComponent) HandleNextLine() error {
	rc.SelectNextLine()

	return rc.HandleSelect()
}

func (rc *RequestComponent) HandlePrevLine() error {
	rc.SelectPrevLine()

	return rc.HandleSelect()
}

func (rc *RequestComponent) SelectNextLine() {
	rc.moveSelectedLine(1)
}

func (rc *RequestComponent) SelectPrevLine() {
	rc.moveSelectedLine(-1)
}

func (rc *RequestComponent) moveSelectedLine(delta int) {
	rc.SetSelectedLineIdx(rc.SelectedIdx + delta)
}

func (rc *RequestComponent) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if rc.List.Len() > 0 {
		clampedValue = Clamp(value, 0, rc.List.Len()-1)
	}

	rc.SelectedIdx = clampedValue
}

func (rc *RequestComponent) HandleNextReqTab() error {
	if rc.RequestContext == nil {
		return nil
	}

	rc.RequestContext.HandleNextReqTab()

	return rc.HandleSelect()
}

func (rc *RequestComponent) HandlePrevReqTab() error {
	if rc.RequestContext == nil {
		return nil
	}

	rc.RequestContext.HandlePrevReqTab()

	return rc.HandleSelect()
}

func (rc *RequestComponent) HandleNextResTab() error {
	if rc.RequestContext == nil {
		return nil
	}

	rc.RequestContext.HandleNextResTab()

	return rc.HandleSelect()
}

func (rc *RequestComponent) HandlePrevResTab() error {
	if rc.RequestContext == nil {
		return nil
	}

	rc.RequestContext.HandlePrevResTab()

	return rc.HandleSelect()
}
