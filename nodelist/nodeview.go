package nodelist

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var spacing = "    "

type searchFunctionType = func(*Node) bool

type nodeView struct {
	node          *Node
	parent        int
	children      []int
	prefix        string
	isHighlighted bool
}

// View presents a subset view of master based on filters or splits
type View struct {
	nodes []nodeView
}

// NewView stuff
// checks and updates children
func NewView(nodes []*Node) (View, error) {
	if len(nodes) == 0 || nodes[0].key != "Root" {
		return View{}, fmt.Errorf("New views must at least have a Root node")
	}
	nodeViews := make([]nodeView, len(nodes))
	for index, node := range nodes {
		nodeViews[index] = nodeView{node, 0, []int{}, "", false}
	}
	v := View{nodeViews}
	v.updateNodeRelationships()
	return v, nil
}

// Size returns number of nodes in view
func (v View) Size() int {
	return len(v.nodes)
}

// GetNodes returns formatted keys of nodes from start to start+num
func (v View) GetNodes(start, num int) string {
	var nodes string
	for index := start; index < start+num && index < len(v.nodes); index++ {
		if v.nodes[index].prefix == "" {
			v.updatePrefix(index)
		}
		nodes += v.nodes[index].prefix + v.nodes[index].node.GetNode() + "\n"
	}
	return strings.TrimRight(nodes, "\n")
}

// GetJSON returns formated JSON for nodeIndex. The JSON output can be offset and
// number of lines returned limited using the offset and num inputs
func (v View) GetJSON(nodeIndex, offset, num int) string {
	return v.getJSON(nodeIndex, 0, offset, &num)
}

// GetNodesMatching searches entire view for matches of matchtype to regex. Set equal to false to invert result
func (v View) GetNodesMatching(regex *regexp.Regexp, matchType MatchType, equal bool) []int {
	searchFunction := getSearchFunction(matchType, regex, equal)
	return v.getChildrenMatching(0, -1, searchFunction)
}

// GetRelativesMatching searches nodes relative to nodeIndex in similar fashion to GetNodesMatching.
// relativeStartLevel defines how many levels above nodeIndex the search should start from
// and depth defines how many levels of children from relativeStartLevel should be searched.
// To search a particular parent, set depth to zero, otherwise that parent is ignored
func (v View) GetRelativesMatching(nodeIndex, relativeStartLevel, depth int, regex *regexp.Regexp, matchType MatchType, equal bool) []int {
	var matchedIndices []int
	var startIndex int
	for index := nodeIndex; index >= 0 && relativeStartLevel >= 0; index = v.nodes[index].parent {
		startIndex = index
		relativeStartLevel--
	}
	searchFunction := getSearchFunction(matchType, regex, equal)
	if depth == 0 && searchFunction(v.nodes[startIndex].node) {
		matchedIndices = append(matchedIndices, startIndex)
	}
	matchedIndices = append(matchedIndices, v.getChildrenMatching(startIndex, depth, searchFunction)...)
	return matchedIndices
}

// Filter returns a new view with the defined node indices along with their parents
func (v View) Filter(nodeIndices []int) (View, error) {
	finalIndices := v.appendAndSortParentIndices(nodeIndices)
	var nodes = make([]*Node, 0, len(finalIndices))
	for _, index := range finalIndices {
		nodes = append(nodes, v.nodes[index].node)
	}
	if len(nodes) == 0 {
		nodes = append(nodes, &Node{"Root", "", 0})
	}
	return NewView(nodes)
}

// Highlight clears current highlight and applies highlight to node pointed to by nodeIndices
func (v *View) Highlight(nodeIndices []int) {
	for index := range v.nodes {
		v.nodes[index].isHighlighted = false
	}
	for _, index := range nodeIndices {
		v.nodes[index].isHighlighted = true
	}
}

// FindNextHighlight will return a new offset to show next highlight
func (v View) FindNextHighlight(nodeIndex, startOffset int) (int, error) {
	numTotalChildren := v.getLastChild(nodeIndex) - nodeIndex
	if startOffset < numTotalChildren {
		offset := (startOffset + 1) % numTotalChildren
		looped := false
		for {
			if v.nodes[nodeIndex+offset].isHighlighted {
				return offset, nil
			}
			offset = (offset + 1) % numTotalChildren
			if looped {
				break
			} else if offset == startOffset {
				looped = true
			}
		}
	}
	return 0, fmt.Errorf("No highlighted nodes in nodeIndex: %d", nodeIndex)
}

/////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS

