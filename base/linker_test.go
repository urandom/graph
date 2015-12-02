package base

import (
	"testing"

	"github.com/urandom/graph"
)

func TestLinker(t *testing.T) {
	var l1 graph.Linker = NewLinker()
	var l2 graph.Linker = NewLinker()

	if l1.Node().Id() == l2.Node().Id() {
		t.Fatalf("Nodes shouldn't be equal: %v, %v\n", l1.Node(), l2.Node())
	}

	c := l1.Connector(graph.OutputName, graph.OutputType)
	if c == nil || c.Name() != graph.OutputName || c.Type() != graph.OutputType {
		t.Fatalf("Unexpected connnector %v from %v\n", c, l1)
	}

	if n, o := c.Target(); n != nil || o != nil {
		t.Fatalf("Connector %v of %v shouldn't have a target\n", c, l1)
	}

	c = l1.Connector("aux")
	if c != nil {
		t.Fatalf("Connector %v of %v shouldn't exist\n", c, l1)
	}

	if n, o := l1.Connection(l1.Connector(graph.OutputName, graph.OutputType)); n != nil || o != nil {
		t.Fatalf("Connector %v of %v shouldn't have a target\n", l1.Connector(graph.OutputName, graph.OutputType), l1)
	}

	l1.Link(l2)

	c = l1.Connector(graph.OutputName, graph.OutputType)
	if n, o := c.Target(); n.Id() != l2.Node().Id() || o != l2.Connector(graph.InputName) {
		t.Fatalf("Connector %v of %v should have %v as a target\n", c, l1, l2)
	}

	c = l2.Connector(graph.InputName)
	if n, o := c.Target(); n.Id() != l1.Node().Id() || o != l1.Connector(graph.OutputName, graph.OutputType) {
		t.Fatalf("Connector %v of %v should have %v as a target\n", c, l2, l1)
	}

	if n, o := l2.Connection(); n.Id() != l1.Node().Id() || o != l1.Connector(graph.OutputName, graph.OutputType) {
		t.Fatalf("Connector %v of %v should have %v as a target\n", c, l2, l1)
	}

	if n, o := l2.Connection(l2.Connector(graph.OutputName, graph.OutputType)); n != nil || o != nil {
		t.Fatalf("Connection to %v via %v from %v shouldn't exist\n", n, o, l1)
	}
}
