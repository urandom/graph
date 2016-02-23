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
	if len(connectors) != 2 {
		t.Fatalf("Expected 2 connectors, got %d", len(connectors))
	}

	for _, c := range connectors {
		switch c.Name() {
		case graph.OutputName:
			target, _ := c.Target()
			if _, ok := target.Node().(saveNode); !ok {
				t.Fatalf("Unknown node type %T", target)
			}
		case "ref":
		default:
			t.Fatalf("Only %s and %s are expected\n", graph.OutputName, "ref")
		}

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
	if len(connectors) != 2 {
		t.Fatalf("Expected 2 connectors, got %d", len(connectors))
	}

	for _, c := range connectors {
		switch c.Name() {
		case graph.OutputName:
			target, ic := c.Target()
			if _, ok := target.Node().(saveNode); !ok {
				t.Fatalf("Unknown node type %T", target)
			}

			if ic.Name() != graph.InputName {
				t.Fatalf("Expected %s, got %s\n", graph.InputName, ic.Name())
			}

			c2 := roots[1].Connector(graph.OutputName, graph.OutputType)
			if c2 == nil {
				t.Fatalf("Expected a default output connector")
			}

			target, ic = c2.Target()
			if _, ok := target.Node().(saveNode); !ok {
				t.Fatalf("Unknown node type %T", target)
			}

			if ic.Name() != "dup" {
				t.Fatalf("Expected %s, got %s\n", "dup", ic.Name())
			}
		case "ref":
		default:
			t.Fatalf("Only %s and %s are expected\n", graph.OutputName, "ref")
		}

	}

}

func TestProcessJSONBranch(t *testing.T) {
	roots, err := graph.ProcessJSON(testBranch, nil)
	if err != nil {
		t.Fatalf("processing testBranch: %v", err)
	}

	connectors := roots[0].Connectors(graph.OutputType)
	if len(connectors) != 2 {
		t.Fatalf("Expected 2 connectors, got %d", len(connectors))
	}

	for _, c := range connectors {
		switch c.Name() {
		case graph.OutputName:
			target, ic := c.Target()
			if _, ok := target.Node().(saveNode); !ok {
				t.Fatalf("Unknown node type %T", target)
			}

			if ic.Name() != graph.InputName {
				t.Fatalf("Expected %s, got %s\n", graph.InputName, ic.Name())
			}
		case "ref":
			target, ic := c.Target()
			if _, ok := target.Node().(passNode); !ok {
				t.Fatalf("Unknown node type %T", target)
			}

			if ic.Name() != graph.InputName {
				t.Fatalf("Expected %s, got %s\n", graph.InputName, ic.Name())
			}

			tc := target.Connector(graph.OutputName, graph.OutputType)
			if tc == nil {
				t.Fatalf("Expected default output connector\n")
			}

			st, sc := tc.Target()
			if _, ok := st.Node().(saveNode); !ok {
				t.Fatalf("Unknown node type %T", st)
			}

			if sc.Name() != "dup" {
				t.Fatalf("Expected %s, got %s\n", "dup", sc.Name())
			}
		default:
			t.Fatalf("Only %s and %s are expected\n", graph.OutputName, "ref")
		}
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

type passNode struct {
	graph.Node
}

func init() {
	graph.RegisterLinker("Load", func(opts json.RawMessage) (graph.Linker, error) {
		var o loadOptions

		if err := json.Unmarshal([]byte(opts), &o); err != nil {
			return nil, fmt.Errorf("constructing Load: %v", err)
		}

		l := base.NewLinkerNode(loadNode{
			Node: base.NewNode(),
			opts: o,
		})

		ref := base.NewOutputConnector("ref")
		l.OutputConnectors[ref.Name()] = ref

		return l, nil
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

	graph.RegisterLinker("Pass", func(opts json.RawMessage) (graph.Linker, error) {
		return base.NewLinkerNode(passNode{
			Node: base.NewNode(),
		}), nil
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
		"Output": {
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
		"Output": {
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
		"Output": {
			"ReferenceTo": 1,
			"Input": "dup"
		}
	}
}
`
	testBranch = `
{
	"Name": "Load",
	"Options": {
		"Path": "1"
	},
	"Outputs": {
		"ref": {
			"Name": "Pass",
			"Outputs": {
				"Output": {
					"ReferenceTo": 1,
					"Input": "dup"
				}
			}
		},
		"Output": {
			"Name": "Save",
			"ReferenceId": 1,
			"Options": {
				"Path": "2"
			}
		}
	}
}
	`
)
