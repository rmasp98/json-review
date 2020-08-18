package nodelist

import (
	"regexp"
	"strconv"
	"strings"
)

// MatchType indicates what will be used to match (key, value or both)
type MatchType int8

const (
	// ANY a
	ANY MatchType = iota
	// KEY a
	KEY
	// VALUE a
	VALUE
)

func (mt MatchType) String() string {
	return [...]string{"Any", "Key", "Value"}[mt]
}

// Node stored information about node is JSON data
type Node struct {
	key   string
	value string
	level int
}

// NewNode stuff
func NewNode(key, value string, level int) Node {
	return Node{key, value, level}
}

// GetJSON returns formatted JSON for the node. If full is false, the key is excluded
func (n Node) GetJSON(full bool) string {
	if (len(n.key) > 1 && n.key[0:2] == "[]") || !full {
		return n.value
	}
	return strconv.Quote(n.key) + ": " + n.value
}

// GetNode returns the key for the node
func (n Node) GetNode() string {
	return strings.TrimLeft(n.key, "[]")
}

// GetLevel returns the nodes level
func (n Node) GetLevel() int {
	return n.level
}

var brackets = map[string]string{"{": "}", "[": "]"}

// GetCloseBracket returns the correct close bracket if map or array, otherwise returns empty
func (n Node) GetCloseBracket() string {
	return brackets[n.value]
}

// Match returns true if regex matches either key or value
func (n Node) Match(r *regexp.Regexp) bool {
	return n.MatchKey(r) || n.MatchValue(r)
}

// MatchKey returns true if regex matches key
func (n Node) MatchKey(r *regexp.Regexp) bool {
	return r.MatchString(n.key)
}

// MatchValue returns true if regex matches value
func (n Node) MatchValue(r *regexp.Regexp) bool {
	return r.MatchString(strings.Trim(n.value, "\""))
}

// UpdateValue allows parser to insert values into map/array nodes
func (n *Node) UpdateValue(value string) {
	n.value = value
}