func (v View) getJSON(nodeIndex, level, offset int, num *int) string {
	var finalJSON string
	if *num != 0 {
		JSON := v.nodes[nodeIndex].node.GetJSON(level > 0)
		if JSON != "" && nodeIndex >= offset {
			finalJSON = strings.Repeat(spacing, level) + JSON
			*num--
		}
		var childrenJSON string
		for _, childIndex := range v.nodes[nodeIndex].children {
			childJSON := v.getJSON(childIndex, level+1, offset, num)
			if childJSON != "" {
				childrenJSON += childJSON + ",\n"
			}
		}
		if finalJSON != "" && childrenJSON != "" {
			finalJSON += "\n"
		}
		finalJSON += strings.TrimRight(childrenJSON, ",\n")

		closeBracket := v.nodes[nodeIndex].node.GetCloseBracket()
		if closeBracket != "" && *num != 0 && v.getLastChild(nodeIndex) >= offset {
			finalJSON += "\n" + strings.Repeat(spacing, level) + closeBracket
			*num--
		}
	}
	return finalJSON
}

func (v View) getChildrenMatching(nodeIndex int, levels int, searchFunction searchFunctionType) []int {
	var matchedIndices []int
	if levels != 0 {
		for _, childIndex := range v.nodes[nodeIndex].children {
			if searchFunction(v.nodes[childIndex].node) {
				matchedIndices = append(matchedIndices, childIndex)
			}
			matchedIndices = append(matchedIndices, v.getChildrenMatching(childIndex, levels-1, searchFunction)...)
		}
	}
	return matchedIndices
}

func (v View) getLastChild(nodeIndex int) int {
	children := v.nodes[nodeIndex].children
	lastIndex := nodeIndex
	for len(children) > 0 {
		lastIndex = children[len(children)-1]
		children = v.nodes[lastIndex].children
	}
	return lastIndex
}

func (v *View) updatePrefix(index int) {
	if index > 0 {
		parentIndex := v.nodes[index].parent
		if parentIndex > 0 && v.nodes[parentIndex].prefix == "" {
			v.updatePrefix(parentIndex)
		}
		newPrefix := convertParentPrefix(v.nodes[parentIndex].prefix)

		if v.isLastInLevel(index, parentIndex) {
			newPrefix += "└──"
		} else {
			newPrefix += "├──"
		}

		v.nodes[index].prefix = newPrefix
	}
}

func (v View) isLastInLevel(nodeIndex, parentIndex int) bool {
	siblings := v.nodes[parentIndex].children
	if siblings[len(siblings)-1] == nodeIndex {
		return true
	}
	return false
}

func (v *View) updateNodeRelationships() {
	for index := range v.nodes {
		parentIndex := v.getParentIndex(index)
		v.nodes[index].parent = parentIndex
		if parentIndex >= 0 {
			v.nodes[parentIndex].children = append(v.nodes[parentIndex].children, index)
		}
	}
}

func (v View) getParentIndex(nodeIndex int) int {
	targetLevel := v.nodes[nodeIndex].node.GetLevel()
	for index := nodeIndex - 1; index >= 0; index-- {
		currentLevel := v.nodes[index].node.GetLevel()
		if currentLevel == targetLevel {
			return v.nodes[index].parent
		} else if currentLevel < targetLevel {
			return index
		}
	}
	return -1
}

func (v View) appendAndSortParentIndices(nodeIndices []int) []int {
	var fullIndices = make(map[int]struct{}, len(v.nodes))
	for _, index := range nodeIndices {
		for i := index; i >= 0; i = v.nodes[i].parent {
			if _, ok := fullIndices[i]; ok {
				break
			}
			fullIndices[i] = struct{}{}
		}
	}
	finalIndices := make([]int, 0, len(fullIndices))
	for index := range fullIndices {
		finalIndices = append(finalIndices, index)
	}
	sort.Ints(finalIndices)
	return finalIndices
}

func getSearchFunction(matchType MatchType, r *regexp.Regexp, equal bool) searchFunctionType {
	if matchType == KEY {
		return func(node *Node) bool { return node.MatchKey(r) == equal }
	} else if matchType == VALUE {
		return func(node *Node) bool { return node.MatchValue(r) == equal }
	}
	return func(node *Node) bool { return node.Match(r) == equal }
}

var prefixConvert = map[rune]rune{
	'─': ' ',
	'│': '│',
	'├': '│',
	'└': ' ',
	' ': ' ',
}

func convertParentPrefix(prefix string) string {
	var output string
	for _, char := range prefix {
		output += string(prefixConvert[char])
	}
	return output
}
