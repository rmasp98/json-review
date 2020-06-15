package nodelist

import "fmt"

// MasterNodeList contains the master copy of the full NodeList
type MasterNodeList struct {
	nodes      []Node
	loadStatus error
}

// NewMasterNodeList stuff
func NewMasterNodeList(jsonData []byte, blocking bool) (MasterNodeList, error) {
	m := MasterNodeList{[]Node{}, fmt.Errorf("Incomplete")}
	parser := NewParser(&m.nodes, m.updateLoadStatus)
	return m, parser.Parse(jsonData, blocking)
}

// LoadStatus returns an error if not loaded fully or nil if completed loading
func (m MasterNodeList) LoadStatus() error {
	return m.loadStatus
}

// GetNodeView returns a View of the masterlist, also returning loadStatus, which will
// be nil if fully loaded. It may also return empty view with a View parsing error
func (m MasterNodeList) GetNodeView() (View, error) {
	nodes := make([]*Node, len(m.nodes))
	for index := range m.nodes {
		nodes[index] = &m.nodes[index]
	}
	view, err := NewView(nodes)
	if err != nil {
		return View{}, err
	}
	return view, m.loadStatus
}

func (m *MasterNodeList) updateLoadStatus(err error) {
	m.loadStatus = err
}
