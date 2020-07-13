package search_test

import (
	"kube-review/mocks"
	"kube-review/search"
	"reflect"
	"testing"
)

func TestBasicCommandExecutesCorrectFunctions(t *testing.T) {
	expression, _ := search.NewExpression("FindNodes(\"test\")")
	mock := mocks.NodeListMock{}
	expression.Execute(&mock)
	expected := []string{"GetNodesMatching"}
	actual := mock.Calls
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexExpressionExecutesCorrectFunctions(t *testing.T) {
	expression, _ := search.NewExpression("FindNodes(\"test\", output=nodes) + FindRelative(nodes, \"test\")")
	mock := mocks.NodeListMock{}
	mock.Returns = [][]int{[]int{1}}
	expression.Execute(&mock)
	expected := []string{"GetNodesMatching", "GetRelativesMatching"}
	actual := mock.Calls
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestComplexExpressionReturnsCorrectOutput(t *testing.T) {
	expression, _ := search.NewExpression("FindNodes(\"test\", output=nodes) + FindRelative(nodes, \"test\")")
	mock := mocks.NodeListMock{}
	mock.Returns = [][]int{[]int{1, 5}, []int{1, 2}, []int{3, 6}}
	actual := expression.Execute(&mock)
	expected := []int{1, 2, 3, 5, 6}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestBracketExpressionReturnsCorrectOutput(t *testing.T) {
	expression, _ := search.NewExpression("FindNodes(\"test\") - (FindNodes(\"test\") + FindNodes(\"test\"))")
	mock := mocks.NodeListMock{}
	mock.Returns = [][]int{[]int{1, 2, 3, 4, 5}, []int{1, 2}, []int{3, 6}}
	actual := expression.Execute(&mock)
	expected := []int{4, 5}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestHintsReturnFunctionSignatures(t *testing.T) {
	actual := search.GetExpressionHints("")
	expected := []string{"FindNodes(regex, matchType, equal, output)", "FindRelative(nodes, regex, relativeStart, depth, matchType, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestHintsReturnsReleventFunctionSignature(t *testing.T) {
	actual := search.GetExpressionHints("findn")
	expected := []string{"FindNodes(regex, matchType, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestHintsWorksWithComplexPrefix(t *testing.T) {
	actual := search.GetExpressionHints("( ( ( ( findN")
	expected := []string{"FindNodes(regex, matchType, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestHintsReturnsOperatorsAfterCloseBracket(t *testing.T) {
	actual := search.GetExpressionHints("findnodes(\"test\")")
	expected := search.Operators
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestReturnsRegexInformationInFunction(t *testing.T) {
	actual := search.GetExpressionHints("FindNodes(")
	expected := []string{"FindNodes(\033[1;31mregex (quoted regex string)\033[0m, matchType, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestReturnsMatchTypeInformationInFunction(t *testing.T) {
	actual := search.GetExpressionHints("FindNodes(\"test\",")
	expected := []string{"ANY", "KEY", "VALUE", "FindNodes(regex, \033[1;31mmatchType (attribute to match against)\033[0m, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestReturnsEqualsInformationInFunction(t *testing.T) {
	actual := search.GetExpressionHints("FindRelative(nodes, \"test\", 0, 1, ANY,")
	expected := []string{"true", "false", "FindRelative(nodes, regex, relativeStart, depth, matchType, \033[1;31mequal (should match be equal or not equal to regex)\033[0m, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestReturnsDepthDetailInFindRelativeHint(t *testing.T) {
	actual := search.GetExpressionHints("FindRelative(nodes, \"test\", 0,")
	expected := []string{"FindRelative(nodes, regex, relativeStart, \033[1;31mdepth (number of levels search should go down)\033[0m, matchType, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestHighlightsCorrectFieldIfKwargsPresent(t *testing.T) {
	actual := search.GetExpressionHints("FindRelative(nodes=nodes, o")
	expected := []string{"FindRelative(nodes, regex, relativeStart, depth, matchType, equal, \033[1;31moutput (variable that holds matched nodes. If exists, append to previous result)\033[0m)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestInsertHintAppendsCorrectOperator(t *testing.T) {
	actual := search.InsertSelectedExpressionHint("FindNodes(\"test\")", 1)
	expected := "FindNodes(\"test\") | "
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCorrectlyInsertsArgumentHint(t *testing.T) {
	actual := search.InsertSelectedExpressionHint("FindNodes(\"test\",", 2)
	expected := "FindNodes(\"test\",VALUE"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCorrectlyInsertsKeywordArgumentHint(t *testing.T) {
	actual := search.InsertSelectedExpressionHint("FindNodes(equal=", 0)
	expected := "FindNodes(equal=true"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestDoesNothingForFunctionHintInArgumentHints(t *testing.T) {
	actual := search.InsertSelectedExpressionHint("FindNodes(", 0)
	expected := "FindNodes("
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestInsertHintsInsertCorrectFunctionUpToBracket(t *testing.T) {
	actual := search.InsertSelectedExpressionHint("Find", 1)
	expected := "FindRelative("
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestGetHintsPreviousFunctions(t *testing.T) {
	actual := search.GetExpressionHints("FindNodes(\"test\") + ")
	expected := []string{"FindNodes(regex, matchType, equal, output)", "FindRelative(nodes, regex, relativeStart, depth, matchType, equal, output)"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestCanInsertSecondFunction(t *testing.T) {
	actual := search.InsertSelectedExpressionHint("FindNodes(\"test\") + Fin", 1)
	expected := "FindNodes(\"test\") + FindRelative("
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}
