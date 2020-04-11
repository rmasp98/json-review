package mocks

import (
	"kube-review/jsontree"
)

// NodeListMock creates a mock for NodeList
type NodeListMock struct {
	Calls []string
	Args  [][]interface{}
}

// GetNodesMatching is a mock function
func (n *NodeListMock) GetNodesMatching(regex string, matchType jsontree.MatchType, equal bool) []int {
	n.Calls = append(n.Calls, "GetNodesMatching")
	args := make([]interface{}, 3)
	args[0] = regex
	args[1] = matchType
	args[2] = equal
	n.Args = append(n.Args, args)
	return []int{1, 2, 5}
}

// GetChildrenMatching is a mock function
func (n *NodeListMock) GetChildrenMatching(nodeIndex int, regex string, matchType jsontree.MatchType, invert bool, allChildren bool) []int {
	n.Calls = append(n.Calls, "GetChildrenMatching")
	args := make([]interface{}, 5)
	args[0] = nodeIndex
	args[1] = regex
	args[2] = matchType
	args[3] = invert
	args[4] = allChildren
	n.Args = append(n.Args, args)
	return []int{1, 6}
}

// GetParentChildrenMatching is a mock function
func (n *NodeListMock) GetParentChildrenMatching(nodeIndex int, regex string, matchType jsontree.MatchType, invert bool, allParents bool) []int {
	n.Calls = append(n.Calls, "GetParentChildrenMatching")
	args := make([]interface{}, 5)
	args[0] = nodeIndex
	args[1] = regex
	args[2] = matchType
	args[3] = invert
	args[4] = allParents
	n.Args = append(n.Args, args)
	return []int{2, 7}
}

// ApplyFilter is a mock function
func (n *NodeListMock) ApplyFilter(nodes []int) error {
	n.Calls = append(n.Calls, "ApplyFilter")
	args := make([]interface{}, 1)
	args[0] = nodes
	n.Args = append(n.Args, args)
	return nil
}

// ApplyHighlight is a mock function
func (n *NodeListMock) ApplyHighlight(nodes []int) error {
	n.Calls = append(n.Calls, "ApplyHighlight")
	args := make([]interface{}, 1)
	args[0] = nodes
	n.Args = append(n.Args, args)
	return nil
}

// FindNextHighlightedNode is a mock function
func (n *NodeListMock) FindNextHighlightedNode() error {
	n.Calls = append(n.Calls, "FindNextHighlightedNode")
	args := make([]interface{}, 0)
	n.Args = append(n.Args, args)
	return nil
}
