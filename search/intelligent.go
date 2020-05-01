package search

import (
	"regexp"
	"strings"
)

// Intelligent stuff
type Intelligent struct {
	commands     []Command
	commandIndex int
}

// NewIntelligent stuff
func NewIntelligent(input string) (Intelligent, error) {
	commands, err := Parse(input)
	if err != nil {
		return Intelligent{}, err
	}
	return Intelligent{commands, 0}, nil
}

// Execute stuff
func (i *Intelligent) Execute(nodeList sNodeList, fMode FunctionEnum) error {
	matchedNodes := i.executeCommands([]int{}, nodeList)
	nodeList.ApplyFilter(matchedNodes)
	return nil
}

// GetCommands returns commands (needed for testing)
func (i Intelligent) GetCommands() []Command {
	return i.commands
}

func (i *Intelligent) executeCommands(currentIndices []int, nodeList sNodeList) []int {
	var nodeIndices []int
	for index := i.commandIndex; index < len(i.commands); index++ {
		command := i.commands[index]
		list := command.RunConitional(currentIndices, nodeList)
		if command.HasOpenBracket() {
			i.commandIndex = index + 1
			list = i.executeCommands(nodeIndices, nodeList)
			index = i.commandIndex
		}
		if command.Operator != "" {
			nodeIndices = command.RunOperation(nodeIndices, list)
		} else {
			nodeIndices = append(nodeIndices, list...)
		}
		if command.HasCloseBracket() {
			i.commandIndex = index
			return nodeIndices
		}
	}

	return nodeIndices
}

// GetIntelligentHints returns hints for controls
func GetIntelligentHints(input string, cursorPos int) []string {
	startIndex, isControl := getInterestingSubstringStart(input[:cursorPos])
	interestring := strings.Trim(input[startIndex:cursorPos], " ")
	if isControl {
		return getControlHints(interestring)
	}
	return operators
}

func getInterestingSubstringStart(input string) (int, bool) {
	controlStartIndex := getControlSubstringStart(input)
	operatorStartIndex := getOperatorSubstringStart(input)
	if controlStartIndex < operatorStartIndex {
		return operatorStartIndex, false
	}
	return controlStartIndex, true
}

func getControlHints(regex string) []string {
	var hints []string
	var altHints []string
	if r, err := regexp.Compile("(?i)" + regex); err == nil {
		for _, control := range controls {
			if r.MatchString(control) {
				hints = append(hints, control)
			}
		}
	}
	return append(altHints, hints...)
}

func getControlSubstringStart(input string) int {
	highestIndex := strings.LastIndex(input, "(")
	for _, operator := range operators {
		if newIndex := strings.LastIndex(input, operator); newIndex > highestIndex {
			highestIndex = newIndex
		}
	}
	return highestIndex + 1
}

func getOperatorSubstringStart(input string) int {
	bracketIndex := strings.LastIndex(input, ")")
	quoteIndex := strings.LastIndex(input, "\"")
	if bracketIndex < quoteIndex && strings.Count(input, "\"")%2 == 0 {
		return quoteIndex + 1
	}
	return bracketIndex + 1
}
