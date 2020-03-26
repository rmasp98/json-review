package jsontree

import (
	"encoding/json"
	"log"
	"strings"
)

const spacing = "   "

type NodeList struct {
	nodes      []TreeNode
	topNode    int
	activeNode int
}

// NewNodeList requires a valid json string and returns a NodeList object
func NewNodeList(jsonData string) (NodeList, error) {
	var data interface{}
	json.Unmarshal([]byte(jsonData), &data)
	treeNodes, err := CreateTreeNodes(data, 1)
	if err != nil {
		return NodeList{}, err
	}

	nodes := []TreeNode{TreeNode{"Root", getNodeValue(treeNodes), "", 0, true}}
	nodes = append(nodes, treeNodes...)
	nodeList := NodeList{nodes, 0, 0}
	for index, _ := range nodeList.nodes {
		nodeList.updatePrefix(index)
	}
	return nodeList, nil
}

// Size returns the number of visible nodes
func (n NodeList) Size() int {
	size := 0
	for _, node := range n.nodes {
		if node.IsVisible {
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
		if n.nodes[index].IsVisible {
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
	var output string
	if num > 0 {
		output = n.nodes[n.activeNode].GetJSON(0) + "\n"
		num--
		baseLevel := n.nodes[n.activeNode].Level
		for index := n.activeNode + 1; index < len(n.nodes) && n.nodes[index].Level > baseLevel && num > 0; index++ {
			levelDiff := n.nodes[index].Level - baseLevel
			output += n.nodes[index].GetJSON(levelDiff)
			num--
			output += n.outputAnyCloseBrackets(index, &num, baseLevel) + "\n"
			if num > 0 && index == len(n.nodes)-1 {
				output += braces[n.nodes[0].value]
			}
		}
	}
	return strings.TrimRight(output, "\n")
}

// MoveTopNode changes the start position (topNode) of what GetNodes returns
// relative to its current position
func (n *NodeList) MoveTopNode(offset int) {
	visibleActiveIndex := n.getVisibleIndex(n.activeNode)
	var step int
	step, offset = getDirectionAndAbs(offset)
	for index := n.topNode; index < len(n.nodes) && index >= 0 && offset >= 0; index = index + step {
		if n.nodes[index].IsVisible {
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
		n.activeNode = newNodeIndex
	}
}

// ExpandActiveNode makes all children of active node visible
func (n *NodeList) ExpandActiveNode() {
	n.alterNodesVisibility(
		n.activeNode+1, n.getLevelEndIndex(n.activeNode+1), true,
	)
}

// CollapseActiveNode makes all nodes on level with active node and below
// invisible and returns the visible index of the parent node
func (n *NodeList) CollapseActiveNode() int {
	parentIndex := n.getParentIndex(n.activeNode)
	n.alterNodesVisibility(
		parentIndex+1,
		n.getLevelEndIndex(n.activeNode),
		false,
	)
	if parentIndex < n.topNode {
		n.topNode = parentIndex
	}
	log.Print("Collapse")
	n.activeNode = parentIndex
	return n.getVisibleIndex(n.activeNode)
}

func (n *NodeList) Search(regex string) {
	// r := regexp.MustCompile(regex)

}

// Internal functions/////////////////////////////////////////////////////

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
		newPrefix := convertParentPrefix(n.nodes[n.getParentIndex(index)].Prefix)

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

func (n NodeList) isLastOnLevel(index int) bool {
	if index == len(n.nodes)-1 {
		return true
	}
	targetLevel := n.nodes[index].Level
	for _, node := range n.nodes[index+1:] {
		if node.IsVisible && node.Level == targetLevel {
			return false
		} else if node.Level < targetLevel {
			return true
		}
	}
	return true
}

func (n NodeList) outputAnyCloseBrackets(index int, num *int, baseLevel int) string {
	nodeLevel := n.nodes[index].Level
	var output string
	if index < len(n.nodes)-1 {
		if nodeLevel > n.nodes[index+1].Level {
			for i := nodeLevel; i > n.nodes[index+1].Level && i > baseLevel; i-- {
				if *num <= 0 {
					return output
				}
				output += "\n" + strings.Repeat(spacing, i-1-baseLevel) + n.getParentType(i, index)
				(*num)--
			}
			if n.nodes[index+1].Level > baseLevel {
				output += ","
			}
		} else if nodeLevel == n.nodes[index+1].Level {
			output += ","
		}
	}
	return output
}

var braces = map[string]string{"{": "}", "[": "]"}

func (n NodeList) getParentType(level, index int) string {
	for i := index; i > 0; i-- {
		node := n.nodes[i]
		newLevel := node.Level
		if newLevel == level-1 && (n.nodes[i].value == "{" || n.nodes[i].value == "[") {
			return braces[n.nodes[i].value]
		}
	}
	return "ERROR"
}

func (n *NodeList) alterNodesVisibility(startIndex, endIndex int, visible bool) {
	for index := startIndex; index < endIndex+1; index++ {
		n.nodes[index].IsVisible = visible
	}
}

func (n NodeList) getNodeIndex(visibleIndex int) int {
	for index := n.topNode; index < len(n.nodes); index++ {
		if n.nodes[index].IsVisible {
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
		if n.nodes[index].IsVisible {
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