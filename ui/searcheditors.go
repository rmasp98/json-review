package ui

import (
	"kube-review/jsontree"
	"kube-review/search"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

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
	basicEditor(v, key, ch, mod)
	if key == gocui.KeyEnter {
		line, _ := v.Line(0)
		e.nodeList.Search(line)
	} else if key == gocui.KeyCtrlQ {
		GetWindow().UpdateEditor(SEARCH, NewQueryEditor(e.nodeList))
	}
	search, _ := v.Line(0)
	GetWindow().UpdateViewContent(SEARCH, search)
}

// QueryEditor stuff
type QueryEditor struct {
	nodeList  *jsontree.NodeList
	queryList search.QueryList
}

// NewQueryEditor stuff
func NewQueryEditor(nodeList *jsontree.NodeList) *QueryEditor {
	ql := search.NewQueryList()
	ql.Load("querylist.json")
	return &QueryEditor{nodeList, ql}
}

const (
	red       = "\x1b[0;31m"
	white     = "\x1b[0;97m"
	bold      = "\033[1m"
	unbold    = "\033[0m"
	delimiter = " - "
)

// Edit stuff
func (e *QueryEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {

	if _, lineNum := v.Cursor(); lineNum == 0 {
		basicEditor(v, key, ch, mod)
	} else {
		switch {
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyEnter:
			line, _ := v.Line(lineNum)
			v.Clear()
			search := strings.Split(line, delimiter)
			v.Write([]byte(search[0]))
			v.SetCursor(len(search[0]), 0)
		}

	}

	if key == gocui.KeyEnter {
		line, _ := v.Line(0)
		search := e.queryList.GetRegex(line)
		e.nodeList.Search(search)
	} else if key == gocui.KeyCtrlQ {
		search, _ := v.Line(0)
		v.Highlight = false
		v.SetCursor(len(search), 0)
		GetWindow().UpdateViewContent(SEARCH, search)
		GetWindow().UpdateEditor(SEARCH, NewSearchEditor(e.nodeList))
		return
	}

	search, _ := v.Line(0)
	queryDetails := e.getPossibleQueries(search)
	GetWindow().UpdateViewContent(SEARCH, search+queryDetails)

	if _, lineNum := v.Cursor(); lineNum == 0 {
		v.Highlight = false
	} else {
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
}

func (e QueryEditor) getPossibleQueries(search string) string {
	var queryDetails string
	r, err := regexp.Compile(search)
	if err == nil && search != "" {
		for _, name := range e.queryList.GetNames() {
			if name == search {
				return ""
			} else if r.MatchString(name) {
				queryDetails += "\n"
				queryDetails += white + bold + name
				queryDetails += delimiter
				queryDetails += red + e.queryList.GetDescription(name)
				queryDetails += unbold + white
			}
		}
		if queryDetails == "" {
			return "\nNo Matched Queries"
		}
	}
	return queryDetails
}

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
	}
}
