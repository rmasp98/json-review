package search

import (
	"kube-review/jsontree"
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
	if input != "" {
		for _, hint := range s.getPossibleHints(input, cursorPos) {
			output += "\n" + hint
			if s.queryMode == QUERY {
				output += " - " + redBold + s.ql.GetDescription(hint)
			}
			output += reset
		}
	}
	return output
}

// InsertSelectedHint returns the search string with hint inserted and the position
// of the end of inserted hint
func (s Search) InsertSelectedHint(input string, cursorPos int, index int) (string, int) {
	hint := s.getSelectedHint(input, cursorPos, index)
	replaceStart := 0
	if s.queryMode == INTELLIGENT {
		replaceStart, _ = getInterestingSubstringStart(input[:cursorPos])
	}
	space := " "
	if replaceStart == 0 || input[replaceStart-1] == '(' {
		space = ""
	}
	return input[:replaceStart] + space + hint + input[cursorPos:], len(input[:replaceStart] + space + hint)
}

// Execute runs a search based on input, QueryMode and searchMode
func (s Search) Execute(input string, nodeList sNodeList) error {
	qMode := s.queryMode
	if qMode == QUERY {
		input, qMode = s.ql.GetQuery(input)
	}

	if qMode == INTELLIGENT {
		intelligent, err := NewIntelligent(input)
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

func (s Search) getSelectedHint(input string, cursorPos int, index int) string {
	hints := s.getPossibleHints(input, cursorPos)
	if index < len(hints) && index >= 0 {
		return hints[index]
	}
	return ""
}

func (s Search) getPossibleHints(input string, cursorPos int) []string {
	if s.queryMode == QUERY {
		var queries []string
		if input != "" {
			r, _ := regexp.Compile(input)
			for _, query := range s.ql.GetNames() {
				if r.MatchString(query) {
					queries = append(queries, query)
				}
				if input == query {
					return []string{}
				}
			}
		}
		return queries
	} else if s.queryMode == INTELLIGENT {
		return GetIntelligentHints(input, cursorPos)
	}
	return []string{}
}

func (s Search) executeRegex(regex string, nodeList sNodeList) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	matchNodes := nodeList.GetNodesMatching(r, jsontree.ANY, true)
	if s.functionMode == FILTER {
		return nodeList.ApplyFilter(matchNodes)
	} else if s.functionMode == FIND {
		if err := nodeList.ApplyHighlight(matchNodes); err != nil {
			return err
		}
		nodeList.FindNextHighlightedNode()
	}
	return nil
}
