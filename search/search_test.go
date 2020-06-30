package search_test

import (
	"kube-review/search"
	"testing"
)

func getQueryList() *search.QueryList {
	var ql = search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	return &ql
}

// tests
// TODO: test the execute paths

func TestGetHintsReturnsNothingInRegexMode(t *testing.T) {
	s := search.NewSearch(search.REGEX, getQueryList())
	actual := s.GetHints("test", 4)
	if actual != "" {
		t.Errorf("Expected empty string but got '%s'", actual)
	}
}

func TestGetHintsReturnsStringOfHints(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	actual := s.GetHints("test", 4)
	expected := "\ntest1more - \033[1;31mFirst test (REGEX)\033[0m\ntest2less - \033[1;31mSecond test (Expression)\033[0m"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestInsertReturnsCorrectResponseForQuery(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	input, cursorPos := s.InsertSelectedHint("test", 4, 1)
	if input != "test2less" && cursorPos != 9 {
		t.Errorf("Expected 'test2less' and 9 but got '%s' and %d", input, cursorPos)
	}
}

func TestInsertReturnsCorrectResponseForExpression(t *testing.T) {
	s := search.NewSearch(search.EXPRESSION, getQueryList())
	input, cursorPos := s.InsertSelectedHint("Find", 4, 0)
	if input != "FindNodes(" && cursorPos != 10 {
		t.Errorf("Expected 'FindNodes(' and 10 but got '%s' and %d", input, cursorPos)
	}
}

func TestInsertReturnsInputAfterCursorPositionForExpression(t *testing.T) {
	s := search.NewSearch(search.EXPRESSION, getQueryList())
	input, cursorPos := s.InsertSelectedHint("Find + FindNodes(\"test\")", 4, 0)
	if input != "FindNodes( + FindNodes(\"test\")" && cursorPos != 10 {
		t.Errorf("Expected 'FindNodes(' and 10 but got '%s' and %d", input, cursorPos)
	}
}

var queryInfo = []string{"Regex", "Expression", "Query"}
var methodInfo = []string{"Find", "Filter"}

func TestGetModeInfoReturnsCorrect(t *testing.T) {
	s := search.NewSearch(search.REGEX, getQueryList())
	for _, query := range queryInfo {
		for _, method := range methodInfo {
			actual := s.GetModeInfo()
			expected := query + "-" + method
			if actual != expected {
				t.Errorf("Expected '%s' but got '%s'", expected, actual)
			}
			s.ToggleSearchMode()
		}
		s.ToggleQueryMode()
	}
}
