package jsontree

import (
	"encoding/json"
	"fmt"
	"kube-review/utils"
	"log"
	"regexp"
	"strings"
)

const spacing = "   "

// NodeList contains JSON data and ways to extract useful UI inforamtion from it
type NodeList struct {
	nodes      []TreeNode
	topNode    int
	activeNode int
	jsonStart  int
	findNode   int
}

// NewNodeList requires a valid json string and returns a NodeList object
func NewNodeList(jsonData string) (NodeList, error) {
	var data interface{}
	json.Unmarshal([]byte(jsonData), &data)
	treeNodes, err := CreateTreeNodes(data, 1)
	if err != nil {
		return NodeList{}, err
	}

	nodes := []TreeNode{TreeNode{"Root", getNodeValue(treeNodes), "", 0, 0, []int{}, true, false, false}}
	nodes = append(nodes, treeNodes...)
	nodeList := NodeList{nodes, 0, 0, 0, -1}
	for index := range nodeList.nodes {
		if index != 0 {
			parentIndex := nodeList.getParentIndex(index)
			nodeList.nodes[index].Parent = parentIndex
			nodeList.nodes[parentIndex].Children = append(nodeList.nodes[parentIndex].Children, index)
		}
	}
	for index := range nodeList.nodes {
		if index != 0 {
			nodeList.updatePrefix(index)
		}
	}
	return nodeList, nil
}

// Size returns the number of visible nodes
func (n NodeList) Size() int {
	size := 0
	for _, node := range n.nodes {
		if node.IsVisible() {
			size++
		}
	}
	return size
}

// Clear resets any affect on nodes
func (n *NodeList) Clear() {
	for index := range n.nodes {
		n.nodes[index].Clear()
	}
}

// TODO: Move IsFiltered and IsExpanded into functions that are checked either
// through IsVisible function or via other functions

// FilterAll sets all (but root) nodes to filtered
func (n *NodeList) FilterAll() {
	for index := 1; index < len(n.nodes); index++ {
		n.nodes[index].Filter(true)
	}
}

// GetNodes returns a formated string list of visible nodes from topNode
// and is only num long
func (n NodeList) GetNodes(num int) string {
	var output string
	for index := n.topNode; num > 0 && index < len(n.nodes); index++ {
		if n.nodes[index].IsVisible() {
			n.updatePrefix(index)
			output += n.nodes[index].GetNode() + "\n"
			num--
		}
	}
	return strings.TrimRight(output, "\n")
}

// GetJSON returns a formatted json string, num lines long for
// the active node. Fields can be hidden using the Filter function
// Passing -1 will remove limit on returned lines
func (n NodeList) GetJSON(num int) string {
	return n.getJSON(n.activeNode, 0, &num)
}

// MoveTopNode changes the start position (topNode) of what GetNodes returns
// relative to its current position
func (n *NodeList) MoveTopNode(offset int) {
	visibleActiveIndex := n.getVisibleIndex(n.activeNode)
	var step int
	step, offset = getDirectionAndAbs(offset)
	for index := n.topNode; index < len(n.nodes) && index >= 0 && offset >= 0; index = index + step {
		if n.nodes[index].IsVisible() {
			n.topNode = index
			offset--
		}
	}
	n.SetActiveNode(visibleActiveIndex)
}

// SetActiveNode lets NodeList know the highlighted node in editor.
// This is the node that all actions will be performed on
// visibleIndex should be index visible to the user
func (n *NodeList) SetActiveNode(visibleIndex int) {
	newNodeIndex := n.getNodeIndex(visibleIndex)
	if newNodeIndex != -1 {
		n.setActiveNode(newNodeIndex)
	}
}

// MoveJSONPosition offsets the Json shown in DISPLAY
func (n *NodeList) MoveJSONPosition(offset int) {
	oldVisiblePosition := n.getVisibleIndex(n.jsonStart)
	newPosition := n.getNodeIndex(oldVisiblePosition + offset)
	if newPosition < n.activeNode {
		n.jsonStart = n.activeNode
	} else {
		lastChild := n.getLastChild(n.activeNode)
		if newPosition > n.getLastChild(n.activeNode) {
			n.jsonStart = lastChild
		} else {
			n.jsonStart = newPosition
		}
	}
}

