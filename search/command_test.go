package search_test

import (
	"kube-review/mocks"
	"kube-review/nodelist"
	"kube-review/search"
	"reflect"
	"regexp"
	"testing"
)

func TestReturnsTrueIfHasOpenBracket(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "+", "(")
	actual := command.HasOpenBracket()
	if actual != true {
		t.Errorf("Expected true but got %t", actual)
	}
}

func TestReturnsTrueIfHasCloseBracket(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "+", ")")
	actual := command.HasCloseBracket()
	if actual != true {
		t.Errorf("Expected true but got %t", actual)
	}
}

func TestOperationReturnsConcatedArraysForNoOperation(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "", ")")
	expected := []int{1, 2, 3, 4, 5}
	actual := command.RunOperation([]int{1, 3, 4}, []int{2, 3, 5})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsConcatedArraysForPlus(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "+", ")")
	expected := []int{1, 2, 3, 4, 5}
	actual := command.RunOperation([]int{1, 3, 4}, []int{2, 3, 5})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsSubtractionForMinus(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "-", ")")
	expected := []int{1, 2}
	actual := command.RunOperation([]int{1, 2, 3}, []int{3, 4, 5})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsIntersectionForVirtBar(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "|", ")")
	expected := []int{2, 3}
	actual := command.RunOperation([]int{1, 2, 3}, []int{2, 3, 4})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsAllIfRightNotEmptyforAmpersand(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "&&", ")")
	expected := []int{1, 2, 3, 4, 5}
	actual := command.RunOperation([]int{1, 3, 4}, []int{2, 3, 5})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsNothingIfRightEmptyforAmpersand(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "&&", ")")
	expected := []int{}
	actual := command.RunOperation([]int{1, 2, 3}, []int{})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsNothingIfLeftEmptyforAmpersand(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "&&", ")")
	expected := []int{}
	actual := command.RunOperation([]int{}, []int{4, 5, 6})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsLeftIfRightNotEmptyforLeftArrow(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "<-", ")")
	expected := []int{1, 2, 3}
	actual := command.RunOperation([]int{1, 2, 3}, []int{4, 5, 6})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsNothingIfRightEmptyforLeftArrow(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "<-", ")")
	expected := []int{}
	actual := command.RunOperation([]int{1, 2, 3}, []int{})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsRightIfLeftNotEmptyforRightArrow(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "->", ")")
	expected := []int{4, 5, 6}
	actual := command.RunOperation([]int{1, 2, 3}, []int{4, 5, 6})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsNothingIfLeftEmptyforRightArrow(t *testing.T) {
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "output", "->", ")")
	expected := []int{}
	actual := command.RunOperation([]int{}, []int{4, 5, 6})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestFindFunctionCallsCorrectFunction(t *testing.T) {
	mock := mocks.NodeListMock{}
	command := search.NewCommand(search.CMDFINDNODES, map[string]string{}, "", "", "")
	command.RunFunction([]int{}, &mock)
	actual := mock.Calls[0]
	expected := "GetNodesMatching"
	if actual != expected {
		t.Errorf("Expected call to '%s' but got '%s'", expected, actual)
	}
}

func TestFindFunctionCallsParsesInputCorrectly(t *testing.T) {
	mock := mocks.NodeListMock{}
	input := map[string]string{"regex": "test", "matchType": "KEY", "equal": "false"}
	command := search.NewCommand(search.CMDFINDNODES, input, "", "", "")
	command.RunFunction([]int{}, &mock)
	actual := mock.Args[0]
	if actual[0].(*regexp.Regexp).String() != "test" || actual[1].(nodelist.MatchType) != nodelist.KEY || actual[2].(bool) != false {
		t.Errorf("Expected '[Test Key false]' but got '%v'", actual)
	}
}

func TestFindRelativeFunctionCallsCorrectFunction(t *testing.T) {
	mock := mocks.NodeListMock{}
	command := search.NewCommand(search.CMDFINDRELATIVE, map[string]string{}, "", "", "")
	command.RunFunction([]int{1}, &mock)
	actual := mock.Calls[0]
	expected := "GetRelativesMatching"
	if actual != expected {
		t.Errorf("Expected call to '%s' but got '%s'", expected, actual)
	}
}

func TestFindRelativeGetsCalledForEachElementofInput(t *testing.T) {
	mock := mocks.NodeListMock{}
	command := search.NewCommand(search.CMDFINDRELATIVE, map[string]string{}, "", "", "")
	command.RunFunction([]int{1, 2, 3, 4, 5}, &mock)
	actual := len(mock.Calls)
	expected := 5
	if actual != expected {
		t.Errorf("Expected call to '%d' but got '%d'", expected, actual)
	}
}

func TestFindRelativeCallParsesBaseInputCorrectly(t *testing.T) {
	mock := mocks.NodeListMock{}
	input := map[string]string{"regex": "test", "matchType": "Value", "equal": "false"}
	command := search.NewCommand(search.CMDFINDRELATIVE, input, "", "", "")
	command.RunFunction([]int{1}, &mock)
	actual := mock.Args[0]
	if actual[3].(*regexp.Regexp).String() != "test" || actual[4].(nodelist.MatchType) != nodelist.VALUE || actual[5].(bool) != false {
		t.Errorf("Expected '[Test Value false]' but got '%v'", actual)
	}
}

func TestFindRelativeCallParsesSearchLocationProperly(t *testing.T) {
	mock := mocks.NodeListMock{}
	input := map[string]string{"relativeStartLevel": "2", "depth": "5", "regex": "test"}
	command := search.NewCommand(search.CMDFINDRELATIVE, input, "", "", "")
	command.RunFunction([]int{1}, &mock)
	actual := mock.Args[0]
	if actual[0].(int) != 1 || actual[1].(int) != 2 || actual[2].(int) != 5 {
		t.Errorf("Expected '[1 2 5 test Any true]' but got '%v'", actual)
	}
}

func TestFindRelativeRunsWithDefaultvaluesIfNotProvided(t *testing.T) {
	mock := mocks.NodeListMock{}
	input := map[string]string{"regex": "test"}
	command := search.NewCommand(search.CMDFINDRELATIVE, input, "", "", "")
	command.RunFunction([]int{1}, &mock)
	actual := mock.Args[0]
	if actual[1].(int) != 0 || actual[2].(int) != 1 || actual[4].(nodelist.MatchType) != nodelist.ANY || actual[5].(bool) != true {
		t.Errorf("Expected '[1 0 1 test Any true]' but got '%v'", actual)
	}
}

func TestFindRelativeReturnsAnOrderedUnionOfIndices(t *testing.T) {
	mock := mocks.NodeListMock{}
	mock.Returns = [][]int{[]int{1, 3, 5, 6}, []int{2, 4, 6}}
	input := map[string]string{"regex": "test"}
	command := search.NewCommand(search.CMDFINDRELATIVE, input, "", "", "")
	expected := []int{1, 2, 3, 4, 5, 6}
	_, actual := command.RunFunction([]int{1, 2}, &mock)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}
