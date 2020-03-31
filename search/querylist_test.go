package search_test

import (
	"kube-review/search"
	"reflect"
	"testing"
)

func TestReturnsAnErrorIfFailedToLoadFile(t *testing.T) {
	ql := search.NewQueryList()
	err := ql.Load("nofile")
	if err == nil {
		t.Errorf("This was supposed to return an error")
	}
}

func TestSuccessfullyLoadsFile(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json")
	actual := ql.GetNames()
	expected := []string{"test1", "test2"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v but instead got %v", expected, actual)
	}
}

func TestCanGetRegexFromName(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json")
	actual := ql.GetRegex("test1")
	if actual != "[a-z]{5}" {
		t.Errorf("Expected '[a-z]{5}' but instead got '%s'", actual)
	}
}

func TestCanGetDescriptionFromName(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json")
	actual := ql.GetDescription("test1")
	if actual != "First test" {
		t.Errorf("Expected 'First test' but instead got '%s'", actual)
	}
}

func TestCanAddToList(t *testing.T) {
	ql := search.NewQueryList()
	ql.Add("test3", "TestRegex", "Third test")
	actual := ql.GetRegex("test3")
	if actual != "TestRegex" {
		t.Errorf("Expected 'TestRegex' but instead got '%s'", actual)
	}
}

func TestCanRemoveFromTheList(t *testing.T) {
	ql := search.NewQueryList()
	ql.Load("../testdata/querylist-test.json")
	ql.Remove("test1")
	actual := ql.GetRegex("test1")
	if actual != "" {
		t.Errorf("Expected empty string but instead got '%s'", actual)
	}
}
