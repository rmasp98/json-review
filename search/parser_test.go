package search_test

import (
	"kube-review/search"
	"reflect"
	"testing"
)

func TestParseBasicConditionReturnsNoError(t *testing.T) {
	_, actual := search.Parse("Any==\"Test\"")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseReturnsErrorIfControlInvalid(t *testing.T) {
	_, actual := search.Parse("Nothing==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseReturnsErrorIfConditionInvalid(t *testing.T) {
	_, actual := search.Parse("Any=\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseReturnsErrorIfQuotesInvalid(t *testing.T) {
	_, actual := search.Parse("Any==\"Test")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseReturnsNoErrorForBasicBrackets(t *testing.T) {
	_, actual := search.Parse("(Key==\"Test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseReturnsErrorIfNotMatchingBrackets(t *testing.T) {
	_, actual := search.Parse("(Any==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseReturnsErrorIfEndsInOperator(t *testing.T) {
	_, actual := search.Parse("Any==\"Test\"&&")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseReturnsErrorIfNoOperator(t *testing.T) {
	_, actual := search.Parse("Any==\"Test\"Any==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParsePicksUpErrorsAfterFirstCondition(t *testing.T) {
	_, actual := search.Parse("Any==\"Test\"&&Nothing==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseCanHandleMultipleBrackets(t *testing.T) {
	_, actual := search.Parse("Value==\"test\"+((AnyChildHasValue==\"Test\")+HasParent==\"Test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseCanHandleSpaces(t *testing.T) {
	_, actual := search.Parse(" Key == \"test\" + ( ( ParentHasChildAny == \"Test\" ) - ChildHasKey == \"Test\" ) ")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseDetectsParentChildCommandsAtBeginningOfInput(t *testing.T) {
	_, actual := search.Parse("ParentHasChildKey==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseAllowsParentChildCommandsWithinBracketsAfterMainCommands(t *testing.T) {
	_, actual := search.Parse("Any == \"Test\" + (AnyParentHasChildKey==\"Test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseDetectsParentChildCommandsNotInBracketsFollowingMainCommands(t *testing.T) {
	_, actual := search.Parse("Any == \"Test\" + ParentHasChildKey==\"Test\"")
	if actual == nil {
		t.Errorf("Expected Error but got no error")
	}
}

func TestParseAllowsCaseInsensitiveControls(t *testing.T) {
	_, actual := search.Parse(" keY == \"test\" + ( ( paRenthaScHildAnY == \"Test\" ) - ChildHasKey == \"Test\" ) ")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParserReturnsErrorForInvalidRegex(t *testing.T) {
	_, actual := search.Parse("Key==\"*\"")
	if actual == nil {
		t.Errorf("Expected error but got no error")
	}
}

func TestBasicCommandCorrectlyParsed(t *testing.T) {
	commands, _ := search.Parse("Any==\"Test\"")
	expected := search.Command{"Any", true, "Test", "", ""}
	actual := commands[0]
	if actual != expected {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexCommandParsesCorrectly(t *testing.T) {
	actual, _ := search.Parse("((Any==\"Test\")+Key!=\"Test2\")")
	expected := []search.Command{
		search.Command{"", false, "", "", "("},
		search.Command{"", false, "", "", "("},
		search.Command{"Any", true, "Test", "", ")"},
		search.Command{"Key", false, "Test2", "+", ")"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}
