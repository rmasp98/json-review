package search

import (
	"fmt"
	"regexp"

	"github.com/rmasp98/kube-review/nodelist"
)

// Search stuff
type Search struct {
	queryMode    QueryEnum
	functionMode FunctionEnum
	ql           *QueryList
}

// NewSearch DWISOTT
func NewSearch(queryMode QueryEnum, queryList *QueryList) Search {
	return Search{queryMode, FIND, queryList}
}

// GetHints stuff
func (s Search) GetHints(input string, cursorPos int) string {
	output := ""
	for _, hint := range s.getPossibleHints(input, cursorPos) {
		output += "\n" + hint
	}
	return output
}

// GetHints2 stuff
func (s Search) GetHints2(input string) []string {
	if input == "" {
		return []string{}
	}
	return s.getPossibleHints(input, len(input))
}

// InsertSelectedHint returns the search string with hint inserted and the position
// of the end of inserted hint
func (s Search) InsertSelectedHint(input string, cursorPos int, index int) (string, int) {
	if s.queryMode == EXPRESSION {
		hint := InsertSelectedExpressionHint(input[:cursorPos], index)
		return hint + input[cursorPos:], len(hint)
	} else if s.queryMode == QUERY {
		hint := s.ql.InsertHint(input, index)
		return hint, len(hint)
	}
	return input, cursorPos
}

// Execute runs a search based on input, QueryMode and searchMode
func (s Search) Execute(input string, nodeList sNodeList) error {
	nodeList.ResetView()
	matchedNodes, err := s.getMatchedNodes(input, nodeList)
	if err != nil {
		return err
	}
	if s.functionMode == FILTER {
		return nodeList.Filter(matchedNodes)
	} else if s.functionMode == FIND {
		nodeList.Highlight(matchedNodes)
		return nodeList.FindNextHighlight()
	}
	return fmt.Errorf("Invalid search type. Should be Filter or Find")
}

// ToggleQueryMode switches between regex and query mode
func (s *Search) ToggleQueryMode() {
	s.queryMode = (s.queryMode + 1) % 3
}

// ToggleSearchMode toggles through each search mode
func (s *Search) ToggleSearchMode() {
	s.functionMode = (s.functionMode + 1) % 2
}

// GetModeInfo returns the search and function type for UI title
func (s Search) GetModeInfo() string {
	return s.queryMode.String() + "-" + s.functionMode.String()
}

func (s Search) getPossibleHints(input string, cursorPos int) []string {
	if s.queryMode == QUERY {
		return s.ql.GetHints(input[:cursorPos])
	} else if s.queryMode == EXPRESSION {
		return GetExpressionHints(input[:cursorPos])
	}
	return []string{}
}

func (s Search) getMatchedNodes(input string, nodeList sNodeList) ([]int, error) {
	qMode := s.queryMode
	regex := input
	if qMode == QUERY {
		if regex, qMode = s.ql.GetQuery(input); regex == "" {
			return nil, fmt.Errorf("'%s' is not a valid query", input)
		}
	}

	if qMode == EXPRESSION {
		expression, err := NewExpression(regex)
		if err != nil {
			return nil, err
		}
		// use output to find/filter
		return expression.Execute(nodeList), nil
	}
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return nodeList.GetNodesMatching(r, nodelist.ANY, true), nil
}
