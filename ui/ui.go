package ui

import (
	"kube-review/nodelist"
	"kube-review/search"
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
	// VIEW a
	VIEW
)

func (ve ViewEnum) String() string {
	return [...]string{"Panel", "Search", "Display", "Help", "View"}[ve]
}

// Help stuff
func (ve ViewEnum) Help() string {
	return [...]string{
		" | E: Expand Node | C: Collapse Node",             //PANEL
		" | Ctrl+Q: Toggle Query Mode | Ctrl+N: Find Next", //SEARCH
		"", //DISPLAY
		"Ctrl+C: Exit  | Tab: Next View | Ctrl+R: Reset View | Ctrl+T: Split View | Ctrl+Y: Change View | Ctrl+S: Save ", //HELP
		"", //VIEW
	}[ve]
}

// Run is the entry point for the curses UI interface
func Run(nodeList *nodelist.NodeList, queryList *search.QueryList) error {
	var err error
	cui, err = NewCursesUI(nodeList, queryList)
	if err != nil {
		return err
	}
	return cui.Run()
}