// ExpandActiveNode makes all children of active node visible
func (n *NodeList) ExpandActiveNode() {
	n.alterNodesExpandedness(
		n.activeNode+1, n.getLevelEndIndex(n.activeNode+1), true,
	)
}

// CollapseActiveNode makes all nodes on level with active node and below
// invisible and returns the visible index of the parent node
func (n *NodeList) CollapseActiveNode() int {
	parentIndex := n.nodes[n.activeNode].Parent
	n.alterNodesExpandedness(
		parentIndex+1,
		n.getLevelEndIndex(n.activeNode),
		false,
	)
	if parentIndex < n.topNode {
		n.topNode = parentIndex
	}
	log.Print("Collapse")
	n.setActiveNode(parentIndex)

	return n.getVisibleIndex(n.activeNode)
}

// GetNodesMatching return node indicies matching regex in matchType
func (n NodeList) GetNodesMatching(regex *regexp.Regexp, matchType MatchType, equal bool) []int {
	var matchList []int
	for index, node := range n.nodes {
		if node.Match(regex, matchType) == equal {
			matchList = append(matchList, index)
		}
	}
	return matchList
}

// GetChildrenMatching stuff
func (n NodeList) GetChildrenMatching(nodeIndex int, regex *regexp.Regexp, matchType MatchType, equal bool, recursive bool) []int {
	var matchList = make([]int, len(n.nodes))
	var matchIndex = 0
	for _, index := range n.nodes[nodeIndex].Children {
		if n.nodes[index].Match(regex, matchType) == equal {
			matchList[matchIndex] = index
			matchIndex++
		}
		if recursive {
			for _, childIndex := range n.GetChildrenMatching(index, regex, matchType, equal, recursive) {
				matchList[matchIndex] = childIndex
				matchIndex++
			}
		}
	}
	return matchList[:matchIndex]
}

// GetParentChildrenMatching stuff
func (n NodeList) GetParentChildrenMatching(nodeIndex int, regex *regexp.Regexp, matchType MatchType, equal bool, recursive bool) []int {
	var matchList []int
	for index := n.nodes[nodeIndex].Parent; ; index = n.nodes[index].Parent {
		matchList = append(matchList, n.GetChildrenMatching(index, regex, matchType, equal, false)...)
		if !recursive || index == 0 {
			break
		}
	}
	return matchList
}

// ApplyFilter stuff
func (n *NodeList) ApplyFilter(nodes []int) error {
	n.FilterAll()
	for _, index := range nodes {
		if index > 0 && index < len(n.nodes) {
			n.nodes[index].Filter(false)
			for i := n.nodes[index].Parent; i != 0; i = n.nodes[i].Parent {
				n.nodes[i].Filter(false)
			}
		}
	}
	return nil
}

// ApplyHighlight sets IsHighlighted on all nodes in nodes
func (n *NodeList) ApplyHighlight(nodes []int) error {
	for index := range n.nodes {
		n.nodes[index].IsHighlighted = false
	}
	if len(n.nodes) > len(nodes) {
		for _, index := range nodes {
			if index > 0 && index < len(n.nodes) {
				n.nodes[index].IsHighlighted = true
			}
		}
	}
	return nil
}

// FindNextHighlightedNode moves GetJSON to the next matched node
func (n *NodeList) FindNextHighlightedNode() error {
	levelSize := n.getLevelEndIndex(n.activeNode+1) - n.activeNode
	startOffset := n.jsonStart - n.activeNode
	offset := startOffset
	for {
		index := n.activeNode + offset
		if n.nodes[index].IsHighlighted && index != n.findNode {
			n.jsonStart = index
			n.findNode = index
			return nil
		}

		offset = (offset + 1) % levelSize
		if offset == startOffset {
			break
		}
	}
	return fmt.Errorf("No matches in this node")
}

