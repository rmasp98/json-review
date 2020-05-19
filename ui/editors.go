package ui

import (
	"kube-review/jsontree"
	"kube-review/search"

	"github.com/awesome-gocui/gocui"
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
		cursorPos, _ := v.Cursor()
		line, _ := v.Line(0)
		if cursorPos != len(line) {
			v.EditDelete(false)
		}
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
	s               search.Search
	nodeList        *jsontree.NodeList
	searchCursorPos int
}

// NewSearchEditor stuff
func NewSearchEditor(nodeList *jsontree.NodeList) *SearchEditor {
	return &SearchEditor{search.NewSearch(search.REGEX, "querylist.json"), nodeList, 0}
}

// Edit stuff
func (e *SearchEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if _, lineNum := v.Cursor(); lineNum == 0 {
		e.searchCursorPos, _ = v.Cursor()
		basicEditor(v, key, ch, mod)
	} else {
		switch {
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyEnter:
			input, _ := v.Line(0)
			searchLine, newCursorPos := e.s.InsertSelectedHint(input, e.searchCursorPos, lineNum-1)
			v.Clear()
			v.Write([]byte(searchLine))
			v.SetCursor(newCursorPos, 0)
			v.Highlight = false
			return
		}
	}

	if key == gocui.KeyEnter {
		input, _ := v.Line(0)
		v.Clear()
		v.Write([]byte(input))
		if err := e.s.Execute(input, e.nodeList); err != nil {
			v.Write([]byte("\n" + err.Error()))
			return
		}
		return
	} else if key == gocui.KeyCtrlQ {
		e.s.ToggleQueryMode()
		clearInput(v)
		cui.UpdateViewTitle(SEARCH, "Search: Mode="+e.s.GetModeInfo())
		return
	} else if key == gocui.KeyCtrlF {
		e.s.ToggleSearchMode()
		clearInput(v)
		cui.UpdateViewTitle(SEARCH, "Search: Mode="+e.s.GetModeInfo())
		return
	} else if key == gocui.KeyCtrlN {
		e.nodeList.FindNextHighlightedNode()
	} else if key == gocui.KeyCtrlR {
		e.nodeList.Clear()
	}

	if cursorPos, y := v.Cursor(); y == 0 {
		input, _ := v.Line(0)
		hints := e.s.GetHints(input, cursorPos)
		v.Clear()
		v.Write([]byte(input + hints))
	}

	_, lineNum := v.Cursor()
	v.Highlight = lineNum != 0
	cui.SetCursor(lineNum == 0)
}

func clearInput(v *gocui.View) {
	v.Clear()
	v.SetCursor(0, 0)
}

// NodesEditor is the editor for the PANEL view
type NodesEditor struct {
	nodeList *jsontree.NodeList
}

// NewNodesEditor creates a new nodesEditor object
func NewNodesEditor(nodeList *jsontree.NodeList) *NodesEditor {
	return &NodesEditor{nodeList}
}

// Edit provides repsonses to given key presses for PANEL
func (e *NodesEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyArrowUp:
		_, oldY := v.Cursor()
		v.MoveCursor(0, -1, false)
		_, newY := v.Cursor()
		if oldY == newY {
			e.nodeList.MoveTopNode(-1)
		}
		e.nodeList.SetActiveNode(newY)
	case key == gocui.KeyArrowDown:
		_, oldY := v.Cursor()
		v.MoveCursor(0, 1, false)
		_, newY := v.Cursor()
		if oldY == newY {
			e.nodeList.MoveTopNode(1)
		}
		e.nodeList.SetActiveNode(newY)
	case key == gocui.KeyArrowRight:
		x, y := v.Origin()
		v.SetOrigin(x+1, y)
	case key == gocui.KeyArrowLeft:
		x, y := v.Origin()
		v.SetOrigin(x-1, y)
	case ch == 'e':
		e.nodeList.ExpandActiveNode()
	case ch == 'c':
		newY := e.nodeList.CollapseActiveNode()
		v.SetCursor(0, newY)
	case key == gocui.KeyPgup:
		e.nodeList.MoveTopNode(-25)
	case key == gocui.KeyPgdn:
		e.nodeList.MoveTopNode(25)
	}
}

// DisplayEditor stuff
type DisplayEditor struct {
	nodeList *jsontree.NodeList
}

// NewDisplayEditor stuff
func NewDisplayEditor(nodeList *jsontree.NodeList) *DisplayEditor {
	return &DisplayEditor{nodeList}
}

// Edit defines response to input for display view
func (e *DisplayEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyArrowUp:
		e.nodeList.MoveJSONPosition(-1)
	case key == gocui.KeyArrowDown:
		e.nodeList.MoveJSONPosition(1)
	case key == gocui.KeyArrowLeft:
		x, _ := v.Origin()
		v.SetOrigin(x-1, 0)
	case key == gocui.KeyArrowRight:
		x, _ := v.Origin()
		v.SetOrigin(x+1, 0)
	case key == gocui.KeyPgup:
		e.nodeList.MoveJSONPosition(-25)
	case key == gocui.KeyPgdn:
		e.nodeList.MoveJSONPosition(25)
	}
}

// SelectPopupEditor stuff
type SelectPopupEditor struct {
	ch chan string
}

// NewSelectPopupEditor stuff
func NewSelectPopupEditor(ch chan string) *SelectPopupEditor {
	return &SelectPopupEditor{ch}
}

// Edit stuff
func (s *SelectPopupEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyArrowUp:
		_, cursorY := v.Cursor()
		_, originY := v.Origin()
		if originY != 0 && cursorY != 0 {
			v.MoveCursor(0, -1, false)
		}
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyEnter:
		_, cursorY := v.Cursor()
		if line, err := v.Line(cursorY); err == nil {
			s.ch <- line
		}
		// case key == gocui.KeyCtrlG:
		// cui.ClosePopup()
	}
}

// WritePopupEditor provides an editor that user can write to and run function on enter
type WritePopupEditor struct {
	ch chan string
}

// NewWritePopupEditor creates a new WritePopupEditor
func NewWritePopupEditor(ch chan string) *WritePopupEditor {
	return &WritePopupEditor{ch}
}

// Edit will allow user to enter information into popup and run function
func (w *WritePopupEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	v.Highlight = false
	if x, y := v.Cursor(); y == 0 {
		v.SetCursor(x, 1)
	}
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		if !atLineBeginning(v) {
			v.EditDelete(true)
		}
	case key == gocui.KeyArrowLeft:
		if !atLineBeginning(v) {
			v.MoveCursor(-1, 0, false)
		}
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyEnter:
		input, _ := v.Line(1)
		w.ch <- input
		return
		// case key == gocui.KeyCtrlG:
		// 	cui.ClosePopup()
	}
}

func atLineBeginning(v *gocui.View) bool {
	xCursor, _ := v.Cursor()
	xOrigin, _ := v.Origin()
	return xCursor == 0 && xOrigin == 0
}

// ConfirmPopupEditor allows user to confirm or reject something
type ConfirmPopupEditor struct {
	ch chan string
}

// NewConfirmPopupEditor stuff
func NewConfirmPopupEditor(ch chan string) *ConfirmPopupEditor {
	return &ConfirmPopupEditor{ch}
}

// Edit only allow user to press y or n
func (c *ConfirmPopupEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch == 'y' && c.ch != nil:
		cui.ClosePopup()
		c.ch <- "y"
	case ch == 'n' && c.ch != nil:
		cui.ClosePopup()
		c.ch <- "n"
	case key == gocui.KeyEnter && c.ch == nil:
		cui.ClosePopup()
	}
}
