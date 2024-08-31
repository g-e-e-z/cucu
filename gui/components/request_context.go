package components

import (
	"github.com/samber/lo"
)

// the active request determines the entire state of the application. Will be used
// to render request parameters/ body/ headers and response
type RequestContext[T any] struct {
	requestTabIdx  int
	responseTabIdx int
	// this function returns the tabs that we can display for an item (the tabs
	// are shown on the request and response views)
	GetUrlTab           func() Tab[T]
	GetRequestInfoTabs  func() []Tab[T]
	GetResponseInfoTabs func() []Tab[T]
	// This tells us whether we need to re-render to the main panel for a given item.
	// This should include the item's ID and if you want to invalidate the cache for
	// some other reason, you can add that to the key as well (e.g. the container's state).
	// GetItemContextCacheKey func(item T) string
}

type Tab[T any] struct {
	// key used as part of the context cache key
	Key string
	// title of the tab, rendered in the respective view
	Title string
	// function to render the content of the tab
	Render func() error // tasks.TaskFunc
}

func (rc *RequestContext[T]) RenderUrl() {
	rc.GetUrlTab().Render()
}

// // Request Tab
func (rc *RequestContext[T]) GetRequestInfoTabTitles() []string {
	return lo.Map(rc.GetRequestInfoTabs(), func(tab Tab[T], _ int) string {
		return tab.Title
	})
}

func (rc *RequestContext[T]) GetCurrentRequestInfoTab() Tab[T] {
	return rc.GetRequestInfoTabs()[rc.requestTabIdx]
}

func (rp *RequestContext[T]) HandleNextReqTab() {
	tabs := rp.GetRequestInfoTabs()

	if len(tabs) == 0 {
		return
	}

	rp.requestTabIdx = (rp.requestTabIdx + 1) % len(tabs)
}

func (rp *RequestContext[T]) HandlePrevReqTab() {
	tabs := rp.GetRequestInfoTabs()

	if len(tabs) == 0 {
		return
	}

	rp.requestTabIdx = (rp.requestTabIdx - 1 + len(tabs)) % len(tabs)
}

// // Response Tab
func (rc *RequestContext[T]) GetResponseInfoTabTitles() []string {
	return lo.Map(rc.GetResponseInfoTabs(), func(tab Tab[T], _ int) string {
		return tab.Title
	})
}

func (rc *RequestContext[T]) GetCurrentResponseInfoTab() Tab[T] {
	return rc.GetResponseInfoTabs()[rc.responseTabIdx]
}

func (rp *RequestContext[T]) HandleNextResTab() {
	tabs := rp.GetResponseInfoTabs()

	if len(tabs) == 0 {
		return
	}

	rp.responseTabIdx = (rp.responseTabIdx + 1) % len(tabs)
}

func (rp *RequestContext[T]) HandlePrevResTab() {
	tabs := rp.GetResponseInfoTabs()

	if len(tabs) == 0 {
		return
	}

	rp.responseTabIdx = (rp.responseTabIdx - 1 + len(tabs)) % len(tabs)
}
