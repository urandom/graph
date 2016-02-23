package graph_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/urandom/graph"
	"github.com/urandom/graph/base"
)

func TestRegisterNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should've been a valid panic")
		}
	}()

	graph.RegisterLinker("test", nil)

	t.Fatalf("Can't register nil constructors")
}

func TestRegisterDup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should've been a valid panic")
		}
	}()

	graph.RegisterLinker("test", func(opts json.RawMessage) (graph.Linker, error) { return nil, nil })
	graph.RegisterLinker("test", func(opts json.RawMessage) (graph.Linker, error) { return nil, nil })
}

func TestRegister(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Shouldn't have been a panic")
		}
	}()

	graph.RegisterLinker("test2", func(opts json.RawMessage) (graph.Linker, error) {
		return nil, nil
	})
}

func TestProcessJSON(t *testing.T) {
	roots, err := graph.ProcessJSON(testData1, nil)
	if err != nil {
		t.Fatalf("processing testData1: %v", err)
	}

	if len(roots) != 1 {
		t.Fatalf("Expected 1 root, got %d", len(roots))
	}

	if _, ok := roots[0].Node().(loadNode); !ok {
		t.Fatalf("Unknown node type %T", roots[0].Node())
	}

	connectors := roots[0].Connectors(graph.OutputType)
	if len(connectors) != 1 {
		t.Fatalf("Expected 1 connector, got %d", len(connectors))
	}

	if connectors[0].Name() != graph.OutputName {
		t.Fatalf("Expected %s, got %s\n", graph.OutputName, connectors[0].Name())
	}

	target, _ := connectors[0].Target()
	if _, ok := target.Node().(saveNode); !ok {
		t.Fatalf("Unknown node type %T", target)
	}
}

func TestProcessJSONTwoRoots(t *testing.T) {
	roots, err := graph.ProcessJSON(testTwoRoots, nil)
	if err != nil {
		t.Fatalf("processing testData1: %v", err)
	}

	if len(roots) != 2 {
		t.Fatalf("Expected 2 roots, got %d", len(roots))
	}

	if _, ok := roots[0].Node().(loadNode); !ok {
		t.Fatalf("Unknown node type %T", roots[0].Node())
	}

	connectors := roots[0].Connectors(graph.OutputType)
	if len(connectors) != 1 {
		t.Fatalf("Expected 1 connector, got %d", len(connectors))
	}

	if connectors[0].Name() != graph.OutputName {
		t.Fatalf("Expected %s, got %s\n", graph.OutputName, connectors[0].Name())
	}

	target, c := connectors[0].Target()
	if _, ok := target.Node().(saveNode); !ok {
		t.Fatalf("Unknown node type %T", target)
	}

	if c.Name() != graph.InputName {
		t.Fatalf("Expected %s, got %s\n", graph.InputName, c.Name())
	}

	connectors = roots[1].Connectors(graph.OutputType)
	if len(connectors) != 1 {
		t.Fatalf("Expected 1 connector, got %d", len(connectors))
	}

	if connectors[0].Name() != graph.OutputName {
		t.Fatalf("Expected %s, got %s\n", graph.OutputName, connectors[0].Name())
	}

	target, c = connectors[0].Target()
	if _, ok := target.Node().(saveNode); !ok {
		t.Fatalf("Unknown node type %T", target)
	}

	if c.Name() != "dup" {
		t.Fatalf("Expected %s, got %s\n", "dup", c.Name())
	}
}

type loadNode struct {
	graph.Node
	opts loadOptions
}
type loadOptions struct {
	Path string
}

type saveNode struct {
	graph.Node
	opts saveOptions
}
type saveOptions struct {
	Path string
}

func init() {
	graph.RegisterLinker("Load", func(opts json.RawMessage) (graph.Linker, error) {
		var o loadOptions

		if err := json.Unmarshal([]byte(opts), &o); err != nil {
			return nil, fmt.Errorf("constructing Load: %v", err)
		}

		return base.NewLinkerNode(loadNode{
			Node: base.NewNode(),
			opts: o,
		}), nil
	})
	graph.RegisterLinker("Save", func(opts json.RawMessage) (graph.Linker, error) {
		var o saveOptions

		if err := json.Unmarshal([]byte(opts), &o); err != nil {
			return nil, fmt.Errorf("constructing Save: %v", err)
		}

		l := base.NewLinkerNode(saveNode{
			Node: base.NewNode(),
			opts: o,
		})

		dup := base.NewInputConnector("dup")
		l.InputConnectors[dup.Name()] = dup

		return l, nil
	})
}

const (
	testData1 = `
{
	"Name": "Load",
	"Options": {
		"Path": "1"
	},
	"Outputs": {
		"output": {
			"Name": "Save",
			"Options": {
				"Path": "2"
			}
		}
	}
}
`
	testTwoRoots = `
{
	"Name": "Load",
	"Options": {
		"Path": "1"
	},
	"Outputs": {
		"output": {
			"Name": "Save",
			"ReferenceId": 1,
			"Options": {
				"Path": "2"
			}
		}
	}
}
{
	"Name": "Load",
	"Options": {
		"Path": "2"
	},
	"Outputs": {
		"output": {
			"ReferenceTo": 1,
			"Input": "dup"
		}
	}
}
`
)
