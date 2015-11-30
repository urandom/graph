package base

import "github.com/urandom/graph"

type Node struct {
	InputConnectors  []graph.Connector
	OutputConnectors []graph.Connector
}

func (n *Node) Connect(target graph.Node, source, sink Connector) {
}

func (n *Node) Link(target graph.Node) {
}

func (n *Node) Connection(source ...graph.Connector) graph.Node {
	s := DefaultOutputConnector
	if len(source) > 0 {
		s = source[0]
	}

	switch s.Type() {
	case graph.Input:
		for _, c := range n.InputConnectors {
			if c == s {
				return c.Target()
			}
		}
	case graph.Output:
		for _, c := range n.OutputConnectors {
			if c == s {
				return c.Target()
			}
		}
	}

	return nil
}
