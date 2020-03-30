package jsontree

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// TreeNode is a JSON node
type TreeNode struct {
	key        string
	value      string
	Prefix     string
	Level      int
	Parent     int
	IsExpanded bool
	IsMatched  bool
}

// GetJSON returns a snippet of json that will be used by NodeList to reconstruct the original JSON
func (t TreeNode) GetJSON(level int) string {
	if t.key == "" || t.key[0:2] == "[]" || level == 0 {
		return strings.Repeat(spacing, level) + t.value
	}
	return strings.Repeat(spacing, level) + strconv.Quote(t.key) + ": " + t.value
}

// GetNode returns the key of node in nice format
func (t TreeNode) GetNode() string {
	return t.Prefix + strings.TrimLeft(t.key, "[]")
}

// Search checks key and value of node agaisnt regex
func (t *TreeNode) Search(r *regexp.Regexp) bool {
	t.IsMatched = r.MatchString(t.key) || r.MatchString(t.value)
	return t.IsMatched
}

// CreateTreeNodes creates an array of TreeNodes that represents incoming JSON data
func CreateTreeNodes(data interface{}, level int) ([]TreeNode, error) {
	switch elem := data.(type) {
	case string:
		return []TreeNode{TreeNode{"", strconv.Quote(elem), "", level, 0, true, true}}, nil
	case float64:
		return []TreeNode{TreeNode{"", strconv.FormatFloat(elem, 'g', -1, 64), "", level, 0, true, true}}, nil
	case bool:
		return []TreeNode{TreeNode{"", strconv.FormatBool(elem), "", level, 0, true, true}}, nil
	case nil:
		return []TreeNode{TreeNode{"", "null", "", level, 0, true, true}}, nil
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
	nodes = append(nodes, TreeNode{key, value, "", level, 0, true, true})
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
