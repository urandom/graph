package graph_test

import (
	"testing"

	"github.com/urandom/graph"
	"github.com/urandom/graph/base"
)

func TestVisitor(t *testing.T) {
	v := graph.NewVisitor()

	l1 := base.NewLinker()
	l2 := base.NewLinker()

	if v.Visited(l1.Node()) {
		t.Fatalf("L1 shouldn't have been visited yet")
	}

	if v.Visited(l2.Node()) {
		t.Fatalf("L2 shouldn't have been visited yet")
	}

	v.Add(l1.Node())

	if !v.Visited(l1.Node()) {
		t.Fatalf("L1 should have been visited")
	}

	if v.Visited(l2.Node()) {
		t.Fatalf("L2 shouldn't have been visited yet")
	}

	v.Add(l2.Node())

	if !v.Visited(l1.Node()) {
		t.Fatalf("L1 should have been visited")
	}

	if !v.Visited(l2.Node()) {
		t.Fatalf("L2 should have been visited")
	}
}
