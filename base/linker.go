package base

import "github.com/urandom/graph"

type Linker struct {
	NodeId           graph.Id
	InputConnectors  map[string]graph.Connector
	OutputConnectors map[string]graph.Connector
}

func NewLinker() *Linker {
	l := &Linker{
		InputConnectors:  make(map[string]graph.Connector),
		OutputConnectors: make(map[string]graph.Connector),
	}

	l.InputConnectors[DefaultInputConnector.Name()] = DefaultInputConnector
	l.InputConnectors[DefaultOutputConnector.Name()] = DefaultOutputConnector

	return l
}

func (l *Linker) Connect(target graph.Linker, source, sink graph.Connector) {
	if l.Connector(source.Name(), source.Type()) == nil ||
		target.Connector(sink.Name(), sink.Type()) == nil {
		return
	}

	source.Connect(target.Node(), sink)
	sink.Connect(l.Node(), source)
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

func (l Linker) Connection(source ...graph.Connector) (graph.Node, graph.Connector) {
	s := DefaultOutputConnector
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

func (l Linker) Id() graph.Id {
	return l.NodeId
}

func (l Linker) Node() graph.Node {
	return l
}
