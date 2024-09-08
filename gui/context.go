package gui

// import (
// 	"github.com/samber/lo"
// )

// the active request determines the entire state of the application. Will be used
// to render request parameters/ body/ headers and response
type RequestContext struct {
	requestTabIdx  int
	responseTabIdx int

	GetUrlTab           func() Tab
	GetRequestInfoTabs  func() []Tab
	GetResponseInfoTabs func() []Tab
}

type Tab struct {
	// key used as part of the context cache key
	Key string
	// title of the tab, rendered in the respective view
	Title string
	// function to render the content of the tab
	Render func() error // tasks.TaskFunc
}
