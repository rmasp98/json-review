package uiold

import (
	"fmt"
	"log"
	"strings"

	"github.com/rmasp98/kube-review/nodelist"
	"github.com/rmasp98/kube-review/search"

	"github.com/awesome-gocui/gocui"
)

// CursesUI is the central UI control point
type CursesUI struct {
	gui       *gocui.Gui
	win       Window
	nodeList  *nodelist.NodeList
	queryList *search.QueryList
}

// NewCursesUI stuff
func NewCursesUI(nodeList *nodelist.NodeList, queryList *search.QueryList) (CursesUI, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return CursesUI{}, err
	}
	gui.Highlight = true
	gui.SelFgColor = gocui.ColorRed
	gui.Cursor = true

	cui := CursesUI{gui, NewWindow(0.2, 1, 3), nodeList, queryList}

	cui.gui.SetManagerFunc(cui.update)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, changeView); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, NewSaveUI(nodeList, queryList).Save); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, cui.splitNodeList); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlY, gocui.ModNone, cui.selectView); err != nil {
		log.Panicln(err)
	}
	gui.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		cui.nodeList.ResetView()
		return nil
	})

	return cui, nil
}

// Run stuff
func (cui CursesUI) Run() error {
	if err := cui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	cui.gui.Close()
	return nil
}

// TriggerUpdate allows non-main threads to trigger an update to the ui
func (cui CursesUI) TriggerUpdate() {
	cui.gui.Update(cui.update)
}

// UpdateViewContent stuff
func (cui *CursesUI) UpdateViewContent(view ViewEnum, content string) {
	if v, err := cui.gui.View(view.String()); err == nil {
		v.Clear()
		fmt.Fprint(v, content)
		return
	}
	log.Printf("Could not update content in %s", view.String())
}

// UpdateViewEditor stuff
func (cui *CursesUI) UpdateViewEditor(view ViewEnum, editor gocui.Editor) {
	if v, err := cui.gui.View(view.String()); err == nil {
		v.Editable = editor != nil
		v.Editor = editor
		return
	}
	log.Printf("Could not update editor in %s", view.String())
}

// UpdateViewTitle stuff
func (cui *CursesUI) UpdateViewTitle(view ViewEnum, title string) {
	if v, err := cui.gui.View(view.String()); err == nil {
		v.Title = title
		return
	}
	log.Printf("Could not update title in %s", view.String())
}

// UpdateHelp adds addedHelp to the base help in the help view
func (cui *CursesUI) UpdateHelp(addedHelp string) {
	if v, err := cui.gui.View(HELP.String()); err == nil {
		v.Clear()
		v.Write([]byte(HELP.Help() + addedHelp))
	}
	log.Printf("Could not update help")
}

// SetCursor allows the cursor to be turned on and off
func (cui *CursesUI) SetCursor(active bool) {
	cui.gui.Cursor = active
}

func (cui CursesUI) update(gui *gocui.Gui) error {
	x, y := cui.gui.Size()
	cui.win.Resize(x, y, cui.getLinesInSearch())
	return cui.setViews()
}

func (cui CursesUI) setViews() error {
	for name, layout := range cui.win.views {
		if view, err := cui.gui.SetView(name.String(), layout.x0, layout.y0, layout.x1, layout.y1, 0); err != nil {
			view.Title = name.String()
			switch name {
			case PANEL:
				view.Highlight = true
				view.Editor = NewNodesEditor(cui.nodeList)
				view.Editable = true
			case DISPLAY:
				view.Editor = NewDisplayEditor(cui.nodeList)
				view.Editable = true
			case SEARCH:
				view.Title = "Search: Mode=Regex-Find"
				view.Editor = NewSearchEditor(cui.nodeList, cui.queryList)
				view.Editable = true
				view.Autoscroll = true
			case HELP:
				view.Write([]byte(HELP.Help()))
			}
		} else {
			switch name {
			case PANEL:
				view.Clear()
				view.Write([]byte(cui.nodeList.GetNodes(layout.y1 - layout.y0)))
			case DISPLAY:
				view.Clear()
				view.Write([]byte(cui.nodeList.GetJSON(layout.y1 - layout.y0)))
			case VIEW:
				view.Clear()
				view.Write([]byte(cui.nodeList.GetCurrentView()))
			}
		}
	}
	return nil
}

