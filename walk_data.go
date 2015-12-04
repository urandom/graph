package graph

type WalkData struct {
	Node       Node
	Connectors []Connector

	done chan struct{}
}

func NewWalkData(n Node, c []Connector, d chan struct{}) WalkData {
	return WalkData{Node: n, Connectors: c, done: d}
}

func (wd WalkData) Close() {
	close(wd.done)
}
