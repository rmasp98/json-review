package ui

import (
	"github.com/rmasp98/kube-review/nodelist"
	"github.com/rmasp98/kube-review/search"

	"github.com/gdamore/tcell"

	"github.com/rmasp98/kube-review/tview"
)

// Run stuff
func Run(nodeList *nodelist.NodeList, queryList *search.QueryList) error {

	nodes := tview.NewTextView()
	nodes.SetText(nodeList.GetNodes(50)).SetWrap(false)
	nodes.SetBorder(true).SetTitle("Nodes")

	// TODO: figure out window size
	json := tview.NewTextView().SetText(nodeList.GetJSON(-1)).SetWrap(false)
	json.SetBorder(true).SetTitle("Json")

	help := tview.NewTextView().SetText("INSERT HELP TEXT")
	help.SetBorder(true).SetTitle("Help")

	searchThing := search.NewSearch(search.EXPRESSION, queryList)

	// Can use set label to define search type
	search := tview.NewInputField().SetFieldBackgroundColor(tcell.ColorDefault).SetLabel("Regex-Filter: ").SetAutocompleteFunc(searchThing.GetHints2)
	search.SetBorder(true).SetTitle("Search")
	search.SetAutocompleteInsertFunc(func(autoCompleteText string, currentText string, cursorPos int, key tcell.Key) string {
		if key == tcell.KeyEnter {
			return autoCompleteText
		}
		return currentText
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(nodes, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(search, 3, 1, true).
				AddItem(json, 0, 1, false),
				0, 3, true),
			0, 1, true).
		AddItem(help, 3, 1, false)

	return tview.NewApplication().SetRoot(flex, true).Run()
}
