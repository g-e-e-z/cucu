package main

import (
	"fmt"
	"log"
	"math"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	maxX -= 1
	maxY -= 1
	if v, err := g.SetView("sidebar", 0, 0, int(math.Floor(.2 * float64(maxX))), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Requests")
	}
	if v, err := g.SetView("request", 1+int(math.Floor(.2 * float64(maxX))), 0, int(math.Floor(.6 * float64(maxX))), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Request")
	}
	if v, err := g.SetView("response", 1+int(math.Floor(.6 * float64(maxX))), 0, maxX, maxY); err != nil {
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
