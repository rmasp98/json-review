package nodelist

// NodeList presents an interface for interacting with nodelists
type NodeList struct {
	master      MasterNodeList
	views       map[string]View
	currentView View
}

// NewNodeList stuff
// parser to create nodelist
func NewNodeList() NodeList {
	return NodeList{}
}

//functions
// SplitViews
// moveTopNode
// setactivenode
// movejsonposition
// GetJSON
// GetNodes
// all search functions
// apply filter
// ApplyHighlight
// FindNexthighlight
