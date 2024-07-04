package main

import (
	"fmt"
	"log"

    "github.com/g-e-e-z/cucu/app"
	"github.com/jroimartin/gocui"
)

func main() {
    app, err := app.NewApp()
	if err == nil {
		err = app.Run()
	}
	app.Close()

	app.Gui.SetManagerFunc(layout)

	if err := app.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := app.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	maxX -= 1
	maxY -= 1
	if v, err := g.SetView("sidebar", 0, 0, maxX/5, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Requests")
	}
	if v, err := g.SetView("request", 1+maxX/5, 0, 2*maxX/3, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Request")
	}
	if v, err := g.SetView("response", 1+2*maxX/3, 0, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Response")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
