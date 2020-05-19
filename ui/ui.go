package ui

import (
	"kube-review/jsontree"
)

var cui CursesUI

// ViewEnum list the possible views available
type ViewEnum int

const (
	// PANEL a
	PANEL ViewEnum = iota
	// SEARCH a
	SEARCH
	// DISPLAY a
	DISPLAY
	// HELP a
	HELP
)

func (ve ViewEnum) String() string {
	return [...]string{"Panel", "Search", "Display", "Help"}[ve]
}

// Help stuff
func (ve ViewEnum) Help() string {
	return [...]string{
		" | E: Expand Node | C: Collapse Node",             //PANEL
		" | Ctrl+Q: Toggle Query Mode | Ctrl+N: Find Next", //SEARCH
		"", //DISPLAY
		"", //HELP
	}[ve]
}

// Run is the entry point for the curses UI interface
func Run(nodeList *jsontree.NodeList) error {
	var err error
	cui, err = NewCursesUI(nodeList)
	if err != nil {
		return err
	}
	return cui.Run()
}
