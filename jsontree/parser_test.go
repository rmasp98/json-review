package jsontree_test

import (
	"kube-review/jsontree"
	"testing"
)

func TestCanParseJSONString(t *testing.T) {
	stringJSON := `"Hello World"`
	json := GetJSONData(stringJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(false)
	if actual != stringJSON {
		t.Errorf("Expecting '%s' but got '%s'", stringJSON, actual)
	}
}

func TestCanParseJSONInt(t *testing.T) {
	intJSON := `1492`
	json := GetJSONData(intJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(false)
	if actual != intJSON {
		t.Errorf("Expecting '%s' but got '%s'", intJSON, actual)
	}
}

func TestCanParseJSONFloat(t *testing.T) {
	floatJSON := `1.3947619`
	json := GetJSONData(floatJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(false)
	if actual != floatJSON {
		t.Errorf("Expecting '%s' but got '%s'", floatJSON, actual)
	}
}

func TestCanParseJSONBool(t *testing.T) {
	boolJSON := `false`
	json := GetJSONData(boolJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(false)
	if actual != boolJSON {
		t.Errorf("Expecting '%s' but got '%s'", boolJSON, actual)
	}
}

func TestCanParseJSONNull(t *testing.T) {
	nullJSON := `null`
	json := GetJSONData(nullJSON)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	actual := nodes[0].GetJSON(false)
	if actual != nullJSON {
		t.Errorf("Expecting '%s' but got '%s'", nullJSON, actual)
	}
}

func TestCanParseJSONMultipleArray(t *testing.T) {
	json := GetJSONData(`["Goodbye","World","Hello","World"]`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	expectedArray := []string{`"Goodbye"`, `"World"`, `"Hello"`, `"World"`}
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(false)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultipleMap(t *testing.T) {
	json := GetJSONData(`{"Goodbye":"World","Hello":"World"}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	actual1 := nodes[0].GetJSON(true)
	expected1 := `"Goodbye":"World"`
	actual2 := nodes[1].GetJSON(true)
	expected2 := `"Hello": "World"`
	if actual1 != expected1 && actual2 != expected2 {
		t.Errorf("Expecting '%s' and '%s' but got '%s' and '%s'", expected1, expected2, actual1, actual2)
	}
}

func TestCanParseJSONMultiLevelArray(t *testing.T) {
	json := GetJSONData(`[{"Goodbye":"Child"},{"Hello":"Adult"}]`)
	nodes, _ := jsontree.CreateTreeNodes(json, 0)
	expectedArray := []string{"{", `"Goodbye": "Child"`, "{", `"Hello": "Adult"`}
	levels := []int{0, 1, 0, 1}
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(levels[index] > 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultiLevelMap(t *testing.T) {
	json := GetJSONData(`{"Goodbye":{"Cruel World":"Test","Hello":"World"}}`)
	nodes, _ := jsontree.CreateTreeNodes(json, 1)
	expectedArray := []string{`{`, `"Cruel World": "Test"`, `"Hello": "World"`}
	levels := []int{0, 1, 1}
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(levels[index] > 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}
