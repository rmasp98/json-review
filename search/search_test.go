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

//"../testdata/querylist-test.json"

func TestLoadsQueryList(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	actual, _ := s.InsertSelectedHint("test", 4, 0)
	if actual != "test1more" {
		t.Errorf("Expected 'test1more' but got '%s'", actual)
	}
}

func TestGetQueryReturnsNothingIfIndexNotExist(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	actual, _ := s.InsertSelectedHint("test", 4, 2)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestGetQueryReturnsNothingIfNotQueryMode(t *testing.T) {
	s := search.NewSearch(search.REGEX, getQueryList())
	actual, _ := s.InsertSelectedHint("test", 4, 0)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestSearchHelpReturnsStringWithMatchingQueriesAndDescriptions(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	actual := s.GetHints("test", 3)
	expected := "\ntest1more - \033[1;31mFirst test (REGEX)\033[0m\ntest2less - \033[1;31mSecond test (INTELLIGENT)\033[0m"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestSearchHelpReturnsNothingIfNotQueryMode(t *testing.T) {
	s := search.NewSearch(search.REGEX, getQueryList())
	actual := s.GetHints("test", 3)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestSearchHelpReturnsNothingIfInputEmpty(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	actual := s.GetHints("", 0)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestSearchHelpReturnsNothingIfMatchesAQueryExactly(t *testing.T) {
	s := search.NewSearch(search.QUERY, getQueryList())
	actual := s.GetHints("test1more", 8)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

var queryInfo = []string{"Regex", "Intelligent", "Query"}
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

func TestInsetHintReplacesBasicSearchInIntelligentMode(t *testing.T) {
	s := search.NewSearch(search.INTELLIGENT, getQueryList())
	actual, _ := s.InsertSelectedHint("An", 2, 1)
	expected := "HasAnyParent"
	if actual != expected {
		t.Errorf("Expecting '%s' but got '%s'", expected, actual)
	}
}

func TestInsetHintLeavesAnyStringAfterCursor(t *testing.T) {
	s := search.NewSearch(search.INTELLIGENT, getQueryList())
	actual, _ := s.InsertSelectedHint("An==\"test\"", 2, 1)
	expected := "HasAnyParent==\"test\""
	if actual != expected {
		t.Errorf("Expecting '%s' but got '%s'", expected, actual)
	}
}

func TestInsetHintLeavesAnyStringBeforeSubstringOfInterest(t *testing.T) {
	s := search.NewSearch(search.INTELLIGENT, getQueryList())
	actual, _ := s.InsertSelectedHint("Any==\"test\" + An==\"test\"", 16, 1)
	expected := "Any==\"test\" + HasAnyParent==\"test\""
	if actual != expected {
		t.Errorf("Expecting '%s' but got '%s'", expected, actual)
	}
}

func TestInsertHintReturnsNewCursorPositionAfterHint(t *testing.T) {
	s := search.NewSearch(search.INTELLIGENT, getQueryList())
	_, actual := s.InsertSelectedHint("Any==\"test\" + An==\"test\"", 16, 1)
	if actual != 26 {
		t.Errorf("Expecting 26 but got '%d'", actual)
	}
}
