package gui

import "github.com/jroimartin/gocui"

type Views struct {
  	Requests    *gocui.View
	RequestInfo *gocui.View
	Response    *gocui.View
}
