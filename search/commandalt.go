package search

import (
	"kube-review/jsontree"
	"regexp"
)

// Find(searchType, "regex", outnodes)
// searchType = any/key/value

// FindRelative(nodes, relation, levels, searchType, outnodes)
// relation = parent/child/sibling
// levels decides number of levels if recursive (siblings will search children)
// TODO: find way to search everything back to main parent

// ForValues(nodes, searchtype, outnodes)
// get values of nodes and perform a find on those

type CommandAlt struct {
	function string
	input    map[string]string
	output   string
	operator string
	bracket  string
}

func (c CommandAlt) RunFunction(input []int, nodeList sNodeList) (string, []int) {
	if c.function == "Find" {
		searchType := c.input["searchType"]
		regex := c.input["regex"]
		r, _ := regexp.Compile(regex)
		result := nodeList.GetNodesMatching(r, getSearchType(searchType), true)
		return c.output, result
	}
	return "", []int{}
}

func (c CommandAlt) RunOperation(left, right []int) []int {
	//do other operations but default to this one
	return append(left, right...)
}

func (c CommandAlt) HasOpenBracket() bool {
	return false
}

func (c CommandAlt) HasCloseBracket() bool {
	return false
}

type IntelligentAlt struct {
	variables map[string][]int
	commands  []CommandAlt
	cmdIndex  int
}

func (i IntelligentAlt) executeCommands(nodeList sNodeList) []int {
	var currentIndices []int
	for cmdIndex := i.cmdIndex; cmdIndex < len(i.commands); cmdIndex++ {
		command := i.commands[cmdIndex]
		input := i.variables[command.input["input"]]
		variable, output := command.RunFunction(input, nodeList)
		if variable != "" {
			i.variables[variable] = append(i.variables[variable], output...)
		}
		if command.HasOpenBracket() {
			i.cmdIndex = cmdIndex + 1
			output = i.executeCommands(nodeList)
			cmdIndex = i.cmdIndex
		}
		currentIndices = command.RunOperation(currentIndices, output)
		if command.HasCloseBracket() {
			i.cmdIndex = cmdIndex
			return currentIndices
		}
	}
	return currentIndices
}

func getSearchType(input string) jsontree.MatchType {
	return jsontree.ANY
}
