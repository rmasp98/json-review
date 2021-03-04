package nodelist_test

import (
	"testing"

	"github.com/rmasp98/kube-review/nodelist"
)

func TestReturnsErrorIfDataNotYetParsed(t *testing.T) {
	m, _ := nodelist.NewMasterNodeList([]byte(fullJson), false)
	if m.LoadStatus() == nil {
		t.Errorf("Expecting 'Incomplete' error but got nothing")
	}
}

func TestReturnNoErrorIfParsingComplete(t *testing.T) {
	m, _ := nodelist.NewMasterNodeList([]byte(fullJson), true)
	if m.LoadStatus() != nil {
		t.Errorf("Expecting no error but got '%s'", m.LoadStatus())
	}
}

func TestReturnsNodeViewForMaster(t *testing.T) {
	m, _ := nodelist.NewMasterNodeList([]byte(fullJson), true)
	actual, err := m.GetNodeView()
	if err != nil || actual.Size() != 17 {
		t.Errorf("Expected size to be 17 but got %d", actual.Size())
	}
}
