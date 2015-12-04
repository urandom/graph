package base

import (
	"sync"

	"github.com/urandom/graph"
)

type Visitor struct {
	sync.RWMutex
	visited map[graph.Id]bool
}

func NewVisitor() Visitor {
	return Visitor{visited: make(map[graph.Id]bool)}
}

func (v Visitor) Add(node graph.Node) bool {
	defer v.Unlock()
	v.Lock()

	if v.visited[node.Id()] {
		return false
	}

	v.visited[node.Id()] = true
	return true
}

func (v Visitor) Visited(node graph.Node) bool {
	defer v.RUnlock()
	v.RLock()

	return v.visited[node.Id()]
}
