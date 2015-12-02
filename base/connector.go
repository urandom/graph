package base

import "github.com/urandom/graph"

type Connector struct {
	targetNode      graph.Node
	targetConnector graph.Connector

	kind graph.ConnectorType
	name string
}

var (
	DefaultInputConnector  graph.Connector = NewInputConnector()
	DefaultOutputConnector graph.Connector = NewOutputConnector()
)

func NewInputConnector(name ...string) *Connector {
	return newConnector(graph.InputType, name...)
}

func NewOutputConnector(name ...string) *Connector {
	return newConnector(graph.OutputType, name...)
}

func newConnector(kind graph.ConnectorType, name ...string) *Connector {
	c := Connector{kind: kind}

	if len(name) > 0 {
		c.name = name[0]
	} else {
		c.name = graph.InputName
	}

	return &c
}

func (c Connector) Type() graph.ConnectorType {
	return c.kind
}

func (c Connector) Name() string {
	return c.name
}

func (c Connector) Target() (graph.Node, graph.Connector) {
	return c.targetNode, c.targetConnector
}

func (c *Connector) Connect(target graph.Node, connector graph.Connector) {
	c.targetNode = target
	c.targetConnector = connector
}

func (c *Connector) Disconnect() {
	c.targetNode = nil
	c.targetConnector = nil
}
