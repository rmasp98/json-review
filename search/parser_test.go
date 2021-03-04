package search_test

import (
	"reflect"
	"testing"

	"github.com/rmasp98/kube-review/search"
)

func TestParseReturnsErrorForInvalidFunction(t *testing.T) {
	_, actual := search.Parse("NotAFunction()")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestParseBasicFindFunctionReturnsNoError(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseReturnsErrorIfNoBracketBeforeFunction(t *testing.T) {
	_, actual := search.Parse("Find")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestParseReturnsErrorIfFindDoesNotHaveRegexArgument(t *testing.T) {
	_, actual := search.Parse("FindNodes( )")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestParseReturnsErrorIfFindHasNoClosingBrackets(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\"")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestParseReturnsNoErrorForBasicBrackets(t *testing.T) {
	_, actual := search.Parse("(FindNodes(\"test\"))")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseReturnsThrowErrorForNoCloseBracket(t *testing.T) {
	_, actual := search.Parse("(FindNodes(\"test\")")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestParseCanHandleMultipleBrackets(t *testing.T) {
	_, actual := search.Parse("(((FindNodes(\"test\"))))")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseReturnsNoErrorForFunctionWithOperator(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\")+FindNodes(\"test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseReturnsErrorForHangingOperator(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\")+")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestParseReturnsErrorIfFunctionsHaveNoOperator(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\")FindNodes(\"test\")")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestCanHandleSimilarOperators(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\")-FindNodes(\"test\")->FindNodes(\"test\")")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestParseCanHandleSpaces(t *testing.T) {
	_, actual := search.Parse(" FindNodes( \"test\") +  FindNodes( \"test\" ) ")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestThrowsErrorForInvalidRegex(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"*\")")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestCanParseAllArgumentsWithoutError(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\", ANY, true, out)")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestCanParseKeywordArgumentsWithoutError(t *testing.T) {
	_, actual := search.Parse("FindNodes(regex=\"test\", matchType=ANY, equal=true, output=out)")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestCanParseThrowsErrorForInvalidKwarg(t *testing.T) {
	_, actual := search.Parse("FindNodes(regex=\"test\", matchType=ANY, equal=true, fakearg=out)")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestCanParseThrowsErrorIfNormalArgFollowsKwarg(t *testing.T) {
	_, actual := search.Parse("FindNodes(regex=\"test\", matchType=ANY, equal=true, out)")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestCanParseAllFindRelativeArgumentsWithoutError(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\", output=nodes) + FindRelative(nodes, \"test\", 0, 1, ANY, true, out)")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestThrowsErrorIfInputNotYetCreated(t *testing.T) {
	_, actual := search.Parse("FindRelative(input, \"test\", 0, 1, ANY, true, out)")
	if actual == nil {
		t.Error("Expected error but got nothing")
	}
}

func TestEverythingButInputOutputIsCaseInsensitive(t *testing.T) {
	_, actual := search.Parse("FindNodes(\"test\", OutPut=nodes) + fiNdreLatiVe(nodes, \"test\", 0, 1, aNy, tRUe, OUt)")
	if actual != nil {
		t.Errorf("Expected no error but got '%s'", actual)
	}
}

func TestBasicCommandCorrectlyParsed(t *testing.T) {
	actual, _ := search.Parse("FindNodes(\"test\", ANY, true, out)")
	expected := search.NewCommand(search.CMDFINDNODES, map[string]string{"regex": "test", "matchType": "ANY", "equal": "true"}, "out", "", "")
	if !reflect.DeepEqual(actual[0], expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual[0])
	}
}

func TestComplexCommandCorrectlyParsed(t *testing.T) {
	actual, _ := search.Parse("((FindNodes(\"test\", Key, output=outnodes)) -> FindRelative(outnodes, \"next test\", 2, 5, KEY, true))")
	expected := []search.Command{
		search.NewCommand(search.CMDNULL, nil, "", "", "("),
		search.NewCommand(search.CMDNULL, nil, "", "", "("),
		search.NewCommand(search.CMDFINDNODES, map[string]string{"regex": "test", "matchType": "Key"}, "outnodes", "", ")"),
		search.NewCommand(search.CMDFINDRELATIVE, map[string]string{"nodes": "outnodes", "regex": "next test", "relativeStart": "2", "depth": "5", "matchType": "KEY", "equal": "true"}, "", "->", ")"),
	}
	for index := range expected {
		if !reflect.DeepEqual(actual[index], expected[index]) {
			t.Errorf("Expected \n'%v' but got \n'%v'", expected[index], actual[index])
		}
	}
}
