package ui

// Window determines size of each view
type Window struct {
	views              map[ViewEnum]layout
	panelRelativeWidth float64
	border             int
	tbBaseBuffer       int
}

// NewWindow stuff
func NewWindow(panelRelativeWidth float64, border, tbBaseBuffer int) Window {
	return Window{map[ViewEnum]layout{
		PANEL:   newLayout(10, 1),
		DISPLAY: newLayout(1, 1),
		SEARCH:  newLayout(1, 2),
		HELP:    newLayout(1, 2),
		VIEW:    newLayout(10, 1),
	}, panelRelativeWidth, border, tbBaseBuffer}
}

// GetDimensions gets the dimensions for a given view
func (w Window) GetDimensions(view ViewEnum) (int, int, int, int) {
	v := w.views[view]
	return v.x0, v.y0, v.x1, v.y1
}

// Resize stuff
func (w *Window) Resize(maxWidth, maxHeight, searchHeight int) error {
	tbBuffer := w.tbBaseBuffer
	if maxHeight < 4*tbBuffer || maxWidth < w.views[SEARCH].minWidth+2*w.border {
		tbBuffer = 0
	}

	panelWidth := int(w.panelRelativeWidth * float64(maxWidth))
	if panelWidth < w.views[PANEL].minWidth || maxHeight < w.views[SEARCH].minHeight+2*w.border {
		panelWidth = 0
	}

	if panelWidth != 0 {
		w.updateViewDimensions(PANEL,
			w.border, tbBuffer+w.border,
			panelWidth, maxHeight-tbBuffer-w.border,
		)
		w.updateViewDimensions(VIEW,
			w.border, w.border,
			panelWidth, tbBuffer,
		)
	} else {
		w.updateViewDimensions(PANEL, 0, 0, 0, 0)
		w.updateViewDimensions(VIEW, 0, 0, 0, 0)
	}
	if tbBuffer != 0 {
		w.updateViewDimensions(SEARCH,
			w.border+panelWidth, w.border,
			maxWidth-w.border, searchHeight,
		)

		w.updateViewDimensions(HELP,
			w.border, maxHeight-tbBuffer,
			maxWidth-w.border, maxHeight-w.border,
		)
	} else {
		w.updateViewDimensions(SEARCH, 0, 0, 0, 0)
		w.updateViewDimensions(HELP, 0, 0, 0, 0)
	}
	if maxWidth > w.views[DISPLAY].minWidth+2*w.border && maxHeight > w.views[DISPLAY].minHeight+2*w.border {
		w.updateViewDimensions(DISPLAY,
			w.border+panelWidth, tbBuffer+w.border,
			maxWidth-w.border, maxHeight-tbBuffer-w.border,
		)
	} else {
		w.updateViewDimensions(DISPLAY, 0, 0, 0, 0)
	}

	// TODO: write tests for these
	// popupWidth := w.views[POPUP].x1 - w.views[POPUP].x0
	// popupHeight := w.views[POPUP].y1 - w.views[POPUP].y0
	// if popupWidth != 0 && popupHeight != 0 && popupWidth < maxWidth && popupHeight < maxHeight {
	// 	w.updateViewDimensions(POPUP,
	// 		maxWidth/2-popupWidth/2,
	// 		maxHeight/2-popupHeight/2,
	// 		maxWidth/2+popupWidth/2+(popupWidth%2),
	// 		maxHeight/2+popupHeight/2+(popupHeight%2),
	// 	)
	// }
	return nil
}

func (w *Window) updateViewDimensions(view ViewEnum, x0, y0, x1, y1 int) {
	v := w.views[view]
	v.x0, v.y0, v.x1, v.y1 = x0, y0, x1, y1
	w.views[view] = v
}

func newLayout(minWidth, minHeight int) layout {
	return layout{minWidth, minHeight, 0, 0, 0, 0}
}

type layout struct {
	minWidth  int
	minHeight int
	x0        int
	y0        int
	x1        int
	y1        int
}
