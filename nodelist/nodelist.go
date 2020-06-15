package nodelist

import (
	"fmt"
	"regexp"
	"strings"
)

// NodeList presents an interface for interacting with nodelists
type NodeList struct {
	master          MasterNodeList
	views           map[string]View
	currentView     View
	currentViewName string
	topNodeIndex    int
	activeNodeIndex int
	jsonViewOffset  int
}

// NewNodeList stuff
// parser to create nodelist
func NewNodeList(jsonData []byte, blocking bool) (NodeList, error) {
	master, err := NewMasterNodeList(jsonData, blocking)
	if err != nil {
		return NodeList{}, err
	}

	if blocking {
		view, err := master.GetNodeView()
		if err != nil {
			return NodeList{}, err
		}
		return NodeList{master, map[string]View{"main": view}, view, "main", 0, 0, 0}, nil
	}
	nodeList := NodeList{master, map[string]View{}, View{}, "", 0, 0, 0}
	//subscribe to master callback
	return nodeList, nil
}

// GetJSON returns a formatted json string, num lines long for
// the active node. Fields can be hidden using the Filter function
// Passing -1 will remove limit on returned lines
func (n NodeList) GetJSON(num int) string {
	return n.currentView.GetJSON(n.activeNodeIndex, n.jsonViewOffset, num)
}

// GetNodes returns a formated string list of visible nodes from topNode
// and is only num nodes long
func (n NodeList) GetNodes(num int) string {
	return n.currentView.GetNodes(n.topNodeIndex, num)
}

// MoveTopNode changes the start position (topNode) of what GetNodes returns
// relative to its current position
func (n *NodeList) MoveTopNode(offset int) {
	n.topNodeIndex += offset
	if n.topNodeIndex < 0 {
		n.topNodeIndex = 0
	} else if n.topNodeIndex >= n.currentView.Size() {
		n.topNodeIndex = n.currentView.Size() - 1
	}
}

// SetActiveNode lets NodeList know the highlighted node in editor.
// This is the node that all actions will be performed on.
// Actual nodeIndex is calculated relative to topNode
func (n *NodeList) SetActiveNode(index int) {
	n.activeNodeIndex = n.topNodeIndex + index
	if n.activeNodeIndex >= n.currentView.Size() {
		n.activeNodeIndex = n.currentView.Size() - 1
	}
	n.jsonViewOffset = 0
}

// MoveJSONView offsets the Json view returned
func (n *NodeList) MoveJSONView(offset int) {
	n.jsonViewOffset += offset
	if n.jsonViewOffset < 0 {
		n.jsonViewOffset = 0
	}
}

// GetNodesMatching searches entire view for matches of matchtype to regex. Set equal to false to invert result
func (n NodeList) GetNodesMatching(regex *regexp.Regexp, matchType MatchType, equal bool) []int {
	return n.currentView.GetNodesMatching(regex, matchType, equal)
}

// GetRelativesMatching searches nodes relative to nodeIndex in similar fashion to GetNodesMatching.
// relativeStartLevel defines how many levels above nodeIndex the search should start from
// and depth defines how many levels of children from relativeStartLevel should be searched.
// To search a particular parent, set depth to zero, otherwise that parent is ignored
func (n NodeList) GetRelativesMatching(nodeIndex, relativeStartLevel, depth int, regex *regexp.Regexp, matchType MatchType, equal bool) []int {
	return n.currentView.GetRelativesMatching(nodeIndex, relativeStartLevel, depth, regex, matchType, equal)
}

// Filter stuff
func (n *NodeList) Filter(nodeIndices []int) error {
	newView, err := n.currentView.Filter(nodeIndices)
	if err != nil {
		return err
	}
	n.currentView = newView
	return nil
}

// Highlight stuff
func (n *NodeList) Highlight(nodeIndices []int) {
	n.currentView.Highlight(nodeIndices)
}

// FindNextHighlight stuff
func (n *NodeList) FindNextHighlight() error {
	newOffset, err := n.currentView.FindNextHighlight(n.activeNodeIndex, n.jsonViewOffset)
	if err != nil {
		return err
	}
	n.jsonViewOffset = newOffset
	return nil
}

// SplitViews stuff
// Split seperates MasterNodeList into different NodeListViews based on seperator
// e.g. "items = kind" will find array items and split based on the value of kind in each element
// If not in items array or does not have kind, will be put into "main" group
// Can choose to wait for loading and split to complete with blocking
func (n *NodeList) SplitViews(separator string) error {
	root, target := parseSeparator(separator)
	split, err := n.views["main"].Split(root, target)
	if err != nil {
		return err
	}
	for name, nodes := range split {
		n.views[name] = nodes
	}

	return nil
}

// ListViews stuff
func (n NodeList) ListViews() []string {
	keys := make([]string, 0, len(n.views))
	for k := range n.views {
		keys = append(keys, k)
	}
	return keys
}

// SetView stuff
func (n *NodeList) SetView(name string) error {
	if view, ok := n.views[name]; ok {
		n.currentView = view
		n.currentViewName = name
		return nil
	}
	return fmt.Errorf("View with name, '%s', does not exist", name)
}

// ResetView stuff
func (n *NodeList) ResetView() {
	n.currentView = n.views[n.currentViewName]
}

func parseSeparator(sep string) ([]string, []string) {
	split := strings.Split(sep, "=")
	if len(split) == 0 {
		return []string{}, []string{}
	} else if len(split) == 1 {
		return []string{}, splitSeparator(split[0])
	}
	return splitSeparator(split[0]), splitSeparator(split[1])
}

func splitSeparator(sep string) []string {
	trimmedSep := strings.Trim(sep, " ")
	return strings.Split(trimmedSep, ".")
}
