package base

import (
	"testing"

	"github.com/urandom/graph"
)

func TestConnector(t *testing.T) {
	var c graph.Connector

	c = NewInputConnector()

	expectedStr := graph.InputName
	if c.Name() != expectedStr {
		t.Fatalf("Expected name %s, got %s\n", expectedStr, c.Name())
	}

	expectedType := graph.InputType
	if c.Type() != expectedType {
		t.Fatalf("Expected type %v, got %v\n", expectedType, c.Type())
	}

	c = NewOutputConnector()

	expectedStr = graph.OutputName
	if c.Name() != expectedStr {
		t.Fatalf("Expected name %s, got %s\n", expectedStr, c.Name())
	}

	expectedType = graph.OutputType
	if c.Type() != expectedType {
		t.Fatalf("Expected type %v, got %v\n", expectedType, c.Type())
	}

	c = NewInputConnector("aux")
	expectedStr = "aux"
	if c.Name() != expectedStr {
		t.Fatalf("Expected name %s, got %s\n", expectedStr, c.Name())
	}

	expectedType = graph.InputType
	if c.Type() != expectedType {
		t.Fatalf("Expected type %v, got %v\n", expectedType, c.Type())
	}

	c = NewOutputConnector("aux")
	expectedStr = "aux"
	if c.Name() != expectedStr {
		t.Fatalf("Expected name %s, got %s\n", expectedStr, c.Name())
	}

	expectedType = graph.OutputType
	if c.Type() != expectedType {
		t.Fatalf("Expected type %v, got %v\n", expectedType, c.Type())
	}

	if n, o := c.Target(); n != nil || o != nil {
		t.Fatalf("Unexpected target")
	}

	l := NewLinker()
	c.Connect(l, l.Connector(graph.InputName))

	if n, o := c.Target(); n != l || o != l.Connector(graph.InputName) {
		t.Fatalf("Unexpected target")
	}

	c.Disconnect()

	if n, o := c.Target(); n != nil || o != nil {
		t.Fatalf("Unexpected target")
	}

	if err := c.Connect(l, NewOutputConnector("test")); err != graph.ErrSameConnectorType {
		t.Fatalf("Expected %v, got %v\n", graph.ErrSameConnectorType, err)
	}
}
