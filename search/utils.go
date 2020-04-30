package search

import "kube-review/jsontree"

// FunctionEnum list the possible views available
type FunctionEnum int

const (
	// FIND a
	FIND FunctionEnum = iota
	// FILTER a
	FILTER
)

func (fe FunctionEnum) String() string {
	return [...]string{"Find", "Filter"}[fe]
}

// QueryEnum list the possible views available
type QueryEnum int

const (
	// REGEX a
	REGEX QueryEnum = iota
	// QUERY a
	QUERY
	// INTELLIGENT a
	INTELLIGENT
)

func (qe QueryEnum) String() string {
	return [...]string{"Regex", "Query", "Intelligent"}[qe]
}

type sNodeList interface {
	GetNodesMatching(regex string, matchType jsontree.MatchType, equal bool) []int
	GetChildrenMatching(nodeIndex int, regex string, matchType jsontree.MatchType, equal bool, recursive bool) []int
	GetParentChildrenMatching(nodeIndex int, regex string, matchType jsontree.MatchType, equal bool, recursive bool) []int
	ApplyFilter(nodes []int) error
	ApplyHighlight(nodes []int) error
	FindNextHighlightedNode() error
}

const (
	redBold   = "\033[1;31m"
	whiteBold = "\033[1;97m"
	reset     = "\033[0m"
)

var controls = []string{
	"Any",
	"Key",
	"Value",
	"HasParent",
	"HasAnyParent",
	"ParentHasChildKey",
	"ParentHasChildValue",
	"ParentHasChildAny",
	"ChildHasKey",
	"ChildHasValue",
	"ChildHasAny",
	"AnyParentHasChildKey",
	"AnyParentHasChildValue",
	"AnyParentHasChildAny",
	"AnyChildHasKey",
	"AnyChildHasValue",
	"AnyChildHasAny",
}

func isMainCommand(command Command) bool {
	if command.Control == controls[0] || // Any
		command.Control == controls[1] || // Key
		command.Control == controls[2] { // Value
		return true
	}
	return false

}

var conditionals = []string{
	"==",
	"!=",
}

var operators = []string{
	"+",  // Union
	"|",  // Intersection
	"-",  // Subtraction of matching elements
	"&&", // Show all if second contains elements
	"<-", // Only show first if second contains elements
	"->", // Only show second elements
}
