package graph

type WalkData struct {
	Node    Node
	Parents []Parent

	done chan struct{}
}

type Parent struct {
	From string
	To   string
	Node Node
}

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

func (wd WalkData) Close() {
	close(wd.done)
}
