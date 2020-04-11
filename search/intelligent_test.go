package search_test

import (
	"kube-review/mocks"
	"kube-review/search"
	"reflect"
	"testing"
)

//tests
// Should not start with anything other than Any, Key, Value
//can return possible keywords
//correctly proccess input

func TestValidateBasicConditionReturnsNoError(t *testing.T) {
	_, actual := search.NewIntelligent("Any==\"Test\"")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestValidateReturnsErrorIfControlInvalid(t *testing.T) {
	_, actual := search.NewIntelligent("Nothing==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidateReturnsErrorIfConditionInvalid(t *testing.T) {
	_, actual := search.NewIntelligent("Any=\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidateReturnsErrorIfQuotesInvalid(t *testing.T) {
	_, actual := search.NewIntelligent("Any==\"Test")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidateReturnsNoErrorForBasicBrackets(t *testing.T) {
	_, actual := search.NewIntelligent("(Any==\"Test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestValidateReturnsErrorIfNotMatchingBrackets(t *testing.T) {
	_, actual := search.NewIntelligent("(Any==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidateReturnsErrorIfEndsInOperator(t *testing.T) {
	_, actual := search.NewIntelligent("Any==\"Test\"&&")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidateReturnsErrorIfNoOperator(t *testing.T) {
	_, actual := search.NewIntelligent("Any==\"Test\"Any==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidatePicksUpErrorsAfterFirstCondition(t *testing.T) {
	_, actual := search.NewIntelligent("Any==\"Test\"&&Nothing==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestValidateCanHandleMultipleBrackets(t *testing.T) {
	_, actual := search.NewIntelligent("((Any==\"Test\")+Any==\"Test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestValidateCanHandleSpaces(t *testing.T) {
	_, actual := search.NewIntelligent(" ( ( Any == \"Test\" ) - Any == \"Test\" ) ")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestBasicCommandCorrectlyParsed(t *testing.T) {
	intelligent, _ := search.NewIntelligent("Any==\"Test\"")
	expected := search.Command{"Any", true, "Test", "", ""}
	actual := intelligent.GetCommands()[0]
	if actual != expected {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexCommandParsesCorrectly(t *testing.T) {
	intelligent, _ := search.NewIntelligent("((Any==\"Test\")+Key!=\"Test2\")")
	expected := []search.Command{
		search.Command{"", false, "", "", "("},
		search.Command{"", false, "", "", "("},
		search.Command{"Any", true, "Test", "", ")"},
		search.Command{"Key", false, "Test2", "+", ")"},
	}
	actual := intelligent.GetCommands()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestBasicCommandExecutesCorrectFunctions(t *testing.T) {
	intelligent, _ := search.NewIntelligent("Any==\"Test\"")
	mock := mocks.NodeListMock{}
	intelligent.Execute(&mock)
	expected := []string{"GetNodesMatching", "ApplyFilter"}
	actual := mock.Calls
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexCommandExecutesCorrectFunctions(t *testing.T) {
	intelligent, _ := search.NewIntelligent("Any==\"Test\" + (ParentHasChildKey==\"Test2\")")
	mock := mocks.NodeListMock{}
	intelligent.Execute(&mock)
	expected := []string{"GetNodesMatching", "GetParentChildrenMatching", "GetParentChildrenMatching", "GetParentChildrenMatching", "ApplyFilter"}
	actual := mock.Calls
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexCommandReturnsCorrectOutput(t *testing.T) {
	intelligent, _ := search.NewIntelligent("Any==\"Test\" + (ParentHasChildKey==\"Test2\")")
	mock := mocks.NodeListMock{}
	intelligent.Execute(&mock)
	expected := []int{1, 2, 5, 2, 7, 2, 7, 2, 7}
	actual := mock.Args[len(mock.Args)-1][0]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSubtractCommandReturnsCorrectOutput(t *testing.T) {
	intelligent, _ := search.NewIntelligent("Any==\"Test\" - Any==\"Test2\"")
	mock := mocks.NodeListMock{}
	intelligent.Execute(&mock)
	actual := mock.Args[2][0].([]int)
	if len(actual) != 0 {
		t.Errorf("Expected empty array but got '%v'", actual)
	}
}

func TestBracketsCommandReturnsCorrectOutput(t *testing.T) {
	intelligent, _ := search.NewIntelligent("Any==\"Test\" - (ParentHasChildKey==\"Test2\" + ChildHasValue==\"Test3\") + Key!=\"Test4\"")
	mock := mocks.NodeListMock{}
	intelligent.Execute(&mock)
	expected := []int{5, 1, 2, 5}
	actual := mock.Args[len(mock.Args)-1][0]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}
