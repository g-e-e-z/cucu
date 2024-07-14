package commands

import (
	"os/exec"
	// "github.com/g-e-e-z/cucu/app"
)

// Command holds all the commands
type Command struct {
	command  func(string, ...string) *exec.Cmd
	getenv   func(string) string
}

// type CommandFunc func(*gocui.Gui, *gocui.View) error
//
// var COMMANDS map[string]func(string, *app.App) CommandFunc = map[string]func(string, *app.App) CommandFunc{
// 	"nextView": func(_ string, a *app.App) CommandFunc {
// 		return a.NextView
// 	},
// 	"prevView": func(_ string, a *app.App) CommandFunc {
// 		return a.PrevView
// 	},
// 	"quit": func(_ string, _ *app.App) CommandFunc {
// 		return quit
// 	},
// 	"openEditor": func(_ string, a *app.App) CommandFunc {
// 		return func(g *gocui.Gui, v *gocui.View) error {
// 			return openEditor(g, v, a.Config.General.Editor)
// 		}
// 	},
// 	"scrollDown": func(_ string, _ *app.App) CommandFunc {
// 		return scrollViewDown
// 	},
// 	"scrollUp": func(_ string, _ *app.App) CommandFunc {
// 		return scrollViewUp
// 	},
// 	"pageDown": func(_ string, _ *app.App) CommandFunc {
// 		return pageDown
// 	},
// 	"pageUp": func(_ string, _ *app.App) CommandFunc {
// 		return pageUp
// 	},
// }
//
// // Return to this when refactoring Views into ViewList
// // moves the cursor up or down by the given amount (up for negative values)
// // func moveSelectedLine(v *gocui.View, delta int) {
// // 	v.SetSelectedLineIdx(v.SelectedIdx + delta)
// // }
// //
// // func selectNextLine(v *gocui.View) {
// //     moveSelectedLine(v, 1)
// // }
// //
// // func selectPrevLine(v *gocui.View) {
// //     moveSelectedLine(v, -1)
// // }
//
//
// // func openEditor(g *gocui.Gui, v *gocui.View, editor string) error {
// // 	file, err := os.CreateTemp(os.TempDir(), "cucu-")
// // 	if err != nil {
// // 		return nil
// // 	}
// // 	defer os.Remove(file.Name())
// //
// // 	val := getViewValue(g, v.Name())
// // 	if val != "" {
// // 		fmt.Fprint(file, val)
// // 	}
// // 	file.Close()
// //
// // 	info, err := os.Stat(file.Name())
// // 	if err != nil {
// // 		return nil
// // 	}
// //
// // 	cmd := exec.Command(editor, file.Name())
// // 	cmd.Stdout = os.Stdout
// // 	cmd.Stdin = os.Stdin
// // 	cmd.Stderr = os.Stderr
// // 	err = cmd.Run()
// // 	// sync termbox to reset console settings
// // 	// this is required because the external editor can modify the console
// // 	defer g.Update(func(_ *gocui.Gui) error {
// // 		termbox.Sync()
// // 		return nil
// // 	})
// // 	if err != nil {
// // 		rv, _ := g.View(ERROR_VIEW)
// // 		rv.Clear()
// // 		fmt.Fprintf(rv, "Editor open error: %v", err)
// // 		return nil
// // 	}
// //
// // 	newInfo, err := os.Stat(file.Name())
// // 	if err != nil || newInfo.ModTime().Before(info.ModTime()) {
// // 		return nil
// // 	}
// //
// // 	newVal, err := os.ReadFile(file.Name())
// // 	if err != nil {
// // 		return nil
// // 	}
// //
// // 	v.SetCursor(0, 0)
// // 	v.SetOrigin(0, 0)
// // 	v.Clear()
// // 	fmt.Fprint(v, strings.TrimSpace(string(newVal)))
// //
// // 	return nil
// // }
