package mocks

import (
	"regexp"

	"github.com/rmasp98/kube-review/nodelist"
)

// NodeListMock creates a mock for NodeList
type NodeListMock struct {
	Calls   []string
	Args    [][]interface{}
	Returns [][]int
}

// GetNodesMatching is a mock function
func (n *NodeListMock) GetNodesMatching(regex *regexp.Regexp, matchType nodelist.MatchType, equal bool) []int {
	n.Calls = append(n.Calls, "GetNodesMatching")
	args := make([]interface{}, 3)
	args[0] = regex
	args[1] = matchType
	args[2] = equal
	n.Args = append(n.Args, args)
	if len(n.Calls)-1 < len(n.Returns) {
		return n.Returns[len(n.Calls)-1]
	}
	return []int{}
}

// GetRelativesMatching is a mock function
func (n *NodeListMock) GetRelativesMatching(nodeIndex, relativeStartLevel, depth int, regex *regexp.Regexp, matchType nodelist.MatchType, equal bool) []int {
	n.Calls = append(n.Calls, "GetRelativesMatching")
	args := make([]interface{}, 6)
	args[0] = nodeIndex
	args[1] = relativeStartLevel
	args[2] = depth
	args[3] = regex
	args[4] = matchType
	args[5] = equal
	n.Args = append(n.Args, args)
	if len(n.Calls)-1 < len(n.Returns) {
		return n.Returns[len(n.Calls)-1]
	}
	return []int{}
}

// Filter is a mock function
func (n *NodeListMock) Filter(nodes []int) error {
	n.Calls = append(n.Calls, "Filter")
	args := make([]interface{}, 1)
	args[0] = nodes
	n.Args = append(n.Args, args)
	return nil
}

// Highlight is a mock function
func (n *NodeListMock) Highlight(nodes []int) {
	n.Calls = append(n.Calls, "Highlight")
	args := make([]interface{}, 1)
	args[0] = nodes
	n.Args = append(n.Args, args)
}

// FindNextHighlight is a mock function
func (n *NodeListMock) FindNextHighlight() error {
	n.Calls = append(n.Calls, "FindNextHighlight")
	args := make([]interface{}, 0)
	n.Args = append(n.Args, args)
	return nil
}

// ResetView is a mock function
func (n *NodeListMock) ResetView() {
	n.Calls = append(n.Calls, "ResetView")
}
