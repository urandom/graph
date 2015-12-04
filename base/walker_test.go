package base

import (
	"math/rand"
	"testing"
	"time"

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

	expected := map[graph.Id]graph.Linker{
		linkers[0].Node().Id():  linkers[0],
		linkers[5].Node().Id():  linkers[5],
		linkers[6].Node().Id():  linkers[6],
		linkers[10].Node().Id(): linkers[10],
	}

	for _, l := range w.roots {
		if _, ok := expected[l.Node().Id()]; !ok {
			t.Fatalf("Unexpected root node: %v\n", l.Node())
		}
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

		switch n.Id() {
		case linkers[2].Node().Id():
			if !v2.Visited(linkers[1].Node()) {
				t.Fatalf("Node 2 depends on 1")
			}

			if !v2.Visited(linkers[3].Node()) {
				t.Fatalf("Node 2 depends on 3")
			}
		case linkers[9].Node().Id():
			if !v2.Visited(linkers[7].Node()) {
				t.Fatalf("Node 9 depends on 7")
			}

			if !v2.Visited(linkers[10].Node()) {
				t.Fatalf("Node 9 depends on 10")
			}
		}

		count++
		time.Sleep(time.Duration((rand.Intn(200-10) + 10)) * time.Millisecond)

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
