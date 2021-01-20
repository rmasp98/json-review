package ui

import (
	"kube-review/nodelist"
	"kube-review/search"

	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

// Run stuff
func Run(nodeList *nodelist.NodeList, queryList *search.QueryList) error {
	nodes := nodeList.GetTViewNodes()
	nodes.SetBorder(true).SetTitle("Nodes")

	// TODO: figure out window size
	json := tview.NewTextView().SetText(nodeList.GetJSON(-1)).SetWrap(false)
	json.SetBorder(true).SetTitle("Json")

	help := tview.NewTextView().SetText("INSERT HELP TEXT")
	help.SetBorder(true).SetTitle("Help")

	// Can use set label to define search type
	search := tview.NewInputField().SetFieldBackgroundColor(tcell.ColorDefault).SetLabel("Regex-Filter: ").SetAutocompleteFunc(func(text string) []string { return []string{"Test1", "Test2"} })
	search.SetBorder(true).SetTitle("Search")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(nodes, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(search, 4, 1, true).
				AddItem(json, 0, 1, false),
				0, 3, true),
			0, 1, true).
		AddItem(help, 4, 1, false)

	return tview.NewApplication().SetRoot(flex, true).Run()
}
