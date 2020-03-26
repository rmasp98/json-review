package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
)

// Window contains all information for populating the window
type Window struct {
	views              map[ViewEnum]Layout
	panelRelativeWidth float32
	border             int
	tbBaseBuffer       int
	saveVisible        bool
	height             int
}

var (
	instance     Window
	once         sync.Once
	helpContents = "Ctrl+C: Exit  | Tab: Next View | Ctrl+S: Save"
)

// GetWindow creates new window if note created and returns window
func GetWindow() *Window {
	once.Do(func() {
		instance = Window{map[ViewEnum]Layout{
			PANEL:   Layout{MinSize{10, 1}, Dimensions{0, 0, 0, 0}, "", nil},
			SEARCH:  Layout{MinSize{1, 2}, Dimensions{0, 0, 0, 0}, "", nil},
			DISPLAY: Layout{MinSize{1, 1}, Dimensions{0, 0, 0, 0}, "", nil},
			HELP:    Layout{MinSize{1, 2}, Dimensions{0, 0, 0, 0}, helpContents, nil},
			SAVE:    Layout{MinSize{1, 1}, Dimensions{50, 30, 100, 32}, "", nil}},
			0.2, 1, 3, false, 0}
	})
	return &instance
}

// ShowSaveView a
func (w *Window) ShowSaveView(visible bool) {
	w.saveVisible = visible
}

// Height returns height of view for use in GetNodes and GetJson
func (w Window) Height(view ViewEnum) int {
	return w.GetDimensions(view).Y1 - w.GetDimensions(view).Y0 - 1
}

// GetDimensions gets the dimensions for a given view
func (w Window) GetDimensions(view ViewEnum) Dimensions {
	return w.views[view].dim
}

// GetContent returns the content of a view
func (w Window) GetContent(view ViewEnum) string {
	return w.views[view].content
}

// Resize updates views based on the current window size
func (w *Window) Resize(maxX, maxY int) {
	w.height = maxY
	tbBuffer := w.tbBaseBuffer
	if maxY < 4*tbBuffer || maxX < w.views[SEARCH].min.width+2*w.border {
		tbBuffer = 0
	}

	panelWidth := int(w.panelRelativeWidth * float32(maxX))
	if panelWidth < w.views[PANEL].min.width || maxY < w.views[SEARCH].min.height+2*w.border {
		panelWidth = 0
	}

	if panelWidth != 0 {
		w.updateViewDimensions(PANEL, Dimensions{
			w.border, w.border,
			panelWidth, maxY - tbBuffer - w.border,
		})
	} else {
		w.updateViewDimensions(PANEL, Dimensions{0, 0, 0, 0})
	}
	if tbBuffer != 0 {
		w.updateViewDimensions(SEARCH, Dimensions{
			w.border + panelWidth, w.border,
			maxX - w.border, tbBuffer,
		})
		w.updateViewDimensions(HELP, Dimensions{
			w.border, maxY - tbBuffer,
			maxX - w.border, maxY - w.border,
		})
	} else {
		w.updateViewDimensions(SEARCH, Dimensions{0, 0, 0, 0})
		w.updateViewDimensions(HELP, Dimensions{0, 0, 0, 0})
	}
	if maxX > w.views[DISPLAY].min.width+2*w.border && maxY > w.views[DISPLAY].min.height+2*w.border {
		w.updateViewDimensions(DISPLAY, Dimensions{
			w.border + panelWidth, tbBuffer + w.border,
			maxX - w.border, maxY - tbBuffer - w.border,
		})
	} else {
		w.updateViewDimensions(DISPLAY, Dimensions{0, 0, 0, 0})
	}
}

// SetViews passes all the view information to gocui
func (w Window) SetViews(gui GoCui) error {
	for name, view := range w.views {
		if view.dim != (Dimensions{0, 0, 0, 0}) {
			if gView, err := gui.SetView(name.String(), view.dim.X0, view.dim.Y0, view.dim.X1, view.dim.Y1); err != gocui.ErrUnknownView {
				gView.Title = name.String()
				if view.content != "" {
					gView.Clear()
					fmt.Fprint(gView, view.content)
					updateCursorPosition(gView, view.content)
				}
				if view.editor != nil {
					gView.Editable = true
					gView.Editor = view.editor
					//Bodge to get good highlighting on panel
					if name == PANEL {
						gView.Highlight = true
						gView.SelBgColor = gocui.ColorGreen
						gView.SelFgColor = gocui.ColorBlack
					}
				}
			}
		}
	}
	if !w.saveVisible {
		gui.SetViewOnBottom(SAVE.String())
	} else {
		gui.SetCurrentView(SAVE.String())
		gui.SetViewOnTop(SAVE.String())
	}
	return nil
}

func (w *Window) updateViewDimensions(view ViewEnum, dim Dimensions) {
	newView := w.views[view]
	newView.dim = dim
	w.views[view] = newView
}

// UpdateViewContent updates the content of the view
func (w *Window) UpdateViewContent(view ViewEnum, content string) {
	newView := w.views[view]
	newView.content = content
	w.views[view] = newView
}

// UpdateEditor sets the editor of the view
func (w *Window) UpdateEditor(view ViewEnum, editor gocui.Editor) {
	newView := w.views[view]
	newView.editor = editor
	w.views[view] = newView
}

func updateCursorPosition(view *gocui.View, content string) {
	_, yCursor := view.Cursor()
	_, yOrigin := view.Origin()
	if strings.Count(content, "\n") < yCursor+yOrigin {
		view.SetCursor(0, 0)
		view.SetOrigin(0, strings.Count(content, "\n"))
	}
}

//////////////////////////////////////////////////////////////////////////
// View data

// MinSize defines the smallest size a view can be before it disappears
type MinSize struct {
	width, height int
}

// Dimensions defines the current dimensions of the view
type Dimensions struct {
	X0, Y0, X1, Y1 int
}

// Layout contains size and content information for a view
type Layout struct {
	min     MinSize
	dim     Dimensions
	content string
	editor  gocui.Editor
}
