package graph

// WalkData represents the data that will be sent through the walk channel
type WalkData struct {
	// Node is the current node being visited
	Node Node
	// Parents contains the Parents of the node
	Parents []Parent

	done chan struct{}
}

// Parent is a simple representation of the connection between the node and its
// parent
type Parent struct {
	// From is the name of the parent's connector
	From ConnectorName
	// To is the name of the connector of the current node's linker
	To ConnectorName
	// Node is the parent
	Node Node
}

// NewWalkData creates a new data object. Used by the walker
func NewWalkData(n Node, conns []Connector, d chan struct{}) WalkData {
	parents := []Parent{}
	for _, c := range conns {
		if t, o := c.Target(); t != nil {
			parents = append(parents,
				Parent{From: o.Name(), To: c.Name(), Node: t.Node()})
		}
	}
	return WalkData{Node: n, Parents: parents, done: d}
}

// Close notifies the walker that any operation done using the information of
// the node is complete and it can proceed to its descendants
func (wd WalkData) Close() {
	close(wd.done)
}
