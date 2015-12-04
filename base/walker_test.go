package base

import (
	"testing"

	"github.com/urandom/graph"
)

func TestWalker(t *testing.T) {
	linkers := setupGraph()

	w := NewWalker(linkers[0])

	expectedInt := 12
	if w.Total() != expectedInt {
		t.Fatalf("Expected %v, got %v\n", expectedInt, w.Total())
	}

	expectedInt = 4
	if len(w.roots) != expectedInt {
		t.Fatalf("Expected %v, got %v\n", expectedInt, len(w.roots))
	}

	if w.roots[0].Node().Id() != linkers[0].Node().Id() {
		t.Fatalf("Root %#v doesnt match %#v\n", w.roots[0], linkers[0])
	}

	if w.roots[1].Node().Id() != linkers[5].Node().Id() {
		t.Fatalf("Root %#v doesnt match %#v\n", w.roots[1], linkers[5])
	}

	if w.roots[2].Node().Id() != linkers[6].Node().Id() {
		t.Fatalf("Root %#v doesnt match %#v\n", w.roots[2], linkers[6])
	}

	if w.roots[3].Node().Id() != linkers[10].Node().Id() {
		t.Fatalf("Root %#v doesnt match %#v\n", w.roots[3], linkers[10])
	}

	walker := w.Walk()

	v1 := NewVisitor()
	v2 := NewVisitor()
	for _, l := range linkers {
		v1.Add(l.Node())
	}

	count := 0
	for wd := range walker {
		n := wd.Node

		if !v1.Visited(n) {
			t.Fatalf("Node %#v should be from the original set\n", n)
		}

		if !v2.Add(n) {
			t.Fatalf("Node %#v should be new\n", n)
		}
		count++

		wd.Close()
	}

	expectedInt = len(linkers)
	if count != expectedInt {
		t.Fatalf("Expected %v, got %v\n", expectedInt, count)
	}
}

func setupGraph() []graph.Linker {
	counter = 0
	linkers := make([]graph.Linker, 12)

	for i := range linkers {
		l := NewLinker()

		switch i {
		case 1:
			l.Connect(linkers[i-1], l.Connector(graph.InputName), linkers[i-1].Connector(graph.OutputName, graph.OutputType))
		case 2:
			c := NewInputConnector("aux")
			l.InputConnectors[c.Name()] = c

			l.Connect(linkers[i-1], l.Connector(graph.InputName), linkers[i-1].Connector(graph.OutputName, graph.OutputType))
		case 3:
			c := NewInputConnector("aux")
			l.InputConnectors[c.Name()] = c

			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector("aux", graph.InputType))
		case 4:
			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector(graph.InputName))
		case 5:
			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector(graph.InputName))
		case 6:
			l.Connect(linkers[3], l.Connector(graph.OutputName, graph.OutputType), linkers[3].Connector("aux"))
		case 7:
			c := NewOutputConnector("dup")
			l.OutputConnectors[c.Name()] = c

			l.Connect(linkers[2], l.Connector(graph.InputName), linkers[2].Connector(graph.OutputName, graph.OutputType))
		case 8:
			l.Connect(linkers[i-1], l.Connector(graph.InputName), linkers[i-1].Connector(graph.OutputName, graph.OutputType))
		case 9:
			c := NewInputConnector("aux")
			l.InputConnectors[c.Name()] = c

			l.Connect(linkers[i-2], l.Connector("aux"), linkers[i-2].Connector("dup", graph.OutputType))
		case 10:
			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector(graph.InputName))
		case 11:
			l.Connect(linkers[i-2], l.Connector(graph.InputName), linkers[i-2].Connector(graph.OutputName, graph.OutputType))
		}

		linkers[i] = l
	}

	return linkers
}