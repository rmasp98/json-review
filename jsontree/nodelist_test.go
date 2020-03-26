package jsontree_test

import (
	"encoding/json"
	"kube-review/jsontree"
	"reflect"
	"testing"
)

func TestReturnCorrectNumNodesForFullJson(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	if nodeList.Size() != 17 {
		t.Errorf("Expected '17' but got '%d'", nodeList.Size())
	}
}

func TestGetNodesReturnsCorrectFormat(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	actual := nodeList.GetNodes(17)
	if actual != fullNodes {
		t.Errorf("Expected out:\n%s\n\nActual output:\n%s", fullNodes, actual)
	}
}

func TestGetNodesReturnsOnlyRequiredNumberOfNodes(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	actual := nodeList.GetNodes(1)
	if actual != "Root" {
		t.Errorf("Expected out:\n'Root'\n\nActual output:\n%s", actual)
	}
}

func TestMoveTopNodeChangesGetNodesResponse(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(5)
	actual := nodeList.GetNodes(1)
	if actual != "│  │     ├──Acronym" {
		t.Errorf("Expected out:\n'│  │     ├──Acronym'\n\nActual output:\n%s", actual)
	}
}

func TestMoveTopNodeReturnsLastNodeIfMovedTooFar(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(20)
	actual := nodeList.GetNodes(1)
	if actual != "└──title" {
		t.Errorf("Expected out:\n'└──title'\n\nActual output:\n%s", actual)
	}
}

func TestMoveTopNodeReturnsFirstNodeIfMovedTooFarBack(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(-2)
	actual := nodeList.GetNodes(1)
	if actual != "Root" {
		t.Errorf(`Expected "%s" but got "%s"`, fullJson, actual)
	}
}

func TestGetJSONReturnsCorrectFormat(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	actual := nodeList.GetJSON(23)
	if !compareJSONStrings(actual, fullJson) {
		t.Errorf(`Expected "%s" but got "%s"`, fullJson, actual)
	}

}

func TestGetJsonReturnsOnlyRequestedNumberOfLines(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	actual := nodeList.GetJSON(1)
	if actual != "{" {
		t.Errorf(`Expected "{" but got "%s"`, actual)
	}
}

func TestGetJSONReturnsJSONForActiveNode(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(7)
	actual := nodeList.GetJSON(3)
	expected := "[\n   \"GML\",\n   \"XML\""
	if actual != expected {
		t.Errorf("Expected out:\n'%s'\n\nActual output:\n%s",
			expected, actual,
		)
	}
}

func TestSetActiveDoesNothingIfTooSmall(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(-1)
	actual := nodeList.GetJSON(1)
	if actual != "{" {
		t.Errorf(`Expected "{" but got "%s"`, actual)
	}
}

func TestSetActiveNodeDoesNothingIfTooBig(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(30)
	actual := nodeList.GetJSON(1)
	if actual != "{" {
		t.Errorf(`Expected "{" but got "%s"`, actual)
	}
}

func TestCanCollapseActiveNode(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(7)
	nodeList.CollapseActiveNode()
	if nodeList.Size() != 13 {
		t.Errorf("Expected 13 nodes but got '%d'", nodeList.Size())
	}
}

func TestCanExpandCollapsedNodes(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(7)
	nodeList.CollapseActiveNode()
	nodeList.SetActiveNode(6)
	nodeList.ExpandActiveNode()
	if nodeList.Size() != 17 {
		t.Errorf("Expected 13 nodes but got '%d'", nodeList.Size())
	}
}

func TestCollapseNodeUpdatesActiveNodeToParent(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(7)
	nodeList.CollapseActiveNode()
	nodeList.ExpandActiveNode()
	if nodeList.Size() != 17 {
		t.Errorf("Expected 13 nodes but got '%d'", nodeList.Size())
	}
}

func TestEnsureActiveNodeIsRelativeToTopNode(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(2)
	nodeList.SetActiveNode(2)
	actual := nodeList.GetJSON(1)
	if actual != `"ISO 8879:1986"` {
		t.Errorf("Expected out:\n'%s'\n\nActual output:\n%s",
			`"ISO 8879:1986"`, actual,
		)
	}
}

func TestGetNodesOnlyReturnsVisibleNode(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.CollapseActiveNode()
	actual := nodeList.GetNodes(17)
	if actual != "Root" {
		t.Errorf("Expected 'Root' nodes but got '%s'", actual)
	}
}

func TestCollapseNodeReturnsVisibleIndexOfParent(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(2)
	nodeList.SetActiveNode(2)
	actual := nodeList.CollapseActiveNode()
	if actual != 1 {
		t.Errorf("Expected '1' nodes but got '%d'", actual)
	}
}

func TestSetActiveNodeAccountsForNonVisibleNodes(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(8)
	nodeList.CollapseActiveNode()
	nodeList.SetActiveNode(8)
	nodeList.CollapseActiveNode()
	actual := nodeList.Size()
	if actual != 13 {
		t.Errorf("Expected '13' nodes but got '%d'", actual)
	}
}

func TestMoveTopNodeAccountsForInvisibleNodes(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.SetActiveNode(4)
	nodeList.CollapseActiveNode()
	nodeList.MoveTopNode(5)
	actual := nodeList.GetNodes(1)
	if actual != "└──title" {
		t.Errorf("Expected '└──title' nodes but got '%s'", actual)
	}
}

func TestUpdateTopNodeUpdatesActiveNode(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(2)
	nodeList.CollapseActiveNode()
	actual := nodeList.Size()
	if actual != 3 {
		t.Errorf("Expected '3' nodes but got '%d'", actual)
	}
}

func TestCollapseActiveNodeMovesTopNodeIfRequired(t *testing.T) {
	nodeList, _ := jsontree.NewNodeList(fullJson)
	nodeList.MoveTopNode(1)
	nodeList.CollapseActiveNode()
	actual := nodeList.GetNodes(1)
	if actual != "Root" {
		t.Errorf("Expected 'Root' nodes but got '%s'", actual)
	}
}

//////////////////////////////////////////////////////////////////////////
//  TEST DATA

var fullJson = `{
	"GlossDiv": {
		"title": "S",
		"GlossList": {
			"GlossEntry": {
				"ID": "SGML",
				"SortAs": "SGML",
				"GlossTerm": "Standard Generalized Markup Language",
				"Acronym": "SGML",
				"Abbrev": "ISO 8879:1986",
				"GlossDef": {
					"para": "A meta-markup language, used to create markup languages such as DocBook.",
					"GlossSeeAlso": ["GML", "XML"]
				},
				"GlossSee": "markup"
			}
		}
	},
	"title": "example glossary"
}`

var fullNodes = `Root
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
