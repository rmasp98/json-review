package jsontree_test

import (
	"encoding/json"
	"kube-review/jsontree"
	"regexp"
	"testing"
)

func TestCanParseJSONString(t *testing.T) {
	stringJSON := `"Hello World"`
	json := GetJSONData(stringJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(0)
	if actual != stringJSON {
		t.Errorf("Expecting '%s' but got '%s'", stringJSON, actual)
	}
}

func TestCanParseJSONInt(t *testing.T) {
	intJSON := `1492`
	json := GetJSONData(intJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(0)
	if actual != intJSON {
		t.Errorf("Expecting '%s' but got '%s'", intJSON, actual)
	}
}

func TestCanParseJSONFloat(t *testing.T) {
	floatJSON := `1.3947619`
	json := GetJSONData(floatJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(0)
	if actual != floatJSON {
		t.Errorf("Expecting '%s' but got '%s'", floatJSON, actual)
	}
}

func TestCanParseJSONBool(t *testing.T) {
	boolJSON := `false`
	json := GetJSONData(boolJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(0)
	if actual != boolJSON {
		t.Errorf("Expecting '%s' but got '%s'", boolJSON, actual)
	}
}

func TestCanParseJSONNull(t *testing.T) {
	nullJSON := `null`
	json := GetJSONData(nullJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(0)
	if actual != nullJSON {
		t.Errorf("Expecting '%s' but got '%s'", nullJSON, actual)
	}
}

func TestCanParseJSONMultipleArray(t *testing.T) {
	json := GetJSONData(`["Goodbye","World","Hello","World"]`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	expectedArray := []string{`"Goodbye"`, `"World"`, `"Hello"`, `"World"`}
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultipleMap(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World","Hello":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	actual1 := nodes[0].GetJSON(1)
	expected1 := `   "Goodbye":"World"`
	actual2 := nodes[1].GetJSON(1)
	expected2 := `   "Hello": "World"`
	if actual1 != expected1 && actual2 != expected2 {
		t.Errorf("Expecting '%s' and '%s' but got '%s' and '%s'", expected1, expected2, actual1, actual2)
	}
}

func TestCanParseJSONMultiLevelArray(t *testing.T) {
	json := GetJSONData(`[{"Goodbye":"Child"},{"Hello":"Adult"}]`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	expectedArray := []string{"{", `   "Goodbye": "Child"`, "{", `   "Hello": "Adult"`}
	levels := []int{0, 1, 0, 1}
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(levels[index])
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultiLevelMap(t *testing.T) {
	json := GetJSONData(`{"Goodbye":{"Cruel World":"Test","Hello":"World"}}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	expectedArray := []string{`{`, `   "Cruel World": "Test"`, `   "Hello": "World"`}
	levels := []int{0, 1, 1}
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(levels[index])
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

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
	actual := nodes[0].GetJSON(1)
	expected := "   \033[41m\"Goodbye\": \"World\"\033[0m"
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
	actual := nodes[0].GetJSON(0)
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
