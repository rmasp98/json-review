package search

import (
	"kube-review/nodelist"
	"regexp"
	"strings"
)

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
	// EXPRESSION a
	EXPRESSION
	// QUERY a
	QUERY
)

func (qe QueryEnum) String() string {
	return [...]string{"Regex", "Expression", "Query"}[qe]
}

// CmdFunc is enum for possible command functions
type CmdFunc int

const (
	//CMDNULL a
	CMDNULL CmdFunc = iota
	// CMDFINDNODES a
	CMDFINDNODES
	// CMDFINDRELATIVE a
	CMDFINDRELATIVE
)

func (cf CmdFunc) String() string {
	return [...]string{"Null", "FindNodes", "FindRelative"}[cf]
}

func (cf CmdFunc) template() []argTemplate {
	return [...][]argTemplate{[]argTemplate{}, findArgs, findRelArgs}[cf]
}

type sNodeList interface {
	GetNodesMatching(regex *regexp.Regexp, matchType nodelist.MatchType, equal bool) []int
	GetRelativesMatching(nodeIndex, relativeStartLevel, depth int, regex *regexp.Regexp, matchType nodelist.MatchType, equal bool) []int
	Filter(nodes []int) error
	Highlight(nodes []int)
	FindNextHighlight() error
	ResetView()
}

const (
	redBold   = "\033[1;31m"
	whiteBold = "\033[1;97m"
	reset     = "\033[0m"
)

var conditionals = []string{
	"==",
	"!=",
}

// Operators lists possible operators in expression search
var Operators = []string{
	// Operators should not match with anything below it (e.g. '-' conflicts with '->')
	"&&", // Show all if second contains elements
	"<-", // Only show first if both contain elements
	"->", // Only show second elements
	"+",  // Union
	"|",  // Intersection
	"-",  // Subtraction of matching elements
}

type argTemplate struct {
	name        string
	argType     string
	description string
}

func getArgIndexByName(name string, argTemp []argTemplate) int {
	for index, arg := range argTemp {
		if len(name) > 0 && len(name) <= len(arg.name) && strings.EqualFold(name, arg.name[:len(name)]) {
			return index
		}
	}
	return -1
}

var findArgs = []argTemplate{
	argTemplate{"regex", "regex", "quoted regex string"},
	argTemplate{"matchType", "MatchType", "attribute to match against"},
	argTemplate{"equal", "bool", "should match be equal or not equal to regex"},
	argTemplate{"output", "output", "variable that holds matched nodes. If exists, append to previous result"},
}

var findRelArgs = []argTemplate{
	argTemplate{"nodes", "input", "nodes to run against. Must be output of previous function call"},
	argTemplate{"regex", "regex", "quoted regex string"},
	argTemplate{"relativeStart", "int", "levels above nodes, that search should begin"},
	argTemplate{"depth", "int", "number of levels search should go down"},
	argTemplate{"matchType", "MatchType", "attribute to match against"},
	argTemplate{"equal", "bool", "should match be equal or not equal to regex"},
	argTemplate{"output", "output", "variable that holds matched nodes. If exists, append to previous result"},
}
