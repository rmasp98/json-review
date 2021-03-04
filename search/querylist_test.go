package search_test

import (
	"reflect"
	"testing"

	"github.com/rmasp98/kube-review/search"
)

func TestReturnsAnErrorIfFailedToLoadFile(t *testing.T) {
	ql := search.NewQueryList()
	err := ql.Load("nofile", "")
	if err == nil {
		t.Errorf("This was supposed to return an error")
	}
}

func TestGetNamesReturnsListInAlphabeticalOrder(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	expected := []string{"test1more", "test2less"}
	for i := 0; i < 100; i++ {
		actual := ql.GetNames()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expected %v but instead got %v", expected, actual)
			break
		}
	}
}

func TestCanGetQueryFromName(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	actual, _ := ql.GetQuery("test1more")
	if actual != "[a-z]{5}" {
		t.Errorf("Expected '[a-z]{5}' but instead got '%s'", actual)
	}
}

func TestCanGetQueryTypeFromName(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	_, actual := ql.GetQuery("test2less")
	if actual != search.EXPRESSION {
		t.Errorf("Expected 'Intelligent' but instead got '%v'", actual)
	}
}

func TestCanGetDescriptionFromName(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	actual := ql.GetDescription("test1more")
	if actual != "First test (REGEX)" {
		t.Errorf("Expected 'First test (REGEX)' but instead got '%s'", actual)
	}
}

func TestCanAddToList(t *testing.T) {
	ql := search.NewQueryList()
	ql.Add("test3", "TestRegex", "Third test", search.REGEX)
	actual, _ := ql.GetQuery("test3")
	if actual != "TestRegex" {
		t.Errorf("Expected 'TestRegex' but instead got '%s'", actual)
	}
}

func TestAddReturnsErrorForIncorrectQueryType(t *testing.T) {
	ql := search.NewQueryList()
	err := ql.Add("test3", "TestRegex", "Third test", search.QUERY)
	if err == nil {
		t.Errorf("Expected an error but instead got nothing")
	}
}

func TestCanRemoveFromTheList(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	ql.Remove("test1more")
	actual, _ := ql.GetQuery("test1more")
	if actual != "" {
		t.Errorf("Expected empty string but instead got '%s'", actual)
	}
}

func TestGetHintsReturnsStringListWithMatchingQueriesAndDescriptions(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	actual := ql.GetHints("test")
	expected := []string{"test1more - \033[1;31mFirst test (REGEX)\033[0m", "test2less - \033[1;31mSecond test (Expression)\033[0m"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSearchHelpReturnsNothingIfMatchesAQueryExactly(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	actual := ql.GetHints("test1more")
	if len(actual) != 0 {
		t.Errorf("Expected nothing but got '%v'", actual)
	}
}

func TestInsertHintReturnsNameOfSelectedHint(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json", "queryschema.json")
	actual := ql.InsertHint("test", 1)
	expected := "test2less"
	if actual != expected {
		t.Errorf("Expected '%s' but instead got '%s'", expected, actual)
	}
}
