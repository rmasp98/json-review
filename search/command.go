package search

import (
	"kube-review/nodelist"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ForValues(nodes, searchtype, outnodes)
// get values of nodes and perform a find on those

// Command stuff
type Command struct {
	function CmdFunc
	input    map[string]string
	output   string
	operator string
	bracket  string
}

// NewCommand stuff
func NewCommand(function CmdFunc, input map[string]string, output string, operator string, bracket string) Command {
	return Command{function, input, output, operator, bracket}
}

// RunFunction stuff
func (c Command) RunFunction(input []int, nodeList sNodeList) (string, []int) {
	if r, matchType, equal, err := c.processBaseInputs(); err == nil {
		if c.function == CMDFINDNODES {
			return c.output, nodeList.GetNodesMatching(r, matchType, equal)
		} else if c.function == CMDFINDRELATIVE {
			var list []int
			relativeStartLevel, depth := c.processRelativeInputs()
			for _, index := range input {
				list = append(list, nodeList.GetRelativesMatching(index, relativeStartLevel, depth, r, matchType, equal)...)
			}
			return c.output, orderedUnion(list, []int{})
		}
	}
	return "", []int{}
}

// RunOperation stuff
func (c Command) RunOperation(left, right []int) []int {
	switch c.operator {
	case "":
		return orderedUnion(left, right)
	case "+":
		return orderedUnion(left, right)
	case "-":
		return subtract(left, right)
	case "|":
		return intersection(left, right)
	case "&&":
		if len(left) > 0 && len(right) > 0 {
			return orderedUnion(left, right)
		}
	case "<-":
		if len(right) > 0 {
			return left
		}
	case "->":
		if len(left) > 0 {
			return right
		}
	}
	return []int{}
}

// HasOpenBracket stuff
func (c Command) HasOpenBracket() bool {
	return c.bracket == "("
}

// HasCloseBracket stuff
func (c Command) HasCloseBracket() bool {
	return c.bracket == ")"
}

// GetInputName returns the name of expected input that should be output by another function call
func (c Command) GetInputName() string {
	return c.input["nodes"]
}

func (c Command) processBaseInputs() (*regexp.Regexp, nodelist.MatchType, bool, error) {
	matchType := getMatchType(c.input["matchType"])
	var equal = false
	if c.input["equal"] == "" || strings.EqualFold(c.input["equal"], "true") {
		equal = true
	}
	r, err := regexp.Compile(c.input["regex"])
	return r, matchType, equal, err
}

func (c Command) processRelativeInputs() (int, int) {
	relativeStartLevel, _ := strconv.Atoi(c.input["relativeStart"])
	var depth = 1
	if temp, ok := c.input["depth"]; ok {
		depth, _ = strconv.Atoi(temp)
	}
	return relativeStartLevel, depth
}

func getMatchType(input string) nodelist.MatchType {
	if strings.EqualFold(input, nodelist.KEY.String()) {
		return nodelist.KEY
	} else if strings.EqualFold(input, nodelist.VALUE.String()) {
		return nodelist.VALUE
	}
	return nodelist.ANY
}

func subtract(left, right []int) []int {
	var result []int
	for _, elemLeft := range left {
		matched := false
		for _, elemRight := range right {
			if elemLeft == elemRight {
				matched = true
			}
		}
		if !matched {
			result = append(result, elemLeft)
		}
	}
	return result
}

func intersection(left, right []int) []int {
	var result []int
	for _, elemLeft := range left {
		for _, elemRight := range right {
			if elemLeft == elemRight {
				result = append(result, elemLeft)
			}
		}
	}
	return result
}

// TODO: this is fairly central so find more efficient way
func orderedUnion(left, right []int) []int {
	unique := make(map[int]struct{}, len(left)+len(right))
	for _, index := range left {
		unique[index] = struct{}{}
	}
	for _, index := range right {
		unique[index] = struct{}{}
	}
	result := make([]int, 0, len(unique))
	for index := range unique {
		result = append(result, index)
	}
	sort.Ints(result)
	return result
}
