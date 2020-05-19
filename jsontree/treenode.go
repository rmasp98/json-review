package jsontree

import (
	"fmt"
	"regexp"
	"sort"
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

// TreeNode is a JSON node
type TreeNode struct {
	key           string
	value         string
	Prefix        string
	Level         int
	Parent        int
	Children      []int
	isExpanded    bool
	isFiltered    bool
	IsHighlighted bool
}

// IsVisible returns false if anything is hiding this node
func (t TreeNode) IsVisible() bool {
	return t.Level == 0 || (t.isExpanded && !t.isFiltered)
}

// SetExpanded defines if a node has been expanded or collapsed
func (t *TreeNode) SetExpanded(isExpanded bool) {
	t.isExpanded = isExpanded
}

// Filter sets the isfiletered flag
func (t *TreeNode) Filter(isFiltered bool) {
	t.isFiltered = isFiltered
}

// GetJSON returns a snippet of json that will be used by NodeList to reconstruct the original JSON
func (t TreeNode) GetJSON(full bool) string {
	if !t.isFiltered {
		var out string
		if t.key == "" || t.key[0:2] == "[]" || !full {
			out = t.value
		} else {
			out = strconv.Quote(t.key) + ": " + t.value
		}
		if t.IsHighlighted {
			return "\033[41m" + out + "\033[0m"
		}
		return out
	}
	return ""
}

// GetNode returns the key of node in nice format
func (t TreeNode) GetNode() string {
	return t.Prefix + strings.TrimLeft(t.key, "[]")
}

var brackets = map[string]string{"{": "}", "[": "]"}

// GetEnding returns any close brackets or a comma if not last on level
func (t TreeNode) GetEnding(lastOnLevel bool) string {
	if !t.isFiltered {
		if !lastOnLevel {
			return brackets[t.value] + ","
		}
		return brackets[t.value]
	}
	return ""
}

// GetCloseBracket stuff
func (t TreeNode) GetCloseBracket() string {
	if !t.isFiltered {
		return brackets[t.value]
	}
	return ""
}

// Clear sets IsFiltered and IsHighlighted flags to false and IsExpanded to true
func (t *TreeNode) Clear() {
	t.isFiltered = false
	t.IsHighlighted = false
	t.isExpanded = true
}

// Match returns true is regex matches key or value of node
func (t TreeNode) Match(r *regexp.Regexp, matchType MatchType) bool {
	return (matchType == ANY || matchType == KEY) && r.MatchString(t.key) ||
		(matchType == ANY || matchType == VALUE) && r.MatchString(t.value)

}

// CreateTreeNodes creates an array of TreeNodes that represents incoming JSON data
func CreateTreeNodes(data interface{}, level int) ([]TreeNode, error) {
	switch elem := data.(type) {
	case string:
		return []TreeNode{TreeNode{"", strconv.Quote(elem), "", level, 0, []int{}, true, false, false}}, nil
	case float64:
		return []TreeNode{TreeNode{"", strconv.FormatFloat(elem, 'g', -1, 64), "", level, 0, []int{}, true, false, false}}, nil
	case bool:
		return []TreeNode{TreeNode{"", strconv.FormatBool(elem), "", level, 0, []int{}, true, false, false}}, nil
	case nil:
		return []TreeNode{TreeNode{"", "null", "", level, 0, []int{}, true, false, false}}, nil
	case map[string]interface{}:
		return newMapNode(elem, level)
	case []interface{}:
		return newArrayNode(elem, level)
	}
	return []TreeNode{}, fmt.Errorf("Incorrectly formated Json")
}

func newMapNode(data map[string]interface{}, level int) ([]TreeNode, error) {
	var returnNodes []TreeNode
	for _, key := range getOrderedMapKeys(data) {
		nodes, err := processNode(key, data[key], level)
		if err != nil {
			return []TreeNode{}, err
		}
		returnNodes = append(returnNodes, nodes...)
	}
	return returnNodes, nil
}

func newArrayNode(data []interface{}, level int) ([]TreeNode, error) {
	var returnNodes []TreeNode
	for index, childInterface := range data {
		nodes, err := processNode("[]"+strconv.Itoa(index), childInterface, level)
		if err != nil {
			return []TreeNode{}, err
		}
		returnNodes = append(returnNodes, nodes...)
	}
	return returnNodes, nil
}

func processNode(key string, childInterface interface{}, level int) ([]TreeNode, error) {
	var nodes []TreeNode
	childNodes, createError := CreateTreeNodes(childInterface, level+1)
	if createError != nil {
		return []TreeNode{}, createError
	}

	value := getNodeValue(childNodes)
	nodes = append(nodes, TreeNode{key, value, "", level, 0, []int{}, true, false, false})
	if value == "{" || value == "[" {
		nodes = append(nodes, childNodes...)
	}
	return nodes, nil
}

func getNodeValue(childNodes []TreeNode) string {
	if len(childNodes) > 0 {
		if childNodes[0].key == "" {
			return childNodes[0].value
		} else if childNodes[0].key[0:2] == "[]" {
			return "["
		} else {
			return "{"
		}
	}
	return ""
}

func getOrderedMapKeys(data map[string]interface{}) []string {
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
