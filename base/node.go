package base

import "github.com/urandom/graph"

// Node is a base implementation of a graph.Node
type Node struct {
	NodeId graph.Id
}

const maxUint = ^uint64(0)

var counter uint64 = 0

// NewNode creates a new node with an incremental id
func NewNode() Node {
	return Node{NodeId: nextId()}
}

func (n Node) Id() graph.Id {
	return n.NodeId
}

func nextId() graph.Id {
	id := graph.Id(counter)

	if counter == maxUint {
		counter = 0
	} else {
		counter++
	}

	return id
}
