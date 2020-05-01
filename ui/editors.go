package ui

import (
	"kube-review/jsontree"

	"github.com/jroimartin/gocui"
)

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

func getCurrentIndex(v *gocui.View) int {
	_, yCursor := v.Cursor()
	_, yOrigin := v.Origin()
	return yCursor + yOrigin
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

// PopupWriteEditor provides an editor that user can write to and run function on enter
type PopupWriteEditor struct {
	dialog   string
	function func(string)
}

// NewPopupWriteEditor creates a new PopupWriteEditor
func NewPopupWriteEditor(dialog string, function func(string)) *PopupWriteEditor {
	return &PopupWriteEditor{dialog, function}
}

// Edit will allow user to enter information into popup and run function
func (p *PopupWriteEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		if x, _ := v.Cursor(); x > 0 {
			v.EditDelete(true)
		}
	case key == gocui.KeyArrowLeft:
		if x, _ := v.Cursor(); x > 0 {
			v.MoveCursor(-1, 0, false)
		}
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyEnter:
		input, _ := v.Line(1)
		p.function(input)
		return
	}
	input, _ := v.Line(1)
	GetWindow().UpdateViewContent(POPUP, p.dialog+input)
}

// PopupConfirmEditor allows user to confirm or reject something
type PopupConfirmEditor struct {
	function func()
}

// NewPopupConfirmEditor stuff
func NewPopupConfirmEditor(function func()) *PopupConfirmEditor {
	return &PopupConfirmEditor{function}
}

// Edit only allow user to press y or n
func (p *PopupConfirmEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch == 'y':
		p.function()
	case ch == 'n':
		GetWindow().ShowPopupView(false, nil)
	}
}
