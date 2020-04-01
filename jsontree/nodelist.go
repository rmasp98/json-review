package jsontree

import (
	"encoding/json"
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
	regex      string
}

// NewNodeList requires a valid json string and returns a NodeList object
func NewNodeList(jsonData string) (NodeList, error) {
	var data interface{}
	json.Unmarshal([]byte(jsonData), &data)
	treeNodes, err := CreateTreeNodes(data, 1)
	if err != nil {
		return NodeList{}, err
	}

	nodes := []TreeNode{TreeNode{"Root", getNodeValue(treeNodes), "", 0, 0, true, true}}
	nodes = append(nodes, treeNodes...)
	nodeList := NodeList{nodes, 0, 0, 0, ""}
	for index := range nodeList.nodes {
		nodeList.nodes[index].Parent = nodeList.getParentIndex(index)
		nodeList.updatePrefix(index)
	}
	return nodeList, nil
}

// Size returns the number of visible nodes
func (n NodeList) Size() int {
	size := 0
	for index := range n.nodes {
		if n.isVisible(index) {
			size++
		}
	}
	return size
}

// GetNodes returns a formated string list of visible nodes from topNode
// and is only num long
func (n NodeList) GetNodes(num int) string {
	var output string
	for index := n.topNode; num > 0 && index < len(n.nodes); index++ {
		if n.isVisible(index) {
			n.updatePrefix(index)
			output += n.nodes[index].GetNode() + "\n"
			num--
		}
	}
	return strings.TrimRight(output, "\n")
}

// GetJSON returns a formatted json string, num lines long for
// the active node. Fields can be hidden using the Search function
func (n NodeList) GetJSON(num int) string {
	var finalJSON string
	if num > 0 {
		baseLevel := n.nodes[n.activeNode].Level
		for index := n.jsonStart; index < len(n.nodes); {
			levelDiff := n.nodes[index].Level - baseLevel
			if n.nodes[index].IsMatched {
				finalJSON += "\n" + n.nodes[index].GetJSON(levelDiff)
				num--
			}
			finalJSON += n.getNodeEndings(index, baseLevel, &num)

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
		if n.isVisible(index) {
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

// Search stuff
func (n *NodeList) Search(regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	for index := 1; index < len(n.nodes); index++ {
		if n.nodes[index].Search(r) {
			for i := n.nodes[index].Parent; i != 0; i = n.nodes[i].Parent {
				n.nodes[i].IsMatched = true
			}
		}
	}
	// Bodge to prevent everything from disappearing if content too small
	n.topNode = 0
	n.setActiveNode(0)
	return nil
}

// Internal functions/////////////////////////////////////////////////////

func (n NodeList) isVisible(nodeIndex int) bool {
	return nodeIndex == 0 || (n.nodes[nodeIndex].IsExpanded && n.nodes[nodeIndex].IsMatched)
}

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
		if nodeLevel == targetLevel && n.isVisible(index) {
			return false
		} else if nodeLevel < targetLevel {
			return true
		}
	}
	return true
}

var brackets = map[string]string{"{": "}", "[": "]"}

// getNodeEndings returns any closed brackets or commas required after nodeIndex
func (n NodeList) getNodeEndings(nodeIndex, baseLevel int, count *int) string {
	var endings string
	if nodeIndex+1 < len(n.nodes) {
		currentLevel := n.nodes[nodeIndex].Level
		nextLevel := n.nodes[nodeIndex+1].Level
		// Close all brackets
		for level := currentLevel; level > nextLevel && level > baseLevel; level-- {
			if *count <= 0 {
				break
			}
			levelParent := n.getLevelParent(nodeIndex, level)
			if n.nodes[levelParent].IsMatched {
				endings += "\n" + strings.Repeat(spacing, level-baseLevel-1) + brackets[n.nodes[levelParent].value]
				*count--
			}
		}
		// If finish closing brackets and not reached end must be another node on same level as last bracket
		if endings != "" && nextLevel > baseLevel {
			endings += ","
		}
		if currentLevel == nextLevel && n.nodes[nodeIndex].IsMatched && n.nodes[nodeIndex+1].IsMatched {
			endings += ","
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
		n.nodes[index].IsExpanded = visible
	}
}

func (n NodeList) getNodeIndex(visibleIndex int) int {
	for index := n.topNode; index < len(n.nodes); index++ {
		if n.isVisible(index) {
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
		if n.isVisible(index) {
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
