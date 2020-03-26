package jsontree

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TreeNode struct {
	key       string
	value     string
	Prefix    string
	Level     int
	IsVisible bool
}

func (t TreeNode) GetJSON(level int) string {
	if t.key == "" || t.key[0:2] == "[]" || level == 0 {
		return strings.Repeat(spacing, level) + t.value
	}
	return strings.Repeat(spacing, level) + strconv.Quote(t.key) + ": " + t.value
}

func (t TreeNode) GetNode() string {
	return t.Prefix + strings.TrimLeft(t.key, "[]")
}

func CreateTreeNodes(data interface{}, level int) ([]TreeNode, error) {
	switch elem := data.(type) {
	case string:
		return []TreeNode{TreeNode{"", strconv.Quote(elem), "", level, true}}, nil
	case float64:
		return []TreeNode{TreeNode{"", strconv.FormatFloat(elem, 'g', -1, 64), "", level, true}}, nil
	case bool:
		return []TreeNode{TreeNode{"", strconv.FormatBool(elem), "", level, true}}, nil
	case nil:
		return []TreeNode{TreeNode{"", "null", "", level, true}}, nil
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
	nodes = append(nodes, TreeNode{key, value, "", level, true})
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
