package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
)

type CommandFunc func(*gocui.Gui, *gocui.View) error

var COMMANDS map[string]func(string, *App) CommandFunc = map[string]func(string, *App) CommandFunc{
	"nextView": func(_ string, a *App) CommandFunc {
		return a.NextView
	},
	"prevView": func(_ string, a *App) CommandFunc {
		return a.PrevView
	},
	"quit": func(_ string, _ *App) CommandFunc {
		return quit
	},
	"openEditor": func(_ string, a *App) CommandFunc {
		return func(g *gocui.Gui, v *gocui.View) error {
			return openEditor(g, v, a.config.General.Editor)
		}
	},
	"scrollDown": func(_ string, _ *App) CommandFunc {
		return scrollViewDown
	},
	"scrollUp": func(_ string, _ *App) CommandFunc {
		return scrollViewUp
	},
	"pageDown": func(_ string, _ *App) CommandFunc {
		return pageDown
	},
	"pageUp": func(_ string, _ *App) CommandFunc {
		return pageUp
	},
}

// Return to this when refactoring Views into ViewList
// moves the cursor up or down by the given amount (up for negative values)
// func moveSelectedLine(v *gocui.View, delta int) {
// 	v.SetSelectedLineIdx(v.SelectedIdx + delta)
// }
//
// func selectNextLine(v *gocui.View) {
//     moveSelectedLine(v, 1)
// }
//
// func selectPrevLine(v *gocui.View) {
//     moveSelectedLine(v, -1)
// }

func scrollView(v *gocui.View, dy int) error {
	v.Autoscroll = false
	ox, oy := v.Origin()
	if oy+dy < 0 {
		dy = -oy
	}
	if _, err := v.Line(dy); dy > 0 && err != nil {
		dy = 0
	}
	v.SetOrigin(ox, oy+dy)
	return nil
}

func scrollViewUp(_ *gocui.Gui, v *gocui.View) error {
	return scrollView(v, -1)
}

func scrollViewDown(_ *gocui.Gui, v *gocui.View) error {
	return scrollView(v, 1)
}

func pageUp(_ *gocui.Gui, v *gocui.View) error {
	_, height := v.Size()
	scrollView(v, -height*2/3)
	return nil
}

func pageDown(_ *gocui.Gui, v *gocui.View) error {
	_, height := v.Size()
	scrollView(v, height*2/3)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func openEditor(g *gocui.Gui, v *gocui.View, editor string) error {
	file, err := os.CreateTemp(os.TempDir(), "cucu-")
	if err != nil {
		return nil
	}
	defer os.Remove(file.Name())

	val := getViewValue(g, v.Name())
	if val != "" {
		fmt.Fprint(file, val)
	}
	file.Close()

	info, err := os.Stat(file.Name())
	if err != nil {
		return nil
	}

	cmd := exec.Command(editor, file.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	// sync termbox to reset console settings
	// this is required because the external editor can modify the console
	defer g.Update(func(_ *gocui.Gui) error {
		termbox.Sync()
		return nil
	})
	if err != nil {
		rv, _ := g.View(ERROR_VIEW)
		rv.Clear()
		fmt.Fprintf(rv, "Editor open error: %v", err)
		return nil
	}

	newInfo, err := os.Stat(file.Name())
	if err != nil || newInfo.ModTime().Before(info.ModTime()) {
		return nil
	}

	newVal, err := os.ReadFile(file.Name())
	if err != nil {
		return nil
	}

	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
	v.Clear()
	fmt.Fprint(v, strings.TrimSpace(string(newVal)))

	return nil
}
