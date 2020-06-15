package nodelist_test

import (
	"fmt"
	"kube-review/nodelist"
	"testing"
)

func TestCanParseJSONString(t *testing.T) {
	stringJSON := `"Hello World"`
	nodes := []nodelist.Node{}
	parser := nodelist.NewParser(&nodes, nil)
	parser.Parse([]byte(stringJSON), true)
	actual := nodes[0].GetJSON(false)
	if actual != stringJSON {
		t.Errorf("Expecting '%s' but got '%s'", stringJSON, actual)
	}
}

func TestCanParseJSONInt(t *testing.T) {
	intJSON := `1492`
	nodeList := getNodeList(intJSON, nil)
	actual := nodeList[0].GetJSON(false)
	if actual != intJSON {
		t.Errorf("Expecting '%s' but got '%s'", intJSON, actual)
	}
}

func TestCanParseJSONFloat(t *testing.T) {
	floatJSON := `1.3947619`
	nodeList := getNodeList(floatJSON, nil)
	actual := nodeList[0].GetJSON(false)
	if actual != floatJSON {
		t.Errorf("Expecting '%s' but got '%s'", floatJSON, actual)
	}
}

func TestCanParseJSONBool(t *testing.T) {
	boolJSON := `false`
	nodeList := getNodeList(boolJSON, nil)
	actual := nodeList[0].GetJSON(false)
	if actual != boolJSON {
		t.Errorf("Expecting '%s' but got '%s'", boolJSON, actual)
	}
}

func TestCanParseJSONNull(t *testing.T) {
	nullJSON := `null`
	nodeList := getNodeList(nullJSON, nil)
	actual := nodeList[0].GetJSON(false)
	if actual != nullJSON {
		t.Errorf("Expecting '%s' but got '%s'", nullJSON, actual)
	}
}

func TestCanParseJSONMultipleArray(t *testing.T) {
	expectedArray := []string{"[", `"Goodbye"`, `"World"`, `"Hello"`, `"World"`}
	nodes := getNodeList(`["Goodbye","World","Hello","World"]`, nil)
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(index != 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultipleMap2(t *testing.T) {
	expectedResults := []string{"{", `"Goodbye": "World"`, `"Hello": "World"`}
	nodes := getNodeList(`{"Goodbye":"World","Hello":"World"}`, nil)
	for index, expected := range expectedResults {
		actual := nodes[index].GetJSON(index != 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultiLevelArray(t *testing.T) {
	expectedArray := []string{"[", "{", `"Goodbye": "Child"`, "{", `"Hello": "Adult"`}
	nodes := getNodeList(`[{"Goodbye":"Child"},{"Hello":"Adult"}]`, nil)
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(index != 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestCanParseJSONMultiLevelMap(t *testing.T) {
	expectedArray := []string{`{`, `"Goodbye": {`, `"Cruel World": "Test"`, `"Hello": "World"`}
	nodes := getNodeList(`{"Goodbye":{"Cruel World":"Test","Hello":"World"}}`, nil)
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(index != 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestEnsureMapsReturnInAlphabeticalOrder(t *testing.T) {
	expectedArray := []string{`{`, `"a": 1`, `"b": 2`, `"c": 3`}
	nodes := getNodeList(`{"c":3, "a":1,"b":2}`, nil)
	for index, expected := range expectedArray {
		actual := nodes[index].GetJSON(index != 0)
		if actual != expected {
			t.Errorf("Expecting '%s' but got '%s'", expected, actual)
		}
	}
}

func TestParseThrowsAnErrorIfJsonInvalid(t *testing.T) {
	nodes := []nodelist.Node{}
	parser := nodelist.NewParser(&nodes, nil)
	actual := parser.Parse([]byte(`{"c"=3, "a":1,"b":2}`), true)
	if actual == nil {
		t.Errorf("Expecting an error but got nothing")
	}
}

func TestCallsCallbackWithNilIfSuccessful(t *testing.T) {
	actual := fmt.Errorf("Test Error")
	getNodeList(`{"c":3, "a":1,"b":2}`, func(err error) { actual = err })
	if actual != nil {
		t.Errorf("Expecting no error but got '%s'", actual)
	}
}

// HELPER FUNCTIONS AND DATA //////////////////////////////////////////////

func getNodeList(jsonData string, callback func(error)) []nodelist.Node {
	nodes := []nodelist.Node{}
	parser := nodelist.NewParser(&nodes, callback)
	parser.Parse([]byte(jsonData), true)
	return nodes
}
