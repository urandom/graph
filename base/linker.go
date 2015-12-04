package base

import "github.com/urandom/graph"

type Linker struct {
	InputConnectors  map[string]graph.Connector
	OutputConnectors map[string]graph.Connector

	Data graph.Node
}

func NewLinker() *Linker {
	l := &Linker{
		Data:             NewNode(),
		InputConnectors:  make(map[string]graph.Connector),
		OutputConnectors: make(map[string]graph.Connector),
	}

	input := NewInputConnector()
	output := NewOutputConnector()
	l.InputConnectors[input.Name()] = input
	l.OutputConnectors[output.Name()] = output

	return l
}

func (l *Linker) Connect(target graph.Linker, source, sink graph.Connector) {
	source = l.Connector(source.Name(), source.Type())
	sink = target.Connector(sink.Name(), sink.Type())

	if source == nil || sink == nil {
		return
	}

	source.Connect(target, sink)
	sink.Connect(l, source)
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

func (l Linker) Connector(name string, kind ...graph.ConnectorType) graph.Connector {
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

	var connectorsOfType map[string]graph.Connector
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
