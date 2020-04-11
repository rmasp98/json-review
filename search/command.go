package search

import (
	"kube-review/jsontree"
	"regexp"
)

// Command stuff
type Command struct {
	Control  string
	Equal    bool
	Regex    string
	Operator string
	Bracket  string
}

// RunConitional stuff
func (c Command) RunConitional(indices []int, nodeList sNodeList) []int {
	var outIndices []int
	if matched, _ := regexp.MatchString("ParentHasChild", c.Control); matched {
		for _, index := range indices {
			outIndices = append(outIndices, nodeList.GetParentChildrenMatching(index, c.Regex, getMatchType(c.Control), c.Equal, getRecursion(c.Control))...)
		}
	} else if matched, _ := regexp.MatchString("ChildHas", c.Control); matched {
		for _, index := range indices {
			outIndices = append(outIndices, nodeList.GetChildrenMatching(index, c.Regex, getMatchType(c.Control), c.Equal, getRecursion(c.Control))...)
		}
	} else if c.Control != "" {
		outIndices = nodeList.GetNodesMatching(c.Regex, getMatchType(c.Control), c.Equal)
	}
	return outIndices
}

// RunOperation stuff
func (c Command) RunOperation(left, right []int) []int {
	switch c.Operator {
	case "+":
		return append(left, right...)
	case "-":
		return subtract(left, right)
	case "&&":
		if len(left) > 0 && len(right) > 0 {
			return append(left, right...)
		}
	}
	return []int{}
}

// HasOpenBracket stuff
func (c Command) HasOpenBracket() bool {
	return c.Bracket == "("
}

// HasCloseBracket stuff
func (c Command) HasCloseBracket() bool {
	return c.Bracket == ")"
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

func getMatchType(control string) jsontree.MatchType {
	for i := 0; i < 3; i++ {
		if matched, _ := regexp.MatchString(jsontree.MatchType(i).String()+"$", control); matched {
			return jsontree.MatchType(i)
		}
	}
	return jsontree.ANY
}

func getRecursion(control string) bool {
	matched, _ := regexp.MatchString("^Any", control)
	return matched
}
