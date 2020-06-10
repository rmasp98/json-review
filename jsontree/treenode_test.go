package jsontree_test

import (
	"encoding/json"
	"kube-review/jsontree"
	"regexp"
	"testing"
)

func TestReturnsPrefixPlusKeyInGetNodes(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	nodes[0].Prefix = "Test - "
	actual := nodes[0].GetNode()
	expected := "Test - Goodbye"
	if actual != expected {
		t.Errorf("Expecting '%s' but got '%s'", expected, actual)
	}
}

func TestReturnsCorrectLevelForEachNode(t *testing.T) {
	json := GetJSONData(`{"Goodbye":{"Cruel World":"Test","Hello":"World"}}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 5)
	expectedLevels := []int{5, 6, 6}
	for index, expected := range expectedLevels {
		actual := nodes[index].Level
		if actual != expected {
			t.Errorf("Expecting '%d' but got '%d'", expected, actual)
		}
	}
}

func TestGetNodeForArrayRemovesBrackets(t *testing.T) {
	json := GetJSONData(`["Goodbye","World","Hello","World"]`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	expectedNodes := []string{"0", "1"}
	for index, expected := range expectedNodes {
		actual := nodes[index].GetNode()
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestMatchReturnsFalseIfNoMatchInKey(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	r := regexp.MustCompile("World")
	actual := nodes[0].Match(r, jsontree.KEY)
	if actual != false {
		t.Errorf("Expecting 'false' but got '%t'", actual)
	}
}

func TestMatchReturnsTrueIfMatchInValue(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	r := regexp.MustCompile("World")
	actual := nodes[0].Match(r, jsontree.VALUE)
	if actual != true {
		t.Errorf("Expecting 'true' but got '%t'", actual)
	}
}

func TestGetJsonGivesCorrectResponseIfHighlighted(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	nodes[0].IsHighlighted = true
	actual := nodes[0].GetJSON(true)
	expected := "\033[41m\"Goodbye\": \"World\"\033[0m"
	if actual != expected {
		t.Errorf("Expecting '%s' but got '%s'", expected, actual)
	}
}

func TestGetEndingReturnsCorrectEndingIfLast(t *testing.T) {
	json := GetJSONData(`{"Goodbye":{"Horrible":"World","Cruel":["World"]}}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	actual1 := nodes[0].GetEnding(true)
	actual2 := nodes[1].GetEnding(true)
	actual3 := nodes[2].GetEnding(true)
	if actual1 != "}" || actual2 != "]" || actual3 != "" {
		t.Errorf("Expecting '}', ']' and '' but got '%s', '%s' and '%s", actual1, actual2, actual3)
	}
}

func TestGetEndingReturnsCloseBracketAndIfValueOpenBracketAndNotLast(t *testing.T) {
	json := GetJSONData(`{"Goodbye":{"Horrible":"World","Cruel":["World"]}}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	actual1 := nodes[0].GetEnding(false)
	actual2 := nodes[1].GetEnding(false)
	actual3 := nodes[2].GetEnding(false)
	if actual1 != "}," || actual2 != "]," || actual3 != "," {
		t.Errorf("Expecting '},', '],' and ',' but got '%s', '%s' and '%s", actual1, actual2, actual3)
	}
}

func TestGetEndingReturnsNothingIfFiltered(t *testing.T) {
	json := GetJSONData(`{"Goodbye":{"Horrible":"World","Cruel":["World"]}}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	nodes[0].Filter(true)
	actual1 := nodes[0].GetEnding(false)
	nodes[1].Filter(true)
	actual2 := nodes[1].GetEnding(false)
	nodes[2].Filter(true)
	actual3 := nodes[2].GetEnding(false)
	if actual1 != "" || actual2 != "" || actual3 != "" {
		t.Errorf("Expecting '', '' and '' but got '%s', '%s' and '%s", actual1, actual2, actual3)
	}
}

func TestJsonReturnsNothingIfFiltered(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	nodes[0].Filter(true)
	actual := nodes[0].GetJSON(false)
	if actual != "" {
		t.Errorf("Expecting nothing but got '%s'", actual)
	}
}

// HELPER FUNCTIONS AND DATA //////////////////////////////////////////////

func GetJSONData(jsonData string) interface{} {
	var data interface{}
	json.Unmarshal([]byte(jsonData), &data)
	return data
}
