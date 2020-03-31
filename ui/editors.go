package ui

import (
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
	// 	x, y := v.Origin()
	// 	v.SetOrigin(x, y-1)
	case key == gocui.KeyArrowDown:
		e.nodeList.MoveJSONPosition(1)
		// 	x, y := v.Origin()
		// 	v.SetOrigin(x, y+1)
		// case key == gocui.KeyArrowLeft:
		// 	x, y := v.Origin()
		// 	v.SetOrigin(x-1, y)
		// case key == gocui.KeyArrowRight:
		// 	x, y := v.Origin()
		// 	v.SetOrigin(x+1, y)
		// case key == gocui.KeyPgup:
		// 	x, y := v.Origin()
		// 	if y-25 > 0 {
		// 		v.SetOrigin(x, y-25)
		// 	} else {
		// 		v.SetOrigin(x, 0)
		// 	}
		// case key == gocui.KeyPgdn:
		// 	x, y := v.Origin()
		// 	v.SetOrigin(x, y+25)
	}
}

func saveEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// switch {
	// case ch != 0 && mod == 0:
	// 	v.EditWrite(ch)
	// case key == gocui.KeySpace:
	// 	v.EditWrite(' ')
	// case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
	// 	v.EditDelete(true)
	// case key == gocui.KeyArrowLeft:
	// 	v.MoveCursor(-1, 0, false)
	// case key == gocui.KeyArrowRight:
	// 	v.MoveCursor(1, 0, false)
	// case key == gocui.KeyEnter:
	// 	filename, _ := v.Line(0)
	// 	err := save(filename)
	// 	v.Clear()
	// 	if err != nil {
	// 		fmt.Fprintf(v, "Failed to write to \"%s\"", filename)
	// 	} else {
	// 		fmt.Fprintf(v, "Sucessfully saved to \"%s\"", filename)
	// 	}
	// 	go closeSave(v)
	// }
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
