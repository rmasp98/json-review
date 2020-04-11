package ui

import (
	"kube-review/jsontree"
	"kube-review/search"

	"github.com/jroimartin/gocui"
)

// TODO: increase functionality of the basic editing

func basicEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyHome:
		v.SetCursor(0, 0)
	case key == gocui.KeyEnd:
		line, _ := v.Line(0)
		v.SetCursor(len(line), 0)
	}
}

// SearchEditor stuff
type SearchEditor struct {
	s        search.Search
	nodeList *jsontree.NodeList
}

// NewSearchEditor stuff
func NewSearchEditor(nodeList *jsontree.NodeList) *SearchEditor {
	return &SearchEditor{search.NewSearch(search.REGEX, "querylist.json"), nodeList}
}

// Edit stuff
func (e *SearchEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	updateTitle(v, &e.s)
	if _, lineNum := v.Cursor(); lineNum == 0 {
		basicEditor(v, key, ch, mod)
	} else {
		switch {
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyEnter:
			input, _ := v.Line(0)
			name := e.s.GetQueryName(input, lineNum-1)
			v.Clear()
			v.Write([]byte(name))
			v.SetCursor(len(name), 0)
		}
	}

	if key == gocui.KeyEnter {
		input, _ := v.Line(0)
		e.s.Execute(input, e.nodeList)
	} else if key == gocui.KeyCtrlQ {
		e.s.ToggleQueryMode()
		clearInput(v)
		updateTitle(v, &e.s)
		return
	} else if key == gocui.KeyCtrlF {
		e.s.ToggleSearchMode()
		clearInput(v)
		updateTitle(v, &e.s)
		return
	} else if key == gocui.KeyCtrlN {
		e.nodeList.FindNextHighlightedNode()
	}

	input, _ := v.Line(0)
	queryDetails := e.s.GetHints(input)
	GetWindow().UpdateViewContent(SEARCH, input+queryDetails)

	_, lineNum := v.Cursor()
	v.Highlight = lineNum != 0
}

func clearInput(v *gocui.View) {
	v.Clear()
	v.SetCursor(0, 0)
	GetWindow().UpdateViewContent(SEARCH, "")
}

func updateTitle(v *gocui.View, s *search.Search) {
	v.Title = "Search: Mode=" + s.GetModeInfo()
}
