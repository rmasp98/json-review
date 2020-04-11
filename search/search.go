package search

import (
	"kube-review/jsontree"
	"regexp"
)

// Search stuff
type Search struct {
	queryMode    QueryEnum
	functionMode FunctionEnum
	ql           QueryList
}

// NewSearch stuff
func NewSearch(queryMode QueryEnum, qlFile string) Search {
	ql := NewQueryList()
	ql.Load(qlFile)
	return Search{queryMode, FIND, ql}
}

// GetHints stuff
func (s Search) GetHints(input string) string {
	output := ""
	if s.queryMode == QUERY && input != "" {
		for _, query := range s.getPossibleQueries(input) {
			if input == query {
				return ""
			}
			output += "\n" + whiteBold + query + " - " + redBold + s.ql.GetDescription(query) + reset
		}
	}
	return output
}

// Execute runs a search based on input, QueryMode and searchMode
func (s Search) Execute(input string, nodeList sNodeList) error {
	regex := input
	if s.queryMode == QUERY {
		regex = s.ql.GetRegex(input)
	} else if s.queryMode == INTELLIGENT {
		intelligent, err := NewIntelligent(input)
		if err != nil {
			return err
		}
		intelligent.Execute(nodeList)
		return nil
	}
	switch s.functionMode {
	case FIND:
		matchNodes := nodeList.GetNodesMatching(regex, jsontree.ANY, true)
		nodeList.ApplyHighlight(matchNodes)
		nodeList.FindNextHighlightedNode()
	case FILTER:
		matchNodes := nodeList.GetNodesMatching(regex, jsontree.ANY, true)
		nodeList.ApplyFilter(matchNodes)
	}
	return nil
}

// GetQueryName get query of index in list of matched querys to user input
func (s Search) GetQueryName(input string, index int) string {
	if s.queryMode == QUERY {
		queries := s.getPossibleQueries(input)
		if index < len(queries) && index >= 0 {
			return queries[index]
		}
	}
	return ""
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
	var mode = s.queryMode.String()
	if s.queryMode != INTELLIGENT {
		mode += "-" + s.functionMode.String()
	}
	return mode
}

func (s Search) getPossibleQueries(input string) []string {
	var queries []string
	r, _ := regexp.Compile(input)
	for _, query := range s.ql.GetNames() {
		if r.MatchString(query) {
			queries = append(queries, query)
		}
	}
	return queries
}
