package base

import "github.com/urandom/graph"

// Connector is a base implementation of a graph.Connector
type Connector struct {
	targetLinker    graph.Linker
	targetConnector graph.Connector

	kind graph.ConnectorType
	name graph.ConnectorName
}

// NewInputConnector creates an input connector with the specified name. If no
// name is given, the default input name is used.
func NewInputConnector(name ...graph.ConnectorName) *Connector {
	return newConnector(graph.InputType, name...)
}

// NewOutputConnector creates an output connector with the specified name. If no
// name is given, the default output name is used.
func NewOutputConnector(name ...graph.ConnectorName) *Connector {
	n := graph.OutputName
	if len(name) > 0 {
		n = name[0]
	}
	return newConnector(graph.OutputType, n)
}

func newConnector(kind graph.ConnectorType, name ...graph.ConnectorName) *Connector {
	c := Connector{kind: kind}

	if len(name) > 0 {
		c.name = name[0]
	} else {
		if kind == graph.InputType {
			c.name = graph.InputName
		} else {
			c.name = graph.OutputName
		}
	}

	return &c
}

func (c Connector) Type() graph.ConnectorType {
	return c.kind
}

func (c Connector) Name() graph.ConnectorName {
	return c.name
}

func (c Connector) Target() (graph.Linker, graph.Connector) {
	return c.targetLinker, c.targetConnector
}

func (c *Connector) Connect(target graph.Linker, connector graph.Connector) error {
	if connector == nil {
		return graph.ErrInvalidConnector
	}

	if connector.Type() == c.Type() {
		return graph.ErrSameConnectorType
	}

	c.targetLinker = target
	c.targetConnector = connector
	return nil
}

func (c *Connector) Disconnect() {
	c.targetLinker = nil
	c.targetConnector = nil
}
