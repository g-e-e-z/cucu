package gui

import (
	"fmt"

	"github.com/g-e-e-z/cucu/utils"
	"github.com/jesseduffield/gocui"
	lcu "github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) CurrentView() *gocui.View {
	return gui.g.CurrentView()
}

func (gui *Gui) GetUrlView() *gocui.View {
	return gui.Views.Url
}

func (gui *Gui) GetRequestInfoView() *gocui.View {
	return gui.Views.RequestInfo
}

func (gui *Gui) GetResponseInfoView() *gocui.View {
	return gui.Views.ResponseInfo
}

func (gui *Gui) IsCurrentView(view *gocui.View) bool {
	return view == gui.CurrentView()
}

func (gui *Gui) FocusY(selectedY int, lineCount int, v *gocui.View) {
	gui.focusPoint(0, selectedY, lineCount, v)
}

func (gui *Gui) RenderErrorString(s string) {
	_ = gui.renderString(gui.g, "response", s)
}

func (gui *Gui) allViews() []*gocui.View {
	return []*gocui.View{
		gui.Views.Requests,
		gui.Views.Url,
		gui.Views.RequestInfo,
		gui.Views.ResponseInfo,
	}
}

// renderString resets the origin of a view and sets its content
func (gui *Gui) renderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			return nil // return gracefully if view has been deleted
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		return gui.setViewContent(v, s)
	})
	return nil
}

func (gui *Gui) nextView(g *gocui.Gui, v *gocui.View) error {
	viewNames := gui.viewNames()
	var focusedViewName string
	if v == nil || v.Name() == viewNames[len(viewNames)-1] {
		focusedViewName = viewNames[0]
	} else {
		viewName := v.Name()
		for i := range viewNames {
			if viewName == viewNames[i] {
				focusedViewName = viewNames[i+1]
				break
			}
			if i == len(viewNames)-1 {
				gui.Log.Info("not in list of views")
				return nil
			}
		}
	}
	focusedView, err := g.View(focusedViewName)
	if err != nil {
		panic(err)
	}
    // Don't leave a view on editable
    if v.Editable {
        gui.handleToggleEdit(g, v)
    }
	// gui.resetMainView()
	// return gui.switchFocus(focusedView)
	if _, err := gui.g.SetCurrentView(focusedView.Name()); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) previousView(g *gocui.Gui, v *gocui.View) error {
	viewNames := gui.viewNames()
	var focusedViewName string
	if v == nil || v.Name() == viewNames[0] {
		focusedViewName = viewNames[len(viewNames)-1]
	} else {
		viewName := v.Name()
		for i := range viewNames {
			if viewName == viewNames[i] {
				focusedViewName = viewNames[i-1]
				break
			}
			if i == len(viewNames)-1 {
				gui.Log.Info("not in list of views")
				return nil
			}
		}
	}
	focusedView, err := g.View(focusedViewName)
	if err != nil {
		panic(err)
	}
    // Don't leave a view on editable
    if v.Editable {
        gui.handleToggleEdit(g, v)
    }
	// gui.resetMainView()
	// return gui.switchFocus(focusedView)
	if _, err := gui.g.SetCurrentView(focusedView.Name()); err != nil {
		return err
	}
	return nil
}

// if the cursor down past the last item, move it to the last line
// nolint:unparam
func (gui *Gui) focusPoint(selectedX int, selectedY int, lineCount int, v *gocui.View) {
	if selectedY < 0 || selectedY > lineCount {
		return
	}
	ox, oy := v.Origin()
	originalOy := oy
	cx, cy := v.Cursor()
	originalCy := cy
	_, height := v.Size()

	ly := lcu.Max(height-1, 0)

	windowStart := oy
	windowEnd := oy + ly

	if selectedY < windowStart {
		oy = lcu.Max(oy-(windowStart-selectedY), 0)
	} else if selectedY > windowEnd {
		oy += (selectedY - windowEnd)
	}

	if windowEnd > lineCount-1 {
		shiftAmount := (windowEnd - (lineCount - 1))
		oy = lcu.Max(oy-shiftAmount, 0)
	}

	if originalOy != oy {
		_ = v.SetOrigin(ox, oy)
	}

	cy = selectedY - oy
	if originalCy != cy {
		_ = v.SetCursor(cx, selectedY-oy)
	}
}

func (gui *Gui) setViewContent(v *gocui.View, s string) error {
	v.Clear()
	fmt.Fprint(v, gui.cleanString(s))
	return nil
}

func (gui *Gui) cleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return utils.NormalizeLinefeeds(output)
}
