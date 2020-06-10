package nodelist_test

import (
	"kube-review/nodelist"
	"regexp"
	"testing"
)

func TestCorrectFormatReturnedForGetJSON(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	actual := node.GetJSON(true)
	expected := "\"key\": \"value\""
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestOnlyValueReturnedForArrayGetJSON(t *testing.T) {
	node := nodelist.NewNode("[]0", "\"value\"", 0)
	actual := node.GetJSON(true)
	expected := "\"value\""
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestReturnsOnlyValueIfFullIsFalse(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	actual := node.GetJSON(false)
	expected := "\"value\""
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestGetNodeRemovesBracketsFromArrayKey(t *testing.T) {
	node := nodelist.NewNode("[]0", "\"value\"", 0)
	actual := node.GetNode()
	expected := "0"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestGetCloseBracketReturnsNothingForNormalNodes(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	actual := node.GetCloseBracket()
	expected := ""
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestGetCloseBracketReturnsBraceForMaps(t *testing.T) {
	node := nodelist.NewNode("key", "{", 0)
	actual := node.GetCloseBracket()
	expected := "}"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestGetCloseBracketReturnsSquareForArrays(t *testing.T) {
	node := nodelist.NewNode("key", "[", 0)
	actual := node.GetCloseBracket()
	expected := "]"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestMatchKeyReturnsTrueForSuccessfulMatch(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("key")
	actual := node.MatchKey(r)
	expected := true
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}

func TestMatchKeyReturnsFalseForUnsuccessfulMatch(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("value")
	actual := node.MatchKey(r)
	expected := false
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}

func TestMatchValueReturnsTrueForSuccessfulMatch(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("^value$")
	actual := node.MatchValue(r)
	expected := true
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}

func TestMatchValueReturnsFalseForUnsuccessfulMatch(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("key")
	actual := node.MatchValue(r)
	expected := false
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}

func TestMatchReturnsTrueForSuccessfulMatchInKey(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("key")
	actual := node.Match(r)
	expected := true
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}

func TestMatchReturnsTrueForSuccessfulMatchInValue(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("value")
	actual := node.Match(r)
	expected := true
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}

func TestMatchReturnsFalseForUnsuccessfulMatch(t *testing.T) {
	node := nodelist.NewNode("key", "\"value\"", 0)
	r := regexp.MustCompile("test")
	actual := node.MatchValue(r)
	expected := false
	if actual != expected {
		t.Errorf("Expected '%t' but got '%t'", expected, actual)
	}
}
