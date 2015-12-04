package graph_test

import (
	"fmt"
	"math/rand"

	"github.com/urandom/graph"
	"github.com/urandom/graph/base"
)

type Processor interface {
	Process(wd graph.WalkData, output chan<- int)
	Result() int
}

type RandomNumberNode struct {
	graph.Node
	result int
}

type MultiplyNode struct {
	graph.Node
	result int
}

type SummingNode struct {
	graph.Node
	result int
}

func (n *RandomNumberNode) Process(wd graph.WalkData, output chan<- int) {
	n.result = rand.Intn(50-10) + 10

	wd.Close()
}

func (n RandomNumberNode) Result() int {
	return n.result
}

func (n *MultiplyNode) Process(wd graph.WalkData, output chan<- int) {
	parent := wd.Parents[0]

	if p, ok := parent.Node.(Processor); ok {
		n.result = p.Result()*rand.Intn(10-1) + 1
	}

	wd.Close()
}

func (n MultiplyNode) Result() int {
	return n.result
}

func (n *SummingNode) Process(wd graph.WalkData, output chan<- int) {
	for _, parent := range wd.Parents {
		if p, ok := parent.Node.(Processor); ok {
			n.result += p.Result()
		}
	}

	wd.Close()
	output <- n.result
}

func (n SummingNode) Result() int {
	return n.result
}

func Example() {
	root := CreateGraph()

	walker := graph.NewWalker(root)
	data := walker.Walk()

	output := make(chan int)

	for wd := range data {
		if p, ok := wd.Node.(Processor); ok {
			go p.Process(wd, output)
		} else {
			wd.Close()
		}
	}

	select {
	case r := <-output:
		fmt.Println(r)
	}
}

func CreateGraph() graph.Linker {
	linkers := make([]graph.Linker, 4)

	for i := range linkers {
		l := base.NewLinker()

		switch i {
		case 0:
			l.Data = &RandomNumberNode{Node: l.Data}
		case 1:
			l.Data = &RandomNumberNode{Node: l.Data}
		case 2:
			l.Data = &MultiplyNode{Node: l.Data}
			l.Connect(linkers[0], l.Connector(graph.InputName), linkers[0].Connector(graph.OutputName, graph.OutputType))
		case 3:
			c := base.NewInputConnector("aux")
			l.InputConnectors[c.Name()] = c
			l.Data = &SummingNode{Node: l.Data}

			l.Connect(linkers[1], l.Connector(graph.InputName), linkers[1].Connector(graph.OutputName, graph.OutputType))
			l.Connect(linkers[2], l.Connector(c.Name()), linkers[2].Connector(graph.OutputName, graph.OutputType))
		}

		linkers[i] = l
	}

	return linkers[0]
}
