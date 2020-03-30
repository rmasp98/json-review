package ui

import (
	"kube-review/jsontree"
	"log"

	"github.com/jroimartin/gocui"
)

// GoCui is an interface to the gocui.Gui struct
type GoCui interface {
	MainLoop() error
	Close()
	SetView(string, int, int, int, int) (*gocui.View, error)
	View(string) (*gocui.View, error)
	SetCurrentView(string) (*gocui.View, error)
	SetViewOnTop(string) (*gocui.View, error)
	SetViewOnBottom(string) (*gocui.View, error)
	DeleteView(string) error
	Size() (int, int)
}

////////////////////////////////////////////////////////////////////////

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
	// SAVE a
	SAVE
)

func (ve ViewEnum) String() string {
	return [...]string{"Panel", "Search", "Display", "Help", "Filename"}[ve]
}

////////////////////////////////////////////////////////////////////////

// CursesUI stuff
type CursesUI struct {
	gui GoCui
}

// NewCursesUI stuff
func NewCursesUI(json *jsontree.NodeList) CursesUI {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	gui.Highlight = true
	gui.SelFgColor = gocui.ColorRed
	gui.Cursor = true

	cui := CursesUI{gui}

	gui.SetManagerFunc(func(gui *gocui.Gui) error {
		x, y := gui.Size()
		GetWindow().Resize(x, y)
		GetWindow().UpdateViewContent(DISPLAY, json.GetJSON(y))
		GetWindow().UpdateViewContent(PANEL, json.GetNodes(y))
		return GetWindow().SetViews(gui)
	})

	y, _ := gui.Size()
	GetWindow().UpdateViewContent(DISPLAY, json.GetJSON(y))
	GetWindow().UpdateViewContent(PANEL, json.GetNodes(y))
	GetWindow().UpdateEditor(PANEL, NewNodesEditor(json))
	GetWindow().UpdateEditor(SEARCH, NewSearchEditor(json))
	GetWindow().UpdateEditor(DISPLAY, NewDisplayEditor(json))
	GetWindow().UpdateEditor(SAVE, gocui.EditorFunc(saveEditor))

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, changeView); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		GetWindow().ShowSaveView(true)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	return cui
}

// Run stuff
func (cui CursesUI) Run() {
	if err := cui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	cui.gui.Close()
}

// Quit stuff
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

var screen = 0

func changeView(g *gocui.Gui, v *gocui.View) error {
	screen = (screen + 1) % 3
	if ViewEnum(screen) == PANEL || ViewEnum(screen) == DISPLAY {
		g.Cursor = false
	} else {
		g.Cursor = true
	}
	g.SetCurrentView(ViewEnum(screen).String())

	return nil
}
