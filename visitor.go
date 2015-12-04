package graph

import "sync"

type Visitor struct {
	sync.RWMutex
	visited map[Id]bool
}

func NewVisitor() Visitor {
	return Visitor{visited: make(map[Id]bool)}
}

func (v Visitor) Add(node Node) bool {
	defer v.Unlock()
	v.Lock()

	if v.visited[node.Id()] {
		return false
	}

	v.visited[node.Id()] = true
	return true
}

func (v Visitor) Visited(node Node) bool {
	defer v.RUnlock()
	v.RLock()

	return v.visited[node.Id()]
}
