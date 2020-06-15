package nodelist_test

import (
	"encoding/json"
	"kube-review/nodelist"
	"reflect"
	"regexp"
	"testing"
)

func TestReturnsErrorIfNotRootNode(t *testing.T) {
	_, err := nodelist.NewView([]*nodelist.Node{})
	if err == nil {
		t.Errorf("Expected an error but got nothing")
	}
}

func TestCanGetSizeOfView(t *testing.T) {
	nodes := createNodes(fullNodesRaw)
	view, _ := nodelist.NewView(nodes)
	actual := view.Size()
	expected := len(nodes)
	if actual != expected {
		t.Errorf("Expected %d but got %d", expected, actual)
	}
}

func TestCanGetFullNodeListFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetNodes(0, 17)
	expected := fullNodes
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanGetSubsetOfNodesFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetNodes(0, 2)
	expected := "Root\n├──GlossDiv"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanGetOffsetListOfNodesFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetNodes(5, 2)
	expected := "│  │     ├──Acronym\n│  │     ├──GlossDef"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanHandleRquestOfNonExistentNodes(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetNodes(0, 50)
	expected := fullNodes
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanGetFullJSONFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetJSON(0, 0, 23)
	expected := fullJson
	if !compareJSONStrings(actual, expected) {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanGetSubsetJSONFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetJSON(0, 0, 2)
	expected := "{\n    \"GlossDiv\": {"
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanGetOffsetJSONFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetJSON(2, 4, 2)
	expected := "        \"GlossDef\": {\n            \"GlossSeeAlso\": ["
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanGetOutOfRangeJSONFromView(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	actual := view.GetJSON(0, 0, 50)
	expected := fullJson
	if !compareJSONStrings(actual, expected) {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestSearchReturnsCorrectIndicesForMatches(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	r := regexp.MustCompile("Gloss")
	expected := []int{1, 2, 3, 6, 7, 11, 12}
	actual := view.GetNodesMatching(r, nodelist.ANY, true)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSearchWithNoMatchesReturnsEmptyArray(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	r := regexp.MustCompile("I will not match")
	matches := view.GetNodesMatching(r, nodelist.ANY, true)
	if len(matches) > 0 {
		t.Errorf("Expected empty array but return %d elements", len(matches))
	}
}

func TestSearchSpecificallyInKeyReturnsCorrectMatches(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	r := regexp.MustCompile("G")
	expected := []int{5, 8, 12, 13, 14}
	actual := view.GetNodesMatching(r, nodelist.VALUE, true)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSearchRelativeMatchesSelfIfChildLevelZero(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	r := regexp.MustCompile("Entry")
	expected := []int{3}
	actual := view.GetRelativesMatching(12, 1, 0, r, nodelist.KEY, true)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSearchRelativeDoesNotIncludeSelfIfChildLevelNotZero(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	r := regexp.MustCompile("Entry")
	actual := view.GetRelativesMatching(3, 0, 1, r, nodelist.KEY, true)
	if len(actual) > 0 {
		t.Errorf("Expected empty array but got '%v'", actual)
	}
}

func TestSearchRelativeReturnsIndexForMatchingChildren(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	r := regexp.MustCompile("Gloss")
	expected := []int{2, 3}
	actual := view.GetRelativesMatching(1, 0, 2, r, nodelist.KEY, true)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestFilterReturnsRootForEmptyNodes(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	newView, _ := view.Filter([]int{})
	expected := "Root"
	actual := newView.GetNodes(0, 2)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestFilterReturnsParentsAsWell(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	newView, _ := view.Filter([]int{4, 5, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	expected := fullNodes
	actual := newView.GetNodes(0, 17)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestFindHighlightDoesNotReturnStartOffsetIfOtherNodesAvailable(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	view.Highlight([]int{2, 5})
	expected := 2
	actual, _ := view.FindNextHighlight(0, 5)
	if actual != expected {
		t.Errorf("Expected %d but got %d", expected, actual)
	}
}

func TestFindHighlightReturnsErrorIfNoHighlightFound(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	view.Highlight([]int{2})
	_, actual := view.FindNextHighlight(3, 0)
	if actual == nil {
		t.Errorf("Expected an error but got nothing")
	}
}

func TestNewHighlightRemovesPreviousHighlights(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(fullNodesRaw))
	view.Highlight([]int{2})
	view.Highlight([]int{3})
	expected := 3
	actual, _ := view.FindNextHighlight(0, 0)
	if actual != expected {
		t.Errorf("Expected %d but got %d", expected, actual)
	}
}

// Now figure out how to test this...
// Test returning errors for failed parts

func TestSplitReturnsErrorIfRootNotValid(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(splitNodesRaw))
	_, actual := view.Split([]string{"path", "to", "non-root"}, []string{"path", "to", "target"})
	if actual == nil {
		t.Error("Expected an error but got none")
	}
}

func TestSplitReturnsErrorIfTargetNotValid(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(splitNodesRaw))
	_, actual := view.Split([]string{"path", "to", "root"}, []string{"path", "to", "non-target"})
	if actual == nil {
		t.Error("Expected an error but got none")
	}
}

func TestSplitReturnsCorrectNamesForSplits(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(splitNodesRaw))
	split, _ := view.Split([]string{"path", "to", "root"}, []string{"path", "to", "target"})
	actual := []string{}
	for key := range split {
		actual = append(actual, key)
	}
	expected := []string{"Hello", "Goodbye"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestSplitReturnsCorrectViews(t *testing.T) {
	view, _ := nodelist.NewView(createNodes(splitNodesRaw))
	split, _ := view.Split([]string{"path", "to", "root"}, []string{"path", "to", "target"})
	expected := "\"Test\""
	actual := split["Hello"].GetJSON(12, 0, 1)
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

// Test data
var (
	fullJson = `{
	"GlossDiv": {
		"GlossList": {
			"GlossEntry": {
				"Abbrev": "ISO 8879:1986",
				"Acronym": "SGML",
				"GlossDef": {
					"GlossSeeAlso": [
						"GML",
						"XML"
					],
					"para": "A meta-markup language, used to create markup languages such as DocBook."
				},
				"GlossSee": "markup",
				"GlossTerm": "Standard Generalized Markup Language",
				"ID": "SGML",
				"SortAs": "SGML"
			}
		},
		"title": "S"
	},
	"title": "example glossary"
}`
	fullNodes = `Root
├──GlossDiv
│  ├──GlossList
│  │  └──GlossEntry
│  │     ├──Abbrev
│  │     ├──Acronym
│  │     ├──GlossDef
│  │     │  ├──GlossSeeAlso
│  │     │  │  ├──0
│  │     │  │  └──1
│  │     │  └──para
│  │     ├──GlossSee
│  │     ├──GlossTerm
│  │     ├──ID
│  │     └──SortAs
│  └──title
└──title`
	fullNodesRaw = []nodelist.Node{
		nodelist.NewNode("Root", "{", 0),
		nodelist.NewNode("GlossDiv", "{", 1),
		nodelist.NewNode("GlossList", "{", 2),
		nodelist.NewNode("GlossEntry", "{", 3),
		nodelist.NewNode("Abbrev", "\"ISO 8879:1986\"", 4),
		nodelist.NewNode("Acronym", "\"SGML\"", 4),
		nodelist.NewNode("GlossDef", "{", 4),
		nodelist.NewNode("GlossSeeAlso", "[", 5),
		nodelist.NewNode("[]0", "\"GML\"", 6),
		nodelist.NewNode("[]1", "\"XML\"", 6),
		nodelist.NewNode("para", "\"A meta-markup language, used to create markup languages such as DocBook.\"", 5),
		nodelist.NewNode("GlossSee", "\"markup\"", 4),
		nodelist.NewNode("GlossTerm", "\"Standard Generalized Markup Language\"", 4),
		nodelist.NewNode("ID", "\"SGML\"", 4),
		nodelist.NewNode("SortAs", "\"SGML\"", 4),
		nodelist.NewNode("title", "\"S\"", 2),
		nodelist.NewNode("title", "\"example glossary\"", 1),
	}

	splitNodesRaw = []nodelist.Node{
		nodelist.NewNode("Root", "{", 0),
		nodelist.NewNode("path", "{", 1),
		nodelist.NewNode("to", "{", 2),
		nodelist.NewNode("root", "[", 3),
		nodelist.NewNode("[]0", "{", 4),
		nodelist.NewNode("path", "{", 5),
		nodelist.NewNode("to", "{", 6),
		nodelist.NewNode("target", "\"Hello\"", 7),
		nodelist.NewNode("[]1", "{", 4),
		nodelist.NewNode("path", "{", 5),
		nodelist.NewNode("to", "{", 6),
		nodelist.NewNode("target", "\"Goodbye\"", 7),
		nodelist.NewNode("[]2", "{", 4),
		nodelist.NewNode("path", "{", 5),
		nodelist.NewNode("to", "{", 6),
		nodelist.NewNode("target", "\"Hello\"", 7),
		nodelist.NewNode("other", "\"Test\"", 7),
	}
)

func createNodes(nl []nodelist.Node) []*nodelist.Node {
	var nodes []*nodelist.Node
	for index := range nl {
		nodes = append(nodes, &nl[index])
	}
	return nodes
}

func compareJSONStrings(json1 string, json2 string) bool {
	var check1, check2 interface{}
	json.Unmarshal([]byte(json1), &check1)
	json.Unmarshal([]byte(json2), &check2)
	if reflect.DeepEqual(check1, check2) {
		return true
	} else {
		return false
	}
}
