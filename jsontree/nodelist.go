package jsontree

import (
	"encoding/json"
	"fmt"
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

	nodes := []TreeNode{TreeNode{"Root", getNodeValue(treeNodes), "", 0, 0, true, false, false}}
	nodes = append(nodes, treeNodes...)
	nodeList := NodeList{nodes, 0, 0, 0, -1}
	for index := range nodeList.nodes {
		nodeList.nodes[index].Parent = nodeList.getParentIndex(index)
		nodeList.updatePrefix(index)
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
func (n NodeList) GetJSON(num int) string {
	var finalJSON string
	if num > 0 {
		baseLevel := n.nodes[n.activeNode].Level
		for index := n.jsonStart; index < len(n.nodes); {
			levelDiff := n.nodes[index].Level - baseLevel
			if nodeJSON := n.nodes[index].GetJSON(levelDiff); nodeJSON != "" {
				finalJSON += "\n" + nodeJSON
				num--
			}
			finalJSON += n.getNodeEndings(index, baseLevel)

			if num > 0 && index == len(n.nodes)-1 {
				finalJSON += "\n" + brackets[n.nodes[0].value]
			}
			index++
			if num <= 0 || (index < len(n.nodes) && n.nodes[index].Level <= baseLevel) {
				break
			}
		}
	}
	return strings.Trim(finalJSON, "\n")
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
		n.jsonStart = newPosition
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
func (n NodeList) GetNodesMatching(regex string, matchType MatchType, equal bool) []int {
	var matchList []int
	r, _ := regexp.Compile(regex)
	for index, node := range n.nodes {
		if node.Match(r, matchType) == equal {
			matchList = append(matchList, index)
		}
	}
	return matchList
}

// GetChildrenMatching stuff
func (n NodeList) GetChildrenMatching(nodeIndex int, regex string, matchType MatchType, equal bool, recursive bool) []int {
	return []int{}
}

// GetParentChildrenMatching stuff
func (n NodeList) GetParentChildrenMatching(nodeIndex int, regex string, matchType MatchType, equal bool, recursive bool) []int {
	return []int{}
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

// TODO: Move to treeNode suppling parent prefix and bool isLastOnLevel

func (n *NodeList) updatePrefix(index int) {
	if index > 0 {
		newPrefix := convertParentPrefix(n.nodes[n.nodes[index].Parent].Prefix)

		if n.isLastOnLevel(index) {
			newPrefix += "└──"
		} else {
			newPrefix += "├──"
		}

		n.nodes[index].Prefix = newPrefix
	}
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

func (n NodeList) isLastOnLevel(currentIndex int) bool {
	if currentIndex == len(n.nodes)-1 {
		return true
	}
	targetLevel := n.nodes[currentIndex].Level
	for index := currentIndex + 1; index < len(n.nodes); index++ {
		nodeLevel := n.nodes[index].Level
		if nodeLevel == targetLevel && n.nodes[index].IsVisible() {
			return false
		} else if nodeLevel < targetLevel {
			return true
		}
	}
	return true
}

// getNodeEndings returns any closed brackets or commas required after nodeIndex
func (n NodeList) getNodeEndings(nodeIndex, baseLevel int) string {
	var endings string
	if nodeIndex+1 < len(n.nodes) {
		currentLevel := n.nodes[nodeIndex].Level
		nextLevel := n.nodes[nodeIndex+1].Level
		if nextLevel <= currentLevel {
			endings = n.nodes[nodeIndex].GetEnding(n.isLastOnLevel(nodeIndex))
		}
		for level := currentLevel; level > nextLevel && level > baseLevel; level-- {
			levelParent := n.getLevelParent(nodeIndex, level)
			nodeEnding := n.nodes[levelParent].GetEnding(n.isLastOnLevel(levelParent))
			if nodeEnding != "" {
				endings += "\n" + strings.Repeat(spacing, level-baseLevel-1) + nodeEnding
			}
		}
	}
	return endings
}

func (n NodeList) getLevelParent(nodeIndex, level int) int {
	for index := nodeIndex; index > 0; index-- {
		if n.nodes[index].Level == level-1 {
			return index
		}
	}
	return 0
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

func getDirectionAndAbs(offset int) (int, int) {
	if offset < 0 {
		return -1, -offset
	}
	return 1, offset
}
