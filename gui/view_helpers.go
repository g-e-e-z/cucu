package gui

import (
	"fmt"

	"github.com/g-e-e-z/cucu/utils"
	lcu "github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jroimartin/gocui"
	"github.com/spkg/bom"
)

func (gui *Gui) CurrentView() *gocui.View {
	return gui.g.CurrentView()
}

func (gui *Gui) GetUrlView() *gocui.View {
	return gui.Views.Url
}

func (gui *Gui) GetParamsView() *gocui.View {
	return gui.Views.Params
}

func (gui *Gui) IsCurrentView(view *gocui.View) bool {
	return view == gui.CurrentView()
}

func (gui *Gui) FocusY(selectedY int, lineCount int, v *gocui.View) {
	gui.focusPoint(0, selectedY, lineCount, v)
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


