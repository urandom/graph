package graph_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/urandom/graph"
	"github.com/urandom/graph/base"
)

func TestWalker(t *testing.T) {
	linkers := setupGraph()

	w := graph.NewWalker(linkers[0])

	expectedInt := 12
	if w.Total() != expectedInt {
		t.Fatalf("Expected %v, got %v\n", expectedInt, w.Total())
	}

	expectedInt = 4
	roots := w.RootNodes()
	if len(roots) != expectedInt {
		t.Fatalf("Expected %v, got %v\n", expectedInt, len(roots))
	}

	expected := map[graph.Id]graph.Linker{
		linkers[0].Node().Id():  linkers[0],
		linkers[5].Node().Id():  linkers[5],
		linkers[6].Node().Id():  linkers[6],
		linkers[10].Node().Id(): linkers[10],
	}

	for _, n := range roots {
		if _, ok := expected[n.Id()]; !ok {
			t.Fatalf("Unexpected root node: %v\n", n)
		}
	}

	walker := w.Walk()

	v1 := graph.NewVisitor()
	v2 := graph.NewVisitor()
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
	linkers := make([]graph.Linker, 12)

	for i := range linkers {
		l := base.NewLinker()

		switch i {
		case 1:
			l.Connect(linkers[i-1], l.Connector(graph.InputName), linkers[i-1].Connector(graph.OutputName, graph.OutputType))
		case 2:
			c := base.NewInputConnector("aux")
			l.InputConnectors[c.Name()] = c

			l.Connect(linkers[i-1], l.Connector(graph.InputName), linkers[i-1].Connector(graph.OutputName, graph.OutputType))
		case 3:
			c := base.NewInputConnector("aux")
			l.InputConnectors[c.Name()] = c

			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector("aux", graph.InputType))
		case 4:
			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector(graph.InputName))
		case 5:
			l.Connect(linkers[i-1], l.Connector(graph.OutputName, graph.OutputType), linkers[i-1].Connector(graph.InputName))
		case 6:
			l.Connect(linkers[3], l.Connector(graph.OutputName, graph.OutputType), linkers[3].Connector("aux"))
		case 7:
			c := base.NewOutputConnector("dup")
			l.OutputConnectors[c.Name()] = c

			l.Connect(linkers[2], l.Connector(graph.InputName), linkers[2].Connector(graph.OutputName, graph.OutputType))
		case 8:
			l.Connect(linkers[i-1], l.Connector(graph.InputName), linkers[i-1].Connector(graph.OutputName, graph.OutputType))
		case 9:
			c := base.NewInputConnector("aux")
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
