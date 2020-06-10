package nodelist

// MasterNodeList contains the master copy of the full NodeList
type MasterNodeList struct {
	nodes        []Node
	loadComplete bool
}

// NewMasterNodeList stuff
func NewMasterNodeList(jsonData []byte, blocking bool) MasterNodeList {
	m := MasterNodeList{[]Node{}, false}
	parser := NewParser(&m.nodes, nil)
	if blocking {
		parser.WaitForComplete()
	}
	return m
}

// GetNode returns node at index. No safety so verify index within Size
// Undefined behaviour if accessed before load completed (verify with LoadComplete)
func (m MasterNodeList) GetNode(index int) *Node {
	return &m.nodes[index]
}

// Size returns size of nodelist
func (m MasterNodeList) Size() int {
	return len(m.nodes)
}

// LoadComplete defines if nodelist has been completely loaded
func (m MasterNodeList) LoadComplete() bool {
	return m.loadComplete
}

// Split seperates MasterNodeList into different NodeListViews based on seperator
// e.g. "items = kind" will find array items and split based on the value of kind in each element
// If not in items array or does not have kind, will be put into "main" group
// Can choose to wait for loading and split to complete with blocking
func (m MasterNodeList) Split(seperator string, out *map[string]View, blocking bool) {
	// goroutine to continually update out until loadcomplete
}
