package graph

import "sync"

// Visitor keeps track of which nodes have already been visited during a walk
type Visitor struct {
	sync.RWMutex
	visited map[Id]bool
}

// NewVisitor creates a new visitor
func NewVisitor() *Visitor {
	return &Visitor{visited: make(map[Id]bool)}
}

// Add marks a node as visited. If the node has already been visited, it will
// return false. A write lock is used during this operation
func (v *Visitor) Add(node Node) bool {
	defer v.Unlock()
	v.Lock()

	if v.visited[node.Id()] {
		return false
	}

	v.visited[node.Id()] = true
	return true
}

// Visited returns whether a node has already been visited. A read lock is used
// when checking
func (v *Visitor) Visited(node Node) bool {
	defer v.RUnlock()
	v.RLock()

	return v.visited[node.Id()]
}
