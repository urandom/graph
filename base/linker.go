package base

import "github.com/urandom/graph"

// Linker provides a base implementation of graph.Linker
type Linker struct {
	// InputConnectors is a map of the input connectors using their names
	InputConnectors map[graph.ConnectorName]graph.Connector
	// OutputConnectors is a map of the output connectors using their names
	OutputConnectors map[graph.ConnectorName]graph.Connector

	// Data is the underlying Node
	Data graph.Node
}

// NewLinker creates a new linker with a node and adds the default input and
// output connectors
func NewLinker() *Linker {
	return NewLinkerNode(NewNode())
}

func NewLinkerNode(node graph.Node) *Linker {
	l := &Linker{
		Data:             node,
		InputConnectors:  make(map[graph.ConnectorName]graph.Connector),
		OutputConnectors: make(map[graph.ConnectorName]graph.Connector),
	}

	input := NewInputConnector()
	output := NewOutputConnector()
	l.InputConnectors[input.Name()] = input
	l.OutputConnectors[output.Name()] = output

	return l
}

func (l *Linker) Connect(target graph.Linker, source, sink graph.Connector) error {
	if source == nil || sink == nil {
		return graph.ErrInvalidConnector
	}

	source = l.Connector(source.Name(), source.Type())
	sink = target.Connector(sink.Name(), sink.Type())

	if source == nil || sink == nil {
		return graph.ErrInvalidConnector
	}

	if err := source.Connect(target, sink); err != nil {
		return err
	}

	if err := sink.Connect(l, source); err != nil {
		source.Disconnect()
		return err
	}

	return nil
}

func (l *Linker) Disconnect(source graph.Connector) {
	if _, tc := source.Target(); tc != nil {
		tc.Disconnect()
	}

	source.Disconnect()
}

func (l *Linker) Link(target graph.Linker) {
	l.Connect(target, l.Connector(graph.OutputName, graph.OutputType), target.Connector(graph.InputName))
}

func (l *Linker) Unlink() {
	c := l.Connector(graph.OutputName, graph.OutputType)

	if t, _ := c.Target(); t != nil {
		l.Disconnect(c)
	}
}

func (l Linker) Connector(name graph.ConnectorName, kind ...graph.ConnectorType) graph.Connector {
	t := graph.InputType
	if len(kind) > 0 {
		t = kind[0]
	}

	switch t {
	case graph.InputType:
		return l.InputConnectors[name]
	case graph.OutputType:
		return l.OutputConnectors[name]
	}

	return nil
}

func (l Linker) Connectors(kind ...graph.ConnectorType) []graph.Connector {
	t := graph.InputType
	if len(kind) > 0 {
		t = kind[0]
	}

	var connectorsOfType map[graph.ConnectorName]graph.Connector
	connectors := []graph.Connector{}

	switch t {
	case graph.InputType:
		connectorsOfType = l.InputConnectors
	case graph.OutputType:
		connectorsOfType = l.OutputConnectors
	}

	for _, v := range connectorsOfType {
		connectors = append(connectors, v)
	}

	return connectors
}

func (l Linker) Connection(source ...graph.Connector) (graph.Linker, graph.Connector) {
	s := l.InputConnectors[graph.InputName]
	if len(source) > 0 {
		s = source[0]
	}

	switch s.Type() {
	case graph.InputType:
		for _, c := range l.InputConnectors {
			if c == s {
				return c.Target()
			}
		}
	case graph.OutputType:
		for _, c := range l.OutputConnectors {
			if c == s {
				return c.Target()
			}
		}
	}

	return nil, nil
}

func (l Linker) Node() graph.Node {
	return l.Data
}
