package ui

import (
	"fmt"
	"kube-review/jsontree"
	"os"
	"time"

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
		e.nodeList.ExpandActiveNode()
	case key == gocui.KeyArrowLeft:
		newY := e.nodeList.CollapseActiveNode()
		v.SetCursor(0, newY)
	// case ch == 'e':
	// 	e.nodeList.ExpandAllNodes()
	// case ch == 'c':
	// 	e.nodeList.CollapseAllNodes()
	// 	v.SetCursor(0, 0)
	case key == gocui.KeyPgup:
		e.nodeList.MoveTopNode(-25)
	case key == gocui.KeyPgdn:
		e.nodeList.MoveTopNode(25)
	}

	// This prevents scrolling as we will be handling that
	// v.SetOrigin(0, 0)

	GetWindow().UpdateViewContent(PANEL, e.nodeList.GetNodes(GetWindow().Height(PANEL)))
	GetWindow().UpdateViewContent(DISPLAY, e.nodeList.GetJSON(GetWindow().Height(DISPLAY)))
}

func getCurrentIndex(v *gocui.View) int {
	_, yCursor := v.Cursor()
	_, yOrigin := v.Origin()
	return yCursor + yOrigin
}

// SearchEditor is the editor for the search view
type SearchEditor struct {
	nodeList *jsontree.NodeList
}

// NewSearchEditor stuff
func NewSearchEditor(nodeList *jsontree.NodeList) *SearchEditor {
	return &SearchEditor{nodeList}
}

// TODO: increase functionality of the search editing

// Edit stuff
func (e *SearchEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
	// line, _ := v.Line(0)
	// e.nodeList.Search(line)
	GetWindow().UpdateViewContent(PANEL, e.nodeList.GetNodes(GetWindow().Height(PANEL)))
	// GetWindow().UpdateViewContent(DISPLAY, e.nodeList.GetJSON(0))
}

// Edit defines response to input for display view
func displayEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyArrowUp:
		x, y := v.Origin()
		v.SetOrigin(x, y-1)
	case key == gocui.KeyArrowDown:
		x, y := v.Origin()
		v.SetOrigin(x, y+1)
	case key == gocui.KeyArrowLeft:
		x, y := v.Origin()
		v.SetOrigin(x-1, y)
	case key == gocui.KeyArrowRight:
		x, y := v.Origin()
		v.SetOrigin(x+1, y)
	case key == gocui.KeyPgup:
		x, y := v.Origin()
		if y-25 > 0 {
			v.SetOrigin(x, y-25)
		} else {
			v.SetOrigin(x, 0)
		}
	case key == gocui.KeyPgdn:
		x, y := v.Origin()
		v.SetOrigin(x, y+25)
	}
}

func saveEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyEnter:
		filename, _ := v.Line(0)
		err := save(filename)
		v.Clear()
		if err != nil {
			fmt.Fprintf(v, "Failed to write to \"%s\"", filename)
		} else {
			fmt.Fprintf(v, "Sucessfully saved to \"%s\"", filename)
		}
		go closeSave(v)
	}
}

func save(filename string) error {
	file, openError := os.Create(filename)
	defer file.Close()
	if openError != nil {
		return openError
	}
	_, writeError := file.Write([]byte(GetWindow().GetContent(DISPLAY)))
	if writeError != nil {
		return writeError
	}
	return nil
}

func closeSave(v *gocui.View) {
	time.Sleep(1 * time.Second)
	GetWindow().ShowSaveView(false)
	v.Clear()
	v.SetCursor(0, 0)
}
