package nodelist_test

import (
	"kube-review/nodelist"
	"reflect"
	"sort"
	"testing"
)

// Tests
// splitviews
// listviews
// setview

func TestMoveTopNodeOffsetsGetNodes(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.MoveTopNode(1)
	expected := "├──GlossDiv"
	actual := nl.GetNodes(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestMoveTopNodeCannotGoBelowZero(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.MoveTopNode(-1)
	expected := "Root"
	actual := nl.GetNodes(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestMoveTopNodeCannotGoAboveFinalNode(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.MoveTopNode(50)
	expected := "└──title"
	actual := nl.GetNodes(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestSetActiveNodeUpdatesJSONOutput(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.SetActiveNode(4)
	expected := `"ISO 8879:1986"`
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestSetActiveNodeIsRelativeToTopNode(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.MoveTopNode(2)
	nl.SetActiveNode(2)
	expected := `"ISO 8879:1986"`
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestSetActiveNodeCannotBeAboveFinalIndex(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.SetActiveNode(20)
	expected := `"example glossary"`
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestMoveJSONViewDoesWhatItSays(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.MoveJSONView(4)
	expected := `                "Abbrev": "ISO 8879:1986"`
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestJSONOffsetResetOnNewActiveNode(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.MoveJSONView(4)
	nl.SetActiveNode(5)
	expected := `"SGML"`
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestFilterRemovesUndefinedNodes(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.Filter([]int{3})
	nl.SetActiveNode(3)
	expected := "{\n}"
	actual := nl.GetJSON(2)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanFindNextHighlight(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(fullJson), true)
	nl.Highlight([]int{5})
	nl.FindNextHighlight()
	expected := `                "Acronym": "SGML"`
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanSplitNodesBasedOnInputString(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(splitNodesBase), true)
	nl.SplitViews("path.to.root = path.to.target")
	expected := []string{"Goodbye", "Hello", "main"}
	actual := nl.ListViews()
	sort.Strings(actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestCanSplitNodesWithNoRoot(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(splitNodesArray), true)
	nl.SplitViews("path.to.target")
	expected := []string{"Goodbye", "Hello", "main"}
	actual := nl.ListViews()
	sort.Strings(actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSplitreturnsErrorIfInvalid(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(splitNodesArray), true)
	err := nl.SplitViews("")
	if err == nil {
		t.Errorf("Expected error but got none")
	}
}

func TestSetViewReturnsErrorIfViewDoesNotExist(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(splitNodesBase), true)
	actual := nl.SetView("Not a View")
	if actual == nil {
		t.Errorf("Expected error but got none")
	}
}

func TestCanSetViewToASplitView(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(splitNodesBase), true)
	nl.SplitViews("path.to.root = path.to.target")
	nl.SetView("Hello")
	nl.SetActiveNode(11)
	expected := "\"Test\""
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestCanResetView(t *testing.T) {
	nl, _ := nodelist.NewNodeList([]byte(splitNodesBase), true)
	nl.SplitViews("path.to.root = path.to.target")
	nl.SetView("Hello")
	nl.Filter([]int{})
	nl.ResetView()
	nl.SetActiveNode(11)
	expected := "\"Test\""
	actual := nl.GetJSON(1)
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

// Test data

var splitNodesBase = `{
	"path": {
		"to": {
			"root": [
				{
					"path": {
						"to": {
							"target": "Hello"
						}
					}
				},
				{
					"path": {
						"to": {
							"target": "Goodbye"
						}
					}
				},
				{
					"path": {
						"to": {
							"target": "Hello",
							"Other": "Test"
						}
					}
				}
			]
		}
	}
}`

var splitNodesArray = `[
	{
		"path": {
			"to": {
				"target": "Hello"
			}
		}
	},
	{
		"path": {
			"to": {
				"target": "Goodbye"
			}
		}
	}
]`
