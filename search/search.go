package search

import (
	"kube-review/nodelist"
	"regexp"
)

// Search stuff
type Search struct {
	queryMode    QueryEnum
	functionMode FunctionEnum
	ql           *QueryList
}

// NewSearch stuff
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
	qMode := s.queryMode
	if qMode == QUERY {
		input, qMode = s.ql.GetQuery(input)
	}

	if qMode == EXPRESSION {
		intelligent, err := NewExpression(input)
		if err != nil {
			return err
		}
		return intelligent.Execute(nodeList, s.functionMode)
	}
	return s.executeRegex(input, nodeList)
}

// ToggleQueryMode switches between regex and query mode
func (s *Search) ToggleQueryMode() {
	s.queryMode = (s.queryMode + 1) % 3
}

// ToggleSearchMode toggles through each search mode
func (s *Search) ToggleSearchMode() {
	s.functionMode = (s.functionMode + 1) % 2
}

// GetModeInfo stuff
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

func (s Search) executeRegex(regex string, nodeList sNodeList) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	matchNodes := nodeList.GetNodesMatching(r, nodelist.ANY, true)
	if s.functionMode == FILTER {
		return nodeList.Filter(matchNodes)
	} else if s.functionMode == FIND {
		nodeList.Highlight(matchNodes)
		if err := nodeList.FindNextHighlight(); err != nil {
			return err
		}
	}
	return nil
}