func (cui CursesUI) splitNodeList(g *gocui.Gui, v *gocui.View) error {
	var ch = make(chan string)
	cui.CreatePopup("Split Nodes", "Define the string used to split the nodes:\n", NewWritePopupEditor(ch), true, false, true)
	go func(ch chan string, nodeList *nodelist.NodeList) {
		splitString := <-ch
		nodeList.SplitViews(splitString)
		cui.ClosePopup()
	}(ch, cui.nodeList)
	return nil
}

func (cui CursesUI) selectView(g *gocui.Gui, v *gocui.View) error {
	var ch = make(chan string)
	content := "Choose the nodelist view:"
	for _, view := range cui.nodeList.ListViews() {
		content += "\n" + view
	}
	cui.CreatePopup("Select View", content, NewSelectPopupEditor(ch), false, true, true)
	go func(ch chan string, nodeList *nodelist.NodeList) {
		view := <-ch
		nodeList.SetView(view)
		cui.ClosePopup()
	}(ch, cui.nodeList)
	return nil
}

var lastView string

// CreatePopup stuff
func (cui *CursesUI) CreatePopup(title, content string, editor gocui.Editor, cursor bool, highlight bool, nextLine bool) error {
	if view := cui.gui.CurrentView(); view != nil {
		lastView = view.Name()
	} else {
		lastView = SEARCH.String()
	}
	winWidth, winHeight := cui.gui.Size()
	x0, y0, x1, y1 := determinePopupDimensions(content, winWidth, winHeight)
	if view, err := cui.gui.SetView("Popup", x0, y0, x1, y1, 0); err != nil {
		view.Title = title
		view.Editable = editor != nil
		view.Editor = editor
		view.Highlight = highlight
		view.Write([]byte(content))
		if nextLine {
			view.SetCursor(0, 1)
		}
		cui.SetCursor(cursor)
		if _, errCurrent := cui.gui.SetCurrentView("Popup"); errCurrent != nil {
			log.Println(errCurrent.Error())
		}
		if _, errTop := cui.gui.SetViewOnTop("Popup"); errTop != nil {
			log.Println(errTop.Error())
		}
		cui.TriggerUpdate()
	} else {
		return fmt.Errorf("Popup already exists")
	}
	return nil
}

// ClosePopup stuff
func (cui *CursesUI) ClosePopup() {
	cui.gui.DeleteView("Popup")
	if _, err := cui.gui.SetCurrentView(lastView); err != nil {
		cui.gui.SetCurrentView(SEARCH.String())
	}
	cui.TriggerUpdate()
}

func (cui CursesUI) getLinesInSearch() int {
	if v, err := cui.gui.View(SEARCH.String()); err == nil {
		if lines := len(v.BufferLines()); lines != 0 {
			return 2 + lines
		}
	}
	return 3
}

var screenID = 0

func changeView(g *gocui.Gui, v *gocui.View) error {
	screenID = (screenID + 1) % 3
	screen := ViewEnum(screenID)
	if screen == SEARCH {
		cui.SetCursor(true)
	} else {
		cui.SetCursor(false)
	}
	if _, err := g.SetCurrentView(screen.String()); err != nil {
		return err
	}
	if _, err := g.SetViewOnTop(screen.String()); err != nil {
		return err
	}
	cui.UpdateHelp(screen.Help())
	return nil
}

func determinePopupDimensions(content string, winWidth, winHeight int) (int, int, int, int) {
	splitText := strings.Split(content, "\n")
	height := len(splitText) + 1
	width := 40
	for _, t := range splitText {
		if len(t) > width {
			width = len(t)
		}
	}
	if width != 0 && height != 0 && width < winWidth && height < winHeight {
		return winWidth/2 - width/2,
			winHeight/2 - height/2,
			winWidth/2 + width/2 + (width % 2) + 1,
			winHeight/2 + height/2 + (height % 2)
	}
	return 0, 0, 0, 0
}
