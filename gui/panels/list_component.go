package panels

import (
	"fmt"

	"github.com/g-e-e-z/cucu/commands"
	"github.com/g-e-e-z/cucu/utils"
	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spkg/bom"
)

type IGui interface {
	IsCurrentView(*gocui.View) bool
	GetUrlView() *gocui.View
	GetParamsView() *gocui.View
	GetResponseView() *gocui.View
	FocusY(selectedLine int, itemCount int, view *gocui.View)
	Update(func() error)
}

// This should be one more level of generalization
type ListComponent struct {
	Log  *logrus.Entry
	View *gocui.View
	ListPanel[*commands.Request]

	Gui            IGui
	NoItemsMessage string
}

func (self *ListComponent) GetView() *gocui.View {
	return self.View
}

func (self *ListComponent) GetRequestStrings(request *commands.Request) []string {
	return []string{request.Method, request.Name}
}

func (self *ListComponent) HandleSelect() error {
	item, err := self.GetSelectedItem(self.NoItemsMessage)
	if err != nil {
		if err.Error() != self.NoItemsMessage {
			return err
		}

		if self.NoItemsMessage != "" {
			self.Log.Warn(self.NoItemsMessage)
		}

		return nil
	}

	self.Refocus()

	return self.renderContext(item)
}

func (self *ListComponent) Rerender() error {
	self.Gui.Update(func() error {
		self.View.Clear()
		table := lo.Map(self.GetItems(), func(req *commands.Request, index int) []string {
			return self.GetRequestStrings(req)
		})
		renderedTable, err := utils.RenderComponent(table)
		if err != nil {
			return err
		}
		fmt.Fprint(self.View, renderedTable)

		// TODO: Find work around to get this back in/ evalute if its problematic being commented out: Figure out all callers
		// if self.Gui.IsCurrentView(self.View) {
		// 	return self.HandleSelect()
		// }
		// return nil
		return self.HandleSelect()
	})

	return nil
}

func (self *ListComponent) renderContext(request *commands.Request) error {
	// if self.ContextState == nil {
	// 	return nil
	// }

	// In lazydocker this is the tab names in the main view, will use something similar in future iterations
	// key := self.ContextState.GetCurrentContextKey(item)
	// if !self.Gui.ShouldRefresh(key) {
	// 	return nil
	// }

	// mainView := self.Gui.GetMainView()
	// mainView.Tabs = self.ContextState.GetMainTabTitles()
	// mainView.TabIndex = self.ContextState.mainTabIdx

	// task := self.ContextState.GetCurrentMainTab().Render(item)
	// return self.Gui.QueueTask(task)

	// TODO: Don't write directly, this whole block is questionable, TextArea etc..
	urlView := self.Gui.GetUrlView()
	urlView.ClearTextArea()
	output := string(bom.Clean([]byte(request.Url)))
	s := utils.NormalizeLinefeeds(output)
	urlView.TextArea.TypeString(s)
	urlView.SetCursor(len(s), 0)
	fmt.Fprint(urlView, s)

	paramsView := self.Gui.GetParamsView()
	paramsView.Clear()
	params, err := request.GetParams()
	if err != nil {
		return err
	}
	table := utils.MapToSlice(utils.ValuesToMap(params))
	renderedTable, err := utils.RenderComponent(table)
	fmt.Fprint(paramsView, renderedTable)

	responseView := self.Gui.GetResponseView()
	responseView.Clear()
	fmt.Fprint(responseView, request.ResponseBody)

	return nil
}

// Keybinding related
func (rp *ListComponent) Refocus() {
	rp.Gui.FocusY(rp.SelectedIdx, rp.List.Len(), rp.View)
}

func (rp *ListComponent) HandleNextLine() error {
	rp.SelectNextLine()

	return rp.HandleSelect()
}

func (rp *ListComponent) HandlePrevLine() error {
	rp.SelectPrevLine()

	return rp.HandleSelect()
}

func (rp *ListComponent) SelectNextLine() {
	rp.moveSelectedLine(1)
}

func (rp *ListComponent) SelectPrevLine() {
	rp.moveSelectedLine(-1)
}

func (rp *ListComponent) moveSelectedLine(delta int) {
	rp.SetSelectedLineIdx(rp.SelectedIdx + delta)
}

func (rp *ListComponent) SetSelectedLineIdx(value int) {
	clampedValue := 0
	if rp.List.Len() > 0 {
		clampedValue = Clamp(value, 0, rp.List.Len()-1)
	}

	rp.SelectedIdx = clampedValue
}

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
