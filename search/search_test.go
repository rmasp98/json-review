package search_test

import (
	"kube-review/search"
	"testing"
)

var qlFile = "../testdata/querylist-test.json"

func TestLoadsQueryList(t *testing.T) {
	search := search.NewSearch(search.QUERY, qlFile)
	actual := search.GetQueryName("test", 0)
	if actual != "test1more" {
		t.Errorf("Expected 'test1more' but got '%s'", actual)
	}
}

func TestGetQueryReturnsNothingIfIndexNotExist(t *testing.T) {
	search := search.NewSearch(search.QUERY, qlFile)
	actual := search.GetQueryName("test", 2)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestGetQueryReturnsNothingIfNotQueryMode(t *testing.T) {
	search := search.NewSearch(search.REGEX, qlFile)
	actual := search.GetQueryName("test", 0)
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestSearchHelpReturnsStringWithMatchingQueriesAndDescriptions(t *testing.T) {
	search := search.NewSearch(search.QUERY, qlFile)
	actual := search.GetHints("test")
	expected := "\n\033[1;97mtest1more - \033[1;31mFirst test\033[0m\n\033[1;97mtest2less - \033[1;31mSecond test\033[0m"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestSearchHelpReturnsNothingIfNotQueryMode(t *testing.T) {
	search := search.NewSearch(search.REGEX, qlFile)
	actual := search.GetHints("test")
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestSearchHelpReturnsNothingIfInputEmpty(t *testing.T) {
	search := search.NewSearch(search.QUERY, qlFile)
	actual := search.GetHints("")
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

func TestSearchHelpReturnsNothingIfMatchesAQueryExactly(t *testing.T) {
	search := search.NewSearch(search.QUERY, qlFile)
	actual := search.GetHints("test1more")
	if actual != "" {
		t.Errorf("Expected nothing but got '%s'", actual)
	}
}

var queryInfo = []string{"Regex", "Query", "Intelligent"}
var methodInfo = []string{"Find", "Filter"}

func TestGetModeInfoReturnsCorrect(t *testing.T) {
	search := search.NewSearch(search.REGEX, qlFile)
	for _, query := range queryInfo {
		for _, method := range methodInfo {
			actual := search.GetModeInfo()
			var expected string
			if query == "Intelligent" {
				expected = query
			} else {
				expected = query + "-" + method
			}
			if actual != expected {
				t.Errorf("Expected '%s' but got '%s'", expected, actual)
			}
			search.ToggleSearchMode()
		}
		search.ToggleQueryMode()
	}
}