// Save writes content of GetJSON to file
func (n NodeList) Save(filename string) error {
	return utils.Save(filename, n.GetJSON(-1), true)
}

// Internal functions/////////////////////////////////////////////////////

func (n *NodeList) setActiveNode(nodeIndex int) {
	n.activeNode = nodeIndex
	n.jsonStart = nodeIndex
}

func (n NodeList) getParentIndex(nodeIndex int) int {
	currentLevel := n.nodes[nodeIndex].Level
	for index := nodeIndex; index >= 0; index-- {
		if n.nodes[index].Level < currentLevel {
			return index
		}
	}
	return 0
}

func (n NodeList) getJSON(nodeIndex, level int, num *int) string {
	if *num != 0 {
		var finalJSON string
		JSON := n.nodes[nodeIndex].GetJSON(level > 0)
		if JSON != "" && nodeIndex >= n.jsonStart {
			finalJSON = strings.Repeat(spacing, level) + JSON
			*num--
		}
		for pos, childIndex := range n.nodes[nodeIndex].Children {
			childJSON := n.getJSON(childIndex, level+1, num)
			if childJSON != "" {
				if finalJSON == "" {
					finalJSON += childJSON
				} else {
					finalJSON += "\n" + childJSON
					if pos < len(n.nodes[nodeIndex].Children)-1 {
						finalJSON += ","
					}
				}
			}
		}
		closeBracket := n.nodes[nodeIndex].GetCloseBracket()

		if closeBracket != "" && *num != 0 && n.getLastChild(nodeIndex) >= n.jsonStart {
			finalJSON += "\n" + strings.Repeat(spacing, level) + closeBracket
			*num--
		}
		return finalJSON
	}
	return ""
}

// Get the last index in the current level including any children
func (n NodeList) getLevelEndIndex(nodeIndex int) int {
	currentLevel := n.nodes[nodeIndex].Level
	for index := nodeIndex; index < len(n.nodes); index++ {
		if n.nodes[index].Level < currentLevel {
			return index - 1
		}
	}
	return len(n.nodes) - 1
}

func (n *NodeList) updatePrefix(index int) {
	if index > 0 {
		parentIndex := n.nodes[index].Parent
		newPrefix := convertParentPrefix(n.nodes[parentIndex].Prefix)

		if n.isLastInLevel(index, parentIndex) {
			newPrefix += "└──"
		} else {
			newPrefix += "├──"
		}

		n.nodes[index].Prefix = newPrefix
	}
}

func (n NodeList) isLastInLevel(nodeIndex, parentIndex int) bool {
	siblings := n.nodes[parentIndex].Children
	for i := len(siblings) - 1; i >= 0; i-- {
		if !n.nodes[siblings[i]].isFiltered {
			if siblings[i] == nodeIndex {
				return true
			}
			return false
		}
	}
	return false
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

//Sorry for the name...
func (n *NodeList) alterNodesExpandedness(startIndex, endIndex int, visible bool) {
	for index := startIndex; index < endIndex+1; index++ {
		n.nodes[index].SetExpanded(visible)
	}
}

func (n NodeList) getNodeIndex(visibleIndex int) int {
	for index := n.topNode; index < len(n.nodes); index++ {
		if n.nodes[index].IsVisible() {
			if visibleIndex == 0 {
				return index
			}
			visibleIndex--
		}
	}
	return -1
}

func (n NodeList) getVisibleIndex(nodeIndex int) int {
	visibleIndex := 0
	for index := n.topNode; index < nodeIndex; index++ {
		if n.nodes[index].IsVisible() {
			visibleIndex++
		}
	}
	return visibleIndex
}

func (n NodeList) getLastChild(nodeIndex int) int {
	children := n.nodes[nodeIndex].Children
	lastIndex := nodeIndex
	for len(children) > 0 {
		lastIndex = children[len(children)-1]
		children = n.nodes[lastIndex].Children
	}
	return lastIndex
}

func getDirectionAndAbs(offset int) (int, int) {
	if offset < 0 {
		return -1, -offset
	}
	return 1, offset
}
