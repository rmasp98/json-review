package search_test

import (
	"kube-review/jsontree"
	"kube-review/mocks"
	"kube-review/search"
	"reflect"
	"testing"
)

func TestReturnsTrueIfHasOpenBracket(t *testing.T) {
	command := search.Command{"", true, "", "", "("}
	actual := command.HasOpenBracket()
	if actual != true {
		t.Errorf("Expected true but got %t", actual)
	}
}

func TestReturnsTrueIfHasCloseBracket(t *testing.T) {
	command := search.Command{"", true, "", "", ")"}
	actual := command.HasCloseBracket()
	if actual != true {
		t.Errorf("Expected true but got %t", actual)
	}
}

func TestOperationReturnsExpectedForAmpersand(t *testing.T) {
	command := search.Command{"", true, "", "&&", ""}
	actual := command.RunOperation([]int{1, 2}, []int{3, 4})
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsExpectedForAmpersandIfNoneInInArray(t *testing.T) {
	command := search.Command{"", true, "", "&&", ""}
	actual := command.RunOperation([]int{1, 2}, []int{})
	expected := []int{}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsExpectedForPlus(t *testing.T) {
	command := search.Command{"", true, "", "+", ""}
	actual := command.RunOperation([]int{1, 2}, []int{3, 4})
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestOperationReturnsExpectedForMinus(t *testing.T) {
	command := search.Command{"", true, "", "-", ""}
	actual := command.RunOperation([]int{1, 2, 3, 4}, []int{2, 4})
	expected := []int{1, 3}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

type functionData struct {
	name       string
	function   string
	regex      string
	regexIndex int
	matchType  jsontree.MatchType
	mtIndex    int
	condition  bool
	condIndex  int
	all        bool
}

var controlFunctions = []functionData{
	functionData{"Any", "GetNodesMatching", "test", 0, jsontree.ANY, 1, false, 2, false},
	functionData{"Key", "GetNodesMatching", "test", 0, jsontree.KEY, 1, false, 2, false},
	functionData{"Value", "GetNodesMatching", "test", 0, jsontree.VALUE, 1, false, 2, false},
	functionData{"ParentHasChildKey", "GetParentChildrenMatching", "test", 1, jsontree.KEY, 2, false, 3, false},
	functionData{"ParentHasChildValue", "GetParentChildrenMatching", "test", 1, jsontree.VALUE, 2, false, 3, false},
	functionData{"ParentHasChildAny", "GetParentChildrenMatching", "test", 1, jsontree.ANY, 2, false, 3, false},
	functionData{"ChildHasKey", "GetChildrenMatching", "test", 1, jsontree.KEY, 2, false, 3, false},
	functionData{"ChildHasValue", "GetChildrenMatching", "test", 1, jsontree.VALUE, 2, false, 3, false},
	functionData{"ChildHasAny", "GetChildrenMatching", "test", 1, jsontree.ANY, 2, false, 3, false},
	functionData{"AnyParentHasChildKey", "GetParentChildrenMatching", "test", 1, jsontree.KEY, 2, false, 3, true},
	functionData{"AnyParentHasChildValue", "GetParentChildrenMatching", "test", 1, jsontree.VALUE, 2, false, 3, true},
	functionData{"AnyParentHasChildAny", "GetParentChildrenMatching", "test", 1, jsontree.ANY, 2, false, 3, true},
	functionData{"AnyChildHasKey", "GetChildrenMatching", "test", 1, jsontree.KEY, 2, false, 3, true},
	functionData{"AnyChildHasValue", "GetChildrenMatching", "test", 1, jsontree.VALUE, 2, false, 3, true},
	functionData{"AnyChildHasAny", "GetChildrenMatching", "test", 1, jsontree.ANY, 2, false, 3, true},
	// TODO: Implement these
	// "HasParent":              "",
	// "HasAnyParent":           "",
}

func TestConditionalRunsAllCommandsOnNodes(t *testing.T) {
	for _, fd := range controlFunctions {
		mock := mocks.NodeListMock{}
		cmd := search.Command{fd.name, true, fd.regex, "", ""}
		cmd.RunConitional([]int{1}, &mock)
		if mock.Calls[0] != fd.function {
			t.Errorf("Expected call for '%s' but got none", fd.name)
		}
	}
}

func TestConditionalRunsWithRegexInCommand(t *testing.T) {
	for _, fd := range controlFunctions {
		mock := mocks.NodeListMock{}
		cmd := search.Command{fd.name, true, fd.regex, "", ""}
		cmd.RunConitional([]int{1}, &mock)
		if mock.Args[0][fd.regexIndex] != fd.regex {
			t.Errorf("Expected '%s' but got '%s'", fd.regex, mock.Args[0][fd.regexIndex])
		}
	}
}

func TestConditionalPassesConditionToFunction(t *testing.T) {
	for _, fd := range controlFunctions {
		mock := mocks.NodeListMock{}
		cmd := search.Command{fd.name, fd.condition, fd.regex, "", ""}
		cmd.RunConitional([]int{1}, &mock)
		if mock.Args[0][fd.condIndex] != fd.condition {
			t.Errorf("Expected '%t' for '%s' but got '%t'", fd.condition, fd.name, mock.Args[0][fd.condIndex])
		}
	}
}

func TestConditionalRunsWithMatchTypeInCommand(t *testing.T) {
	for _, fd := range controlFunctions {
		mock := mocks.NodeListMock{}
		cmd := search.Command{fd.name, true, fd.regex, "", ""}
		cmd.RunConitional([]int{1}, &mock)
		if mock.Args[0][fd.mtIndex] != fd.matchType {
			t.Errorf("Expected '%s' but got '%s'", fd.matchType, mock.Args[0][fd.mtIndex])
		}
	}
}

func TestConditionalPassesArrayToDependantFunctions(t *testing.T) {
	for _, fd := range controlFunctions[3:] {
		mock := mocks.NodeListMock{}
		cmd := search.Command{fd.name, true, fd.regex, "", ""}
		cmd.RunConitional([]int{5, 10}, &mock)
		if mock.Args[0][0] != 5 || mock.Args[1][0] != 10 {
			t.Errorf("Expected '5' and '10' but got '%d' and '%d'", mock.Args[0][0], mock.Args[1][0])
		}
	}
}

func TestAnyConditionalsPassesTrueToDependantFunctions(t *testing.T) {
	for _, fd := range controlFunctions[3:] {
		mock := mocks.NodeListMock{}
		cmd := search.Command{fd.name, true, fd.regex, "", ""}
		cmd.RunConitional([]int{1}, &mock)
		if mock.Args[0][4] != fd.all {
			t.Errorf("Expected '%t' for '%s'  but got '%t'", fd.all, fd.name, mock.Args[0][4])
		}
	}
}

func TestBasicCommandReturnsExpectedNodeIndices(t *testing.T) {
	mock := mocks.NodeListMock{}
	cmd := search.Command{"Any", true, "Test", "", ""}
	actual := cmd.RunConitional([]int{}, &mock)
	expected := []int{1, 2, 5}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexCommandReturnsExpectedNodeIndices(t *testing.T) {
	mock := mocks.NodeListMock{}
	cmd := search.Command{"AnyChildHasKey", false, "Test", "", ""}
	actual := cmd.RunConitional([]int{1, 2}, &mock)
	expected := []int{1, 6, 1, 6}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}
